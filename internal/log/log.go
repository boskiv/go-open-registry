package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var Logger = logrus.New()

// Fields wraps logrus.Fields, which is a map[string]interface{}
type Fields logrus.Fields

// SetLogLevel setting minimal log level
func SetLogLevel(level logrus.Level) {
	Logger.Level = level
}

// SetLogFormatter setting output formatter
func SetLogFormatter(formatter logrus.Formatter) {
	Logger.Formatter = formatter
}

// Debug logs a message at level Debug on the standard Logger.
func Debug(args ...interface{}) {
	if Logger.Level >= logrus.DebugLevel {
		entry := Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args...)
	}
}

// DebugWithFields logs a message with fields at level Debug on the standard Logger.
func DebugWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.DebugLevel {
		entry := Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(l)
	}
}

// Info logs a message at level Info on the standard Logger.
func Info(args ...interface{}) {
	if Logger.Level >= logrus.InfoLevel {
		entry := Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// InfoWithFields logs a message with fields at level Info on the standard Logger.
func InfoWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.InfoLevel {
		entry := Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(l)
	}
}

// Warn logs a message at level Warn on the standard Logger.
func Warn(args ...interface{}) {
	if Logger.Level >= logrus.WarnLevel {
		entry := Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

// WarnWithFields logs a message with fields at level Warn on the standard Logger.
func WarnWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.WarnLevel {
		entry := Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(l)
	}
}

// Error logs a message at level Error on the standard Logger.
func Error(args ...interface{}) {
	if Logger.Level >= logrus.ErrorLevel {
		entry := Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// ErrorWithFields logs a message with fields at level Error on the standard Logger.
func ErrorWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.ErrorLevel {
		entry := Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(l)
	}
}

// Fatal logs a message at level Fatal on the standard Logger.
func Fatal(args ...interface{}) {
	if Logger.Level >= logrus.FatalLevel {
		entry := Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

// FatalWithFields logs a message with fields at level Fatal on the standard Logger.
func FatalWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.FatalLevel {
		entry := Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(l)
	}
}

// Panic logs a message at level Panic on the standard Logger
func Panic(args ...interface{}) {
	entry := Logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo(2)
	entry.Panic(args...)
}

// PanicWithFields logs a message with fields at level Panic on the standard Logger.
func PanicWithFields(l interface{}, f Fields) {
	entry := Logger.WithFields(logrus.Fields(f))
	entry.Data["file"] = fileInfo(2)
	entry.Panic(l)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
