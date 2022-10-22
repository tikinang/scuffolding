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
) *Handler {
	// TODO(mpavlicek): register every component to runner with reflection on application maybe?
	r := &Handler{
		log:  log,
		web:  web,
		api:  api,
		site: site,
	}
	runner.Register(r)
	return r
}

func (r *Handler) Before(ctx context.Context) error {
	guest := r.web.Router().Group("/")
	{
		guest.GET("/", skelet.HtmlGuestGet(r.site.Index, "index.gohtml"))
		guest.POST("/do-art", skelet.RestGuest(r.api.DoArt))
	}
	return nil
}
