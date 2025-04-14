package storage

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/storage" //nolint:depguard
	"github.com/stretchr/testify/assert"
)

// MockEvent для тестирования.
type MockEvent struct {
	ID      string
	UserID  string
	Title   string
	StartAt time.Time
	EndAt   time.Time
}

// Глобальные переменные для тестов.
var (
	mockEvent1 = MockEvent{
		ID:      "1",
		Title:   "1",
		UserID:  "1",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(1 * time.Hour),
	}
	mockEvent2 = MockEvent{
		ID:      "2",
		Title:   "2",
		UserID:  "2",
		StartAt: time.Now().Add(2 * time.Hour),
		EndAt:   time.Now().Add(3 * time.Hour),
	}
)

func TestAddEvent(t *testing.T) {
	inMemoryStorage := NewInMemoryStorage()

	// Успешное добавление
	err := inMemoryStorage.AddEvent(&storage.Event{
		ID:      mockEvent1.ID,
		Title:   mockEvent1.Title,
		UserID:  mockEvent1.UserID,
		StartAt: mockEvent1.StartAt,
		EndAt:   mockEvent1.EndAt,
	})
	assert.NoError(t, err)
	assert.Contains(t, inMemoryStorage.events, mockEvent1.ID)

	// Попытка добавить событие с пересекающимся временем
	err = inMemoryStorage.AddEvent(&storage.Event{
		ID:      "3",
		Title:   mockEvent1.Title,
		UserID:  mockEvent1.UserID,
		StartAt: mockEvent1.StartAt.Add(30 * time.Minute),
		EndAt:   mockEvent1.EndAt.Add(30 * time.Minute),
	})
	assert.ErrorIs(t, err, storage.ErrDateBusy)

	// Попытка добавить событие с некорректными данными
	invalidEvent := MockEvent{
		ID:      "4",
		Title:   "4",
		UserID:  "4",
		StartAt: time.Now().Add(1 * time.Hour),
		EndAt:   time.Now(),
	}
	err = inMemoryStorage.AddEvent(&storage.Event{
		ID:      invalidEvent.ID,
		Title:   invalidEvent.Title,
		UserID:  invalidEvent.UserID,
		StartAt: invalidEvent.StartAt,
		EndAt:   invalidEvent.EndAt,
	})
	assert.Error(t, err)
}

func TestUpdateEvent(t *testing.T) {
	inMemoryStorage := NewInMemoryStorage()
	// Добавляем первое событие
	err := inMemoryStorage.AddEvent(&storage.Event{
		ID:      mockEvent1.ID,
		Title:   mockEvent1.Title,
		UserID:  mockEvent1.UserID,
		StartAt: mockEvent1.StartAt,
		EndAt:   mockEvent1.EndAt,
	})
	assert.NoError(t, err)

	// Успешное обновление
	updatedEvent := MockEvent{
		ID:      mockEvent1.ID,
		Title:   mockEvent1.Title,
		UserID:  mockEvent1.UserID,
		StartAt: time.Now().Add(5 * time.Hour),
		EndAt:   time.Now().Add(6 * time.Hour),
	}
	err = inMemoryStorage.UpdateEvent(&storage.Event{
		ID:      updatedEvent.ID,
		Title:   updatedEvent.Title,
		UserID:  updatedEvent.UserID,
		StartAt: updatedEvent.StartAt,
		EndAt:   updatedEvent.EndAt,
	})
	assert.NoError(t, err)
	event, _ := inMemoryStorage.GetEventByID("1")
	assert.Equal(t, updatedEvent.StartAt, event.StartAt)

	// Добавляем второе событие
	err = inMemoryStorage.AddEvent(&storage.Event{
		ID:      mockEvent2.ID,
		Title:   mockEvent2.Title,
		UserID:  mockEvent2.UserID,
		StartAt: mockEvent2.StartAt,
		EndAt:   mockEvent2.EndAt,
	})
	assert.NoError(t, err)

	// Попытка обновить событие с пересекающимся временем
	err = inMemoryStorage.UpdateEvent(&storage.Event{
		ID:      "1",
		Title:   mockEvent2.Title,
		UserID:  mockEvent2.UserID,
		StartAt: mockEvent2.StartAt.Add(-30 * time.Minute),
		EndAt:   mockEvent2.EndAt.Add(-30 * time.Minute),
	})
	assert.ErrorIs(t, err, storage.ErrDateBusy)
}

