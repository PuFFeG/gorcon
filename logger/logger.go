package logger

import (
    "log"
    "os"
    "sync"
)

var (
    logFile   *os.File
    logFileMu sync.Mutex
)

// InitLogFile инициализирует файл логов. Эта функция должна быть вызвана при запуске приложения.
func InitLogFile(path string) error {
    logFileMu.Lock()
    defer logFileMu.Unlock()

    if logFile != nil {
        // Если файл логов уже инициализирован, закрываем его перед созданием нового
        logFile.Close()
    }

    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return err
    }
    logFile = file
    return nil
}

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
func NewLogger(level LogLevel) *Logger {
    return &Logger{
        level:  level,
        logger: log.New(logFile, "", log.LstdFlags|log.Lshortfile),
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
func NewErrorLogger() *Logger {
    return NewLogger(Error)
}

// NewWarningLogger создает новый экземпляр Logger с уровнем Warning
func NewWarningLogger() *Logger {
    return NewLogger(Warning)
}

// NewInfoLogger создает новый экземпляр Logger с уровнем Info
func NewInfoLogger() *Logger {
    return NewLogger(Info)
}
