package skelet

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
)

func RegisterTestRunner[T any](
	t *testing.T,
	flavors FlavorProvider,
	bones ...any,
) *T {
	// TODO(mpavlicek): maybe use deadline from testing
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	t.Cleanup(cancel)

	runner, skelet, err := AssembleRunner[T](
		nil,
		"test",
		flavors,
		bones...,
	)
	assert.Nil(t, err)

	err = runner.RunBefore(ctx)
	assert.Nil(t, err)

	t.Cleanup(runner.RunAfter)

	return skelet
}

func ApiOk[Out, In any](t *testing.T, h http.Handler, method, url string, in In) Out {
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
