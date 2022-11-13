package pipeline_test

import (
	"authstore/internal/common/http/handler"
	"authstore/pkg/pipeline"
	"authstore/tests/mocks/pkg/mock_pipeline"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func newPipline() *pipeline.Pipeline {
	return pipeline.New()
}

func newHandleContext() *handler.HandleContext {
	return handler.NewHandleContext(
		httptest.NewRecorder(),
		&http.Request{},
		make(httprouter.Params, 0),
	)
}
func TestPipeline_pipe(t *testing.T) {
	p := newPipline()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandle := mock_pipeline.NewMockHandle(ctrl)
	mockHandle2 := mock_pipeline.NewMockHandle(ctrl)

	p.Pipe(mockHandle)
	p.Pipe(mockHandle2)

	assert.NotNil(t, newHandleContext())
	assert.NotEmpty(t, p.Handlers)
	assert.Equal(t, 2, p.Handlers.Length())
}

func TestPipeline_runWithDefaultHandle(t *testing.T) {
	p := newPipline()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hctx := newHandleContext()
	mockHandle := mock_pipeline.NewMockHandle(ctrl)
	mockHandle.EXPECT().Handle(hctx, nil).Return(nil)

	err := p.Handle(hctx, mockHandle)

	assert.NoError(t, err)
}
func TestPipeline_runWithNilDefaultHandle(t *testing.T) {
	p := newPipline()
	err := p.Handle(newHandleContext(), nil)
	assert.Error(t, err)
}
func TestPipeline_runWithManyHandlers(t *testing.T) {
	p := newPipline()

	mock1 := mockMiddleware1{}
	mock2 := mockMiddleware2{}
	last := mockMiddleware3{}

	p.Pipe(mock1)
	p.Pipe(mock2)
	hctx := newHandleContext()
	err := p.Handle(hctx, last)
	assert.NoError(t, err)
	assert.Equal(t, "suka", hctx.HttpContext.W.Header().Get("fuck"))
	assert.Equal(t, "suka", hctx.HttpContext.W.Header().Get("gandon"))
	assert.Equal(t, "fuck", hctx.HttpContext.W.Header().Get("pidor"))
}
func TestPipeline_doubleRun(t *testing.T) {
	p := newPipline()

	mock1 := mockMiddleware1{}
	mock2 := mockMiddleware2{}
	last := mockMiddleware3{}

	p.Pipe(mock1)
	p.Pipe(mock2)
	hctx := newHandleContext()
	err := p.Handle(hctx, last)
	assert.NoError(t, err)
	assert.Equal(t, "suka", hctx.HttpContext.W.Header().Get("fuck"))
}

type mockMiddleware1 struct{}

func (m mockMiddleware1) Handle(hctx *handler.HandleContext, next pipeline.Handle) error {
	hctx.W.Header().Add("fuck", "suka")
	return next.Handle(hctx, nil)
}

type mockMiddleware2 struct{}

func (m mockMiddleware2) Handle(hctx *handler.HandleContext, next pipeline.Handle) error {
	hctx.W.Header().Add("pidor", "fuck")
	return next.Handle(hctx, nil)
}

type mockMiddleware3 struct{}

func (m mockMiddleware3) Handle(hctx *handler.HandleContext, next pipeline.Handle) error {
	hctx.W.Header().Add("gandon", "suka")
	return nil
}
