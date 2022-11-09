package middleware_test

import (
	"authstore/internal/common/http/handler"
	"authstore/internal/common/http/middleware"
	"authstore/tests/stubs/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapterMiddleware(t *testing.T) {
	stubLogger := logger.Logger{}
	m := middleware.NewMiddleware(stubLogger)
	test := struct {
		CalledHandler    bool
		HandleContextArg *handler.HandleContext
	}{
		CalledHandler:    false,
		HandleContextArg: &handler.HandleContext{},
	}
	var testHandle handler.Handle = func(hc *handler.HandleContext) error {
		test.HandleContextArg = hc
		test.CalledHandler = true
		return nil
	}

	h := m.AdapterMiddleware(testHandle)
	w := httptest.NewRecorder()
	r := &http.Request{}
	h(w, r, make(httprouter.Params, 1))
	require.True(t, test.CalledHandler, "test handler must be called")
	assert.NotNil(t, test.HandleContextArg)
}
