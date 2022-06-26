package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"html/template"
	"icbaat/pkg/shared/tikigo/logger"
	"icbaat/pkg/shared/tikigo/skelet"
	"io/fs"
	"net/http"
	"time"
)

type Config struct {
	ListenPort int32
	Mode       string
}

func DefaultConfig() Config {
	return Config{
		ListenPort: 1999,
		Mode:       gin.ReleaseMode,
	}
}

type Handler struct {
	config Config
	log    *logger.Handler

	engine *gin.Engine
	server *http.Server
}

func New(
	runner *skelet.Runner,
	config Config,
	log *logger.Handler,
) (r *Handler) {
	defer func() { runner.Register(r) }()
	return &Handler{
		config: config,
		log:    log,
	}
}

func (r *Handler) Router() gin.IRouter {
	return r.engine
}

func (r *Handler) SetHTML(
	staticPath string,
	staticFs fs.FS,
	tmpl *template.Template,
	funcMap template.FuncMap,
) {
	r.engine.StaticFS(staticPath, http.FS(staticFs))
	r.engine.SetHTMLTemplate(tmpl)
	r.engine.SetFuncMap(funcMap)
}

func (r *Handler) Before(ctx context.Context) error {

	gin.SetMode(r.config.Mode)
	r.engine = gin.New()

	r.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", r.config.ListenPort),
		Handler:           r.engine.Handler(),
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 3,
		WriteTimeout:      time.Second * 15,
	}

	return nil
}

func (r *Handler) Run(ctx context.Context, done chan<- error) {
	defer close(done)

	shutdown := make(chan error)
	go func() {
		defer close(shutdown)

		r.log.Debugf("listening on %s", r.server.Addr)
		if err := r.server.ListenAndServe(); err != nil {
			shutdown <- errors.WithStack(err)
		}
	}()
	<-ctx.Done()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := r.server.Shutdown(timeoutCtx); err != nil {
		r.log.WithError(errors.WithStack(err)).Error("shutdown error")
	}

	if err := <-shutdown; err != nil && !errors.Is(err, http.ErrServerClosed) {
		r.log.WithError(err).Error("listen and serve error")
	}
}
