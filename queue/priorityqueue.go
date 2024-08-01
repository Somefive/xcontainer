package queue

import (
	"sync"
)

type indexedItem[T any] struct {
	value T   // The value of the item; arbitrary.
	index int // The index of the item in the heap.
}

// PriorityQueue the priority queue interface
type PriorityQueue[T any] interface {
	Len() int
	Push(T)
	Pop() T
	Top() T
}

type priorityQueue[T any] struct {
	mu         sync.RWMutex
	items      []*indexedItem[T]
	comparator func(T, T) bool
}

// NewPriorityQueue create a thread-safe priority queue with given comparator. The comparator must not be nil.
func NewPriorityQueue[T any](data []T, comparator func(T, T) bool) PriorityQueue[T] {
	items := make([]*indexedItem[T], len(data))
	for i := 0; i < len(data); i++ {
		items[i] = &indexedItem[T]{value: data[i], index: i}
	}
	in := &priorityQueue[T]{
		items:      items,
		comparator: comparator,
	}
	for i := len(data) - 1; i >= 0; i-- {
		in.fix(i)
	}
	return in
}

// Len return the length of the queue
func (in *priorityQueue[T]) Len() int {
	return len(in.items)
}

func (in *priorityQueue[T]) swap(i, j int) {
	in.items[i], in.items[j] = in.items[j], in.items[i]
	in.items[i].index, in.items[j].index = i, j
}

// Push add a value to the queue
func (in *priorityQueue[T]) Push(value T) {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.items = append(in.items, &indexedItem[T]{value, len(in.items)})
	in.fix(len(in.items) - 1)
}

// Pop return the top item of the queue
func (in *priorityQueue[T]) Pop() T {
	in.mu.Lock()
	defer in.mu.Unlock()
	item := in.items[0].value
	in.swap(0, len(in.items)-1)
	in.items = in.items[:len(in.items)-1]
	if len(in.items) > 0 {
		in.fix(0)
	}
	return item
}

// Top return the top item of the queue
func (in *priorityQueue[T]) Top() T {
	in.mu.RLock()
	defer in.mu.RUnlock()
	return in.items[0].value
}

func (in *priorityQueue[T]) fix(index int) {
	i := index
	for { // down fix
		left, right, top := 2*i+1, 2*i+2, i
		if left < len(in.items) && in.comparator(in.items[top].value, in.items[left].value) {
			top = left
		}
		if right < len(in.items) && in.comparator(in.items[top].value, in.items[right].value) {
			top = right
		}
		if i == top {
			break
		}
		in.swap(i, top)
		i = top
	}
	if i != index {
		return
	}
	for { // up fix
		parent := (i - 1) >> 1
		if i == parent || !in.comparator(in.items[i].value, in.items[i].value) {
			break
		}
		in.swap(i, parent)
		i = parent
	}
}
