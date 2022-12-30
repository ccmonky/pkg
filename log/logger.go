package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

// Logger log interface
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// NopLogger is a logger which discards all logs.
type NopLogger struct{}

func (NopLogger) Error(msg string, keysAndValues ...interface{}) {}
func (NopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (NopLogger) Debug(msg string, keysAndValues ...interface{}) {}

// StdLoggerAdapter is a logger implementation, which sends al logs to provided standard output.
type StdLogger struct {
	*log.Logger
	prefix string
	out    io.Writer
	level  Level
}

type StdLoggerOption func(*StdLogger)

func WithOut(out io.Writer) StdLoggerOption {
	return func(logger *StdLogger) {
		logger.out = out
	}
}

func WithPrefix(prefix string) StdLoggerOption {
	return func(logger *StdLogger) {
		logger.prefix = prefix
	}
}

func WithLevel(level Level) StdLoggerOption {
	return func(logger *StdLogger) {
		logger.level = level
	}
}

// NewStdLogger creates StdLoggerAdapter which sends al logs to stderr.
func NewStdLogger(opts ...StdLoggerOption) Logger {
	stdLogger := &StdLogger{
		prefix: "[ccmlog] ",
		out:    os.Stderr,
		level:  InfoLevel,
	}
	for _, opt := range opts {
		opt(stdLogger)
	}
	stdLogger.Logger = log.New(stdLogger.out, stdLogger.prefix, log.LstdFlags|log.Lshortfile)
	return stdLogger
}

func (logger *StdLogger) Debug(msg string, keysAndValues ...interface{}) {
	if logger.level <= DebugLevel {
		logger.Printf("debug: %s, %v", msg, keysAndValues)
	}
}

func (logger *StdLogger) Info(msg string, keysAndValues ...interface{}) {
	if logger.level <= InfoLevel {
		logger.Printf("info: %s, %v", msg, keysAndValues)
	}
}

func (logger *StdLogger) Error(msg string, keysAndValues ...interface{}) {
	if logger.level <= ErrorLevel {
		logger.Printf("error: %s, %v", msg, keysAndValues)
	}
}

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production.
	DebugLevel Level = -1
	// InfoLevel is the default logging priority.
	InfoLevel Level = 0
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel Level = 2
)

// LoadDefault load default logger for specified key
func LoadDefault(ctx context.Context, key any) (Logger, error) {
	switch key.(string) {
	case "", "std":
		return NewStdLogger(), nil
	case "nop":
		return &NopLogger{}, nil
	}
	return nil, fmt.Errorf("no default logger for %v", key)
}

// Interface guard
var (
	_ Logger = (*NopLogger)(nil)
	_ Logger = (*StdLogger)(nil)
)
