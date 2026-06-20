package goutils

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type LoggerHook struct {
	logrus.Hook
}

func (h *LoggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type LoggerInterface interface {
	SetLevel(level string)
	GetLevel() logrus.Level
	SetFormatter(formatter logrus.Formatter)
	SetReportCaller(reportCaller bool)
	SetOutput(output logrus.Hook)
	AddHook(hook logrus.Hook)
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
}

func NewLogger() *Logger {
	logger := logrus.StandardLogger()
	return &Logger{logger}
}

func (l *Logger) SetLevel(level string) {
	if level, err := logrus.ParseLevel(level); err != nil {
		l.Logger.Errorf("Invalid %q provided for log level", level)
		l.Logger.SetLevel(logrus.InfoLevel)
	} else {
		l.Logger.SetLevel(level)
	}
}

func (l *Logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.SetFormatter(formatter)
}

func (l *Logger) SetReportCaller(reportCaller bool) {
	l.Logger.SetReportCaller(reportCaller)
}

func (l *Logger) AddHook(hook LoggerHook) {
	l.Logger.AddHook(hook.Hook)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}
