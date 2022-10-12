package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type (
	LogLevel int
	LogType  int
)

const (
	LogFatal = LogType(0x1)
	LogError = LogType(0x2)
	LogWarn  = LogType(0x4)
	LogInfo  = LogType(0x8)
	LogDebug = LogType(0x10)
)

const (
	LogLevelNone  = LogLevel(0x0)
	LogLevelFatal = LogLevelNone | LogLevel(LogFatal)
	LogLevelError = LogLevelFatal | LogLevel(LogError)
	LogLevelWarn  = LogLevelError | LogLevel(LogWarn)
	LogLevelInfo  = LogLevelWarn | LogLevel(LogInfo)
	LogLevelDebug = LogLevelInfo | LogLevel(LogDebug)
	LogLevelAll   = LogLevelDebug
)

type Logger struct {
	_log         *log.Logger
	level        LogLevel
	highlighting bool
}

var _log = New()

func init() {
	setFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func New() *Logger {
	return NewLogger(os.Stderr, "")
}

func NewLogger(w io.Writer, prefix string) *Logger {
	var level LogLevel
	if l := os.Getenv("LOG_LEVEL"); len(l) != 0 {
		level = stringToLogLevel(os.Getenv("LOG_LEVEL"))
	} else {
		level = LogLevelInfo
	}

	return &Logger{_log: log.New(w, prefix, log.LstdFlags), level: level, highlighting: true}
}

func Fatal(format string, v ...interface{}) {
	_log.fatal(format, v...)
}

func Error(format string, v ...interface{}) {
	_log.error(format, v...)
}

func Warn(format string, v ...interface{}) {
	_log.warn(format, v...)
}

func Debug(format string, v ...interface{}) {
	_log.debug(format, v...)
}

func Info(format string, v ...interface{}) {
	_log.info(format, v...)
}

func setFlags(flags int) {
	_log._log.SetFlags(flags)
}

func stringToLogLevel(level string) LogLevel {
	switch level {
	case "fatal":
		return LogLevelFatal
	case "error":
		return LogLevelError
	case "warn":
		return LogLevelWarn
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	}
	return LogLevelAll
}

func (l *Logger) fatal(format string, v ...interface{}) {
	l.logf(LogFatal, format, v...)
	os.Exit(1)
}

func (l *Logger) error(format string, v ...interface{}) {
	l.logf(LogError, format, v...)
}

func (l *Logger) warn(format string, v ...interface{}) {
	l.logf(LogWarn, format, v...)
}

func (l *Logger) debug(format string, v ...interface{}) {
	l.logf(LogDebug, format, v...)
}

func (l *Logger) info(format string, v ...interface{}) {
	l.logf(LogInfo, format, v...)
}

func (l *Logger) logf(t LogType, format string, v ...interface{}) {
	if l.level|LogLevel(t) != l.level {
		return
	}

	logStr, logColor := logTypeToString(t)
	var s string
	if l.highlighting {
		s = "\033" + logColor + "m[" + logStr + "]" + fmt.Sprintf(format, v...) + "\033[0m"
	} else {
		s = "[" + logStr + "]" + fmt.Sprintf(format, v...)
	}
	_ = l._log.Output(4, s)
}

func logTypeToString(t LogType) (string, string) {
	switch t {
	case LogFatal:
		return "fatal", "[0;31"
	case LogError:
		return "error", "[0;31"
	case LogWarn:
		return "warn", "[0;33"
	case LogDebug:
		return "debug", "[0;36"
	case LogInfo:
		return "info", "[0;37"
	}
	return "unknown", "[0;37"
}
