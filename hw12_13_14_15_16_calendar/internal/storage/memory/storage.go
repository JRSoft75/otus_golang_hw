package storage

import (
	"sync"
	"time"

	"github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/storage" //nolint:depguard
)

// InMemoryStorage реализует интерфейс Storage, храня события в памяти.
type InMemoryStorage struct {
	mu     sync.RWMutex
	events map[string]*storage.Event
}

// NewInMemoryStorage создает новый экземпляр InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		events: make(map[string]*storage.Event),
	}
}

// AddEvent добавляет новое событие в хранилище.
func (s *InMemoryStorage) AddEvent(event *storage.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isTimeOverlapping(event.StartAt, event.EndAt) {
		return storage.ErrDateBusy
	}

	s.events[event.ID] = event
	return nil
}

// UpdateEvent обновляет существующее событие в хранилище.
func (s *InMemoryStorage) UpdateEvent(event *storage.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return storage.ErrEventNotFound
	}

	if s.isTimeOverlapping(event.StartAt, event.EndAt, event.ID) {
		return storage.ErrDateBusy
	}

	s.events[event.ID] = event
	return nil
}

// DeleteEvent удаляет событие из хранилища.
func (s *InMemoryStorage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

// GetEventsByTimeRange возвращает список событий за указанный временной диапазон.
func (s *InMemoryStorage) GetEventsByTimeRange(startTime, endTime time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []*storage.Event
	for _, event := range s.events {
		if event.StartAt.Before(endTime) && event.EndAt.After(startTime) {
			events = append(events, event)
		}
	}
	return events, nil
}

// GetEventByID возвращает событие по его ID.
func (s *InMemoryStorage) GetEventByID(id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return nil, storage.ErrEventNotFound
	}
	return event, nil
}

// Проверяет, пересекается ли указанное время с другими событиями.
func (s *InMemoryStorage) isTimeOverlapping(start, end time.Time, excludeID ...string) bool {
	for id, event := range s.events {
		if len(excludeID) > 0 && id == excludeID[0] {
			continue
		}
		if !(end.Before(event.StartAt) || start.After(event.EndAt)) {
			return true
		}
	}
	return false
}
