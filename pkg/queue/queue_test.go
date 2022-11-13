package queue_test

import (
	"authstore/pkg/queue"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Length(t *testing.T) {
	q := queue.New()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	assert.Equal(t, 3, q.Length())

	q = queue.New()
	q.Enqueue(1)
	assert.Equal(t, 1, q.Length())
}
func TestQueue_IsEmpty(t *testing.T) {
	q := queue.New()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	assert.False(t, q.IsEmpty())
	q = queue.New()
	assert.True(t, q.IsEmpty())
}
func TestQueue_withIntItems(t *testing.T) {
	q := queue.New()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	assert.Equal(t, 1, q.Dequeue())
	assert.Equal(t, 2, q.Dequeue())
	assert.Equal(t, 3, q.Dequeue())
}

func TestQueue_withStrItems(t *testing.T) {
	q := queue.New()
	q.Enqueue("fuck")
	q.Enqueue("suka")
	q.Enqueue("aaaa")

	assert.Equal(t, "fuck", q.Dequeue())
	assert.Equal(t, "suka", q.Dequeue())
	assert.Equal(t, "aaaa", q.Dequeue())
}

func TestQueue_withMultitypeItems(t *testing.T) {
	q := queue.New()
	q.Enqueue("fuck")
	q.Enqueue(1)
	q.Enqueue("aaaa")

	assert.Equal(t, "fuck", q.Dequeue())
	assert.Equal(t, 1, q.Dequeue())
	assert.Equal(t, "aaaa", q.Dequeue())
}
func TestQueue_withInterfaceItems(t *testing.T) {
	q := queue.New()
	q.Enqueue(queueItemImplementation1{})
	q.Enqueue(queueItemImplementation2{})

	queueItem1, ok1 := q.Dequeue().(queueItemInterface1)
	queueItem2, ok2 := q.Dequeue().(queueItemInterface1)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, "suka", queueItem1.TestFunc())
	assert.Equal(t, "fuck", queueItem2.TestFunc())
}
func TestQueue_withInterfaceItemsWithError(t *testing.T) {
	q := queue.New()
	q.Enqueue(queueItemImplementation1{})
	q.Enqueue(queueItemImplementation2{})

	queueItem1, ok1 := q.Dequeue().(queueItemInterface1)
	queueItem2, ok2 := q.Dequeue().(queueItemInterface2)
	assert.True(t, ok1)
	assert.False(t, ok2)
	assert.Equal(t, "suka", queueItem1.TestFunc())
	assert.Nil(t, queueItem2)
}

type queueItemInterface2 interface {
	AAAA() string
}
type queueItemInterface1 interface {
	TestFunc() string
}

type queueItemImplementation1 struct{}

func (q queueItemImplementation1) TestFunc() string {
	return "suka"
}

type queueItemImplementation2 struct{}

func (q queueItemImplementation2) TestFunc() string {
	return "fuck"
}
