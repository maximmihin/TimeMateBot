package use_cases

import (
	"fmt"
	"time"
	"timeMate_bot/entities"
	"timeMate_bot/storage"
)

// todo сделать проверку на случай, если придет nil интерфейс

type StackUserChanges interface {
	AddTagId(tagId int64) error
	AddEventId(eventId int64) error
	Reset() error
	Undo(db *storage.Store) error
}

func AddNewEvent(db *storage.Store, lastMove StackUserChanges, tagName, date, comment string) error {

	lastEvent, err := db.GetLastEvent()
	if err != nil {
		return err
	}
	dateLastEvent, err := time.Parse("2006-01-02 15:04:05", lastEvent.Date)
	if err != nil {
		return err
	}
	dateAddingEvent, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return err
	}
	if dateLastEvent.After(dateAddingEvent) {
		return fmt.Errorf("нельзя добавлять события раньше последнего добавленного")
	}

	event, err := entities.NewEventBuilder().
		Tag(tagName).
		Date(date).
		Comment(comment).
		Build()
	if err != nil {
		return err
	}

	// пытаемся создать эвент
	eventId, err := db.CreateEvent(event)
	if err != nil {
		// если эвент не создался, из-за отсутствия тега - создаем тег и снова создаем эвент
		if err == storage.ErrTagNotExist {
			tagId, err := db.CreateTag(entities.Tag{Tag: tagName})
			if err != nil {
				return err
			}
			_ = lastMove.AddTagId(tagId)
			eventId, err = db.CreateEvent(event)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	_ = lastMove.AddEventId(eventId)

	return nil
}
