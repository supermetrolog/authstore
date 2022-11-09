package pipeline_test

import (
	"authstore/internal/common/http/handler"
	"authstore/internal/common/pipeline"
	mock_handle "authstore/tests/mocks/pipeline"
	"authstore/tests/stubs/logger"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func newPipline() *pipeline.Pipeline {
	return pipeline.NewPipline(
		logger.Logger{},
		httptest.NewRecorder(),
		&http.Request{},
		make(httprouter.Params, 0),
	)
}
func TestPipeline_pipe(t *testing.T) {
	p := newPipline()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandle := mock_handle.NewMockHandle(ctrl)
	mockHandle2 := mock_handle.NewMockHandle(ctrl)

	p.Pipe(mockHandle)
	p.Pipe(mockHandle2)

	assert.NotNil(t, p.HandleContext)
	assert.NotEmpty(t, p.Handlers)
	assert.Equal(t, 2, p.Handlers.Length())
}

func TestPipeline_runWithDefaultHandle(t *testing.T) {
	p := newPipline()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandle := mock_handle.NewMockHandle(ctrl)
	mockHandle.EXPECT().Handle(p.HandleContext, nil).Return(nil)

	err := p.Handle(p.HandleContext, mockHandle)

	assert.NoError(t, err)
}
func TestPipeline_runWithNilDefaultHandle(t *testing.T) {
	p := newPipline()

	err := p.Handle(p.HandleContext, nil)
	assert.Error(t, err)
}
func TestPipeline_runWithManyHandlers(t *testing.T) {
	p := newPipline()

	mock1 := mockMiddleware1{}
	mock2 := mockMiddleware2{}
	last := mockMiddleware3{}

	p.Pipe(mock1)
	p.Pipe(mock2)

	err := p.Handle(p.HandleContext, last)
	assert.NoError(t, err)
	assert.Equal(t, "suka", p.HandleContext.HttpContext.W.Header().Get("fuck"))
}

// func TestPipeline_runWithManyHandlers(t *testing.T) {
// 	p := newPipline()
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockHandle := mock_handle.NewMockHandle(ctrl)
// 	mockHandle2 := mock_handle.NewMockHandle(ctrl)
// 	mockHandleDefault := mock_handle.NewMockHandle(ctrl)
// 	mockHandle.EXPECT().Handle(p.HandleContext, mockHandle2).DoAndReturn(func(hctx *handler.HandleContext, next pipeline.Handle) error {
// 		return next.Handle(hctx, nil)
// 	})
// 	mockHandle2.EXPECT().Handle(p.HandleContext, nil).DoAndReturn(func(hctx *handler.HandleContext, next pipeline.Handle) error {
// 		return next.Handle(hctx, nil)
// 	})
// 	// gomock.InOrder(

// 	// )
// 	p.Pipe(mockHandle)
// 	p.Pipe(mockHandle2)

// 	err := p.Handle(p.HandleContext, mockHandleDefault)
// 	assert.NoError(t, err)
// }
func TestPipeline_queue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHandle := mock_handle.NewMockHandle(ctrl)
	mockHandle2 := mock_handle.NewMockHandle(ctrl)
	q := pipeline.NewHandlersQueue()
	q.Enqueue(mockHandle)
	q.Enqueue(mockHandle2)

	assert.Equal(t, mockHandle, q.Dequeue())
	assert.Equal(t, mockHandle2, q.Dequeue())

}

type mockMiddleware1 struct{}

func (m mockMiddleware1) Handle(hctx *handler.HandleContext, next pipeline.Handle) error {
	fmt.Println("MD 1")
	hctx.W.Header().Add("fuck", "suka")
	return next.Handle(hctx, nil)
}

type mockMiddleware2 struct{}

func (m mockMiddleware2) Handle(hctx *handler.HandleContext, next pipeline.Handle) error {
	fmt.Println("MD 2")
	hctx.W.Header().Add("pidor", "suka")
	return next.Handle(hctx, nil)
}

type mockMiddleware3 struct{}

func (m mockMiddleware3) Handle(hctx *handler.HandleContext, next pipeline.Handle) error {
	fmt.Println("MD 3")
	hctx.W.Header().Add("gandon", "suka")
	return nil
}
