package entities

import (
	"fmt"
	"time"
	"unicode/utf8"
)

type Event struct {
	ID      int64  `json:"id" csv:"ID"`
	Tag     string `json:"tag" csv:"Tag"`
	Date    string `json:"date" csv:"Date"`
	Comment string `json:"comment,omitempty" csv:"Comment,omitempty"`
}

type EventBuilder struct {
	id      int64
	tag     string
	date    string
	comment string
}

type EventBuilderError struct {
	hasError      bool
	buildErrors   []error
	tagErrors     []error
	dataErrors    []error
	commentErrors []error
}

func (EventBuilderError) Error() string {
	return fmt.Sprintf("лол, моя кастомная ошибка")
}

func NewEventBuilder() *EventBuilder {
	return &EventBuilder{}
}

func (eb *EventBuilder) Tag(tag string) *EventBuilder {
	eb.tag = tag
	return eb
}

func (eb *EventBuilder) Date(date string) *EventBuilder {
	eb.date = date
	return eb
}

func (eb *EventBuilder) Comment(comment string) *EventBuilder {
	eb.comment = comment
	return eb
}

func (eb *EventBuilder) Build() (Event, error) {
	var err EventBuilderError

	// todo добавить проверку тега
	err.dataErrors = isTimeValid(eb.date)
	err.commentErrors = isCommentValid(eb.comment)

	// todo может надо поменять местами? проверки на билд то так и не было
	// хотя, можно вместо этого буля просто добавить ошибку в buildErrors
	if err.tagErrors != nil || err.dataErrors != nil || err.commentErrors != nil {
		err.hasError = true
		// todo  ну и ретерн тут же
	}

	return Event{
		Tag:     eb.tag,
		Date:    eb.date,
		Comment: eb.comment,
	}, nil

}

func isTimeValid(potentialTime string) []error {
	var errs []error

	if potentialTime == "" {
		return []error{fmt.Errorf("дата должна быть обязательно указана")}
	}

	// todo добавить возможность добавлять событие с времнем и датой
	_, err := time.Parse("15:04", potentialTime)
	if err != nil {
		errs = append(errs, fmt.Errorf("не верный формат времени: %s: %s", potentialTime, err))
	}
	return errs
}

func isCommentValid(potentialComment string) []error {
	var errs []error

	if potentialComment == "" {
		return nil
	}

	if !utf8.ValidString(potentialComment) {
		errs = append(errs, fmt.Errorf("строка может быть только в кодировке UTF-8"))
	}
	// todo вынести 256 в константы
	if utf8.RuneCountInString(potentialComment) > 256 {
		errs = append(errs, fmt.Errorf("строка не может быть длиннее 256 символов"))
	}
	return nil
}
