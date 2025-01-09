iters
===

Go library for genric interator operators

There are tools for [maps](./maps/README.md), [slices](./slices/README.md), and [channels](./chans/README.md).

As well as the following top level offerings:

```go
package iters // import "github.com/kendfss/iters"


CONSTANTS

const ErrBadCopy but.Note = "bad copy: should've copied %d but copied %d"

FUNCTIONS

func And[T any](args ...func(T, T) bool) func(T, T) bool
    And returns a predicate that seeks the satisfaction of all arguments

func Eq[T comparable](a, b T) bool
    Eq checks the given arguments for equality

func Lt[T rules.Ordered](a, b T) bool
    Lt checks that the first argument is less than the second argument

func Neq[T comparable](a, b T) bool
    Neq checks the given arguments for inequality

func Not[T any](pred func(T, T) bool) func(T, T) bool
    Not returns the negation of the given predicate

func Or[T any](args ...func(T, T) bool) func(T, T) bool
    Or returns a predicate that seeks the satisfaction of any argument

func Slice[T any](args ...T) []T
    Slice simply returns its arguments in a slice


TYPES

type Queue[T any] struct {
	Buf []T
	*sync.RWMutex
}
    Queue is a threadsafe FIFO data structure

func NewQueue[T any](elems ...T) *Queue[T]
    NewQueue initializes a Queue variadically. If you already have a slice,
    just use the constructor

func (q *Queue[T]) Cast(fn func(T) T) *Queue[T]
    Cast performs an in-place map of q

func (q *Queue[T]) Clone() *Queue[T]
    Clone returns a new queue with the same elements in a new buffer

func (q *Queue[T]) Do(fn func(T)) *Queue[T]
    Do executes some function on each element of q

func (q *Queue[T]) Eq(p *Queue[T], equals func(T, T) bool) bool
    Eq checks two Queues for equality

func (q *Queue[T]) GoCast(fn func(T) T) *Queue[T]
    GoCast performs a concurrent Cast of each element in the Queue

func (q *Queue[T]) GoDo(fn func(T)) *Queue[T]
    GoDo executes some function on each element of q concurrently

func (q *Queue[T]) GoMap(fn func(T) T) *Queue[T]
    GoMap performs a concurrent mapping on the elements of the Queue

func (q *Queue[T]) Len() int
    Len returns the length of the underlying buffer

func (q *Queue[T]) Map(fn func(T) T) *Queue[T]
    Map returns a new Queue whose elements are return values of fn on q

func (q *Queue[T]) Pop() *T
    Pop returns the first element from q returns nil for the empty Queue

func (q *Queue[T]) Push(e T) *Queue[T]
    Push appends an element to q

func (q Queue[T]) String() string
    String returns a text representation of the owner. Specifically it returns
    the type name and pointer of the first element

```
