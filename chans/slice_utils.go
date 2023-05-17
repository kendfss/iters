package chans

// sliceContains reports whether v is present in s.
func sliceContains[E comparable](s []E, v E) bool {
	return sliceIndex(s, v) >= 0
}

// sliceContains reports whether v is present in s.
func sliceContainsFunc[E comparable](eq func(E, E) bool, s []E, v E) bool {
	return sliceIndexFunc(eq, v, s) >= 0
}

// sliceIndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func sliceIndexFunc[E any](eq func(E, E) bool, val E, s []E) int {
	for i, v := range s {
		if eq(v, val) {
			return i
		}
	}
	return -1
}

// sliceIndex returns the index of the first occurrence of v in s,
// or -1 if not present.
func sliceIndex[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

func sliceToChan[T any](slice []T) chan T {
	out := make(chan T)

	go func() {
		for _, e := range slice {
			out <- e
		}
	}()

	return out
}
