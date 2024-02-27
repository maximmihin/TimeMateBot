package storage

import (
	"fmt"
	"testing"
	"time"
	"timeMate_bot/entities"
)

const (
	rootPath = "/Users/gradagas/Desktop/timeMate_bot/"

	dbPath = "/Users/gradagas/Desktop/timeMate_bot/my.db"

	tagsCsv   = "/Users/gradagas/Desktop/timeMate_bot/storage/test_data/csvs/tags.csv"
	eventsCsv = "/Users/gradagas/Desktop/timeMate_bot/storage/test_data/csvs/events.csv"

	dataJson = "/Users/gradagas/Desktop/timeMate_bot/storage/test_data/json/data.json"
)

//func testCsvImport(t *testing.T) {
//
//	db, err := New(dbPath)
//	if err != nil {
//		t.Errorf("%s\n", err)
//	}
//
//	err = db.CsvImport(tagsCsv, eventsCsv)
//	if err != nil {
//		t.Errorf("%s\n", err)
//	}
//}
//
//func TestStore_ImportJson(t *testing.T) {
//	db, err := New(dbPath)
//	if err != nil {
//		t.Errorf("%s\n", err)
//	}
//
//	err = db.ImportJson(dataJson)
//	if err != nil {
//		t.Errorf("%s\n", err)
//	}
//}
//
//func TestStore_ExportJson(t *testing.T) {
//	db, err := New(dbPath)
//	if err != nil {
//		t.Errorf("%s\n", err)
//	}
//
//	_, err = db.ExportJson(rootPath + "storage/test_data/")
//	if err != nil {
//		t.Errorf("%s\n", err)
//	}
//}

func TestInsertEvent(t *testing.T) {
	db, err := New(dbPath)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	eventId, err := db.CreateEvent(entities.Event{
		Tag:  "сон",
		Date: "2024-02-24 00:48:07",
	})
	fmt.Println(eventId)

	eventId, err = db.CreateEvent(entities.Event{
		Tag:  "нереальный тег",
		Date: "2024-02-25 00:48:07",
	})
	if err == ErrTagNotExist {
		fmt.Println(err)

		fmt.Println("круто!")

	}
	fmt.Println(err)
	//if errors.Is(err, sqlite3.ErrNoExtended()){}}
}

func TestStore_GetLastEvent(t *testing.T) {
	db, err := New(dbPath)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	event, err := db.GetLastEvent()
	if err != nil {
		t.Errorf("%s\n", err)
	}
	fmt.Println(event)

}

func TestStore_GetEventsByDate(t *testing.T) {
	db, err := New(dbPath)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	t1, err1 := time.Parse("2006-01-02 15:04:05", "2024-02-13 01:23:33")
	if err1 != nil {
		return
	}
	t2, err2 := time.Parse("2006-01-02 15:04:05", "2024-02-15 21:38:58")
	if err2 != nil {
		return
	}

	events, err := db.GetEventsByDate(t1, t2)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	fmt.Println(events)

}
