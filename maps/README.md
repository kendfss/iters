```go

package maps // import "github.com/kendfss/iters/maps"


// FUNCTIONS

func Clear[K comparable, V any](m map[K]V)
    // Clear removes all entries from m, leaving it empty.

func Clone[K comparable, V any](m map[K]V) map[K]V
    // Clone returns a copy of m. This is a shallow clone: the new keys and values
    // are set using ordinary assignment.

func Copy[K comparable, V any](dst, src map[K]V)
    // Copy copies all key/value pairs in src adding them to dst. When a key in src
    // is already present in dst, the value in dst will be overwritten by the value
    // associated with the key in src.

func Countf[K comparable, T any](fn func(T) K, args ...T) map[K]int
func DeleteFunc[K comparable, V any](m map[K]V, del func(K, V) bool)
    // DeleteFunc deletes any key/value pairs from m for which del returns true.

func Equal[K, V comparable](m1, m2 map[K]V) bool
    // Equal reports whether two maps contain the same key/value pairs. Values are
    // compared using ==.

func EqualFunc[K comparable, V1, V2 any](m1 map[K]V1, m2 map[K]V2, eq func(V1, V2) bool) bool
    // EqualFunc is like Equal, but compares values using eq. Keys are still
    // compared with ==.

func Extend[K comparable, V any](m map[K]V, key K, val V) map[K]V
    // Extend

func Filter[K comparable, V any](pred func(V) bool, m map[K]V) map[K]V
    // FilterKV creates a new map consisting of values which satisfy a predicate

func FilterKV[K comparable, V any](pred func(K, V) bool, m map[K]V) map[K]V
    // FilterKV creates a new map consisting of key-value pairs which satisfy a
    // predicate

func FromKeys[K comparable, V any](fn func(K) V, args ...K) map[K]V
    // FromKeys creates map values by casting keys

func FromKeys2[K comparable, V any](fn func(K) V, args ...K) map[K][]V
    // FromKeys2 creates map values by casting keys the values are kept in an array
    // to avoid collisions

func FromVals[K comparable, V any](fn func(V) K, vals ...V) map[K]V
    // FromVals creates map keys by casting values

func FromVals2[K comparable, V any](fn func(V) K, vals ...V) map[K][]V
    // FromVals2 creates map keys by casting values the values are kept in an array
    // to avoid collisions

func Getter[V any, K comparable](table map[K]V) func(K) V
    // Getter returns a castable operator that fetches the element of that index
    // from a slice. see Get and Cast for more info

func Keys[K comparable, V any](m map[K]V) []K
    // Keys returns the keys of the map m. The keys will be in an indeterminate
    // order.

func Keys2[K comparable, V any](m map[K]V) []K
func Mass[K comparable, V any](m map[K][]V) (out int)
    // Mass computes the number of items in values of a slice-valued map

func Select[K comparable, V any](m map[K]V, keys []K) []V
func Values[K comparable, V any](m map[K]V) []V
    // Values returns the values of the map m. The values will be in an
    // indeterminate order.

func Values2[K comparable, V any](m map[K]V) []V
func Vector[K comparable](args ...K) map[K]int
```