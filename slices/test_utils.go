package slices

import (
	"math/rand"

	"github.com/kendfss/rules"
)

func eqComparable[T comparable](a, b T) bool {
	return a == b
}

func randSign[N rules.Signed | rules.Float | rules.Complex](n N) N {
	var sgn N
	if rand.Intn(2) == 1 {
		sgn = 1
	} else {
		sgn = -1
	}
	return sgn * n
}

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

// equalNaN is like == except that all NaNs are equal.
func equalNaN[T comparable](v1, v2 T) bool {
	isNaN := func(f T) bool { return f != f }
	return v1 == v2 || (isNaN(v1) && isNaN(v2))
}

// offByOne returns true if integers v1 and v2 differ by 1.
func offByOne[E rules.Integer](v1, v2 E) bool {
	return v1 == v2+1 || v1 == v2-1
}

func equalToCmp[T comparable](eq func(T, T) bool) func(T, T) int {
	return func(v1, v2 T) int {
		if eq(v1, v2) {
			return 0
		}
		return 1
	}
}

func cmp[T rules.Ordered](v1, v2 T) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	} else {
		return 0
	}
}

func equalToIndex[T any](f func(T, T) bool, v1 T) func(T) bool {
	return func(v2 T) bool {
		return f(v1, v2)
	}
}
