package pipeline

import (
	"authstore/internal/common/http/handler"
	"authstore/pkg/queue"
	"errors"
)

type Handle interface {
	Handle(hctx *handler.HandleContext, next Handle) error
}

type Pipeline struct {
	Handlers queue.Queue
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) Pipe(handle Handle) {
	p.Handlers.Enqueue(handle)
}
func (p *Pipeline) Handle(hctx *handler.HandleContext, nextDefault Handle) error {
	if nextDefault == nil {
		return errors.New("default Handle can not be nil")
	}
	n := newNext(p.Handlers, nextDefault)
	return n.Next(hctx)
}
