package slices

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"unsafe"

	"github.com/kendfss/but"
	"github.com/kendfss/oprs"
	"github.com/kendfss/oprs/math/real"
	"github.com/kendfss/rules"
)

// Equal reports whether two slices are equal: the same length and all
// elements equal. If the lengths are different, Equal returns false.
// Otherwise, the elements are compared in increasing index order, and the
// comparison stops at the first unequal pair.
// Floating point NaNs are not considered equal.
func Equal[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// EqualFunc reports whether two slices are equal using a comparison
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func EqualFunc[E1, E2 any](eq func(E1, E2) bool, s1 []E1, s2 []E2) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v1 := range s1 {
		v2 := s2[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}

// Compare compares the elements of s1 and s2.
// The elements are compared sequentially, starting at index 0,
// until one element is not equal to the other.
// The result of comparing the first non-matching elements is returned.
// If both slices are equal until one of them ends, the shorter slice is
// considered less than the longer one.
// The result is 0 if s1 == s2, -1 if s1 < s2, and +1 if s1 > s2.
// Comparisons involving floating point NaNs are ignored.
func Compare[E rules.Ordered](s1, s2 []E) int {
	s2len := len(s2)
	for i, v1 := range s1 {
		if i >= s2len {
			return +1
		}
		v2 := s2[i]
		switch {
		case v1 < v2:
			return -1
		case v1 > v2:
			return +1
		}
	}
	if len(s1) < s2len {
		return -1
	}
	return 0
}

// CompareFunc is like Compare but uses a comparison function
// on each pair of elements. The elements are compared in increasing
// index order, and the comparisons stop after the first time cmp
// returns non-zero.
// The result is the first non-zero result of cmp; if cmp always
// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func CompareFunc[E1, E2 any](cmp func(E1, E2) int, s1 []E1, s2 []E2) int {
	s2len := len(s2)
	for i, v1 := range s1 {
		if i >= s2len {
			return +1
		}
		v2 := s2[i]
		if c := cmp(v1, v2); c != 0 {
			return c
		}
	}
	if len(s1) < s2len {
		return -1
	}
	return 0
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](val E, s []E) int {
	return IndexFunc(oprs.Eq[E], val, s)
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[E any](eq func(E, E) bool, val E, s []E) int {
	return IndexPred(oprs.Method(val, eq), s)
}

// IndexPred returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexPred[E any](eq func(E) bool, s []E) int {
	for i, v := range s {
		if eq(v) {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[E comparable](s []E, v E) bool {
	return Index(v, s) >= 0
}

// ContainsFunc reports whether v is present in s, using eq as an equivalence operator.
func ContainsFunc[E any](eq func(E, E) bool, s []E, v E) bool {
	return IndexFunc(eq, v, s) >= 0
}

// ContainsPred reports whether anything in s satisfies the given predicate.
func ContainsPred[E any](pred func(E) bool, s []E) bool {
	return IndexPred(pred, s) >= 0
}

// ContainsAny reports whether any element of args is present in s.
func ContainsAny[E comparable](s []E, args ...E) bool {
	return ContainsAnyFunc(oprs.Eq[E], s, args...)
}

// ContainsAnyFunc reports whether of args is present in s.
func ContainsAnyFunc[E any](eq func(E, E) bool, s []E, args ...E) bool {
	for _, arg := range args {
		if ContainsFunc(eq, s, arg) {
			return true
		}
	}
	return false
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// In the returned slice r, r[i] == v[0].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
func Insert[E any](s []E, i int, args ...E) []E {
	tot := len(s) + len(args)
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(args):], s[i:])
		copy(s2[i:], args)
		return s2
	}
	s2 := make([]E, tot)
	copy(s2, s[:i])
	copy(s2[i:], args)
	copy(s2[i+len(args):], s[i:])
	return s2
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-(j-i)), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
func Delete[E any](s []E, i, j int) []E {
	return append(s[:i], s[j:]...)
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[E any](s []E) []E {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	out := make([]E, len(s))
	copy(out, s)
	return out
	// return append([]E{}, s...)
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice s; it does not create a new slice.
func Compact[E comparable](s []E) []E {
	if len(s) == 0 {
		return s
	}
	i := 1
	last := s[0]
	for _, v := range s[1:] {
		if v != last {
			s[i] = v
			i++
			last = v
		}
	}
	return s[:i]
}

// Compacted clones the slice and runs Compact on said clone
func Compacted[E comparable](s []E) []E {
	c := Clone(s)
	return Compact(c)
}

// CompactFunc is like Compact but uses a comparison function.
func CompactFunc[E any](eq func(E, E) bool, s []E) []E {
	if len(s) == 0 {
		return s
	}
	i := 1
	last := s[0]
	for _, v := range s[1:] {
		if !eq(v, last) {
			s[i] = v
			i++
			last = v
		}
	}
	return s[:i]
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. Grow may modify elements of the
// slice between the length and the capacity. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[E any](s []E, n int) []E {
	return append(s, make([]E, n)...)[:len(s)]
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[E any](s []E) []E {
	return s[:len(s):len(s)]
}

// Returns the index of the shortest slice received at call-time
// -1 if no arguments are passed
func Shortest[E any](args ...[]E) (out int) {
	switch len(args) {
	case 0:
		out--
	case 1:
	default:
		for i, arg := range args {
			if len(arg) < len(args[out]) {
				out = i
			}
		}
	}
	return
}

// Returns the index of the longest slice received at call-time
// -1 if no arguments are passed
func Longest[E any](args ...[]E) (out int) {
	switch len(args) {
	case 0:
		out--
	case 1:
	default:
		for i, arg := range args {
			if len(arg) > len(args[out]) {
				out = i
			}
		}
	}
	return
}

// Cast returns a slice whose values are the result of the
// application of the given function to all elements of the given slice
// it behaves like "map" in languages whose hashtables are called "associative array" or "dictionary"
func Cast[E, V any](f func(E) V, s []E) []V {
	out := make([]V, len(s))
	for i, e := range s {
		out[i] = f(e)
	}
	return out
}

// Filter returns a slice featuring all truthy elements
func Filter(args []bool) (out []bool) {
	for _, e := range args {
		out = append(out, e)
	}
	return out
}

// FilterFunc returns a slice featuring all elements of the incident that satisfy the given predicate
func FilterFunc[E any](f func(E) bool, args []E) (out []E) {
	for _, e := range args {
		if f(e) {
			out = append(out, e)
		}
	}
	return out
}

// Get returns the i'th element from a slice, even if i is negative
// uses the same indexing convention as python lists/tuples
func Get[E any, I rules.Integer](index I, slice []E) E {
	if index >= 0 {
		return slice[index]
	}
	return slice[len(slice)+int(index)]
}

// Reduce returns the outcome of successive applications of
// a function, f, as a binary operator over the slice, s.
func Reduce[E any](f func(E, E) E, s []E) (out E) {
	switch len(s) {
	case 0:
		out = *new(E)
	case 1:
		out = s[0]
	default:
		out = s[0]
		for _, e := range s[1:] {
			out = f(out, e)
		}
	}
	return out
}

// Trot returns the outcome of step-wise applications of
// a function, f, as a binary operator over the slice, s.
// Trot{addition, {1, 2, 3}} == {1, 1, 1}
func Trot[T, O any](operator func(T, T) O, data []T) []O {
	rack := make([]O, len(data)-1)
	for i, datum := range data[1:] {
		rack[i] = operator(data[i], datum)
	}
	return rack
}

// Convolve type-equivalent slices
func Zip[K any](args ...[]K) (out [][]K) {
	min := Shortest(args...)
	l := len(args[min])
	out = make([][]K, l)
	if min > -1 {
		for i := range out {
			out[i] = make([]K, len(args))

			for j, arg := range args {
				out[i][j] = arg[i]
			}
		}
		// for i, arg := range args {
		// 	out[i] = arg[:l]
		// }
	}
	return
}

type (
	LR[L, R any] struct {
		// LR holds two values, Left and Right, of any types.
		Left  L
		Right R
	}
	Pair[T any] struct {
		// Pair holds two values, Left and Right, of any type.
		Left, Right T
	}
)

func (p Pair[T]) Split() (l, r T) {
	return p.Left, p.Right
}

func (lr LR[L, R]) Split() (l L, r R) {
	return lr.Left, lr.Right
}

func (p Pair[T]) From(l, r T) Pair[T] {
	p.Left, p.Right = l, r
	return p
}

func (lr LR[L, R]) From(l L, r R) LR[L, R] {
	lr.Left, lr.Right = l, r
	return lr
}

func (p Pair[T]) L() T {
	return p.Left
}

func (p Pair[T]) R() T {
	return p.Right
}

func (lr LR[L, R]) L() L {
	return lr.Left
}

func (lr LR[L, R]) R() R {
	return lr.Right
}

func ReducePairLeft[T any](op func(l, r T) T, rack []Pair[T]) T {
	return Reduce(op, Cast(Pair[T].L, rack))
}

func ReducePairRight[T any](op func(l, r T) T, rack []Pair[T]) T {
	return Reduce(op, Cast(Pair[T].R, rack))
}

// func (P)

// Convolve pairs of type-distinct slices with a Pair
func Zip2[L, R any](left []L, right []R) (out []LR[L, R]) {
	if len(left) > len(right) {
		out = make([]LR[L, R], len(right))
	} else {
		out = make([]LR[L, R], len(left))
	}
	for i := range out {
		out[i] = LR[L, R]{Left: left[i], Right: right[i]}
	}
	return out
}

// Convolve pairs of type-distinct slices with a closure
func Zip3[L, R any](left []L, right []R) (out []func() (L, R)) {
	if len(left) > len(right) {
		out = make([]func() (L, R), len(right))
	} else {
		out = make([]func() (L, R), len(left))
	}
	for i := range out {
		out[i] = func() (L, R) {
			return left[i], right[i]
		}
	}
	return out
}

// Concatenate slices
func Chain[E any](args ...[]E) (out []E) {
	for _, arg := range args {
		out = append(out, arg...)
	}
	return out
}

// Return r length subsequences of elements from the input
// empty if r > len(slice) || r < 0
//
// The combination tuples are emitted in lexicographic ordering according to
// the order of the input iterable. So, if the input iterable is sorted,
// the combination tuples will be produced in sorted order.
//
// Elements are treated as unique based on their position, not on their value.
// So if the input elements are unique, there will be no repeat values
// in each combination.
// func Combinations[E any, U rules.Unsigned](slice []E, r U) ([][]E, error) {
// 	if r < 0 || int(r) > len(slice) {
// 		return nil, ErrIndex
// 	}
// 	pool := Clone(slice)
// 	n := Len[U](pool)

// 	indices := Upton[U](r)
// 	for {
// 		for ii, i := range Reversed(Upton[U](r)) {
// 			// if indices[i] != uint(i)+n-r {
// 			if indices[i] != (i)+n-r {
// 				break
// 			} else if ii == len(slice)-1 {
// 				indices[i]++
// 				for _, j := range Uptonm[U](i+1, r) {
// 					indices[j] = indices[j-1] + 1
// 				}
// 			}
// 		}
// 	}
// }

// Indices of a slice with given length
// Range[byte](256)
func Upton[O, I rules.Real](stop I) []O {
	return Upto[O](0, stop, 1)
}

// Consecutive ints, including start, smaller than stop, and separated by one
// Uptonm[byte](0, 256)
func Uptonm[O, I rules.Real](start, stop I) []O {
	return Upto[O](start, stop, 1)
}

// Consecutive ints, including start, smaller than stop, and separated by step
// Upto[byte](0, 256, 1)
func Upto[O, I rules.Real](start, stop, step I) []O {
	if stop < start && step >= 0 {
		panic(but.New("start %v exceeds stop %v but step %v is non-negative", start, stop, step))
		// 	return
	}
	out := []O{}
	if start <= stop {
		for i := O(start); i < O(stop); i += O(step) {
			out = append(out, i)
		}
	} else {
		for i := O(start); i > O(stop); i += O(step) {
			out = append(out, i)
		}
	}

	return out
}

// Produce a reversed copy of a slice
func Reversed[E any](slice []E) []E {
	// if len(slice) <= 1 {
	// 	return slice
	// }
	// out := make([]E, len(slice))
	// for i, j := 0, len(slice)-1; j != 0; i, j = i+1, j-1 {
	// 	out[i], out[j] = slice[j], slice[i]
	// }
	// if l := len(slice); l%2 == 1 {
	// 	out[l/2+1] = slice[l/2+1]
	// }
	out := Clone(slice)
	Reverse(out)
	return out
}

// Reverse a slice in place
// func Reverse[[]E ~[]E, E any](slice []E) {
func Reverse[E any](slice []E) {
	if len(slice) > 1 {
		for i, j := 0, len(slice)-1; j > 0; i, j = i+1, j-1 {
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
}

// Swap the elements at a pair of indices (in place)
func Swap[E any](slice []E, i, j int) []E {
	slice[i], slice[j] = slice[j], slice[i]
	return slice
}

// Swap the elements at a pair of indices (copied)
func Swapped[E any](slice []E, i, j int) []E {
	out := append([]E{}, slice...)
	out[i], out[j] = slice[j], slice[i]
	return out
}

// Len returns the length of a slice as the desired type of integer
func Len[I rules.Integer, E any](slice []E) I {
	return I(len(slice))
}

// Check if all elements of a slice are true
func All(args []bool) bool {
	// return AllFunc(oprs.IsTrue, args...)
	switch len(args) {
	case 0:
		return false
	case 1:
		return args[0]
	default:
		return args[0] && All(args[1:])
	}
}

// use a custom predicate to check if all elements of a slice
// have a common property
func AllFunc[E any](pred func(E) bool, slice []E) bool {
	return All(Cast(pred, slice))
}

// Check if any elements of a slice are true
func Any(args ...bool) bool {
	return AnyFunc(oprs.IsTrue, args...)
}

// use a custom predicate to check if any elements of a slice
// have a common property
func AnyFunc[E any](pred func(E) bool, slice ...E) bool {
	for _, e := range slice {
		if pred(e) {
			return true
		}
	}
	return false
}

// Max returns the index of the Maximal value of a slice
func Max[E rules.Ordered](args ...E) (out int) {
	for i, arg := range args {
		if arg > args[out] {
			out = i
		}
	}
	return
}

// Min returns the index of the Minimal value of a slice
func Min[E rules.Ordered](args ...E) (out int) {
	for i, arg := range args {
		if arg < args[out] {
			out = i
		}
	}
	return
}

// Extremal finds the index of a maximum, or minimum, value of a slice
// by passing a key corresponding to greater than or less than
// Extremal[MyType](gt, mySlice...) -> maximal value
// Extremal[MyType](lt, mySlice...) -> minimal value
func Extremal[E any](operator func(E, E) bool, args ...E) (out int) {
	for i, arg := range args {
		if operator(arg, args[out]) {
			out = i
		}
	}
	return out
}

// Deprecated, use Chain
func Union[E any](first []E, rest ...[]E) []E {
	fmt.Fprintln(os.Stderr, "Union is deprecated, use Chain")
	return Chain(first, Chain(rest...))
}

// Deprecated, use Chain
func Flatter[M ~[][]E, E any](arg M) (out []E) {
	fmt.Fprintln(os.Stderr, "Flatter is deprecated, use Chain")
	return Chain(arg...)
}

// Snap breaks a slice into sections of given width
// Snap(2, []int{1, 2, 3, 4}) == [][]int{{1, 2}, {3, 4}}
// Snap(3, []int{1, 2, 3, 4}) == [][]int{{1, 2, 3}, {4}}
// func Snap[[]E ~[]E, E any](arg []E, width int) (out [][]E) {
func Snap[I rules.I, E any](width I, arg []E) (out [][]E) {
	if width == 0 {
		return [][]E{arg[:0], arg[0:]}
	}
	for i_, e := range arg {
		i := I(i_)
		if i%width == 0 {
			out = append(out, []E{})
		}
		ind := len(out) - 1
		out[ind] = append(out[ind], e)
	}
	return out
}

// Split "cuts" the slice at all occurrences of breaker
func Split[E comparable](slice []E, breaker E) [][]E {
	return SplitFunc(oprs.Eq[E], slice, breaker)
}

// SplitFunc "cuts" the slice at all occurrences of breaker
func SplitFunc[E any](eq func(E, E) bool, slice []E, breaker E) [][]E {
	// pred := func(arg E) bool {
	// 	return eq(arg, breaker)
	// }
	pred := oprs.Method(breaker, eq)
	return SplitPred(pred, slice)
}

// SplitPred "cuts" the slice at all elements satisfying some predicate
func SplitPred[E any](pred func(E) bool, slice []E) [][]E {
	out := make([][]E, 1)
	for _, e := range slice {
		if pred(e) {
			out = append(out, []E{})
			continue
		}
		i := len(out) - 1
		out[i] = append(out[i], e)
	}
	return out
}

// SplitAfter "cuts" the slice at all matching elements without discarding them
func SplitAfter[E comparable](slice []E, breaker E) [][]E {
	return SplitAfterFunc(oprs.Eq[E], breaker, slice)
}

// SplitAfterFunc "cuts" the slice at all matching elements without discarding them
func SplitAfterFunc[E any](function func(E, E) bool, breaker E, slice []E) [][]E {
	return SplitAfterPred(oprs.Method(breaker, function), slice)
}

// SplitAfterPred "cuts" the slice at all satisfying elements without discarding them
func SplitAfterPred[E any](function func(E) bool, slice []E) [][]E {
	out := make([][]E, 1)
	for _, e := range slice {
		out[len(out)-1] = append(out[len(out)-1], e)
		if function(e) {
			out = append(out, []E{})
		}
	}
	return out
}

// Deprecated, use Repeat
func Ones[T rules.Integer](count T) []T {
	fmt.Fprintln(os.Stderr, "Ones is deprecated, use Repeat")
	return Repeat(T(1), count)
}

// Convert a slice of any type into one of empty interfaces
func Anify[E any](slice []E) []any {
	out := make([]any, len(slice))
	for i := range out {
		out[i] = slice[i]
	}
	return out
}

// Repeat returns a slice, with length count, of integers initialized to seed
// if you want to repeat an empty slice you should use Tee instead
func Repeat[T any, C rules.Integer](seed T, count C) []T {
	out := make([]T, count)
	for i := range out {
		out[i] = *(*T)(unsafe.Pointer(&seed))
	}
	return out
}

func Extend[T any, C rules.Integer](slice []T, seed T, count C) []T {
	return append(slice, Repeat(seed, count)...)
}

// Movers returns the indices of slice elements that are not equal to their successors
func Movers[E comparable](s []E) (out []int) {
	for i, e := range s[1:] {
		if e != s[i-1] {
			out = append(out, i)
		}
	}
	return out
}

// MoversFunc returns the indices of slice elements that are not equal to their successors
func MoversFunc[E any](eq func(E, E) bool, s []E) (out []int) {
	for i, e := range s[1:] {
		if eq(e, s[i-1]) {
			out = append(out, i)
		}
	}
	return out
}

// Standers returns the indices of the slice elements that are equal to their successors
func Standers[E comparable](s []E) (out []E) {
	for i, e := range s[1:] {
		if e == s[i-1] {
			out = append(out, e)
		}
	}
	return out
}

// StandersFunc returns the indices of the slice elements that are equal to their successors
func StandersFunc[E any](eq func(E, E) bool, s []E) (out []E) {
	for i, e := range s[1:] {
		if eq(e, s[i-1]) {
			out = append(out, e)
		}
	}
	return out
}

func VariadicFilter[E any](adicity int, walk bool, f func(...E) bool, slice []E) (out [][]E) {
	step := 1
	if !walk {
		step = adicity
	}

	for i := 0; i < (len(slice) - adicity); i += step {
		args := make([]E, adicity)
		for j := range args {
			args[j] = slice[i+j]
		}

		if f(args...) {
			out = append(out, args)
		}
	}

	return out
}

// select returns the elements of a slice located at the chosen indices
// note: all indices are wrapped by a modulus equal to the length of the slice
// use StrictSelect to mitigate this behaviour
func Select[E any](slice []E, indices []int) []E {
	out := make([]E, len(indices))
	for i, e := range indices {
		out[i] = slice[e%len(slice)]
	}
	return out
}

// select returns the elements of a slice located at the chosen indices
// note: indices greater than slice length with cause panic
// use Select to mitigate this behaviour
func SelectStrict[E any](slice []E, indices []int) []E {
	out := make([]E, len(indices))
	for _, e := range indices {
		out = append(out, slice[e])
	}
	return out
}

// Prefill prepends some number of elements to a slice
func Prefill[T any](s []T, by uint) []T {
	out := make([]T, by)
	return append(out, s...)
}

// PrefillSeed prepends some number of elements to a slice
func PrefillSeed[T any](slice []T, seed T, by uint) []T {
	out := make([]T, by)
	for i := range out {
		out[i] = seed
	}
	return append(out, slice...)
}

func Cartesian[L, R any](left []L, right []R) []LR[L, R] {
	out := make([]LR[L, R], len(left)*len(right))
	ctr := 0
	for _, l := range left {
		for _, r := range right {
			out[ctr].Left = l
			out[ctr].Right = r
			ctr++
		}
	}
	return out
}

// Count returns the number of occurences of item in rack
func Count[T comparable](item T, rack []T) (out uint) {
	for _, e := range rack {
		if e == item {
			out++
		}
	}
	return
}

// CountFunc returns the number of occurences, with respect to eq == true, of item in rack
func CountFunc[E any](eq func(E, E) bool, item E, rack []E) (out uint) {
	for _, e := range rack {
		if eq(item, e) {
			out++
		}
	}
	return
}

// Indices returns the positions at which item can be found in rack
func Indices[T comparable](item T, rack []T) (out []int) {
	for i, e := range rack {
		if e == item {
			out = append(out, i)
		}
	}
	return
}

// IndicesFunc returns the positions at which item can be found, by eq == true, in rack
func IndicesFunc[T comparable](eq func(T, T) bool, item T, rack []T) (out []int) {
	for i, e := range rack {
		if eq(e, item) {
			out = append(out, i)
		}
	}
	return
}

// Dot returns a dot product analog of left with right.
// Dot({2, 3}, {1, 2}) === {2, 6}
// Dot({2}, {1, 2}) === {2, 0}
// Dot({1, 2}, {2}) === {2, 0}
func Dot[N rules.Num](left, right []N) []N {
	if len(left) > len(right) {
		return Dot(right, left)
	}
	out := append(left, make([]N, len(right)-len(left))...)
	for i, e := range out {
		out[i] = e * right[i]
	}
	return out
}

// DotFunc returns a dot product analog of left with right,
// using mul as a binary operator over the chosen type.
func DotFunc[T any](mul func(T, T) T, left, right []T) []T {
	if len(left) > len(right) {
		return DotFunc(mul, right, left)
	}
	out := append(left, make([]T, len(right)-len(left))...)
	for i, e := range out {
		out[i] = mul(e, right[i])
	}
	return out
}

// Mul returns a dot product analog of left with right.
// Each argument of length 1 is treaded as a scalar
// Mul({2, 3}, {1, 2}) === {2, 6}
// Mul({2}, {1, 2}) === {2, 4}
// Mul({1, 2}, {2}) === {2, 4}, Mul(right, left) if len(left) > len(right)
// If you want ordinary Muls, use Cast(a+b, Zip(left, right))
func Mul[N rules.Num](left, right []N) []N {
	if len(left) > len(right) {
		return Mul(right, left)
	}
	out := append(left, make([]N, len(right)-len(left))...)
	for i, e := range out {
		out[i] = e * right[i]
	}
	if len(left) == len(right) {
		mul := func(rack []N) N {
			return rack[0] * rack[1]
		}
		return Cast(mul, Zip(left, right))
	}
	return out
}

// Rotated returns the index-shifted of a given slice
// as though the operation were taking place on a torus (no elements lost or added)
func Rotated[T any, I rules.I](slice []T, steps I) []T {
	if len(slice) == 0 {
		return make([]T, 0)
	}
	steps %= I(len(slice))
	if steps < 0 {
		steps += I(len(slice))
	}
	return append(slice[steps:], slice[:steps]...)
}

// Send is like Cast but for impure functions
func Send[T any](f func(T), args []T) {
	for _, arg := range args {
		f(arg)
	}
}

// Pointers returns an array of pointers to the values of given slice
// These pointers should not agree with other reference to the data
func Pointers[T any](s []T) []*T {
	out := make([]*T, len(s))
	for i, e := range s {
		out[i] = &e
	}
	return out
}

// Values returns a slice of values to members of a slice of pointers
func Values[T any](s []*T) []T {
	out := make([]T, len(s))
	for i, e := range s {
		out[i] = *e
	}
	return out
}

// Partition uses a function to categorize elements of a slice
func Partition[K comparable, V any](pred func(V) K, slice []V) map[K][]V {
	out := make(map[K][]V)
	for _, e := range slice {
		key := pred(e)
		out[key] = append(out[key], e)
	}
	return out
}

// Deprecated, use Repeat
func Copies[T any, U rules.I](length U, val T) []T {
	fmt.Fprintln(os.Stderr, "Copies is deprecated, use Repeat")
	return Repeat(val, length)
}

// Channel returns
// 		a buffered channel iff cap < 1
// 		an unbuffered channel otherwise
// use an empty slice if you want to reproduce "make(chan T, 0)"
func Channel[T any, Int rules.Integer](slice []T, cap Int) <-chan T {
	var out chan T
	if cap < 1 {
		out = make(chan T)
	} else {
		out = make(chan T, cap)
	}

	go func() {
		defer close(out)
		for _, e := range slice {
			out <- e
		}
	}()
	return out
}

// Consume passes each element of the given channel to the given slice
func Consume[T any](channel chan T) (out []T) {
	for e := range channel {
		out = append(out, e)
	}
	return out
}

// Feed passes each element of the given slice to the given channel
func Feed[T any](channel chan T, slice []T) {
	for _, e := range slice {
		channel <- e
	}
}

// Nest places a slice of T into a new Matrix of T
func Nest[T any](s ...T) [][]T {
	return [][]T{s}
}

// Break an iterable into len(iterable)-length steps of the given length, with each step's starting point one after its predecessor
// example
//	 >>> for i in walks(itertools.count(),2):print(''.join(i))
//	 (0, 1)
//	 (1, 2)
//	 (2, 3)
//	 # etc.
// Inspired by the hyperoperation 16**2[5]2
func Walks[T any, I rules.Integer](length I, slice []T) (out [][]T) {
	tee := Tee(slice, length)
	fmt.Println("tee:", tee)

	// fmt.Println("enum:", )
	for _, p := range Enumerate[int](tee) {
		n, it := p()
		fmt.Println("n, it:", n, it)
		if n == 0 {
			out = append(out, it)
			continue
		}
		out = append(out, it[n:])
	}
	fmt.Println("out:", out)
	return Zip(out...)
}

// Enumerate returns a slice of closures whose return values are tuples of
// elements of the given slice prefixed by their indices
func Enumerate[I rules.Integer, T any](slice []T) []func() (I, T) {
	out := make([]func() (I, T), len(slice))
	for i, e := range slice {
		out[i] = func() (I, T) {
			return I(i), e
		}
	}
	return out
}

// Tee returns a slice of independent slices
func Tee[T any, I rules.Integer](seed []T, count I) [][]T {
	out := make([][]T, count)
	for i, e := range out {
		out[i] = append(e, seed...)
	}
	return out
}

// Make initializes a slice
func Make[T any, I rules.Integer](length I) []T {
	return make([]T, length)
}

// Copy is a wrapper on the builtin copy
func Copy[T any](dst, src []T) int {
	return copy(dst, src)
}

func Permutations[T any](r int, pool []T) (out [][]T) {
	n := len(pool)
	// println(n - repeat)
	r = oprs.Ternary(n <= 0, n, r)
	if r > n {
		return
	}
	indices := Upton[int](n)
	cycles := Upto[int](n, n-r, -1)
	out = append(out, Select(pool, indices[:r]))
	// println("commencing loop")
	// outer:
	// for n > 0 {
	// for ; n > 0; n-- {
	for n != 0 {
		// println(n)
		fmt.Printf("indices: %v\n", indices)
		for pos, i := range Reversed(Upton[int](r)) {
			cycles[i]--
			if cycles[i] == 0 {
				// pos++
				// indices = append(indices[i:], append(indices[i+1:], indices[i:i+1]...)...)
				indices = Skip(pos, indices)
				cycles[i] = n - i
			} else {
				j := cycles[i]
				fmt.Printf("preswap: %v %v\n", indices[i], indices[len(indices)-j])
				indices[i], indices[len(indices)-j] = indices[len(indices)-j], indices[i]
				fmt.Printf("postswap: %v %v\n", indices[i], indices[len(indices)-j])
				out = append(out, Select(pool, indices[:r-1]))
				goto next
			}
		}
		return
	next:
	}
	return
}

// func Permutations[T any](arg []T) (out [][]T) {
// 	// // # If the length of list=0 no permuataions possible
// 	// if len(arg) == 0 {
// 	// 	return [][]T{{}}
// 	// }
// 	// # If the length of list=1, return that element
// 	if len(arg) <= 1 {
// 		return [][]T{arg}
// 	}
// 	for i, m := range arg {
// 		// # Extract list1[i] or m from the list. rem is
// 		// # remaining list
// 		fmt.Printf("i := %v\n", i)
// 		fmt.Printf("arg[:i] := %v\n", arg[:i])
// 		fmt.Printf("arg[i+1:] := %v\n", arg[i+1:])
// 		// rem := append(arg[:i], arg[i+1:]...)
// 		rem := Skip(arg, i)
// 		fmt.Printf("m, rem := %#v, %#v\n", m, rem)
// 		println("")
// 		// # Generating all permutations where m is first
// 		// # element
// 		for _, p := range Permutations(rem) {
// 			out = append(out, append([]T{m}, p...))
// 		}
// 		println("................................................................")

// 	}
// 	return
// }

func Product[T any](repeat int, args ...[]T) [][]T {
	pools := Chain(Tee(args, repeat+1)...)
	out := make([][]T, 1)
	for _, pool := range pools {
		for _, y := range pool {
			for _, x := range out {
				out = append(out, append(x, y))
			}
		}
	}
	return out
}

// Skip returns a version of the slice without the element at the given index
// returns arg if index is greater than len(arg)-1
func Skip[T any, I rules.Int](index I, arg []T) (out []T) {
	if uint64(index) > uint64(len(arg)-1) {
		copy(arg, out)
		return
	}
	for i_, e := range arg {
		i := I(i_)
		if i != index {
			out = append(out, e)
		}
	}

	return out
}

// Return r length subsequences of elements from the input
// empty if r > len(slice) || r < 0
//
// The combination tuples are emitted in lexicographic ordering according to
// the order of the input iterable. So, if the input iterable is sorted,
// the combination tuples will be produced in sorted order.
//
// Elements are treated as unique based on their position, not on their value.
// So if the input elements are unique, there will be no repeat values
// in each combination.
// Combinations('ABCD', 2) --> AB AC AD BC BD CD
// Combinations(range(4), 3) --> 012 013 023 123
func Combinations[T any](pool []T, r int) (out [][]T) {
	n := len(pool)
	for _, indices := range Permutations(r, Upton[int](n)) {
		if Equal(Sorted(indices), indices) {
			out = append(out, Select(pool, indices))
		}
	}
	return out
}

// func Combinations[T any](pool []T, r int) (out [][]T) {
// 	n := len(pool)
// 	if r > n {
// 		return
// 	}
// 	indices := Upton[int](r)
// 	// out = append(out, )
// 	// yield tuple(pool[i] for i in indices)
// 	out = append(out, Select(pool, indices))

// 	// defer func() { out = append(out, Select(pool, indices)) }()
// 	for {
// 		i := 0
// 		if r > 0 {
// 			for _, i_ := range Reversed(Upton[int](r)) {
// 				i = i_
// 				if indices[i_] != i_+n-r {
// 					// println("breaking")
// 					fmt.Printf("breaking!\ti: %v\n", i)
// 					goto rest
// 				}
// 			}
// 			return
// 		}
// 	rest:
// 		// println(i)
// 		indices[i] += 1
// 		for _, j := range Uptonm[int](i+1, r) {
// 			indices[j] = indices[j-1] + 1
// 		}
// 		// yield tuple(pool[i] for i in indices)
// 		out = append(out, Select(pool, indices))
// 		// println("")
// 	}
// }

// PairAll returns a sequence of pairs from a matrix
func PairAll[T any](arg [][]T, shift int) (out []Pair[T]) {
	if !AllFunc(
		oprs.Method(2, oprs.Le[int]),
		Cast(Len[int, T], arg),
	) {
		return nil
	}
	out = make([]Pair[T], len(arg))
	for i, pair := range arg {
		out[i].Left = pair[0+shift]
		out[i].Left = pair[1+shift]
	}
	return out
}

// Reducer returns a castable operator that reduces a slice.
// see Reduce and Cast for more info
func Reducer[T any](op func(T, T) T) func([]T) T {
	return oprs.MethodOp(op, Reduce[T])
}

// Getter returns a castable operator that fetches the element of that index from a slice.
// see Get and Cast for more info
func Getter[T any, I rules.Int](index I) func([]T) T {
	return oprs.Method(index, Get[T, I])
}

// Prefiller returns a castable operator for prefilling slices
// see Prefill and Cast for more info
func Prefiller[T any](by uint) func([]T) []T {
	return func(t []T) []T {
		return Prefill(t, by)
	}
}

// PrefillSeeder returns a castable operator for prefilling slices
// see PrefillSeed and Cast for more info
func PrefillSeeder[T any](seed T, by uint) func([]T) []T {
	return func(t []T) []T {
		return PrefillSeed(t, seed, by)
	}
}

// Caster integrates a castable operator for use with a matrix
// see Cast for more info
func Caster[I, O any](op func(I) O) func([]I) []O {
	return func(arg []I) []O {
		return Cast(op, arg)
	}
}

// Snapper returns a castable operator for snapping slices
// see Snap and Cast for more info
func Snapper[T any, I rules.Int](stride I) func([]T) [][]T {
	return func(t []T) [][]T {
		return Snap(stride, t)
	}
}

// Repeater returns a castable operator for creating slices of given length populated by given value
// see Repeat and Cast for more info
func Repeater[T any, I rules.Int](count I) func(T) []T {
	return func(seed T) []T {
		return Repeat(seed, count)
	}
}

// Choose selects an element of the gicen slice at random
func Choose[T any](arg []T) T {
	return arg[rand.Intn(len(arg))]
}

// Get an element of a slice situated at a point (x,y) when
// the slice is interpreted as
func Getxy[E any](slice []E, stride, x, y int) E {
	return slice[y*stride+x]
}

// ReduceAs applies Reduce after converting a slice of real numbers
// an overflow-safe way for operating on small numbers
func ReduceAs[I, O rules.Real](op func(O, O) O, args ...I) O {
	rack := make([]O, len(args))
	for i, arg := range args {
		rack[i] = O(arg)
	}
	return Reduce(op, rack)
}

func Windows[T any](src []T, size int) (out [][]T) {
	if size > 0 {
		for i := 0; i+size <= len(src); i++ {
			out = append(out, src[i:i+size])
		}
	}
	return out
}

func Resize[T any](s []T, shape ...int) []T {
	dim := Reduce(real.Mul[int], shape)
	switch l := len(s); cmp(dim, l) {
	case 1:
		return s[:dim]
	case -1:
		return append(s, make([]T, l-dim)...)
	default:
		return s
	}
}

// Show prints each element of a slice to a stdout on a new cell
func Show[T any](args ...T) {
	Fshow(os.Stdout, args)
}

// Showf prints each element of a slice to a stdout on a new cell
// using a given format string each time
func Showf[T any](format string, args ...T) {
	Fshowf(os.Stdout, format, args)
}

// Showln prints each element of a slice to a stdout on a new line
func Showln[T any](args ...T) {
	Fshowln(os.Stdout, args)
}

// Fshow prints each element of a slice to a given writer on a new cell
func Fshow[T any](w io.Writer, args []T) {
	for _, arg := range args {
		fmt.Fprint(w, arg)
	}
}

// Fshowf prints each element of a slice to a given writer on a new cell
// using a given format string each time
func Fshowf[T any](w io.Writer, format string, args []T) {
	for _, arg := range args {
		fmt.Fprintf(w, format, arg)
	}
}

// Fshowln prints each element of a slice to a given writer on a new line
func Fshowln[T any](w io.Writer, args []T) {
	for _, arg := range args {
		fmt.Fprintln(w, arg)
	}
}

func Join[T rules.Ordered](s []T, sep T) (out T) {
	for i, e := range s {
		out += e
		if i < len(s)-1 {
			out += sep
		}
	}
	return out
}

func JoinFunc[T any](add func(T, T) T, s []T, sep T) (out T) {
	for i, e := range s {
		out = add(out, e)
		if i < len(s)-1 {
			out = add(out, sep)
		}
	}
	return out
}

// Pairwise(ABCD) -> AB BC CD
func Pairwise[T any](args ...T) [][]T {
	tee := Tee(args, 2)
	a, b := tee[0], tee[1][1:]
	return Zip(a, b)
}

func SortedKey[T any, U rules.Ordered](k func(T) U, s []T) []T {
	key := Key[T, U](k)
	return SortedFunc(key.Lt, s)
}

func Remove[T any, int rules.Int](s *[]T, indices ...int) {
	for ctr := int(0); len(indices) > 0; ctr-- {
		i := indices[0] + ctr
		*s = append((*s)[:i], (*s)[i+1:]...)
	}
}

func Pop[T any, int rules.Int](s []T, i int) LR[T, []T] {
	return LR[T, []T]{
		Left:  s[i],
		Right: append(s[:i], s[i+1:]...),
	}
}

// CastAsync behaves much like cast except that all operations are concurrent
func CastAsync[I, O any](cast func(I) O, args ...I) []O {
	wg := new(sync.WaitGroup)
	wg.Add(len(args))
	out := make([]O, len(args))
	for i, arg := range args {
		go func(i int, arg I) {
			out[i] = cast(arg)
			wg.Done()
		}(i, arg)
	}
	wg.Wait()
	return out
}

// Rcast returns a slice whose values are the result of the
// application of the given function to all elements of the given slice
// it behaves like "map" in languages whose hashtables are called "associative array" or "dictionary"
func Rcast[I any, O any](fs []func([]I) O, s []I) []O {
	out := make([]O, len(fs))
	for i, f := range fs {
		out[i] = f(s)
	}
	return out
}

// Shuffle returns a permutation
func Shuffle[T any](args []T) []T {
	indices := rand.Perm(len(args))
	out := make([]T, len(args))
	// ctr := 0
	for j, i := range indices {
		out[j] = args[i]
		// ctr++
	}
	return out
}

func Deref[T any](arg []*T) []T {
	out := make([]T, len(arg))
	for i, e := range arg {
		out[i] = *e
	}
	return out
}

func Ref[T any](arg []T) []*T {
	out := make([]*T, len(arg))
	for i, e := range arg {
		out[i] = &e
	}
	return out
}
