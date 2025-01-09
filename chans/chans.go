package chans

import (
	"context"
	"fmt"
	"sync"

	"github.com/kendfss/rules"
)

type (
	Readable[T any] interface{ chan T | <-chan T }
	Writable[T any] interface{ chan T | chan<- T }
)

// New initializes an unbuffered channel
func New[T any]() chan T {
	return make(chan T)
}

// Buffered initializes a buffered channel
func Buffered[T any](cap int) chan T {
	return make(chan T, cap)
}

// Inf produces an infinite channel of given capacity using init to create each element/
// Inf handles channel closure
func Inf[T any](init func() T, capacity int) <-chan T {
	out := make(chan T, capacity)
	if init == nil {
		init = func() T { return *new(T) }
	}
	go func() {
		defer close(out)
		for {
			out <- init()
		}
	}()
	return out
}

// Drain consumes a channel to depletion
func Drain[T any](src chan T) {
	for range src {
	}
}

// Filter takes a boolean channel and skips all false values it receives
func Filter(src chan bool) chan bool {
	out := make(chan bool, cap(src))
	go func() {
		for b := range src {
			if b {
				out <- b
			}
		}
	}()
	return out
}

// FilterPred skips any channel receipts that fail to satisfy the given predicate
func FilterPred[T any, channel Readable[T]](pred func(T) bool, src channel) chan T {
	out := make(chan T, cap(src))
	go func() {
		defer close(out)
		for e := range src {
			if pred(e) {
				out <- e
			}
		}
	}()
	return out
}

// Next extracts one receipt from the channel
func Next[T any](src <-chan T) T {
	return <-src
}

// NextSafe extracts one receipt from the channel and informs whether it's closed
func NextSafe[T any, chanT Readable[T]](src chanT) (T, bool) {
	v, open := <-src
	return v, open
}

// PopperDefault returns a popper that yields a default value once the channel is closed
func PopperDefault[T any, chanT Readable[T]](src chanT, defaultVal T) func() T {
	return func() T {
		out, open := <-src
		if !open {
			return defaultVal
		}
		return out
	}
}

// Popper functions return the next channel receipt that satisfies the given predicate.
// If you don't know when the source channel is closed, use NextSafe instead; unless the zero value is known to not satisfy the predicate.
func Popper[T any, channel Readable[T]](pred func(T) bool, src channel) func() T {
	return func() T {
		var (
			out  T
			done bool
		)
		for {
			out, done = <-src
			if done {
				break
			}
			if !pred(out) {
				continue
			}
			break
		}
		return out
	}
}

// PopperSafe functions return the next channel receipt that satisfies the given predicate and tells the user when the channel is closed
func PopperSafe[T any](fn func(T) bool, src chan T) func() (T, bool) {
	return func() (T, bool) {
		var (
			out  T
			open bool
		)
		for {
			out, open = <-src
			if !fn(out) {
				continue
			}
			break
		}
		return out, open
	}
}

// Context returns a channel that closes as soon as either the context is done or the source channel is closed.
// It does not, otherwise, operate on the context object an cannot cancel it
func Context[T any](ctx context.Context, src chan T) chan T {
	out := make(chan T, cap(src))
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, open := <-src:
				if !open {
					return
				}
				out <- v
			}
		}
	}()
	return out
}

// Deprecated: use ReadWrite
func RW[T any](c <-chan T) chan T {
	return ReadWrite(c)
}

// ReadWrite wraps a read-only channel with a read-write one
func ReadWrite[T any](c <-chan T) chan T {
	out := make(chan T, cap(c))
	go func() {
		defer close(out)
		for x := range c {
			out <- x
		}
	}()
	return out
}

// ReadOnly coerces a read-write channel's type to read-only
func ReadOnly[T any](c chan T) <-chan T {
	return c
}

// WriteOnly coerces a read-write channel's type to write-only
func WriteOnly[T any](c chan T) chan<- T {
	return c
}

// Deprecated: use ReadOnly instead
func RO[T any](c chan T) <-chan T { return ReadOnly(c) }

// Deprecated: use WriteOnly instead
func WO[T any](c chan T) chan<- T { return WriteOnly(c) }

// Count returns the number of elements passed through the channel before it is closed (externally)
func Count[T any, chanT Readable[T]](c chanT) (out uint64) {
	for range c {
		out++
	}
	return
}

// Compact removes all duplicates from a channel in constant time
func Compact[T comparable, chanT Readable[T]](src chanT) <-chan T {
	return CompactHash(selfHash[T], src)
}

