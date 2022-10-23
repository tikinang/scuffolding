package skelet

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type WebConfig struct {
	ListenPort int32
	Mode       string
}

func DefaultWebConfig() WebConfig {
	return WebConfig{
		ListenPort: 1999,
		Mode:       gin.ReleaseMode,
	}
}

type Web struct {
	config WebConfig
	log    *Logger

	engine *gin.Engine
	server *http.Server
}

func NewWeb(
	runner *Runner,
	config WebConfig,
	log *Logger,
) *Web {
	r := &Web{
		config: config,
		log:    log,
	}
	runner.Register(r)
	return r
}

func (r *Web) Router() gin.IRouter {
	return r.engine
}

func (r *Web) SetHTML(
	staticPath string,
	staticFs fs.FS,
	tmpl *template.Template,
	funcMap template.FuncMap,
) {
	r.engine.StaticFS(staticPath, http.FS(staticFs))
	r.engine.SetHTMLTemplate(tmpl)
	r.engine.SetFuncMap(funcMap)
}

func (r *Web) Before(_ context.Context) error {
	gin.SetMode(r.config.Mode)
	r.engine = gin.New()
	r.engine.Use(gin.Logger())
	r.engine.Use(gin.Recovery())
	r.engine.Use(cors.Default())

	r.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", r.config.ListenPort),
		Handler:           r.engine.Handler(),
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 3,
		WriteTimeout:      time.Second * 15,
	}

	return nil
}

func (r *Web) Run(ctx context.Context, done chan<- error) {
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

func (r *Web) Handler() http.Handler {
	return r.engine
}
