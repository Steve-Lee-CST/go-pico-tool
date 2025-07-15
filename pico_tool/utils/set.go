package utils

import "sync"

type ISet[T comparable] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Size() int
	IsEmpty() bool
	Clear()
	Items() []T
	Copy() *Set[T]
	Union(other *Set[T]) *Set[T]
	Intersection(other *Set[T]) *Set[T]
	Difference(other *Set[T]) *Set[T]
}

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

func (s *Set[T]) Items() []T {
	items := make([]T, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return items
}

func (s *Set[T]) Copy() *Set[T] {
	newSet := NewSet[T]()
	for item := range s.items {
		newSet.Add(item)
	}
	return newSet
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	unionSet := NewSet[T]()
	for item := range s.items {
		unionSet.Add(item)
	}
	for item := range other.items {
		unionSet.Add(item)
	}
	return unionSet
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	intersectionSet := NewSet[T]()
	for item := range s.items {
		if other.Contains(item) {
			intersectionSet.Add(item)
		}
	}
	return intersectionSet
}

func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	differenceSet := NewSet[T]()
	for item := range s.items {
		if !other.Contains(item) {
			differenceSet.Add(item)
		}
	}
	return differenceSet
}

type ConcurrentSet[T comparable] struct {
	set  Set[T]
	lock sync.RWMutex
}

func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{
		set:  *NewSet[T](),
		lock: sync.RWMutex{},
	}
}

func (cs *ConcurrentSet[T]) Add(item T) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.set.Add(item)
}

func (cs *ConcurrentSet[T]) Remove(item T) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.set.Remove(item)
}

func (cs *ConcurrentSet[T]) Contains(item T) bool {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return cs.set.Contains(item)
}

func (cs *ConcurrentSet[T]) Size() int {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return cs.set.Size()
}

func (cs *ConcurrentSet[T]) IsEmpty() bool {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return cs.set.IsEmpty()
}

func (cs *ConcurrentSet[T]) Clear() {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.set.Clear()
}

func (cs *ConcurrentSet[T]) Items() []T {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return cs.set.Items()
}

func (cs *ConcurrentSet[T]) Copy() *ConcurrentSet[T] {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	return &ConcurrentSet[T]{set: *cs.set.Copy(), lock: sync.RWMutex{}}
}

func (cq *ConcurrentSet[T]) Union(other *ConcurrentSet[T]) *ConcurrentSet[T] {
	unionSet := NewConcurrentSet[T]()
	func() {
		cq.lock.RLock()
		defer cq.lock.RUnlock()
		for item := range cq.set.items {
			unionSet.Add(item)
		}
	}()
	func() {
		other.lock.RLock()
		defer other.lock.RUnlock()
		for item := range other.set.items {
			unionSet.Add(item)
		}
	}()
	return unionSet
}

func (cq *ConcurrentSet[T]) Intersection(other *ConcurrentSet[T]) *ConcurrentSet[T] {
	intersectionSet := NewConcurrentSet[T]()
	otherCopy := other.Copy()

	cq.lock.RLock()
	defer cq.lock.RUnlock()
	for item := range cq.set.items {
		if otherCopy.set.Contains(item) {
			intersectionSet.Add(item)
		}
	}
	return intersectionSet
}

func (cq *ConcurrentSet[T]) Difference(other *ConcurrentSet[T]) *ConcurrentSet[T] {
	differenceSet := NewConcurrentSet[T]()
	otherCopy := other.Copy()

	cq.lock.RLock()
	defer cq.lock.RUnlock()
	for item := range cq.set.items {
		if !otherCopy.set.Contains(item) {
			differenceSet.Add(item)
		}
	}
	return differenceSet
}
