package iters

import (
	"math/rand"
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	plusOne := func(i int) int { return i + 1 }
	nothing := func(int) {}
	buf := ints(314159)
	q := &Queue[int]{buf, new(sync.RWMutex)}
	var p *Queue[int]

	t.Run("eq", func(t *testing.T) {
		p = NewQueue(buf...)
		if !q.Eq(p, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, p.Buf)
		}
	})

	t.Run("clone", func(t *testing.T) {
		p := q.Clone()
		if !q.Eq(p, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, p.Buf)
		}
	})

	t.Run("len", func(t *testing.T) {
		have, want := q.Len(), len(q.Buf)
		if have != want {
			t.Errorf("have %d, want %d", have, want)
		}
	})

	t.Run("GoMap", func(t *testing.T) {
		clone := q.Clone().GoMap(plusOne)
		q := q.GoMap(plusOne)
		if !q.Eq(clone, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, clone.Buf)
		}
	})

	t.Run("map", func(t *testing.T) {
		clone := q.Clone().Map(plusOne)
		q := q.Map(plusOne)
		if !q.Eq(clone, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, clone.Buf)
		}
	})

	t.Run("GoCast", func(t *testing.T) {
		clone := q.Clone().GoCast(plusOne)
		q := q.GoCast(plusOne)
		if !q.Eq(clone, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, clone.Buf)
		}
	})

	t.Run("Cast", func(t *testing.T) {
		clone := q.Clone().Cast(plusOne)
		q := q.Cast(plusOne)
		if !q.Eq(clone, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, clone.Buf)
		}
	})

	t.Run("Do", func(t *testing.T) {
		clone := q.Clone().Do(nothing)
		if !q.Eq(clone, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, clone.Buf)
		}
	})

	t.Run("GoDo", func(t *testing.T) {
		clone := q.Clone().GoDo(nothing)
		if !q.Eq(clone, Eq) {
			t.Errorf("Queue with head @ %p != Que with head @ %p", q.Buf, clone.Buf)
		}
	})

	t.Run("Pop", func(t *testing.T) {
		want := q.Buf[0]
		have := q.Pop()
		if *have != want {
			t.Errorf("have %d, want %d", *have, want)
		}
	})
	t.Run("Push", func(t *testing.T) {
		want := rand.Int()
		have := q.Push(want).Buf[len(q.Buf)-1]
		if have != want {
			t.Errorf("have %d, want %d", have, want)
		}
	})
}
