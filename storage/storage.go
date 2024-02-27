package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
	"timeMate_bot/entities"
	"unicode/utf8"
)

const (
	maxRow  = 100
	someMax = 100
)

type Store struct {
	*sql.DB
}

var ErrTagNotExist = errors.New("tag not exists")

type jsonDump struct {
	Tags   []entities.Tag   `json:"tags"`
	Events []entities.Event `json:"events"`
}

func (s *Store) CreateTag(tag entities.Tag) (int64, error) {
	if s.DB == nil {
		return 0, errors.New("store не был инициализирован")
	}
	stmt := "INSERT INTO Tags (Tag, Comment) VALUES (?, ?)"

	res, err := s.DB.Exec(stmt, tag.Tag, tag.Comment)
	if err != nil {
		return 0, err
	}

	tagId, _ := res.LastInsertId()
	return tagId, nil
}

func (s *Store) CreateTags(tags []entities.Tag) (int, error) {
	// todo добавить проверки на валидность тега перед вставкой
	if s.DB == nil {
		return 0, errors.New("store не был инициализирован")
	}

	switch len(tags) {
	case 0:
		return 0, errors.New("len(tags) == 0")
	case 1:
		_, err := s.CreateTag(tags[0])
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	tampStmt := "INSERT INTO Tags (Tag, Comment) VALUES (?, ?)"

	var valueStrings strings.Builder
	valueArgs := make([]any, 0)

	valueArgs = append(valueArgs, tags[0].Tag, tags[0].Comment)

	j := 1
	for i := 1; j < len(tags); i++ {
		for ; j < len(tags) && j < i*maxRow; j++ {
			valueStrings.WriteString(", (?, ?)")
			valueArgs = append(valueArgs, tags[j].Tag, tags[j].Comment)
		}

		stmt := fmt.Sprintf("%s %s", tampStmt, valueStrings.String())
		_, err := s.DB.Exec(stmt, valueArgs...)
		if err != nil {
			return 0, err
		}

		valueStrings.Reset()
		valueArgs = nil
	}

	return len(valueArgs), nil
}

func (s *Store) GetTagId(tagName string) (int64, error) {

	// todo сделать возврат нескольких ошибок
	if tagName == "" {
		return 0, fmt.Errorf("пустой тег")
	}
	if !utf8.ValidString(tagName) {
		return 0, fmt.Errorf("недопустимая кодировка имени тега: %s - доступна только UTF-8", tagName)
	}
	// todo вынести 64 в константы
	if utf8.RuneCountInString(tagName) > 64 {
		return 0, fmt.Errorf("строка не может быть длиннее 64 символов")
	}

	stmt := "SELECT t.id FROM Tags t WHERE t.Tag = ?"

	rows, err := s.DB.Query(stmt, tagName)
	if err != nil {
		return 0, fmt.Errorf("ошибка при выполнении запроса: %s", err)
	}
	defer rows.Close()

	// todo сделать дополнительную проверку на случай если в базе случайно окажутся теги с одним именем
	var tagId int64
	for rows.Next() {
		err = rows.Scan(&tagId)
		if err != nil {
			return 0, fmt.Errorf("ошибка при сканировании строк : %s", err)
		}
	}
	if tagId == 0 {
		err = sql.ErrNoRows
	}
	return tagId, err
}

func (s *Store) GetAllTags() ([]entities.Tag, error) {
	// TODO сделать ограничение на селект (лимт) - по сути пагинация
	stmtTags := "SELECT t.id, t.tag, t.comment FROM Tags t"

	rows, err := s.DB.Query(stmtTags)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе в бд: %s", err)
	}
	defer rows.Close()

	var tag entities.Tag
	var tags []entities.Tag

	for rows.Next() {
		err = rows.Scan(&tag.ID, &tag.Tag, &tag.Comment)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк : %s", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *Store) DeleteTagById(tagId int64) error {

	if tagId <= 0 {
		return fmt.Errorf("tagId должен быть больше 0")
	}
	stmt := "DELETE FROM Tags WHERE id = ?;"
	_, err := s.DB.Exec(stmt, tagId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateEvent(event entities.Event) (int64, error) {
	if s.DB == nil {
		return 0, errors.New("store не был инициализирован")
	}

	stmt := "INSERT INTO Events (Tag, Date, Comment) VALUES ((SELECT id FROM Tags t WHERE t.Tag = ?), ?, ?)"

	res, err := s.DB.Exec(stmt, event.Tag, event.Date, event.Comment)
	if err != nil {
		var sqlErr sqlite3.Error
		errors.As(err, &sqlErr)
		if sqlErr.Code == sqlite3.ErrConstraint {
			return 0, ErrTagNotExist
		}
		// todo обработать конкретный код ошибки (искать где-то здесь sqlite3.ErrNoExtended)
		return 0, err
	}
	eventId, _ := res.LastInsertId()
	return eventId, nil
}

func (s *Store) GetLastEvent() (event entities.Event, err error) {
	//stmt := "SELECT e.Id, e.Tag, e.Date, e.Comment FROM Events e ORDER BY datetime(Date) DESC LIMIT 1"
	stmt := `
SELECT e.Id, T.Tag, e.Date, e.Comment
FROM Events e
    INNER JOIN Tags T on T.ID = e.Tag
ORDER BY datetime(e.Date) DESC LIMIT 1`

	rows, err := s.DB.Query(stmt)
	if err != nil {
		return entities.Event{}, fmt.Errorf("ошибка при выполнении запроса: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&event.ID, &event.Tag, &event.Date, &event.Comment)
		if err != nil {
			return entities.Event{}, fmt.Errorf("ошибка при сканировании строк : %s", err)
		}
	}
	return event, err
}

func (s *Store) DeleteEventById(eventId int64) error {

	if eventId <= 0 {
		return fmt.Errorf("eventId должен быть больше 0")
	}
	stmt := "DELETE FROM Events WHERE id = ?;"
	_, err := s.DB.Exec(stmt, eventId)
	if err != nil {
		return err
	}
	return nil
}

//func (s *Store) CreateEvents(events []entities.Event) (int, error) {
//	// todo переписать tempStmt - работает со старым форматом
//	if s.DB == nil {
//		return 0, errors.New("store не был инициализирован")
//	}
//	if len(events) == 0 {
//		return 0, errors.New("len(events) == 0")
//	}
//
//	tempStmt := "INSERT INTO Events (Tag, Date, Comment) VALUES (?, ?, ?)"
//
//	var valueStrings strings.Builder
//	valueArgs := make([]any, 0)
//
//	valueArgs = append(valueArgs, events[0].Tag, events[0].Date, events[0].Comment)
//
//	j := 1
//	for i := 1; j < len(events); i++ {
//		for ; j < len(events) && j < i*maxRow; j++ {
//			valueStrings.WriteString(", (?, ?, ?)")
//			valueArgs = append(valueArgs, events[j].Tag, events[j].Date, events[j].Comment)
//		}
//
//		stmt := fmt.Sprintf("%s %s", tempStmt, valueStrings.String())
//		_, err := s.DB.Exec(stmt, valueArgs...)
//		if err != nil {
//			return 0, err
//		}
//
//		valueStrings.Reset()
//		valueArgs = nil
//
//	}
//
//	return len(valueArgs), nil
//}

func (s *Store) GetAllEvents() ([]entities.Event, error) {
	stmtEvents := `
		SELECT e.Id, T.Tag, e.Date, e.Comment
		FROM Events e
    		INNER JOIN Tags T on T.ID = e.Tag`

	rows, err := s.DB.Query(stmtEvents)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе в бд: %s", err)
	}
	defer rows.Close()

	var event entities.Event
	var events []entities.Event

	for rows.Next() {
		err = rows.Scan(&event.ID, &event.Tag, &event.Date, &event.Comment)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк : %s", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Store) GetEventsByDate(startTime time.Time, endTime time.Time) ([]entities.Event, error) {
	stmt := `
SELECT e.ID, T.Tag, e.Date, e.Comment
FROM Events e
    INNER JOIN Tags T on e.Tag = T.ID
WHERE e.Date BETWEEN datetime(?) AND datetime(?);`

	procStartTime := startTime.Format("2006-01-02 15:04:05")
	procEndTime := endTime.Format("2006-01-02 15:04:05")

	rows, err := s.DB.Query(stmt, procStartTime, procEndTime)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе в бд: %s", err)
	}
	defer rows.Close()

	var event entities.Event
	var events []entities.Event

	for rows.Next() {
		err = rows.Scan(&event.ID, &event.Tag, &event.Date, &event.Comment)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк : %s", err)
		}
		events = append(events, event)
	}
	return events, nil

}

//func (s *Store) CsvImport(csvTags, csvEvents string) error {
//	tagsFD, err := os.Open(csvTags)
//	if err != nil {
//		return fmt.Errorf("не получилось открыть файл %s: ошибка: %s\n", csvTags, err)
//	}
//	defer tagsFD.Close()
//
//	eventsFD, err := os.Open(csvEvents)
//	if err != nil {
//		return fmt.Errorf("не получилось открыть файл %s: ошибка: %s\n", csvEvents, err)
//	}
//	defer eventsFD.Close()
//
//	// мапим csv на структуры entities
//	var tags []entities.Tag
//	err = gocsv.Unmarshal(tagsFD, &tags)
//	if err != nil {
//		return fmt.Errorf("не получилось анмаршалить файл %s: ошибка: %s\n", csvTags, err)
//	}
//
//	var events []entities.Event
//	err = gocsv.Unmarshal(eventsFD, &events)
//	if err != nil {
//		return fmt.Errorf("не получилось анмаршалить файл %s: ошибка: %s\n", csvEvents, err)
//	}
//
//	_, err = s.CreateTags(tags)
//	if err != nil {
//		return fmt.Errorf("не получилось созранить в базе tags: %s\n", err)
//	}
//
//	_, err = s.CreateEvents(events)
//	if err != nil {
//		return fmt.Errorf("не получилось созранить в базе events: %s\n", err)
//	}
//
//	return nil
//}
//
//func (s *Store) CsvExport(csvTags, csvEvents string) error {
//
//}

func (s *Store) Close() {
	_ = s.DB.Close()
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return &Store{}, fmt.Errorf("не получтлось открыть базу %s\n", err)
	}

	err = db.Ping()
	if err != nil {
		return &Store{}, fmt.Errorf("не получтлось пингануть базу %s\n", err)
	}

	return &Store{
		db,
	}, nil
}

//
//func (s *Store) ImportJson(jsonPath string) error {
//	//jsonFD, err := os.Open(jsonPath)
//	//if err != nil {
//	//	return fmt.Errorf("не получилось открыть json файл: %s", err)
//	//}
//	//defer jsonFD.Close()
//
//	jsonData, err := os.ReadFile(jsonPath)
//	if err != nil {
//		return fmt.Errorf("не получилось прочитать json файл: %s", err)
//	}
//
//	var dataDump jsonDump
//
//	err = json.Unmarshal(jsonData, &dataDump)
//	if err != nil {
//		return fmt.Errorf("ошибка анмаршалинка json файла: %s", err)
//	}
//
//	_, err = s.CreateTags(dataDump.Tags)
//	if err != nil {
//		return err
//	}
//
//	_, err = s.CreateEvents(dataDump.Events)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// ExportJson export all user data in .json file in directory pathToDir
//func (s *Store) ExportJson(pathToDir string) (string, error) {
//
//	var dump jsonDump
//
//	var err error
//	dump.Tags, err = s.GetAllTags()
//	if err != nil {
//		return "", err
//	}
//
//	dump.Events, err = s.GetAllEvents()
//	if err != nil {
//		return "", err
//	}
//
//	rowDump, err := json.MarshalIndent(dump, "", "  ")
//	if err != nil {
//		return "", err
//	}
//
//	fileName := fmt.Sprintf("%sdata_%s.json", pathToDir, time.Now().Format("2006-01-02_15:04:05"))
//	err = os.WriteFile(fileName, rowDump, fs.ModePerm)
//	if err != nil {
//		return "", fmt.Errorf("не получилось создать файл :%s :%s", fileName, err)
//	}
//	return fileName, nil
//}
