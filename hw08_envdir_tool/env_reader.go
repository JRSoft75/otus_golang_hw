package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	newEnv := make(Environment)
	// Получаем список переменных окружения
	envVars := os.Environ()

	// Создаем мапу для хранения переменных окружения
	envMap := make(map[string]string)

	// Заполняем мапу переменными окружения
	for _, envVar := range envVars {
		// Разделяем строку на ключ и значение
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			envMap[key] = value
		}
	}

	// Получаем список всех файлов в директории testdata
	files, err := os.ReadDir(dir)
	if err != nil {
		return newEnv, fmt.Errorf("ошибка при чтении директории %s: %w", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Пропускаем поддиректории
		}
		if strings.Contains(file.Name(), "=") {
			continue // Пропускаем файл с = в имени
		}

		envValue, _ := getEnv(dir, file.Name())
		if envValue.NeedRemove {
			delete(envMap, file.Name())
		} else {
			envMap[file.Name()] = sanitizeValue(envValue.Value)
		}
	}

	for key, value := range envMap {
		newEnv[key] = EnvValue{
			Value:      value,
			NeedRemove: false,
		}
	}

	return newEnv, err
}

func sanitizeValue(value string) string {
	// Удаляем пробелы и табуляцию в конце строки
	sanitized := strings.TrimRight(value, " \t")

	// Заменяем терминальные нули (0x00) на перевод строки (\n)
	sanitized = strings.ReplaceAll(sanitized, "\000", "\n")
	return sanitized
}

func getEnv(dir, fileName string) (EnvValue, error) {
	result := EnvValue{
		Value:      "",
		NeedRemove: false,
	}
	srcFilePath := filepath.Join(dir, fileName)
	// Открываем файл
	file, err := os.Open(srcFilePath)
	if err != nil {
		return result, fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer func() {
		err := file.Close() // Закрываем файл в конце
		if err != nil {
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return result, fmt.Errorf("ошибка при получении информации о файле %s: %w", srcFilePath, err)
	}

	result.Value = fileName
	if fileInfo.Size() == 0 {
		result.NeedRemove = true
		return result, nil
	}

	// Создаем новый сканер для чтения файла
	scanner := bufio.NewScanner(file)

	// Читаем первую строку
	if scanner.Scan() {
		firstLine := scanner.Text() // Получаем текст первой строки
		result.Value = firstLine
		result.NeedRemove = false
		return result, nil
	} else if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("ошибка при чтении файла %s: %w", srcFilePath, err)
	}

	return result, nil
}
