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

type node[T any] struct {
	data T
	next *node[T]
}

type queue[T any] struct {
	head *node[T]
	tail *node[T]
	size int
}

func NewQueue[T any]() Queue[T] {
	return &queue[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (q *queue[T]) Enqueue(item T) {
	newNode := &node[T]{
		data: item,
		next: nil,
	}

	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}

	q.size++
}

func (q *queue[T]) Dequeue() (T, bool) {
	var zero T

	if q.head == nil {
		return zero, false
	}

	data := q.head.data
	q.head = q.head.next

	if q.head == nil {
		q.tail = nil
	}

	q.size--
	return data, true
}

func (q *queue[T]) Peek() (T, bool) {
	var zero T

	if q.head == nil {
		return zero, false
	}

	return q.head.data, true
}

func (q *queue[T]) IsEmpty() bool {
	return q.head == nil
}

func (q *queue[T]) Size() int {
	return q.size
}

func (q *queue[T]) Clear() {
	q.head = nil
	q.tail = nil
	q.size = 0
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
