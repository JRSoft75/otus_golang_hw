package storage

import "time"

// Storage определяет интерфейс для работы с хранилищем событий.
type Storage interface {
	// AddEvent добавляет новое событие в хранилище.
	AddEvent(event *Event) error

	// UpdateEvent обновляет существующее событие в хранилище.
	UpdateEvent(event *Event) error

	// DeleteEvent удаляет событие из хранилища.
	DeleteEvent(id string) error

	// GetEventsByTimeRange возвращает список событий за указанный временной диапазон.
	GetEventsByTimeRange(startTime, endTime time.Time) ([]*Event, error)

	// GetEventByID возвращает событие по его ID.
	GetEventByID(id string) (*Event, error)
}
