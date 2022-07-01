package di

import (
	"icbaat/pkg/api/api"
	"icbaat/pkg/api/model"
	"icbaat/pkg/api/route"
	"icbaat/pkg/api/site"
	"icbaat/pkg/shared/httpClient"
	"icbaat/pkg/shared/skelet"
)

type Config struct {
	Logger     skelet.LoggerConfig
	HttpClient httpClient.Config
	Web        skelet.WebConfig
	Orm        skelet.OrmConfig
}

// FIXME: mpavlicek - create tikigo config
func DefaultConfig() Config {
	return Config{
		Logger:     skelet.DefaultLoggerConfig(),
		HttpClient: httpClient.DefaultConfig(),
		Web:        skelet.DefaultWebConfig(),
		Orm:        skelet.DefaultOrmConfig(),
	}
}

// FIXME: mpavlicek - can it be better?
func (r Config) GetFlavors() []any {
	return []any{
		func() skelet.LoggerConfig { return r.Logger },
		func() httpClient.Config { return r.HttpClient },
		func() skelet.WebConfig { return r.Web },
		func() skelet.OrmConfig { return r.Orm },
	}
}

type Handler struct {
	HttpClient        *httpClient.Handler
	Web               *skelet.Web
	Routes            *route.Handler
	Api               *api.Handler
	Site              *site.Handler
	Orm               *skelet.Orm
	RepositoryFactory *model.Factory
}

// FIXME: mpavlicek - only things you want and things that needs runner need to be here
func New(
	httpClient *httpClient.Handler,
	web *skelet.Web,
	route *route.Handler,
	api *api.Handler,
	site *site.Handler,
	orm *skelet.Orm,
	repositoryFactory *model.Factory,
) *Handler {
	return &Handler{
		HttpClient:        httpClient,
		Web:               web,
		Routes:            route,
		Api:               api,
		Site:              site,
		Orm:               orm,
		RepositoryFactory: repositoryFactory,
	}
}

func Providers() []any {
	return []any{
		New,
		skelet.NewLogger,
		skelet.NewWeb,
		httpClient.New,
		route.New,
		api.New,
		site.New,
		skelet.NewOrm,
		model.NewFactory,
	}
}
