package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/storage" //nolint:depguard
	_ "github.com/jackc/pgx/stdlib"                                                 // justifying it
)

// SQLStorage реализует интерфейс Storage, используя базу данных SQL.
type SQLStorage struct {
	db *sql.DB
}

// NewSQLStorage создает новый экземпляр SQLStorage.
func NewSQLStorage(db *sql.DB) *SQLStorage {
	return &SQLStorage{db: db}
}

// AddEvent добавляет новое событие в базу данных.
func (s *SQLStorage) AddEvent(event *storage.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	query := `
        INSERT INTO events (id, title, description, user_id, start_at, end_at, notify_before)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	description := ""
	if event.Description != nil {
		description = *event.Description
	}

	var notifyBefore sql.NullInt64
	if event.NotifyBefore != nil {
		notifyBefore.Valid = true
		notifyBefore.Int64 = int64(*event.NotifyBefore)
	}

	_, err := s.db.Exec(query, event.ID, event.Title, description, event.UserID, event.StartAt, event.EndAt, notifyBefore)
	if err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}
	return nil
}

// UpdateEvent обновляет существующее событие в базу данных.
func (s *SQLStorage) UpdateEvent(event *storage.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	query := `
        UPDATE events
        SET title = ?, description = ?, user_id = ?, start_at = ?, end_at = ?, notify_before = ?
        WHERE id = ?
    `
	description := ""
	if event.Description != nil {
		description = *event.Description
	}

	var notifyBefore sql.NullInt64
	if event.NotifyBefore != nil {
		notifyBefore.Valid = true
		notifyBefore.Int64 = int64(*event.NotifyBefore)
	}

	_, err := s.db.Exec(query, event.Title, description, event.UserID, event.StartAt, event.EndAt, notifyBefore, event.ID)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

// DeleteEvent удаляет событие из базы данных.
func (s *SQLStorage) DeleteEvent(id string) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

// GetEventsByTimeRange возвращает список событий за указанный временной диапазон.
func (s *SQLStorage) GetEventsByTimeRange(startTime, endTime time.Time) ([]*storage.Event, error) {
	query := `
        SELECT id, title, description, user_id, start_at, end_at, notify_before
        FROM events
        WHERE start_at < ? AND end_at > ?
    `
	rows, err := s.db.Query(query, endTime, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by time range: %w", err)
	}
	defer rows.Close()

	var events []*storage.Event
	for rows.Next() {
		var event storage.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.UserID,
			&event.StartAt,
			&event.EndAt,
			&event.NotifyBefore)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

// GetEventByID возвращает событие по его ID.
func (s *SQLStorage) GetEventByID(id string) (*storage.Event, error) {
	query := `
        SELECT id, title, description, user_id, start_at, end_at, notify_before
        FROM events
        WHERE id = ?
    `
	row := s.db.QueryRow(query, id)

	var event storage.Event
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.UserID,
		&event.StartAt,
		&event.EndAt,
		&event.NotifyBefore)
	if err == sql.ErrNoRows {
		return nil, storage.ErrEventNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	return &event, nil
}