// CompactFunc removes all duplicates from a channel in linear time
func CompactFunc[T comparable, chanT Readable[T]](eq func(T, T) bool, src chanT) <-chan T {
	out := make(chan T, cap(src))
	go func() {
		marked := []T{}
		pred := func(arg T) bool {
			if sliceContainsFunc(eq, marked, arg) {
				return true
			}
			marked = append(marked, arg)
			return false
		}
		out = FilterPred(pred, src)
	}()
	return out
}

// CompactHash removes all duplicates from a channel in constant time
func CompactHash[T any, H comparable, chanT Readable[T]](hash func(T) H, src chanT) <-chan T {
	out := make(chan T, cap(src))
	go func() {
		marked := map[H]struct{}{}
		pred := func(arg T) bool {
			sum := hash(arg)
			_, found := marked[sum]
			if found {
				return true
			}
			marked[sum] = struct{}{}
			return false
		}
		out = FilterPred(pred, src)
	}()
	return out
}

// Deprecated: use Do instead
func Send[T any, chanT Readable[T]](f func(T), src chanT) {
	Do(f, src)
}

// Do calls a function on every value of a slice
func Do[T any, chanT Readable[T]](f func(T), src chanT) {
	go func() {
		for e := range src {
			f(e)
		}
	}()
}

// Map calls Cast
func Map[I, O any, chanI Readable[I]](f func(I) O, src chanI) <-chan O {
	return Cast(f, src)
}

// Cast calls a pure function on every value of a channel and returns a channel
// containing all the results
func Cast[I, O any, chanI Readable[I]](f func(I) O, src chanI) chan O {
	dst := make(chan O, cap(src))
	go func() {
		for e := range src {
			dst <- f(e)
		}
	}()
	return dst
}

// Chain collects several channels and returns one populated by their content.
// Chain expects the user to control closure of argument channels.
// Chain handles closure of returned channel.
// Returned channel has capacity equal to that of the argument with highest capacity
func Chain[T any, chanT Readable[T]](args ...chanT) <-chan T {
	capacity := 0
	for _, arg := range args {
		capacity = max(cap(arg), capacity)
	}
	out := make(chan T, capacity)
	wg := new(sync.WaitGroup)
	wg.Add(len(args))
	go func() {
		defer close(out)
		for _, channel := range args {
			go func(c chanT) {
				defer wg.Done()
				for e := range c {
					out <- e
				}
			}(channel)
		}
		wg.Wait()
	}()
	return out
}

// ChainCap creates a new channel, with desired capacity, that serves as a frontent for the given arguments
func ChainCap[T any, chanT Readable[T]](capacity int, args ...chanT) <-chan T {
	out := make(chan T, capacity)
	wg := sync.WaitGroup{}
	wg.Add(len(args))
	go func() {
		defer close(out)
		for _, arg := range args {
			go func(arg <-chan T) {
				defer wg.Done()
				for e := range arg {
					out <- e
				}
			}(arg)
		}
		wg.Wait()
	}()
	return out
}

// Extend the first argument with the contents of the successors.
// non blocking, non order-preserving
func Extend[T any, readable Readable[T]](receiver chan T, args ...readable) {
	wg := new(sync.WaitGroup)
	wg.Add(len(args))
	go func() {
		for _, arg := range args {
			go func(arg <-chan T) {
				defer wg.Done()
				for e := range arg {
					receiver <- e
				}
			}(arg)
		}
		wg.Wait()
	}()
}

// Extender functions create extensions of the target channel.
// See Extend for more information
func Extender[T any, readable Readable[T]](target chan T) func(...readable) {
	return func(args ...readable) {
		Extend(target, args...)
	}
}

// Castro casts multiple read-write channels to readonly
func Castro[T any, chanT Readable[T]](buf []chanT) []<-chan T {
	out := make([]<-chan T, len(buf))
	for i, e := range buf {
		out[i] = e
	}
	return out
}

// Lazify converts a slice into a readonly channel
func Lazify[T any](arg []T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, e := range arg {
			out <- e
		}
	}()
	return out
}

// Upto returns an iterator whose content depends on the number of arguments as follows
//
//			# of args 	|| 	behaviour
//		 		 1	 	|| 	stop
//			 	 2	 	|| 	start, stop
//			 	 3	 	|| 	start, stop, step
//	           else 	|| 	error
//
// step defaults to 1
// start defaults to 0
func Upto[T rules.Real](args ...T) (chan T, error) {
	switch len(args) {
	case 1:
		return Upto(0, args[0], 1)
	case 2:
		return Upto(args[0], args[1], 1)
	case 3:
		out := make(chan T)
		go func() {
			start, stop, delta := args[0], args[1], args[2]
			for stop-delta >= start {
				out <- start
				start += delta
			}
		}()
		return out, nil
	case 0:
		return nil, fmt.Errorf("chans.Range: not enough args (%d). want 1, 2, or 3", len(args))
	default:
		return nil, fmt.Errorf("chans.Range: too many args (%d). want 1, 2, or 3", len(args))
	}
}

