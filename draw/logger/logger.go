package logger

import (
    "log"
    "os"
)

// LogLevel определяет уровни логгирования
type LogLevel int

// Определение уровней логгирования
const (
    Error LogLevel = iota
    Warning
    Info
)

// Уровень логгирования в строковом формате
var logLevelStrings = [...]string{"ERROR", "WARNING", "INFO"}

// Logger определяет наш кастомный логгер
type Logger struct {
    level  LogLevel
    logger *log.Logger
}

// NewLogger создает новый экземпляр Logger
func NewLogger(level LogLevel, file *os.File) *Logger {
    return &Logger{
        level:  level,
        logger: log.New(file, "", log.LstdFlags|log.Lshortfile),
    }
}

// Error записывает сообщение в лог с уровнем Error
func (l *Logger) Error(format string, args ...interface{}) {
    l.Log(Error, format, args...)
}

// Warning записывает сообщение в лог с уровнем Warning
func (l *Logger) Warning(format string, args ...interface{}) {
    l.Log(Warning, format, args...)
}

// Info записывает сообщение в лог с уровнем Info
func (l *Logger) Info(format string, args ...interface{}) {
    l.Log(Info, format, args...)
}

// Log записывает сообщение в лог с уровнем указанного типа
func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
    if level <= l.level {
        l.logger.Printf("[%s] "+format, append([]interface{}{logLevelStrings[level]}, args...)...)
    }
}

// NewErrorLogger создает новый экземпляр Logger с уровнем Error
func NewErrorLogger(file *os.File) *Logger {
    return NewLogger(Error, file)
}

// NewWarningLogger создает новый экземпляр Logger с уровнем Warning
func NewWarningLogger(file *os.File) *Logger {
    return NewLogger(Warning, file)
}

// NewInfoLogger создает новый экземпляр Logger с уровнем Info
func NewInfoLogger(file *os.File) *Logger {
    return NewLogger(Info, file)
}

