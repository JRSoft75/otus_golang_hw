package logger

import (
	"fmt"
	"log"
	//"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New(level string) (*Logger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "debug":
		zapLevel = zapcore.DebugLevel
	default:
		return nil, fmt.Errorf("unsupported log level: %s", level)
	}

	config := zap.NewProductionConfig()
	config.Level.SetLevel(zapLevel)
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	return &Logger{zapLogger}, nil
}

// Sync синхронизирует запись логов (например, при завершении работы программы).
func (l *Logger) Sync() {
	if err := l.Logger.Sync(); err != nil {
		log.Printf("failed to sync logger: %v", err)
	}
}

// Info записывает информационное сообщение в лог.
func (l *Logger) Info(msg string) {
	l.Logger.Info(msg)
}

// Error записывает сообщение об ошибке в лог.
func (l *Logger) Error(msg string) {
	l.Logger.Error(msg)
}

// Debug записывает отладочное сообщение в лог.
func (l *Logger) Debug(msg string) {
	l.Logger.Debug(msg)
}
