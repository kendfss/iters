// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"math/bits"

	"github.com/kendfss/rules"
)

// Sort sorts a slice of any ordered type in ascending order.
// Sort may fail to sort correctly when sorting slices of floating-point
// numbers containing Not-a-number (NaN) values.
// Use slices.SortFunc(x, func(a, b float64) bool {return a < b || (math.IsNaN(a) && !math.IsNaN(b))})
// instead if the input may contain NaNs.
func Sort[E rules.Ordered](x []E) {
	n := len(x)
	pdqsortOrdered(x, 0, n, bits.Len(uint(n)))
}

// Sorted sorts a slice of any ordered, type after cloning it, in ascending order.
// Sort may fail to sort correctly when sorting slices of floating-point
// numbers containing Not-a-number (NaN) values.
// Use slices.SortFunc(x, func(a, b float64) bool {return a < b || (math.IsNaN(a) && !math.IsNaN(b))})
// instead if the input may contain NaNs.
func Sorted[E rules.Ordered](x []E) []E {
	y := Clone(x)
	Sort(y)
	return y
}

// SortFunc sorts the slice x in ascending order as determined by the less function.
// This sort is not guaranteed to be stable.
//
// SortFunc requires that less is a strict weak ordering.
// See https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.
func SortFunc[E any](less func(a, b E) bool, x []E) {
	n := len(x)
	pdqsortLessFunc(x, 0, n, bits.Len(uint(n)), less)
}

// SortedFunc sorts a clone of the slice x in ascending order as determined by the less function.
// This sort is not guaranteed to be stable.
//
// SortFunc requires that less is a strict weak ordering.
// See https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.
func SortedFunc[E any](less func(a, b E) bool, x []E) []E {
	y := Clone(x)
	SortFunc(less, y)
	return y
}

// SortKey wraps a Key with a less than (<) function before deferring to SortFunc
// see slices.Key for more info
func SortKey[E any, O rules.Ordered](key func(E) O, arg []E) {
	k := Key[E, O](key)
	SortFunc(k.Lt, arg)
}

// SortStable sorts the slice x while keeping the original order of equal
// elements, using less to compare elements.
func SortStableFunc[E any](less func(a, b E) bool, x []E) {
	stableLessFunc(x, len(x), less)
}

// SortStableKey accepts a measuring key and calls SortStableFunc
func SortStableKey[E any, O rules.Ordered](key func(E) O, data []E) {
	k := Key[E, O](key)
	SortStableFunc(k.Lt, data)
}

// IsSorted reports whether x is sorted in ascending order.
func IsSorted[E rules.Ordered](x []E) bool {
	for i := len(x) - 1; i > 0; i-- {
		if x[i] < x[i-1] {
			return false
		}
	}
	return true
}

// IsSortedFunc reports whether x is sorted in ascending order, with less as the
// comparison function.
func IsSortedFunc[E any](less func(a, b E) bool, x []E) bool {
	for i := len(x) - 1; i > 0; i-- {
		if less(x[i], x[i-1]) {
			return false
		}
	}
	return true
}

// IsSortedKey accepts a measuring key and calls IsSortedFunc
func IsSortedKey[E any, O rules.Ordered](key func(E) O, data []E) bool {
	k := Key[E, O](key)
	return IsSortedFunc(k.Lt, data)
}

// BinarySearch searches for target in a sorted slice and returns the position
// where target is found, or the position where target would appear in the
// sort order; it also returns a bool saying whether the target is really found
// in the slice. The slice must be sorted in increasing order.
func BinarySearch[E rules.Ordered](target E, space []E) (int, bool) {
	// search returns the leftmost position where f returns true, or len(x) if f
	// returns false for all x. This is the insertion position for target in x,
	// and could point to an element that's either == target or not.
	pos := search(len(space), func(i int) bool { return space[i] >= target })
	if pos >= len(space) || space[pos] != target {
		return pos, false
	} else {
		return pos, true
	}
}

// BinarySearchFunc works like BinarySearch, but uses a custom comparison
// function. The slice must be sorted in increasing order, where "increasing" is
// defined by cmp. cmp(a, b) is expected to return an integer comparing the two
// parameters: 0 if a == b, a negative number if a < b and a positive number if
// a > b.
func BinarySearchFunc[E any](cmp func(E, E) int, target E, space []E) (int, bool) {
	pos := search(len(space), func(i int) bool { return cmp(space[i], target) >= 0 })
	if pos >= len(space) || cmp(space[pos], target) != 0 {
		return pos, false
	} else {
		return pos, true
	}
}

// BinarySearchKey accepts a measuring key and calls BinarySearchFunc
func BinarySearchKey[E any, O rules.Ordered](key func(E) O, target E, space []E) (int, bool) {
	k := Key[E, O](key)
	return BinarySearchFunc(k.Cmp, target, space)
}

func search(n int, f func(int) bool) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if !f(h) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

type sortedHint int // hint for pdqsort when choosing the pivot

const (
	unknownHint sortedHint = iota
	increasingHint
	decreasingHint
)

// xorshift paper: https://www.jstatsoft.org/article/view/v008i14/xorshift.pdf
type xorshift uint64

func (r *xorshift) Next() uint64 {
	*r ^= *r << 13
	*r ^= *r >> 17
	*r ^= *r << 5
	return uint64(*r)
}

func nextPowerOfTwo(length int) uint {
	return 1 << bits.Len(uint(length))
}
