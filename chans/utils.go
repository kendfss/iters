package chans

// selfHash simply returns its argument, used for hash accepting functions
func selfHash[T comparable](t T) T {
	return t
}
