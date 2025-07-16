package tools

import (
	"container/list"
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

type queue[T any] list.List

func NewQueue[T any]() Queue[T] {
	return &queue[T]{}
}

func (q *queue[T]) Enqueue(item T) {
	(*list.List)(q).PushBack(item)
}

func (q *queue[T]) Dequeue() (T, bool) {
	if (*list.List)(q).Len() == 0 {
		var zero T
		return zero, false
	}
	elem := (*list.List)(q).Front()
	(*list.List)(q).Remove(elem)
	return elem.Value.(T), true
}

func (q *queue[T]) Peek() (T, bool) {
	if (*list.List)(q).Len() == 0 {
		var zero T
		return zero, false
	}
	front := (*list.List)(q).Front()
	return front.Value.(T), true
}

func (q *queue[T]) IsEmpty() bool {
	return (*list.List)(q).Len() == 0
}

func (q *queue[T]) Size() int {
	return (*list.List)(q).Len()
}

func (q *queue[T]) Clear() {
	(*list.List)(q).Init()
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
