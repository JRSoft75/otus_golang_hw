package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID           string         `json:"id"`           // Уникальный идентификатор события
	Title        string         `json:"title"`        // Заголовок события
	Description  *string        `json:"description"`  // Описание события (опционально)
	UserID       string         `json:"userId"`       // ID пользователя, владельца события
	StartAt      time.Time      `json:"startAt"`      // Дата и время начала события
	EndAt        time.Time      `json:"endAt"`        // Дата и время окончания события
	NotifyBefore *time.Duration `json:"notifyBefore"` // За сколько времени высылать уведомление (опционально)
}

var ErrInvalidEvent = errors.New("invalid event data")

// Validate проверяет корректность данных события.
func (e *Event) Validate() error {
	if e.ID == "" || e.Title == "" || e.UserID == "" {
		return ErrInvalidEvent
	}
	if e.StartAt.After(e.EndAt) {
		return ErrStartAfterEnd
	}
	if e.StartAt.IsZero() || e.EndAt.IsZero() {
		return ErrInvalidEvent
	}
	return nil
}
