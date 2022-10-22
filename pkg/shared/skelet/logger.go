package skelet

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	Level        string
	ReportCaller bool
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:        logrus.DebugLevel.String(),
		ReportCaller: false,
	}
}

type Logger struct {
	config LoggerConfig

	logrus.FieldLogger
}

func NewLogger(
	config LoggerConfig,
) (*Logger, error) {

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, errors.Wrap(err, "parse logrus log level")
	}

	logger := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			DisableColors: true,
		},
		ReportCaller: config.ReportCaller,
		Level:        level,
	}

	return &Logger{
		config:      config,
		FieldLogger: logger,
	}, nil
}

func (r *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		config:      r.config,
		FieldLogger: r.FieldLogger.WithField(key, value),
	}
}

func (r *Logger) WithFields(fields logrus.Fields) *Logger {
	return &Logger{
		config:      r.config,
		FieldLogger: r.FieldLogger.WithFields(fields),
	}
}

func (r *Logger) WithError(err error) *Logger {
	return &Logger{
		config:      r.config,
		FieldLogger: r.FieldLogger.WithError(err),
	}
}