// MustUpto returns an iterator whose behaviour is equivalent to that of Range
func MustUpto[T rules.Real](args ...T) chan T {
	out, err := Upto(args...)
	if err != nil {
		panic(err)
	}
	return out
}

// Deprecated: use Push instead
func Put[T any, int rules.Int, chanT Writable[T]](count int, dst chanT) {
	Push(count, dst)
}

// Push writes "count" copies of zero-initialized "T" instances to "dst"
func Push[T any, int rules.Int, chanT Writable[T]](count int, dst chanT) {
	PushVal(count, dst, *new(T))
}

// PushVal writes a given value to a channel. Useful for spawning go routines before returning
func PushVal[T any, int rules.Int, chanT Writable[T]](count int, dst chanT, val T) {
	for ; count > 0; count-- {
		dst <- val
	}
}

// Deprecated: use Pusher instead
func Putter[T any, chanT Writable[T]](dst chanT) func(T) error {
	return Pusher(dst)
}

// Pusher returns a method of the given channel which sends it a given argument.
// It returns a non nil error when the destination is closed
func Pusher[T any, chanT Writable[T]](dst chanT) func(T) error {
	closed := false
	return func(arg T) error {
		var out error
		defer func() {
			if err := recover(); err != nil {
				out = fmt.Errorf("%s", err)
				closed = true
			}
		}()
		if closed {
			return ErrWriteOnEmpty
		}
		dst <- arg
		return out
	}
}

// Deprecated: use Pop instead
func Get[T any, chanT Readable[T]](count int, src chanT) []T {
	return Pop(count, src)
}

// Pop receives "count" items from "src"
func Pop[T any, chanT Readable[T]](count int, src chanT) []T {
	out := make([]T, count)
	for ; count > 0; count-- {
		out[len(out)-count] = <-src
	}
	return out
}

// Discard ignores "count" items from "src"
func Discard[T any, int rules.Int, chanT Readable[T]](count int, src chanT) {
	for ; count > 0; count-- {
		<-src
	}
}

// Watch feeds dst with items received from src.
// It does not close either of them
func Watch[T any](dst chan T, src <-chan T) {
	for e := range src {
		dst <- e
	}
}

// PredPutter returns a method of the given channel which sends it the given argument
// if, and only if, the argument satisfies the given predicate.
// It returns a non-nil error when the destination is closed
func PredPutter[T any](dst chan T, pred func(T) bool) func(T) error {
	put := Pusher(dst)
	return func(arg T) error {
		if pred(arg) {
			return put(arg)
		}
		return nil
	}
}

// Tee returns a pair of mutually independent copies of the source channel.
// Blocks at each iterative step until both channels have been read from, use Teen to avoid this behaviour
func Tee[T any](src <-chan T) (one, two chan T) {
	one = make(chan T, cap(src))
	two = make(chan T, cap(src))
	go func() {
		defer close(one)
		defer close(two)
		wg := new(sync.WaitGroup)
		for e := range src {
			wg.Add(2)
			go func(t T) { defer wg.Done(); one <- e }(e)
			go func(t T) { defer wg.Done(); two <- e }(e)
			wg.Wait()
		}
	}()
	return one, two
}

// // Teen returns n mutually indipendent copies of the source channel.
// // Teen does not implement channel closure
// func Teen[T any](src <-chan T, n int) []<-chan T {
// 	out := make([]chan T, n)
// 	w := waiter.New(n)
// 	for i := range out {
// 		out[i] = make(chan T, cap(src))
// 	}
// 	go func() {
// 		for e := range src {
// 			w.Add(1)
// 			for i, ch := range out {
// 				go func(i int, ch chan T, t T) { defer w.DoneAt(i); ch <- t }(i, ch, e)
// 			}
// 		}
// 		w.Wait()
// 	}()
// 	return Castro(out)
// }

// Close closes a channel
func Close[T any](ch chan T) { close(ch) }

// Bytes yields each byte in s
func Bytes(s string) <-chan byte {
	out := make(chan byte)
	go func() {
		defer close(out)
		for _, c := range []byte(s) {
			out <- c
		}
	}()
	return out
}

// Runes yields each rune in s
func Runes(s string) <-chan rune {
	out := make(chan rune)
	go func() {
		defer close(out)
		for _, c := range s {
			out <- c
		}
	}()
	return out
}
