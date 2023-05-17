```go
package chans // import "github.com/kendfss/iters/chans"


// VARIABLES

var DefaultCapacity = 0
var ErrUnsatisfied = but.New("Predicate was not satisfied")

// FUNCTIONS

func Cast[I, O any](f func(I) O, ch <-chan I) chan O
    // Cast calls a pure function on every value of a channel and returns a channel
    // containing all the results

func Chain[T any](args ...chan T) <-chan T
    // Chain collects several channels and returns one populated by their content

func Compact[T comparable](ch chan T) chan T
    // remove all duplicates from a channel

func CompactFunc[T comparable](eq func(T, T) bool, ch chan T) chan T
    // remove all duplicates from a channel of a non-comparable type

func Count[T any](c chan T) (out uint64)

func Do[T any](f func(T), ch <-chan T)
    // Send calls a function on every value of a slice

func Extend[T any](receiver chan T, args ...<-chan T)
    // Extend the first argument with the contents of the successors non blocking,
    // non order-preserving

func Extender[T any](target chan T) func(...<-chan T)
func Filter(ch chan bool) chan bool
func FilterPred[T any](pred func(T) bool, ch chan T) chan T
func Get[T any, int rules.Int](count int, ch chan T)
    // Get receives (discards) "count" items from "ch"

func Inf[T any, cap rules.OrderedNumber](init func() T, args ...cap) chan T
func Lazify[T any](arg []T) <-chan T
func Make[T any, cap rules.OrderedNumber](args ...cap) chan T
    // Make creates a buffered channel of given capacity or an unbuffered channel
    // if the capacity is negative

func MustUpto[T rules.Real](args ...T) chan T
    // MustUpto returns an iterator whose behaviour is equivalent to that of Range

func ParseCap[cap rules.OrderedNumber](args ...cap) uint64
    // ParseCap helps you to anticipate the behaviour of functions with a "args
    // ...cap" parameter

func PredPutter[T any](dst chan T, pred func(T) bool) func(T) error
    // PredPutter returns a method of the given channel which sends it the given
    // argument if, and only if, the argument satisfies the given predicate the
    // put-method returns ErrUnsatisfied if the predicate is not satisfied

func Process[T any](c chan T)
    // Process consumes a channel

func Put[T any, int rules.Int](count int, ch chan T)
    // Put writes "count" copies of zero-initialized "T" instances to "ch"

func PutVal[T any](ch chan T, val T)
    // PutVal writes a given value to a channel. Useful for spawning go routines
    // before returning

func Putter[T any](dst chan T) func(T) error
    // Putter returns a method of the given channel which sends it the given
    // argument

func RO[T any](c chan T) <-chan T
    // RO wraps a read-write channel with a read-only one

func RW[T any](c <-chan T) chan T
    // RW wraps a read-only channel with a read-write one

func StepStr[T rules.Char](arg string) chan T
func Upto[T rules.Real](args ...T) (chan T, error)
    // Upto returns an iterator whose content depends on the number of arguments as
    // follows
    //     		# of args 	|| 	behaviour
    //     	 		 1	 	|| 	stop
    //     		 	 2	 	|| 	start, stop
    //     		 	 3	 	|| 	start, stop, step
    //             else 	|| 	error

func Watch[T any](dst, src chan T)
    // Watch feeds dst with items received from src does not close either of them

```