package pipeline

import (
	"authstore/internal/common/http/handler"
	"authstore/pkg/queue"
	"errors"
)

type next struct {
	nextDefault Handle
	Handlers    queue.Queue
}
type nextWrapper struct {
	n *next
}

func (n nextWrapper) Handle(hctx *handler.HandleContext, next Handle) error {
	return n.n.Next(hctx)
}
func newNext(q queue.Queue, nextDefault Handle) next {
	return next{
		Handlers:    q,
		nextDefault: nextDefault,
	}
}
func (n next) Next(hctx *handler.HandleContext) error {
	if n.Handlers.IsEmpty() {
		return n.nextDefault.Handle(hctx, nil)
	}
	current, ok := n.Handlers.Dequeue().(Handle)
	if !ok {
		return errors.New("unknown item in Handlers Queue")
	}
	return current.Handle(hctx, nextWrapper{n: &n})
}
