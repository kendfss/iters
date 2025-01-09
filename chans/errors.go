package chans

import "github.com/kendfss/but"

const (
	ErrUnsatisfied  but.Note = "Predicate was not satisfied"
	ErrWriteOnEmpty but.Note = "cannot write to empty channel"
)
