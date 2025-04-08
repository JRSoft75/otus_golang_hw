package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoggerConfig представляет настройки логгера.
type LoggerConfig struct {
	Level string `yaml:"level"`
}

// DatabaseConfig представляет настройки базы данных.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// Config представляет основную структуру конфигурации сервиса.
type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Storage  struct {
		Type string `yaml:"type"` // "in-memory" или "sql"
	} `yaml:"storage"`
}

func LoadConfig(filePath string) (*Config, error) {
	// Проверка существования файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", filePath)
	}

	// Чтение содержимого файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Парсинг YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
