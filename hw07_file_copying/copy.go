package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Проверяем, что указаны пути к файлам
	if fromPath == "" || toPath == "" {
		return fmt.Errorf("не указаны пути к исходному и целевому файлам")
	}

	// Получаем информацию о исходном файле
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("ошибка при открытии исходного файла: %w", err)
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			fmt.Printf("Ошибка при закрытии исходного файла: %v", err)
		}
	}()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("ошибка при получении информации о файле: %w", err)
	}

	// Проверяем, что offset не превышает размер файла
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Проверяем, является ли файл специальным
	if fileInfo.Mode()&os.ModeNamedPipe != 0 || fileInfo.Mode()&os.ModeDevice != 0 {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	// Устанавливаем смещение
	if _, err := srcFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("ошибка при установке смещения: %w", err)
	}

	// Определяем количество байт для копирования
	bytesToCopy := fileInfo.Size() - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	// Открываем целевой файл для записи
	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("ошибка при создании целевого файла: %w", err)
	}
	defer func() {
		err := dstFile.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии целевого файла: %v", err)
		}
	}()

	// Копируем данные с отображением прогресса
	buffer := make([]byte, 4096)
	var totalCopied int64
	progress := pb.Full.Start64(limit)
	progress.Start()
	for {
		readedBytes, err := srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("ошибка при чтении из исходного файла: %w", err)
		}
		if readedBytes == 0 {
			break // Достигнут конец файла
		}

		// Проверяем, не превышает ли количество скопированных байт лимит
		if totalCopied+int64(readedBytes) > bytesToCopy {
			readedBytes = int(bytesToCopy - totalCopied) // Ограничиваем количество копируемых байт
		}

		// Записываем данные в целевой файл
		if _, err := dstFile.Write(buffer[:readedBytes]); err != nil {
			return fmt.Errorf("ошибка при записи в целевой файл: %w", err)
		}

		totalCopied += int64(readedBytes)

		// Выводим прогресс
		progress.SetCurrent(totalCopied)
		time.Sleep(time.Millisecond * 1000)
		if totalCopied >= bytesToCopy {
			break
		}
	}
	progress.Finish()

	return nil
}
