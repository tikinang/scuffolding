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

	"scuffolding/pkg/shared/skelet"
	"scuffolding/pkg/web/api"
	"scuffolding/pkg/web/di"
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

	h := app.Web.Handler()

	out := RequestApi[api.DoArtIn, api.DoArtOut](t, h, http.MethodPost, "/do-art", api.DoArtIn{
		Id:   "xxx",
		Hash: "bbb",
	})
	t.Log(out.OldHash)
}

func RequestApi[In, Out any](t *testing.T, h http.Handler, method, url string, in In) Out {
	b, err := json.Marshal(in)
	assert.Nil(t, err)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(b))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.EqualValues(t, resp.StatusCode, 200)
	var o Out
	err = json.NewDecoder(resp.Body).Decode(&o)
	assert.Nil(t, err)

	return o
}
