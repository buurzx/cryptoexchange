package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	client *zap.SugaredLogger
}

func New(level string) (*Logger, error) {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("log level %w", err)
	}
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed build logger %w", err)
	}
	defer logger.Sync() // flushes buffer, if any

	sugar := logger.Sugar()

	return &Logger{
		client: sugar,
	}, nil
}

func (l *Logger) Debug(msg ...interface{}) {
	l.client.Debug(msg...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.client.Debugf(format, args...)
}

func (l *Logger) Info(msg ...interface{}) {
	l.client.Info(msg...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.client.Infof(format, args...)
}

func (l *Logger) Warn(msg ...interface{}) {
	l.client.Warn(msg...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.client.Warnf(format, args...)
}

func (l *Logger) Error(msg ...interface{}) {
	l.client.Error(msg...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.client.Errorf(format, args...)
}
