package chans

import (
	"fmt"
	"sync"

	"github.com/kendfss/but"
	"github.com/kendfss/rules"
)

var DefaultCapacity = 0

// Make creates a buffered channel of given capacity
// or an unbuffered channel if the capacity is negative
func Make[T any, cap rules.OrderedNumber](args ...cap) chan T {
	if c := ParseCap(args...); c > 0 {
		return make(chan T, c)
	}
	return make(chan T)

}

// ParseCap helps you to anticipate the behaviour of
// functions with a "args ...cap" parameter
func ParseCap[cap rules.OrderedNumber](args ...cap) uint64 {
	if len(args) > 0 {
		return uint64(args[0])
	}
	return 0
}

func Inf[T any, cap rules.OrderedNumber](init func() T, args ...cap) chan T {
	out := make(chan T, ParseCap(args...))
	if init == nil {
		init = func() T { return *new(T) }
	}
	go func() {
		defer close(out)
		for {

		}
	}()
	return out
}

// Process consumes a channel
func Process[T any](c chan T) {
	for range c {
	}
}

func Filter(ch chan bool) chan bool {
	out := make(chan bool, DefaultCapacity)
	go func() {
		for b := range ch {
			if b {
				out <- b
			}
		}
	}()
	return out
}

func FilterPred[T any](pred func(T) bool, ch chan T) chan T {
	out := make(chan T, DefaultCapacity)
	go func() {
		defer close(out)
		for e := range ch {
			if pred(e) {
				out <- e
			}
		}
	}()
	return out
}

// RW wraps a read-only channel with a read-write one
func RW[T any](c <-chan T) chan T {
	out := make(chan T, cap(c))
	go func() {
		defer close(out)
		for x := range c {
			out <- x
		}
	}()
	return out
}

// RO wraps a read-write channel with a read-only one
func RO[T any](c chan T) <-chan T {
	out := make(chan T, cap(c))
	go func() {
		defer close(out)
		for x := range c {
			out <- x
		}
	}()
	return out
}

func Count[T any](c chan T) (out uint64) {
	for range c {
		out++
	}
	return
}

// remove all duplicates from a channel
func Compact[T comparable](ch chan T) chan T {
	out := make(chan T, DefaultCapacity)
	go func() {
		marked := []T{}
		pred := func(arg T) bool {
			if sliceContains(marked, arg) {
				return true
			}
			marked = append(marked, arg)
			return false
		}
		out = FilterPred(pred, ch)
	}()
	return out
}

// remove all duplicates from a channel of a non-comparable type
func CompactFunc[T comparable](eq func(T, T) bool, ch chan T) chan T {
	out := make(chan T, DefaultCapacity)
	go func() {
		marked := []T{}
		pred := func(arg T) bool {
			if sliceContainsFunc(eq, marked, arg) {
				return true
			}
			marked = append(marked, arg)
			return false
		}
		out = FilterPred(pred, ch)
	}()
	return out
}

// Send calls a function on every value of a slice
func Do[T any](f func(T), ch <-chan T) {
	go func() {
		for e := range ch {
			f(e)
		}
	}()
}

// Cast calls a pure function on every value of a channel and returns a channel
// containing all the results
func Cast[I, O any](f func(I) O, ch <-chan I) chan O {
	out := make(chan O, DefaultCapacity)
	go func() {
		for e := range ch {
			out <- f(e)
		}
	}()
	return out
}

func StepStr[T rules.Char](arg string) chan T {
	out := make(chan T)
	go func() {
		for _, char := range []rune(arg) {
			out <- T(char)
		}
		close(out)
	}()
	return out
}

// Chain collects several channels and returns one populated by their content
func Chain[T any](args ...chan T) <-chan T {
	out := make(chan T)

	go func() {
		wg := new(sync.WaitGroup)
		for _, c := range args {
			wg.Add(1)
			go func(c chan T) {
				defer wg.Done()
				for e := range c {
					out <- e
				}
			}(c)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

// Extend the first argument with the contents of the successors
// non blocking, non order-preserving
func Extend[T any](receiver chan T, args ...<-chan T) {
	wg := new(sync.WaitGroup)
	go func() {
		for _, arg := range args {
			go func(arg <-chan T) {
				wg.Add(1)
				for e := range arg {
					receiver <- e
				}
				wg.Done()
			}(arg)
		}
		wg.Wait()
	}()
}

func Extender[T any](target chan T) func(...<-chan T) {
	return func(args ...<-chan T) {
		Extend(target, args...)
	}
}

func Lazify[T any](arg []T) <-chan T {
	out := make(chan T)
	go func() {
		for _, e := range arg {
			out <- e
		}
	}()
	return out
}

// Upto returns an iterator whose content depends on the number of arguments as follows
// 		# of args 	|| 	behaviour
//	 		 1	 	|| 	stop
//		 	 2	 	|| 	start, stop
//		 	 3	 	|| 	start, stop, step
//         else 	|| 	error
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
	if err == nil {
		return out
	}
	panic(err)
}

// Put writes "count" copies of zero-initialized "T" instances to "ch"
func Put[T any, int rules.Int](count int, ch chan T) {
	for ; count > 0; count-- {
		ch <- *new(T)
	}
}

// PutVal writes a given value to a channel. Useful for spawning go routines before returning
func PutVal[T any](ch chan T, val T) {
	ch <- val
}

// Get receives (discards) "count" items from "ch"
func Get[T any, int rules.Int](count int, ch chan T) {
	for ; count > 0; count-- {
		<-ch
	}
}

// Watch feeds dst with items received from src
// does not close either of them
func Watch[T any](dst, src chan T) {
	for e := range src {
		dst <- e
	}
}

// Putter returns a method of the given channel which sends it the given argument
func Putter[T any](dst chan T) func(T) error {
	return func(arg T) error {
		var out error
		defer func() {
			if err := recover(); err != nil {
				out = err.(error)
			}
		}()

		dst <- arg

		return out
	}
}

// PredPutter returns a method of the given channel which sends it the given argument
// if, and only if, the argument satisfies the given predicate
// the put-method returns ErrUnsatisfied if the predicate is not satisfied
func PredPutter[T any](dst chan T, pred func(T) bool) func(T) error {
	put := Putter(dst)
	return func(arg T) error {
		if pred(arg) {
			return put(arg)
		}
		return ErrUnsatisfied
	}
}

var ErrUnsatisfied = but.New("Predicate was not satisfied")
