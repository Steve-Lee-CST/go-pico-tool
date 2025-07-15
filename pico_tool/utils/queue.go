package utils

import (
	"container/list"
	"sync"
)

type IQueue[T any] interface {
	Enqueue(item T)
	Dequeue() (T, bool)
	IsEmpty() bool
	Size() int
}

var (
	_ IQueue[int] = (*Queue[int])(nil)
	_ IQueue[int] = (*ConcurrentQueue[int])(nil)
)

type Queue[T any] struct {
	items *list.List
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: list.New()}
}

func (q *Queue[T]) Enqueue(item T) {
	q.items.PushBack(item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if q.items.Len() == 0 {
		var zero T
		return zero, false
	}
	front := q.items.Front()
	q.items.Remove(front)
	return front.Value.(T), true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.items.Len() == 0
}

func (q *Queue[T]) Size() int {
	return q.items.Len()
}

type ConcurrentQueue[T any] struct {
	items *list.List
	lock  sync.RWMutex
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		items: list.New(),
		lock:  sync.RWMutex{},
	}
}

func (cq *ConcurrentQueue[T]) Enqueue(item T) {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	cq.items.PushBack(item)
}

func (cq *ConcurrentQueue[T]) Dequeue() (T, bool) {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	if cq.items.Len() == 0 {
		var zero T
		return zero, false
	}
	front := cq.items.Front()
	cq.items.Remove(front)
	return front.Value.(T), true
}

func (cq *ConcurrentQueue[T]) IsEmpty() bool {
	cq.lock.RLock()
	defer cq.lock.RUnlock()
	return cq.items.Len() == 0
}

func (cq *ConcurrentQueue[T]) Size() int {
	cq.lock.RLock()
	defer cq.lock.RUnlock()
	return cq.items.Len()
}
