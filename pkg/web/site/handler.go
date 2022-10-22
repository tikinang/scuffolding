package site

import (
	"context"
	"embed"
	"html/template"
	"icbaat/pkg/shared/skelet"
	"io/fs"

	"github.com/pkg/errors"
)

type Handler struct {
	web *skelet.Web
}

func New(
	runner *skelet.Runner,
	web *skelet.Web,
) *Handler {
	r := &Handler{
		web: web,
	}
	runner.Register(r)
	return r
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
