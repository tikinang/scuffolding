package di

import (
	"icbaat/pkg/shared/httpClient"
	"icbaat/pkg/shared/tikigo/logger"
	"icbaat/pkg/shared/tikigo/web"
	"icbaat/pkg/web/api"
	"icbaat/pkg/web/route"
	"icbaat/pkg/web/site"
)

type Config struct {
	Logger     logger.Config
	HttpClient httpClient.Config
	Web        web.Config
}

func DefaultConfig() Config {
	return Config{
		Logger:     logger.DefaultConfig(),
		HttpClient: httpClient.DefaultConfig(),
		Web:        web.DefaultConfig(),
	}
}

// FIXME: mpavlicek - can it be better?
func (r Config) GetFlavors() []any {
	return []any{
		func() logger.Config { return r.Logger },
		func() httpClient.Config { return r.HttpClient },
		func() web.Config { return r.Web },
	}
}

type Handler struct {
	HttpClient *httpClient.Handler
	Web        *web.Handler
	Routes     *route.Handler
	Api        *api.Handler
	Site       *site.Handler
}

func New(
	httpClient *httpClient.Handler,
	web *web.Handler,
	route *route.Handler,
	api *api.Handler,
	site *site.Handler,
) *Handler {
	return &Handler{
		HttpClient: httpClient,
		Web:        web,
		Routes:     route,
		Api:        api,
		Site:       site,
	}
}

func Providers() []any {
	return []any{
		New,
		logger.New,
		httpClient.New,
		web.New,
		route.New,
		api.New,
		site.New,
	}
}
