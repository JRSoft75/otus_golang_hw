package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	tmpDir := "/tmp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}

	// Путь к директории с исходными файлами
	srcDir := "./testdata"

	// Получаем список всех файлов в директории testdata
	files, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("Ошибка при чтении директории testdata: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Пропускаем поддиректории
		}

		srcFilePath := filepath.Join(srcDir, file.Name())
		dstFilePath := filepath.Join(tmpDir, file.Name())

		t.Run("test coping file: "+file.Name(), func(t *testing.T) {
			err := Copy(srcFilePath, dstFilePath, 0, 0)
			if err != nil {
				t.Fatalf("Ошибка при копировании файла %s: %v", srcFilePath, err)
			}

			srcFileInfo, err := os.Stat(srcFilePath)
			if err != nil {
				t.Fatalf("Ошибка при получении информации о исходном файле %s: %v", srcFilePath, err)
			}

			dstFileInfo, err := os.Stat(dstFilePath)
			if err != nil {
				t.Fatalf("Ошибка при получении информации о целевом файле %s: %v", dstFilePath, err)
			}

			if dstFileInfo.Size() != srcFileInfo.Size() {
				t.Errorf("Размеры файлов не совпадают для %s: ожидается %d, получено %d", file.Name(), srcFileInfo.Size(), dstFileInfo.Size())
			}
			assert.NoError(t, err)

		})

	}

	// Удаляем временные файлы
	for _, file := range files {
		if file.IsDir() {
			continue // Пропускаем поддиректории
		}
		dstFilePath := filepath.Join(tmpDir, file.Name())
		err := os.Remove(dstFilePath)
		if err != nil {
			t.Fatalf("Ошибка при удалении временного файла %s: %v", dstFilePath, err)
		}
	}
}
