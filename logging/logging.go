package logging

import (
	"fmt"
	"io"
	"os"
)

type LogLevel uint8

const (
	LevelFatal LogLevel = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelAll
)

func (x LogLevel) String() string {
	switch x {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARNING"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelAll:
		return "all"
	}
	panic("Unknown level")
}

type Config struct {
	Level LogLevel
}

type Logger struct {
	cfg Config
	w   io.Writer
}

var logger Logger

/// Sets up global logger
func Initialize(cfg Config) {
	logger = NewLogger(cfg)
}

func Get() Logger {
	return logger
}

func NewLogger(cfg Config) Logger {
	return Logger{
		cfg: cfg,
		w:   os.Stderr,
	}
}

func (x Logger) log(lvl LogLevel, format string, args ...any) bool {
	if lvl > x.cfg.Level {
		// fmt.Printf("\t- %q (%d) is greater than %q (%d)\n", lvl, lvl, x.cfg.Level, x.cfg.Level)
		return false
	}

	fmt.Fprintf(x.w, "[%s] %s\n",
		lvl, fmt.Sprintf(format, args...))
	return true
}

func (x Logger) Fatal(format string, args ...any) {
	x.log(LevelFatal, format, args...)
}

func (x Logger) Error(format string, args ...any) {
	x.log(LevelError, format, args...)
}

func (x Logger) Warning(format string, args ...any) {
	x.log(LevelWarning, format, args...)
}

func (x Logger) Info(format string, args ...any) {
	x.log(LevelInfo, format, args...)
}

func (x Logger) Debug(format string, args ...any) {
	x.log(LevelDebug, format, args...)
}
