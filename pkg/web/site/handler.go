package site

import (
	"context"
	"embed"
	"github.com/pkg/errors"
	"html/template"
	"icbaat/pkg/shared/tikigo/skelet"
	"icbaat/pkg/shared/tikigo/web"
	"io/fs"
)

type Handler struct {
	web *web.Handler
}

func New(
	runner *skelet.Runner,
	web *web.Handler,
) (r *Handler) {
	defer func() { runner.Register(r) }()
	return &Handler{
		web: web,
	}
}

//go:embed template public/asset
var content embed.FS

func (r *Handler) Before(ctx context.Context) error {
	tmpl, err := template.New("web").ParseFS(content, "template/*.gohtml")
	if err != nil {
		return errors.Wrap(err, "load templates from templates")
	}
	assets, err := fs.Sub(content, "public")
	if err != nil {
		return errors.WithStack(err)
	}
	r.web.SetHTML("/public", assets, tmpl, nil)
	return nil
}
