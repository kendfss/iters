```go
package slices // import "github.com/kendfss/iters/slices"


VARIABLES

var (
	ErrInsuff = errors.New("Insufficient Elements")
	ErrIndex  = errors.New("slice index out of range")
)

FUNCTIONS

func All(args []bool) bool
    Check if all elements of a slice are true

func AllFunc[E any](pred func(E) bool, slice []E) bool
    use a custom predicate to check if all elements of a slice have a common
    property

func Anify[E any](slice []E) []any
    Convert a slice of any type into one of empty interfaces

func Any(args ...bool) bool
    Check if any elements of a slice are true

func AnyFunc[E any](pred func(E) bool, slice ...E) bool
    use a custom predicate to check if any elements of a slice have a common
    property

func BinarySearch[E rules.Ordered](target E, space []E) (int, bool)
    BinarySearch searches for target in a sorted slice and returns the position
    where target is found, or the position where target would appear in the sort
    order; it also returns a bool saying whether the target is really found in
    the slice. The slice must be sorted in increasing order.

func BinarySearchFunc[E any](cmp func(E, E) int, target E, space []E) (int, bool)
    BinarySearchFunc works like BinarySearch, but uses a custom comparison
    function. The slice must be sorted in increasing order, where "increasing"
    is defined by cmp. cmp(a, b) is expected to return an integer comparing
    the two parameters: 0 if a == b, a negative number if a < b and a positive
    number if a > b.

func BinarySearchKey[E any, O rules.Ordered](key func(E) O, target E, space []E) (int, bool)
    BinarySearchKey accepts a measuring key and calls BinarySearchFunc

func Cast[E, V any](f func(E) V, s []E) []V
    Cast returns a slice whose values are the result of the application of the
    given function to all elements of the given slice it behaves like "map" in
    languages whose hashtables are called "associative array" or "dictionary"

func CastAsync[I, O any](cast func(I) O, args ...I) []O
    CastAsync behaves much like cast except that all operations are concurrent

func Caster[I, O any](op func(I) O) func([]I) []O
    Caster integrates a castable operator for use with a matrix see Cast for
    more info

func Chain[E any](args ...[]E) (out []E)
    Concatenate slices

func Channel[T any, Int rules.Integer](slice []T, cap Int) <-chan T
    Channel returns

        a buffered channel iff cap < 1
        an unbuffered channel otherwise

    use an empty slice if you want to reproduce "make(chan T, 0)"

func Choose[T any](arg []T) T
    Choose selects an element of the gicen slice at random

func Clip[E any](s []E) []E
    Clip removes unused capacity from the slice, returning s[:len(s):len(s)].

func Clone[E any](s []E) []E
    Clone returns a copy of the slice. The elements are copied using assignment,
    so this is a shallow clone.

func Combinations[T any](pool []T, r int) (out [][]T)
    Return r length subsequences of elements from the input empty if r >
    len(slice) || r < 0

    The combination tuples are emitted in lexicographic ordering according
    to the order of the input iterable. So, if the input iterable is sorted,
    the combination tuples will be produced in sorted order.

    Elements are treated as unique based on their position, not on their value.
    So if the input elements are unique, there will be no repeat values
    in each combination. Combinations('ABCD', 2) --> AB AC AD BC BD CD
    Combinations(range(4), 3) --> 012 013 023 123

func Compact[E comparable](s []E) []E
    Compact replaces consecutive runs of equal elements with a single copy.
    This is like the uniq command found on Unix. Compact modifies the contents
    of the slice s; it does not create a new slice.

func CompactFunc[E any](eq func(E, E) bool, s []E) []E
    CompactFunc is like Compact but uses a comparison function.

func Compacted[E comparable](s []E) []E
    Compacted clones the slice and runs Compact on said clone

func Compare[E rules.Ordered](s1, s2 []E) int
    Compare compares the elements of s1 and s2. The elements are compared
    sequentially, starting at index 0, until one element is not equal to the
    other. The result of comparing the first non-matching elements is returned.
    If both slices are equal until one of them ends, the shorter slice is
    considered less than the longer one. The result is 0 if s1 == s2, -1 if
    s1 < s2, and +1 if s1 > s2. Comparisons involving floating point NaNs are
    ignored.

func CompareFunc[E1, E2 any](cmp func(E1, E2) int, s1 []E1, s2 []E2) int
    CompareFunc is like Compare but uses a comparison function on each pair
    of elements. The elements are compared in increasing index order, and the
    comparisons stop after the first time cmp returns non-zero. The result is
    the first non-zero result of cmp; if cmp always returns 0 the result is 0 if
    len(s1) == len(s2), -1 if len(s1) < len(s2), and +1 if len(s1) > len(s2).

func Consume[T any](channel chan T) (out []T)
    Consume passes each element of the given channel to the given slice

func Contains[E comparable](s []E, v E) bool
    Contains reports whether v is present in s.

func ContainsAny[E comparable](s []E, args ...E) bool
    ContainsAny reports whether any element of args is present in s.

func ContainsAnyFunc[E any](eq func(E, E) bool, s []E, args ...E) bool
    ContainsAnyFunc reports whether of args is present in s.

func ContainsFunc[E any](eq func(E, E) bool, s []E, v E) bool
    ContainsFunc reports whether v is present in s, using eq as an equivalence
    operator.

func ContainsPred[E any](pred func(E) bool, s []E) bool
    ContainsPred reports whether anything in s satisfies the given predicate.

func Copies[T any, U rules.I](length U, val T) []T
    Deprecated, use Repeat

func Copy[T any](dst, src []T) int
    Copy is a wrapper on the builtin copy

func Count[T comparable](item T, rack []T) (out uint)
    Count returns the number of occurences of item in rack

func CountFunc[E any](eq func(E, E) bool, item E, rack []E) (out uint)
    CountFunc returns the number of occurences, with respect to eq == true,
    of item in rack

func Delete[E any](s []E, i, j int) []E
    Delete removes the elements s[i:j] from s, returning the modified slice.
    Delete panics if s[i:j] is not a valid slice of s. Delete modifies
    the contents of the slice s; it does not create a new slice. Delete is
    O(len(s)-(j-i)), so if many items must be deleted, it is better to make a
    single call deleting them all together than to delete one at a time.

func Deref[T any](arg []*T) []T
func Dot[N rules.Num](left, right []N) []N
    Dot returns a dot product analog of left with right. Dot({2, 3}, {1,
    2}) === {2, 6} Dot({2}, {1, 2}) === {2, 0} Dot({1, 2}, {2}) === {2, 0}

func DotFunc[T any](mul func(T, T) T, left, right []T) []T
    DotFunc returns a dot product analog of left with right, using mul as a
    binary operator over the chosen type.

func Enumerate[I rules.Integer, T any](slice []T) []func() (I, T)
    Enumerate returns a slice of closures whose return values are tuples of
    elements of the given slice prefixed by their indices

func Equal[E comparable](s1, s2 []E) bool
    Equal reports whether two slices are equal: the same length and all elements
    equal. If the lengths are different, Equal returns false. Otherwise, the
    elements are compared in increasing index order, and the comparison stops at
    the first unequal pair. Floating point NaNs are not considered equal.

func EqualFunc[E1, E2 any](eq func(E1, E2) bool, s1 []E1, s2 []E2) bool
    EqualFunc reports whether two slices are equal using a comparison function
    on each pair of elements. If the lengths are different, EqualFunc returns
    false. Otherwise, the elements are compared in increasing index order,
    and the comparison stops at the first index for which eq returns false.

func Extend[T any, C rules.Integer](slice []T, seed T, count C) []T
func Extremal[E any](operator func(E, E) bool, args ...E) (out int)
    Extremal finds the index of a maximum, or minimum, value of a
    slice by passing a key corresponding to greater than or less than
    Extremal[MyType](gt, mySlice...) -> maximal value Extremal[MyType](lt,
    mySlice...) -> minimal value

func Feed[T any](channel chan T, slice []T)
    Feed passes each element of the given slice to the given channel

func Filter(args []bool) (out []bool)
    Filter returns a slice featuring all truthy elements

func FilterFunc[E any](f func(E) bool, args []E) (out []E)
    FilterFunc returns a slice featuring all elements of the incident that
    satisfy the given predicate

func Flatter[M ~[][]E, E any](arg M) (out []E)
    Deprecated, use Chain

func Fshow[T any](w io.Writer, args []T)
    Fshow prints each element of a slice to a given writer on a new cell

func Fshowf[T any](w io.Writer, format string, args []T)
    Fshowf prints each element of a slice to a given writer on a new cell using
    a given format string each time

func Fshowln[T any](w io.Writer, args []T)
    Fshowln prints each element of a slice to a given writer on a new line

func Get[E any, I rules.Integer](index I, slice []E) E
    Get returns the i'th element from a slice, even if i is negative uses the
    same indexing convention as python lists/tuples

func Getter[T any, I rules.Int](index I) func([]T) T
    Getter returns a castable operator that fetches the element of that index
    from a slice. see Get and Cast for more info

func Getxy[E any](slice []E, stride, x, y int) E
    Get an element of a slice situated at a point (x,y) when the slice is
    interpreted as

func Grow[E any](s []E, n int) []E
    Grow increases the slice's capacity, if necessary, to guarantee space for
    another n elements. After Grow(n), at least n elements can be appended
    to the slice without another allocation. Grow may modify elements of the
    slice between the length and the capacity. If n is negative or too large to
    allocate the memory, Grow panics.

func Index[E comparable](val E, s []E) int
    Index returns the index of the first occurrence of v in s, or -1 if not
    present.

func IndexFunc[E any](eq func(E, E) bool, val E, s []E) int
    IndexFunc returns the first index i satisfying f(s[i]), or -1 if none do.

func IndexPred[E any](eq func(E) bool, s []E) int
    IndexPred returns the first index i satisfying f(s[i]), or -1 if none do.

func Indices[T comparable](item T, rack []T) (out []int)
    Indices returns the positions at which item can be found in rack

func IndicesFunc[T comparable](eq func(T, T) bool, item T, rack []T) (out []int)
    IndicesFunc returns the positions at which item can be found, by eq == true,
    in rack

func Insert[E any](s []E, i int, args ...E) []E
    Insert inserts the values v... into s at index i, returning the modified
    slice. In the returned slice r, r[i] == v[0]. Insert panics if i is out of
    range. This function is O(len(s) + len(v)).

func IsSorted[E rules.Ordered](x []E) bool
    IsSorted reports whether x is sorted in ascending order.

func IsSortedFunc[E any](less func(a, b E) bool, x []E) bool
    IsSortedFunc reports whether x is sorted in ascending order, with less as
    the comparison function.

func IsSortedKey[E any, O rules.Ordered](key func(E) O, data []E) bool
    IsSortedKey accepts a measuring key and calls IsSortedFunc

func Join[T rules.Ordered](s []T, sep T) (out T)
func JoinFunc[T any](add func(T, T) T, s []T, sep T) (out T)
func Len[I rules.Integer, E any](slice []E) I
    Len returns the length of a slice as the desired type of integer

func Longest[E any](args ...[]E) (out int)
    Returns the index of the longest slice received at call-time -1 if no
    arguments are passed

func Make[T any, I rules.Integer](length I) []T
    Make initializes a slice

func Max[E rules.Ordered](args ...E) (out int)
    Max returns the index of the Maximal value of a slice

func Min[E rules.Ordered](args ...E) (out int)
    Min returns the index of the Minimal value of a slice

func Movers[E comparable](s []E) (out []int)
    Movers returns the indices of slice elements that are not equal to their
    successors

func MoversFunc[E any](eq func(E, E) bool, s []E) (out []int)
    MoversFunc returns the indices of slice elements that are not equal to their
    successors

func Mul[N rules.Num](left, right []N) []N
    Mul returns a dot product analog of left with right. Each argument of
    length 1 is treaded as a scalar Mul({2, 3}, {1, 2}) === {2, 6} Mul({2}, {1,
    2}) === {2, 4} Mul({1, 2}, {2}) === {2, 4}, Mul(right, left) if len(left) >
    len(right) If you want ordinary Muls, use Cast(a+b, Zip(left, right))

func Nest[T any](s ...T) [][]T
    Nest places a slice of T into a new Matrix of T

func Ones[T rules.Integer](count T) []T
    Deprecated, use Repeat

func Pairwise[T any](args ...T) [][]T
    Pairwise(ABCD) -> AB BC CD

func Partition[K comparable, V any](pred func(V) K, slice []V) map[K][]V
    Partition uses a function to categorize elements of a slice

func Permutations[T any](r int, pool []T) (out [][]T)
func Pointers[T any](s []T) []*T
    Pointers returns an array of pointers to the values of given slice These
    pointers should not agree with other reference to the data

func Prefill[T any](s []T, by uint) []T
    Prefill prepends some number of elements to a slice

func PrefillSeed[T any](slice []T, seed T, by uint) []T
    PrefillSeed prepends some number of elements to a slice

func PrefillSeeder[T any](seed T, by uint) func([]T) []T
    PrefillSeeder returns a castable operator for prefilling slices see
    PrefillSeed and Cast for more info

func Prefiller[T any](by uint) func([]T) []T
    Prefiller returns a castable operator for prefilling slices see Prefill and
    Cast for more info

func Product[T any](repeat int, args ...[]T) [][]T
func Rcast[I any, O any](fs []func([]I) O, s []I) []O
    Rcast returns a slice whose values are the result of the application of the
    given function to all elements of the given slice it behaves like "map" in
    languages whose hashtables are called "associative array" or "dictionary"

func Reduce[E any](f func(E, E) E, s []E) (out E)
    Reduce returns the outcome of successive applications of a function, f,
    as a binary operator over the slice, s.

func ReduceAs[I, O rules.Real](op func(O, O) O, args ...I) O
    ReduceAs applies Reduce after converting a slice of real numbers an
    overflow-safe way for operating on small numbers

func ReducePairLeft[T any](op func(l, r T) T, rack []Pair[T]) T
func ReducePairRight[T any](op func(l, r T) T, rack []Pair[T]) T
func Reducer[T any](op func(T, T) T) func([]T) T
    Reducer returns a castable operator that reduces a slice. see Reduce and
    Cast for more info

func Ref[T any](arg []T) []*T
func Remove[T any, int rules.Int](s *[]T, indices ...int)
func Repeat[T any, C rules.Integer](seed T, count C) []T
    Repeat returns a slice, with length count, of integers initialized to seed
    if you want to repeat an empty slice you should use Tee instead

func Repeater[T any, I rules.Int](count I) func(T) []T
    Repeater returns a castable operator for creating slices of given length
    populated by given value see Repeat and Cast for more info

func Resize[T any](s []T, shape ...int) []T
func Reverse[E any](slice []E)
    Reverse a slice in place func Reverse[[]E ~[]E, E any](slice []E) {

func Reversed[E any](slice []E) []E
    Produce a reversed copy of a slice

func Rotated[T any, I rules.I](slice []T, steps I) []T
    Rotated returns the index-shifted of a given slice as though the operation
    were taking place on a torus (no elements lost or added)

func Select[E any](slice []E, indices []int) []E
    select returns the elements of a slice located at the chosen indices note:
    all indices are wrapped by a modulus equal to the length of the slice use
    StrictSelect to mitigate this behaviour

func SelectStrict[E any](slice []E, indices []int) []E
    select returns the elements of a slice located at the chosen indices note:
    indices greater than slice length with cause panic use Select to mitigate
    this behaviour

func Send[T any](f func(T), args []T)
    Send is like Cast but for impure functions

func Shortest[E any](args ...[]E) (out int)
    Returns the index of the shortest slice received at call-time -1 if no
    arguments are passed

func Show[T any](args ...T)
    Show prints each element of a slice to a stdout on a new cell

func Showf[T any](format string, args ...T)
    Showf prints each element of a slice to a stdout on a new cell using a given
    format string each time

func Showln[T any](args ...T)
    Showln prints each element of a slice to a stdout on a new line

func Shuffle[T any](args []T) []T
    Shuffle returns a permutation

func Skip[T any, I rules.Int](index I, arg []T) (out []T)
    Skip returns a version of the slice without the element at the given index
    returns arg if index is greater than len(arg)-1

func Snap[I rules.I, E any](width I, arg []E) (out [][]E)
    Snap breaks a slice into sections of given width Snap(2, []int{1, 2, 3,
    4}) == [][]int{{1, 2}, {3, 4}} Snap(3, []int{1, 2, 3, 4}) == [][]int{{1, 2,
    3}, {4}} func Snap[[]E ~[]E, E any](arg []E, width int) (out [][]E) {

func Snapper[T any, I rules.Int](stride I) func([]T) [][]T
    Snapper returns a castable operator for snapping slices see Snap and Cast
    for more info

func Sort[E rules.Ordered](x []E)
    Sort sorts a slice of any ordered type in ascending order. Sort may fail
    to sort correctly when sorting slices of floating-point numbers containing
    Not-a-number (NaN) values. Use slices.SortFunc(x, func(a, b float64) bool
    {return a < b || (math.IsNaN(a) && !math.IsNaN(b))}) instead if the input
    may contain NaNs.

func SortFunc[E any](less func(a, b E) bool, x []E)
    SortFunc sorts the slice x in ascending order as determined by the less
    function. This sort is not guaranteed to be stable.

    SortFunc requires that less is a strict weak ordering. See
    https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.

func SortKey[E any, O rules.Ordered](key func(E) O, arg []E)
    SortKey wraps a Key with a less than (<) function before deferring to
    SortFunc see slices.Key for more info

func SortStableFunc[E any](less func(a, b E) bool, x []E)
    SortStable sorts the slice x while keeping the original order of equal
    elements, using less to compare elements.

func SortStableKey[E any, O rules.Ordered](key func(E) O, data []E)
    SortStableKey accepts a measuring key and calls SortStableFunc

func Sorted[E rules.Ordered](x []E) []E
    Sorted sorts a slice of any ordered, type after cloning it, in ascending
    order. Sort may fail to sort correctly when sorting slices of floating-point
    numbers containing Not-a-number (NaN) values. Use slices.SortFunc(x, func(a,
    b float64) bool {return a < b || (math.IsNaN(a) && !math.IsNaN(b))}) instead
    if the input may contain NaNs.

func SortedFunc[E any](less func(a, b E) bool, x []E) []E
    SortedFunc sorts a clone of the slice x in ascending order as determined by
    the less function. This sort is not guaranteed to be stable.

    SortFunc requires that less is a strict weak ordering. See
    https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.

func SortedKey[T any, U rules.Ordered](k func(T) U, s []T) []T
func Split[E comparable](slice []E, breaker E) [][]E
    Split "cuts" the slice at all occurrences of breaker

func SplitAfter[E comparable](slice []E, breaker E) [][]E
    SplitAfter "cuts" the slice at all matching elements without discarding them

func SplitAfterFunc[E any](function func(E, E) bool, breaker E, slice []E) [][]E
    SplitAfterFunc "cuts" the slice at all matching elements without discarding
    them

func SplitAfterPred[E any](function func(E) bool, slice []E) [][]E
    SplitAfterPred "cuts" the slice at all satisfying elements without
    discarding them

func SplitFunc[E any](eq func(E, E) bool, slice []E, breaker E) [][]E
    SplitFunc "cuts" the slice at all occurrences of breaker

func SplitPred[E any](pred func(E) bool, slice []E) [][]E
    SplitPred "cuts" the slice at all elements satisfying some predicate

func Standers[E comparable](s []E) (out []E)
    Standers returns the indices of the slice elements that are equal to their
    successors

func StandersFunc[E any](eq func(E, E) bool, s []E) (out []E)
    StandersFunc returns the indices of the slice elements that are equal to
    their successors

func Swap[E any](slice []E, i, j int) []E
    Swap the elements at a pair of indices (in place)

func Swapped[E any](slice []E, i, j int) []E
    Swap the elements at a pair of indices (copied)

func Tee[T any, I rules.Integer](seed []T, count I) [][]T
    Tee returns a slice of independent slices

func Trot[T, O any](operator func(T, T) O, data []T) []O
    Trot returns the outcome of step-wise applications of a function, f,
    as a binary operator over the slice, s. Trot{addition, {1, 2, 3}} == {1, 1,
    1}

func Union[E any](first []E, rest ...[]E) []E
    Deprecated, use Chain

func Upto[O, I rules.Real](start, stop, step I) []O
    Consecutive ints, including start, smaller than stop, and separated by step
    Upto[byte](0, 256, 1)

func Upton[O, I rules.Real](stop I) []O
    Indices of a slice with given length Range[byte](256)

func Uptonm[O, I rules.Real](start, stop I) []O
    Consecutive ints, including start, smaller than stop, and separated by one
    Uptonm[byte](0, 256)

func Values[T any](s []*T) []T
    Values returns a slice of values to members of a slice of pointers

func VariadicFilter[E any](adicity int, walk bool, f func(...E) bool, slice []E) (out [][]E)
func Walks[T any, I rules.Integer](length I, slice []T) (out [][]T)
    Break an iterable into len(iterable)-length steps of the given length,
    with each step's starting point one after its predecessor example

        >>> for i in walks(itertools.count(),2):print(''.join(i))
        (0, 1)
        (1, 2)
        (2, 3)
        # etc.

    Inspired by the hyperoperation 16**2[5]2

func Windows[T any](src []T, size int) (out [][]T)
func Zip[K any](args ...[]K) (out [][]K)
    Convolve type-equivalent slices

func Zip3[L, R any](left []L, right []R) (out []func() (L, R))
    Convolve pairs of type-distinct slices with a closure


TYPES

type Key[I any, O rules.Ordered] func(I) O
    Keys are functions that give a notion of size to members of unordered types.
    They are utilities for creating comparison operators on unordered types.
    In mathspeak, they're quite like measures.

func (k Key[I, O]) Cmp(left, right I) int
    Key.Cmp(a, b) is expected to return an integer comparing the two parameters:
    0 if a == b, a negative number if a < b and a positive number if a > b.

func (k Key[I, O]) Eq(left, right I) bool
    Key.Eq checks if left ... right

func (k Key[I, O]) Ge(left, right I) bool
    Key.Ge checks if left ... right

func (k Key[I, O]) Gt(left, right I) bool
    Key.Gt checks if left ... right

func (k Key[I, O]) Le(left, right I) bool
    Key.Le checks if left ... right

func (k Key[I, O]) Lt(left, right I) bool
    Key.Lt checks if left ... right

func (k Key[I, O]) Ne(left, right I) bool
    Key.Ne checks if left ... right

type LR[L, R any] struct {
	// LR holds two values, Left and Right, of any types.
	Left  L
	Right R
}

func Cartesian[L, R any](left []L, right []R) []LR[L, R]

func Pop[T any, int rules.Int](s []T, i int) LR[T, []T]

func Zip2[L, R any](left []L, right []R) (out []LR[L, R])
    Convolve pairs of type-distinct slices with a Pair

func (lr LR[L, R]) From(l L, r R) LR[L, R]

func (lr LR[L, R]) L() L

func (lr LR[L, R]) R() R

func (lr LR[L, R]) Split() (l L, r R)

type Pair[T any] struct {
	// Pair holds two values, Left and Right, of any type.
	Left, Right T
}

func PairAll[T any](arg [][]T, shift int) (out []Pair[T])
    PairAll returns a sequence of pairs from a matrix

func (p Pair[T]) From(l, r T) Pair[T]

func (p Pair[T]) L() T

func (p Pair[T]) R() T

func (p Pair[T]) Split() (l, r T)

```
