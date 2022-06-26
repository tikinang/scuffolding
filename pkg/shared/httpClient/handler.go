package httpClient

import (
	"icbaat/pkg/shared/tikigo/skelet"
	"time"
)

type Config struct {
	Url     string
	Timeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		Url:     "mysql://haha",
		Timeout: time.Second * 5,
	}
}

type Handler struct {
	config Config
}

func New(
	runner *skelet.Runner,
	config Config,
) (r *Handler) {
	defer func() { runner.Register(r) }()
	return &Handler{
		config: config,
	}
}
