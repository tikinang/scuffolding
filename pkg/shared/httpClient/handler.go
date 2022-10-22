package httpClient

import (
	"icbaat/pkg/shared/skelet"
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
) *Handler {
	r := &Handler{
		config: config,
	}
	runner.Register(r)
	return r
}
