package iters

import (
	"math/rand"

	"github.com/kendfss/rules"
)

// ints produces a slice with a given number of pseudo-randomly generated integers
func ints(count int) []int {
	out := make([]int, count)
	for i := range out {
		out[i] = rand.Int()
	}
	return out
}

// Slice simply returns its arguments in a slice
func Slice[T any](args ...T) []T { return args }

// Eq checks the given arguments for equality
func Eq[T comparable](a, b T) bool { return a == b }

// Neq checks the given arguments for inequality
func Neq[T comparable](a, b T) bool { return a != b }

// Lt checks that the first argument is less than the second argument
func Lt[T rules.Ordered](a, b T) bool { return a < b }

// Not returns the negation of the given predicate
func Not[T any](pred func(T, T) bool) func(T, T) bool {
	return func(t1, t2 T) bool {
		return !pred(t1, t2)
	}
}

// And returns a predicate that seeks the satisfaction of all arguments
func And[T any](args ...func(T, T) bool) func(T, T) bool {
	return func(t1, t2 T) bool {
		for _, arg := range args {
			if !arg(t1, t2) {
				return false
			}
		}
		return true
	}
}

// Or returns a predicate that seeks the satisfaction of any argument
func Or[T any](args ...func(T, T) bool) func(T, T) bool {
	return func(t1, t2 T) bool {
		for _, arg := range args {
			if arg(t1, t2) {
				return true
			}
		}
		return false
	}
}