func TestDeleteEvent(t *testing.T) {
	inMemoryStorage := NewInMemoryStorage()
	err := inMemoryStorage.AddEvent(&storage.Event{
		ID:      mockEvent1.ID,
		Title:   mockEvent1.Title,
		UserID:  mockEvent1.UserID,
		StartAt: mockEvent1.StartAt,
		EndAt:   mockEvent1.EndAt,
	})
	assert.NoError(t, err)

	// Успешное удаление
	err = inMemoryStorage.DeleteEvent("1")
	assert.NoError(t, err)
	_, err = inMemoryStorage.GetEventByID("1")
	assert.ErrorIs(t, err, storage.ErrEventNotFound)

	// Попытка удалить несуществующее событие
	err = inMemoryStorage.DeleteEvent("nonexistent")
	assert.ErrorIs(t, err, storage.ErrEventNotFound)
}

func TestGetEventsByTimeRange(t *testing.T) {
	inMemoryStorage := NewInMemoryStorage()
	err := inMemoryStorage.AddEvent(&storage.Event{
		ID:      mockEvent1.ID,
		Title:   mockEvent1.Title,
		UserID:  mockEvent1.UserID,
		StartAt: mockEvent1.StartAt,
		EndAt:   mockEvent1.EndAt,
	})
	assert.NoError(t, err)

	err = inMemoryStorage.AddEvent(&storage.Event{
		ID:      mockEvent2.ID,
		Title:   mockEvent2.Title,
		UserID:  mockEvent2.UserID,
		StartAt: mockEvent2.StartAt,
		EndAt:   mockEvent2.EndAt,
	})
	assert.NoError(t, err)

	// Временной диапазон, охватывающий оба события
	events, err := inMemoryStorage.GetEventsByTimeRange(
		time.Now().Add(-1*time.Hour),
		time.Now().Add(4*time.Hour),
	)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	// Временной диапазон, охватывающий только одно событие
	events, err = inMemoryStorage.GetEventsByTimeRange(
		time.Now().Add(1*time.Hour),
		time.Now().Add(3*time.Hour),
	)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestValidationErrors(t *testing.T) {
	inMemoryStorage := NewInMemoryStorage()

	invalidEvent := MockEvent{
		ID:      "1",
		Title:   "1",
		UserID:  "1",
		StartAt: time.Now().Add(1 * time.Hour),
		EndAt:   time.Now(),
	}
	err := inMemoryStorage.AddEvent(&storage.Event{
		ID:      invalidEvent.ID,
		Title:   invalidEvent.Title,
		UserID:  invalidEvent.UserID,
		StartAt: invalidEvent.StartAt,
		EndAt:   invalidEvent.EndAt,
	})
	assert.ErrorIs(t, err, storage.ErrStartAfterEnd)
}

func TestTimeOverlapErrors(t *testing.T) {
	inMemoryStorage := NewInMemoryStorage()
	err := inMemoryStorage.AddEvent(&storage.Event{
		ID:      mockEvent1.ID,
		Title:   mockEvent1.Title,
		UserID:  mockEvent1.UserID,
		StartAt: mockEvent1.StartAt,
		EndAt:   mockEvent1.EndAt,
	})
	assert.NoError(t, err)

	overlappingEvent := MockEvent{
		ID:      "2",
		Title:   "2",
		UserID:  "2",
		StartAt: mockEvent1.StartAt.Add(30 * time.Minute),
		EndAt:   mockEvent1.EndAt.Add(30 * time.Minute),
	}
	err = inMemoryStorage.AddEvent(&storage.Event{
		ID:      overlappingEvent.ID,
		Title:   overlappingEvent.Title,
		UserID:  overlappingEvent.UserID,
		StartAt: overlappingEvent.StartAt,
		EndAt:   overlappingEvent.EndAt,
	})
	assert.ErrorIs(t, err, storage.ErrDateBusy)
}

func TestConcurrencySafety(t *testing.T) {
	inMemoryStorage := NewInMemoryStorage()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 100) // Задержка для имитации конкурентности
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = inMemoryStorage.AddEvent(&storage.Event{
				ID:      fmt.Sprintf("%d", id),
				Title:   fmt.Sprintf("%d", id),
				UserID:  fmt.Sprintf("%d", id),
				StartAt: time.Now().Add(time.Duration(id) * time.Hour),
				EndAt:   time.Now().Add(time.Duration(id+1) * time.Hour),
			})
		}(i)
	}
	wg.Wait()

	assert.Len(t, inMemoryStorage.events, 10)
}
