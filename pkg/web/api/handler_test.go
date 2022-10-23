package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"icbaat/pkg/shared/skelet"
	"icbaat/pkg/web/api"
	"icbaat/pkg/web/di"
)

func TestHandler_DoArt(t *testing.T) {
	// TODO(mpavlicek): maybe use deadline from testing
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	t.Cleanup(cancel)

	runner, app, err := skelet.AssembleRunner(
		nil,
		"test",
		new(di.Handler),
		di.DefaultConfig(),
		di.Providers()...,
	)
	assert.Nil(t, err)

	err = runner.RunBefore(ctx)
	assert.Nil(t, err)

	t.Cleanup(runner.RunAfter)

	out := RequestApi[api.DoArtIn, api.DoArtOut](t, app.Web.Handler(), http.MethodPost, "/do-art", api.DoArtIn{
		Id:   "xxx",
		Hash: "bbb",
	})
	t.Log(out.OldHash)
}

func RequestApi[In, Out any](t *testing.T, handler http.Handler, method, url string, in In) Out {
	b, err := json.Marshal(in)
	assert.Nil(t, err)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(b))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.EqualValues(t, resp.StatusCode, 200)
	var o Out
	err = json.NewDecoder(resp.Body).Decode(&o)
	assert.Nil(t, err)

	return o
}
