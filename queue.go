package iters

import (
	"fmt"
	"strings"
	"sync"

	"github.com/kendfss/but"
)

// Queue is a threadsafe FIFO data structure
type Queue[T any] struct {
	Buf []T
	*sync.RWMutex
}

// NewQueue initializes a Queue variadically. If you already have a slice, just use the constructor
func NewQueue[T any](elems ...T) *Queue[T] {
	return &Queue[T]{elems, new(sync.RWMutex)}
}

// Push appends an element to q
func (q *Queue[T]) Push(e T) *Queue[T] {
	q.RLock()
	defer q.RUnlock()
	q.Buf = append(q.Buf, e)
	return q
}

// Pop returns the first element from q
// returns nil for the empty Queue
func (q *Queue[T]) Pop() *T {
	q.Lock()
	defer q.Unlock()
	if len(q.Buf) == 0 {
		return nil
	}
	e := q.Buf[0]
	q.Buf = q.Buf[1:]
	return &e
}

// Map returns a new Queue whose elements are return values of fn on q
func (q *Queue[T]) Map(fn func(T) T) *Queue[T] {
	q.RLock()
	defer q.RUnlock()
	buf := make([]T, len(q.Buf))
	for i, e := range q.Buf {
		buf[i] = fn(e)
	}
	return &Queue[T]{buf, new(sync.RWMutex)}
}

// GoMap performs a concurrent mapping on the elements of the Queue
func (q *Queue[T]) GoMap(fn func(T) T) *Queue[T] {
	q.RLock()
	defer q.RUnlock()
	wg := new(sync.WaitGroup)
	wg.Add(len(q.Buf))
	buf := make([]T, len(q.Buf))
	for i, e := range q.Buf {
		go func(i int, t T) { defer wg.Done(); buf[i] = fn(t) }(i, e)
	}
	wg.Wait()
	p := &Queue[T]{buf, new(sync.RWMutex)}
	return p
}

// Cast performs an in-place map of q
func (q *Queue[T]) Cast(fn func(T) T) *Queue[T] {
	q.Lock()
	defer q.Unlock()
	for i, e := range q.Buf {
		q.Buf[i] = fn(e)
	}
	return q
}

// GoCast performs a concurrent Cast of each element in the Queue
func (q *Queue[T]) GoCast(fn func(T) T) *Queue[T] {
	q.Lock()
	defer q.Unlock()
	wg := new(sync.WaitGroup)
	wg.Add(len(q.Buf))
	for i, e := range q.Buf {
		go func(i int, t T) {
			defer wg.Done()
			q.Buf[i] = fn(t)
		}(i, e)
	}
	wg.Wait()
	return q
}

// Do executes some function on each element of q
func (q *Queue[T]) Do(fn func(T)) *Queue[T] {
	q.RLock()
	defer q.RUnlock()
	for _, e := range q.Buf {
		fn(e)
	}
	return q
}

// GoDo executes some function on each element of q concurrently
func (q *Queue[T]) GoDo(fn func(T)) *Queue[T] {
	q.RLock()
	defer q.RUnlock()
	wg := new(sync.WaitGroup)
	wg.Add(len(q.Buf))
	for _, e := range q.Buf {
		go func(t T) { defer wg.Done(); fn(t) }(e)
	}
	wg.Wait()
	return q
}

// Eq checks two Queues for equality
func (q *Queue[T]) Eq(p *Queue[T], equals func(T, T) bool) bool {
	q.RLock()
	defer q.RUnlock()
	p.RLock()
	defer p.RUnlock()
	if len(q.Buf) != len(p.Buf) {
		return false
	}
	for i, e := range q.Buf {
		if !equals(e, p.Buf[i]) {
			return false
		}
	}
	return true
}

// Len returns the length of the underlying buffer
func (q *Queue[T]) Len() int {
	return len(q.Buf)
}

// String returns a text representation of the owner.
// Specifically it returns the type name and pointer of the first element
func (q Queue[T]) String() string {
	q.RLock()
	defer q.RUnlock()
	name := fmt.Sprintf("%T", q)
	parts := strings.Split(name, ".")
	name = parts[len(parts)-1]
	return fmt.Sprintf("%s@%p", name, q.Buf)
}

const ErrBadCopy but.Note = "bad copy: should've copied %d but copied %d"

// Clone returns a new queue with the same elements in a new buffer
func (q *Queue[T]) Clone() *Queue[T] {
	q.RLock()
	defer q.RUnlock()
	buf := make([]T, len(q.Buf))
	n := copy(buf, q.Buf)
	if n != len(q.Buf) {
		panic(ErrBadCopy.Fmt(len(q.Buf), n))
	}
	return &Queue[T]{buf, new(sync.RWMutex)}
}
