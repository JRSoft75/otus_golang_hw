package storage

import "errors"

var (
	ErrDateBusy       = errors.New("данное время уже занято другим событием")
	ErrEventNotFound  = errors.New("событие не найдено")
	ErrInvalidEventID = errors.New("неверный ID события")
	ErrStartAfterEnd  = errors.New("начало события не может быть позже его окончания")
)
