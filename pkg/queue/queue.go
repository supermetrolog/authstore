package queue

type Queue struct {
	handlers []any
}

func New() *Queue {
	return &Queue{}
}
func (q *Queue) IsEmpty() bool {
	return len(q.handlers) == 0
}
func (q *Queue) Length() int {
	return len(q.handlers)
}

func (q *Queue) Enqueue(item any) {
	q.handlers = append(q.handlers, item)
}

func (q *Queue) Dequeue() any {
	if q.IsEmpty() {
		return nil
	}
	h := q.handlers[0]
	q.handlers = q.handlers[1:]
	return h
}
