package slices

import "github.com/kendfss/rules"

// Keys are functions that give a notion of size to members of unordered types.
// They are utilities for creating comparison operators on unordered types.
// In mathspeak, they're quite like measures.
type Key[I any, O rules.Ordered] func(I) O

// Key.Lt checks if left ... right
func (k Key[I, O]) Lt(left, right I) bool {
	return k(left) < k(right)
}

// Key.Le checks if left ... right
func (k Key[I, O]) Le(left, right I) bool {
	return k(left) <= k(right)
}

// Key.Gt checks if left ... right
func (k Key[I, O]) Gt(left, right I) bool {
	return k(left) > k(right)
}

// Key.Ge checks if left ... right
func (k Key[I, O]) Ge(left, right I) bool {
	return k(left) >= k(right)
}

// Key.Eq checks if left ... right
func (k Key[I, O]) Eq(left, right I) bool {
	return k(left) == k(right)
}

// Key.Ne checks if left ... right
func (k Key[I, O]) Ne(left, right I) bool {
	return k(left) != k(right)
}

// Key.Cmp(a, b) is expected to return an integer comparing the two
// parameters: 0 if a == b, a negative number if a < b and a positive number if
// a > b.
func (k Key[I, O]) Cmp(left, right I) int {
	if l, r := k(left), k(right); l == r {
		return 0
	} else if l < r {
		return -1
	}
	return 1
}
