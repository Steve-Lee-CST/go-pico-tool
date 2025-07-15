package tools

import "sync"

var (
	_ ISet[any] = (*Set[any])(nil)
	_ ISet[any] = (*ConcurrentSet[any])(nil)
)

type ISet[T comparable] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Size() int
	IsEmpty() bool
	Clear()
	Items() []T
	Copy() ISet[T]
	Union(other ISet[T]) ISet[T]
	Intersection(other ISet[T]) ISet[T]
	Difference(other ISet[T]) ISet[T]
}

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

func NewSetFromISet[T comparable](src ISet[T]) *Set[T] {
	newSet := NewSet[T]()
	for _, item := range src.Items() {
		newSet.Add(item)
	}
	return newSet
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

func (s *Set[T]) Copy() ISet[T] {
	newSet := NewSet[T]()
	for item := range s.items {
		newSet.Add(item)
	}
	return newSet
}

func (s *Set[T]) Union(other ISet[T]) ISet[T] {
	unionSet := NewSet[T]()
	for item := range s.items {
		unionSet.Add(item)
	}
	for _, item := range other.Items() {
		unionSet.Add(item)
	}
	return unionSet
}

func (s *Set[T]) Intersection(other ISet[T]) ISet[T] {
	intersectionSet := NewSet[T]()
	for item := range s.items {
		if other.Contains(item) {
			intersectionSet.Add(item)
		}
	}
	return intersectionSet
}

func (s *Set[T]) Difference(other ISet[T]) ISet[T] {
	differenceSet := NewSet[T]()
	for item := range s.items {
		if !other.Contains(item) {
			differenceSet.Add(item)
		}
	}
	return differenceSet
}

type ConcurrentSet[T comparable] struct {
	set  *Set[T]
	lock sync.RWMutex
}

func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{
		set:  NewSet[T](),
		lock: sync.RWMutex{},
	}
}

func NewConcurrentSetFromISet[T comparable](other ISet[T]) *ConcurrentSet[T] {
	return &ConcurrentSet[T]{
		set:  NewSetFromISet(other),
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

func (cs *ConcurrentSet[T]) Copy() ISet[T] {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	// return concurrent set in ConcurrentSet.Copy
	return NewConcurrentSetFromISet(cs.set.Copy())
}

func (cs *ConcurrentSet[T]) Union(other ISet[T]) ISet[T] {
	// Create a copy of the other set to ensure thread safety
	otherCopy := other.Copy()

	cs.lock.RLock()
	defer cs.lock.RUnlock()
	// return concurrent set in ConcurrentSet.Union
	return NewConcurrentSetFromISet(cs.set.Union(otherCopy))
}

func (cs *ConcurrentSet[T]) Intersection(other ISet[T]) ISet[T] {
	// Create a copy of the other set to ensure thread safety
	otherCopy := other.Copy()

	cs.lock.RLock()
	defer cs.lock.RUnlock()
	// return concurrent set in ConcurrentSet.Intersection
	return NewConcurrentSetFromISet(cs.set.Intersection(otherCopy))
}

func (cs *ConcurrentSet[T]) Difference(other ISet[T]) ISet[T] {
	// Create a copy of the other set to ensure thread safety
	otherCopy := other.Copy()

	cs.lock.RLock()
	defer cs.lock.RUnlock()
	// return concurrent set in ConcurrentSet.Difference
	return NewConcurrentSetFromISet(cs.set.Difference(otherCopy))
}
