package tools

import (
	"sync"
)

var (
	_ Queue[any] = (*queue[any])(nil)
	_ Queue[any] = (*concurrentQueue[any])(nil)
)

type Queue[T any] interface {
	Enqueue(item T)
	Dequeue() (T, bool)
	Peek() (T, bool)
	IsEmpty() bool
	Size() int
	Clear()
}

type queue[T any] []T

func NewQueue[T any]() Queue[T] {
	q := make(queue[T], 0)
	return &q
}

func (q *queue[T]) Enqueue(item T) {
	*q = append(*q, item)
}

func (q *queue[T]) Dequeue() (T, bool) {
	if len(*q) == 0 {
		var zero T
		return zero, false
	}
	elem := (*q)[0]
	*q = (*q)[1:]
	return elem, true
}

func (q *queue[T]) Peek() (T, bool) {
	if len(*q) == 0 {
		var zero T
		return zero, false
	}
	return (*q)[0], true
}

func (q *queue[T]) IsEmpty() bool {
	return len(*q) == 0
}

func (q *queue[T]) Size() int {
	return len(*q)
}

func (q *queue[T]) Clear() {
	*q = (*q)[:0]
}

type concurrentQueue[T any] struct {
	queue[T]
	rwLock sync.RWMutex
}

func NewConcurrentQueue[T any]() Queue[T] {
	return &concurrentQueue[T]{
		queue: queue[T]{},
	}
}

func (cq *concurrentQueue[T]) Enqueue(item T) {
	cq.rwLock.Lock()
	defer cq.rwLock.Unlock()
	cq.queue.Enqueue(item)
}

func (cq *concurrentQueue[T]) Dequeue() (T, bool) {
	cq.rwLock.Lock()
	defer cq.rwLock.Unlock()
	return cq.queue.Dequeue()
}

func (cq *concurrentQueue[T]) Peek() (T, bool) {
	cq.rwLock.RLock()
	defer cq.rwLock.RUnlock()
	return cq.queue.Peek()
}

func (cq *concurrentQueue[T]) IsEmpty() bool {
	cq.rwLock.RLock()
	defer cq.rwLock.RUnlock()
	return cq.queue.IsEmpty()
}

func (cq *concurrentQueue[T]) Size() int {
	cq.rwLock.RLock()
	defer cq.rwLock.RUnlock()
	return cq.queue.Size()
}

func (cq *concurrentQueue[T]) Clear() {
	cq.rwLock.Lock()
	defer cq.rwLock.Unlock()
	cq.queue.Clear()
}
