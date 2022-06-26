package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Level        string
	ReportCaller bool
}

func DefaultConfig() Config {
	return Config{
		Level:        logrus.DebugLevel.String(),
		ReportCaller: false,
	}
}

type Handler struct {
	config Config
	logrus.FieldLogger
}

func New(
	config Config,
) (*Handler, error) {

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

	return &Handler{
		config:      config,
		FieldLogger: logger,
	}, nil
}

func (r *Handler) WithField(key string, value any) *Handler {
	return &Handler{
		config:      r.config,
		FieldLogger: r.FieldLogger.WithField(key, value),
	}
}

func (r *Handler) WithFields(fields logrus.Fields) *Handler {
	return &Handler{
		config:      r.config,
		FieldLogger: r.FieldLogger.WithFields(fields),
	}
}

func (r *Handler) WithError(err error) *Handler {
	return &Handler{
		config:      r.config,
		FieldLogger: r.FieldLogger.WithError(err),
	}
}
