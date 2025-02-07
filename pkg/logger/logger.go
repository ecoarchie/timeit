package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Interface -.
type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *slog.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string) *Logger {
	var l slog.Level

	switch strings.ToLower(level) {
	case "error":
		l = slog.LevelError
	case "warn":
		l = slog.LevelInfo
	case "info":
		l = slog.LevelInfo
	case "debug":
		l = slog.LevelDebug
	default:
		l = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level:     l,
		AddSource: true,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return &Logger{
		logger: logger,
	}
}

// Debug -.
func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug(message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info(message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn(message, args...)
}

// Error -.
func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Error(message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Error(message, args...)
	os.Exit(1)
}
