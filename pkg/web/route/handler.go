package route

import (
	"context"
	"github.com/tikinang/tikigo/logger"
	"github.com/tikinang/tikigo/skelet"
	"github.com/tikinang/tikigo/web"
	"icbaat/pkg/web/api"
	"icbaat/pkg/web/site"
)

type Handler struct {
	log  *logger.Handler
	web  *web.Handler
	api  *api.Handler
	site *site.Handler
}

func New(
	runner *skelet.Runner,
	log *logger.Handler,
	web *web.Handler,
	api *api.Handler,
	site *site.Handler,
) (r *Handler) {
	// FIXME: mpavlicek - register every component to runner with reflection on application maybe?
	defer func() { runner.Register(r) }()
	return &Handler{
		log:  log,
		web:  web,
		api:  api,
		site: site,
	}
}

func (r *Handler) Before(ctx context.Context) error {
	guest := r.web.Router().Group("/")
	{
		guest.GET("/", web.HtmlGuestGet(r.site.Index, "index.gohtml"))
		guest.GET("/ping", web.RestGuest(r.api.Ping))
	}
	return nil
}
