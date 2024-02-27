package last_move

import (
	"fmt"
	"sync"
	"timeMate_bot/storage"
)

type LastMove struct {
	tagId   int64
	eventId int64
	wg      *sync.WaitGroup
}

func (m *LastMove) AddTagId(tagId int64) error {
	if m == nil {
		return fmt.Errorf("LastMove не был инициализирован")
	}
	if tagId <= 0 {
		return fmt.Errorf("tagId не может быть отрицательным числом: %d", tagId)
	}
	m.tagId = tagId
	return nil
}

func (m *LastMove) AddEventId(eventId int64) error {
	if m == nil {
		return fmt.Errorf("LastMove не был инициализирован")
	}
	if eventId <= 0 {
		return fmt.Errorf("eventId не может быть отрицательным числом: %d", eventId)
	}
	m.eventId = eventId
	return nil
}

func (m *LastMove) Undo(db *storage.Store) error {
	if m == nil {
		return fmt.Errorf("LastMove не был инициализирован")
	}

	// todo - тк ласт мув переедет в слой тг приложения - вызывать удаление тегов будет можно только через юз кейс
	if m.tagId != 0 {
		err := db.DeleteTagById(m.tagId)
		if err != nil {
			return err
		}
	}

	if m.eventId != 0 {
		err := db.DeleteEventById(m.eventId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *LastMove) Reset() error {
	if m == nil {
		return fmt.Errorf("LastMove не был инициализирован")
	}

	m.tagId = 0
	m.eventId = 0
	return nil
}

func NewLastMove() *LastMove {
	return &LastMove{}
}
