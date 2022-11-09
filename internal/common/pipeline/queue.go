package pipeline

type queue struct {
	handlers []Handle
}

func NewHandlersQueue() *queue {
	return &queue{}
}
func (q *queue) IsEmpty() bool {
	return len(q.handlers) == 0
}
func (q *queue) Length() int {
	return len(q.handlers)
}

func (q *queue) Enqueue(h Handle) {
	q.handlers = append(q.handlers, h)
}

func (q *queue) Dequeue() Handle {
	if len(q.handlers) == 0 {
		return nil
	}
	h := q.handlers[0]
	q.handlers = q.handlers[1:]
	return h
}
