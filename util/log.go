package util

import "fmt"

const (
	LOG_LEVEL_INFO = iota
	LOG_LEVEL_DEBUG
	LOG_LEVEL_WAENING
	LOG_LEVEL_ERROR
	LOG_LEVEL_FATAL
)

type Logger struct {
	level int
}

var Log Logger

func init() {
	Log = Logger{}
}
func (l *Logger) SetLevel(level int) {
	if l.level != level {
		l.level = level
	}

}

var globalLevel int

func init() {
	globalLevel = LOG_LEVEL_ERROR
}

func SetGlobalLevel(gl int) {
	if globalLevel != gl {
		globalLevel = gl
	}

}

const (
	COLOR_RED  = ""
	COLOR_BLUE = ""
	COLOR_END  = ""
)

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level > globalLevel {
		format = "[INFO] " + format
		fmt.Printf(format, args...)
	}
}
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level > globalLevel {
		format = "[DEBUG] " + format
		fmt.Printf(format, args...)
	}
}
