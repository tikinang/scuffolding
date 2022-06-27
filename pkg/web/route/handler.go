package route

import (
	"context"
	"icbaat/pkg/shared/skelet"
	"icbaat/pkg/web/api"
	"icbaat/pkg/web/site"
)

type Handler struct {
	log  *skelet.Logger
	web  *skelet.Web
	api  *api.Handler
	site *site.Handler
}

func New(
	runner *skelet.Runner,
	log *skelet.Logger,
	web *skelet.Web,
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
		guest.GET("/", skelet.HtmlGuestGet(r.site.Index, "index.gohtml"))
		guest.POST("/do-art", skelet.RestGuest(r.api.DoArt))
	}
	return nil
}
