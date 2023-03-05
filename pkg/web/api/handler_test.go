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

	runner, skelet, err := skelet.AssembleRunner[di.Di](
		nil,
		"test",
		di.DefaultConfig(),
		di.Providers()...,
	)
	assert.Nil(t, err)

	err = runner.RunBefore(ctx)
	assert.Nil(t, err)

	t.Cleanup(runner.RunAfter)

	h := skelet.Web.Handler()

	out := RequestApi[api.DoArtOut](t, h, http.MethodPost, "/do-art", api.DoArtIn{
		Id:   "xxx",
		Hash: "bbb",
	})
	t.Log(out.OldHash)
}

func RequestApi[Out, In any](t *testing.T, h http.Handler, method, url string, in In) Out {
	b, err := json.Marshal(in)
	assert.Nil(t, err)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(b))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.EqualValues(t, resp.StatusCode, http.StatusOK)
	var o Out
	err = json.NewDecoder(resp.Body).Decode(&o)
	assert.Nil(t, err)

	return o
}
