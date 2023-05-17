package maps

func Keys2[K comparable, V any](m map[K]V) []K {
	out := make([]K, len(m))
	ctr := 0
	for k := range m {
		out[ctr] = k
		ctr++
	}
	return out
}

func Values2[K comparable, V any](m map[K]V) []V {
	out := make([]V, len(m))
	ctr := 0
	for k := range m {
		out[ctr] = m[k]
		ctr++
	}
	return out
}

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
func Keys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Values returns the values of the map m.
// The values will be in an indeterminate order.
func Values[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

// Equal reports whether two maps contain the same key/value pairs.
// Values are compared using ==.
func Equal[K, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// EqualFunc is like Equal, but compares values using eq.
// Keys are still compared with ==.
func EqualFunc[K comparable, V1, V2 any](m1 map[K]V1, m2 map[K]V2, eq func(V1, V2) bool) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || !eq(v1, v2) {
			return false
		}
	}
	return true
}

// Clear removes all entries from m, leaving it empty.
func Clear[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

// func Intersection[K comparable, V any](a, b map[K]V) map[K]V {
// }

func Vector[K comparable](args ...K) map[K]int {
	out := make(map[K]int, len(args))
	for _, arg := range args {
		out[arg]++
	}
	return out
}

func Countf[K comparable, T any](fn func(T) K, args ...T) map[K]int {
	out := make(map[K]int, len(args))
	for _, arg := range args {
		out[fn(arg)]++
	}
	return out
}

// Clone returns a copy of m.  This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func Clone[K comparable, V any](m map[K]V) map[K]V {
	r := make(map[K]V, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r
}

func Select[K comparable, V any](m map[K]V, keys []K) []V {
	out := make([]V, len(keys))
	for i, key := range keys {
		out[i] = m[key]
	}
	return out
}

// Copy copies all key/value pairs in src adding them to dst.
// When a key in src is already present in dst,
// the value in dst will be overwritten by the value associated
// with the key in src.
func Copy[K comparable, V any](dst, src map[K]V) {
	for k, v := range src {
		dst[k] = v
	}
}

// DeleteFunc deletes any key/value pairs from m for which del returns true.
func DeleteFunc[K comparable, V any](m map[K]V, del func(K, V) bool) {
	for k, v := range m {
		if del(k, v) {
			delete(m, k)
		}
	}
}

// FilterKV creates a new map consisting of key-value pairs which satisfy a predicate
func FilterKV[K comparable, V any](pred func(K, V) bool, m map[K]V) map[K]V {
	out := make(map[K]V)
	for k, v := range m {
		if pred(k, v) {
			out[k] = v
		}
	}
	return out
}

// FilterKV creates a new map consisting of values which satisfy a predicate
func Filter[K comparable, V any](pred func(V) bool, m map[K]V) map[K]V {
	out := make(map[K]V)
	for k, v := range m {
		if pred(v) {
			out[k] = v
		}
	}
	return out
}

// Extend
func Extend[K comparable, V any](m map[K]V, key K, val V) map[K]V {
	_, ok := m[key]
	if !ok {
		m[key] = val
	}
	return m
}

// Getter returns a castable operator that fetches the element of that index from a slice.
// see Get and Cast for more info
func Getter[V any, K comparable](table map[K]V) func(K) V {
	return func(k K) V {
		return table[k]
	}
}

// FromKeys creates map values by casting keys
func FromKeys[K comparable, V any](fn func(K) V, args ...K) map[K]V {
	out := map[K]V{}
	for _, arg := range args {
		out[arg] = fn(arg)
	}
	return out
}

// FromKeys2 creates map values by casting keys
// the values are kept in an array to avoid collisions
func FromKeys2[K comparable, V any](fn func(K) V, args ...K) map[K][]V {
	out := map[K][]V{}
	for _, arg := range args {
		v, _ := out[arg]
		out[arg] = append(v, fn(arg))
	}
	return out
}

// FromVals creates map keys by casting values
func FromVals[K comparable, V any](fn func(V) K, vals ...V) map[K]V {
	out := map[K]V{}
	for _, val := range vals {
		out[fn(val)] = val
	}
	return out
}

// FromVals2 creates map keys by casting values
// the values are kept in an array to avoid collisions
func FromVals2[K comparable, V any](fn func(V) K, vals ...V) map[K][]V {
	out := map[K][]V{}
	for _, val := range vals {
		key := fn(val)
		v, _ := out[key]
		out[key] = append(v, val)
	}
	return out
}

// Mass computes the number of items in values of a slice-valued map
func Mass[K comparable, V any](m map[K][]V) (out int) {
	for _, v := range m {
		out += len(v)
	}
	return
}
