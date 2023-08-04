package logger

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Level int64

const (
	NoLog   Level = -1
	Info    Level = 0
	Warning Level = 1
	Debug   Level = 2
	Error   Level = 3
)

func (l Level) String() string {
	switch l {
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Debug:
		return "DEBUG"
	case Error:
		return "ERROR"
	}
	return "unknow"
}

type Logger struct {
	Scope  string
	Level  Level
	logger *log.Entry
}

func (l *Logger) Init() *Logger {
	l.logger = l.logger.WithFields(log.Fields{
		"SCOPE": l.Scope,
		"LEVEL": l.Level,
	})
	switch l.Level {
	case Info:
		l.logger.Logger.SetLevel(log.InfoLevel)
	case Warning:
		l.logger.Logger.SetLevel(log.WarnLevel)
	case Error:
		l.logger.Logger.SetLevel(log.ErrorLevel)
	case Debug:
		l.logger.Logger.SetLevel(log.DebugLevel)
	}
	return l
}
func GetLogger(level Level, scope string) *Logger {
	logger := &Logger{
		Scope:  scope,
		Level:  level,
		logger: log.NewEntry(log.New()),
	}
	logger.Init()
	return logger
}
func (l *Logger) getLogString(args ...interface{}) (string, error) {
	if len(args)%2 != 0 || len(args) == 0 {
		return "", errors.New("args not even")
	}
	res := ""
	for i := 0; i < len(args); i += 2 {
		res += fmt.Sprintf("%v=%v\t", args[i], args[i+1])
	}
	return res, nil
}

func (l *Logger) Info(args ...interface{}) {
	if res, err := l.getLogString(args...); err == nil {
		l.logger.Infoln(res)
	} else {
		l.logger.Infoln(args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if res, err := l.getLogString(args...); err == nil {
		l.logger.Debugln(res)
	} else {
		l.logger.Debugln(args...)
	}

}

func (l *Logger) Warn(args ...interface{}) {
	if res, err := l.getLogString(args...); err == nil {
		l.logger.Warnln(res)
	} else {
		l.logger.Warnln(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if res, err := l.getLogString(args...); err == nil {
		l.logger.Errorln(res)
	} else {
		l.logger.Errorln(args...)
	}

}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatalln(args...)
}
