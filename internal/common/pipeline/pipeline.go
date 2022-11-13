package pipeline

import (
	"authstore/internal/common/http/handler"
	"authstore/internal/common/loggerinterface"
	"errors"
)

type Handle interface {
	Handle(hctx *handler.HandleContext, next Handle) error
}

type Pipeline struct {
	logger      loggerinterface.Logger
	Handlers    queue
	nextDefault Handle
}
type handleWrapper struct {
	pipeline *Pipeline
}

func (f handleWrapper) Handle(hctx *handler.HandleContext, next Handle) error {
	return f.pipeline.Next(hctx, next)
}
func NewPipline(logger loggerinterface.Logger) *Pipeline {
	return &Pipeline{
		logger: logger,
	}
}

func (p *Pipeline) Pipe(handle Handle) {
	p.Handlers.Enqueue(handle)
}
func (p *Pipeline) Handle(hctx *handler.HandleContext, nextDefault Handle) error {
	if nextDefault == nil {
		return errors.New("default Handle can not be nil")
	}
	p.nextDefault = nextDefault
	return p.Next(hctx, nextDefault)
}
func (p Pipeline) Next(hctx *handler.HandleContext, _ Handle) error {
	if p.Handlers.IsEmpty() {
		return p.nextDefault.Handle(hctx, nil)
	}
	current := p.Handlers.Dequeue()
	return current.Handle(hctx, handleWrapper{pipeline: &p})
}
