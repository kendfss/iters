package slices

import "errors"

var (
	ErrInsuff = errors.New("Insufficient Elements")
	ErrIndex  = errors.New("slice index out of range")
)
