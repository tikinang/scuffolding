package api

import (
	_ "embed"
	"icbaat/pkg/shared/tikigo/web"
)

type Handler struct {
	web *web.Handler
}

func New(
	web *web.Handler,
) (r *Handler) {
	return &Handler{
		web: web,
	}
}
