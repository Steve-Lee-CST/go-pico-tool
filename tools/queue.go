package tools

import (
	"container/list"
	"sync"
)

var (
	_ IQueue[any] = (*Queue[any])(nil)
	_ IQueue[any] = (*ConcurrentQueue[any])(nil)
)

type IQueue[T any] interface {
	Enqueue(item T)
	Dequeue() (T, bool)
	IsEmpty() bool
	Size() int
	Clear()
}

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

func (q *Queue[T]) Clear() {
	q.items.Init()
}

type ConcurrentQueue[T any] struct {
	q    *Queue[T]
	lock sync.RWMutex
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		q:    NewQueue[T](),
		lock: sync.RWMutex{},
	}
}

func (cq *ConcurrentQueue[T]) Enqueue(item T) {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	cq.q.Enqueue(item)
}

func (cq *ConcurrentQueue[T]) Dequeue() (T, bool) {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	return cq.q.Dequeue()
}

func (cq *ConcurrentQueue[T]) IsEmpty() bool {
	cq.lock.RLock()
	defer cq.lock.RUnlock()
	return cq.q.IsEmpty()
}

func (cq *ConcurrentQueue[T]) Size() int {
	cq.lock.RLock()
	defer cq.lock.RUnlock()
	return cq.q.Size()
}

func (cq *ConcurrentQueue[T]) Clear() {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	cq.q.Clear()
}
