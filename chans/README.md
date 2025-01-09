iters/chans
===

chans exposes tools for manipulating and debugging channels

```go
package chans // import "github.com/kendfss/iters/chans"


CONSTANTS

const (
	ErrUnsatisfied  but.Note = "Predicate was not satisfied"
	ErrWriteOnEmpty but.Note = "cannot write to empty channel"
)

FUNCTIONS

func Buffered[T any](cap int) chan T
    Buffered initializes a buffered channel

func Bytes(s string) <-chan byte
    Bytes yields each byte in s

func Cast[I, O any, chanI Readable[I]](f func(I) O, src chanI) chan O
    Cast calls a pure function on every value of a channel and returns a channel
    containing all the results

func Castro[T any, chanT Readable[T]](buf []chanT) []<-chan T
    Castro casts multiple read-write channels to readonly

func Chain[T any, chanT Readable[T]](args ...chanT) <-chan T
    Chain collects several channels and returns one populated by their content.
    Chain expects the user to control closure of argument channels. Chain
    handles closure of returned channel. Returned channel has capacity equal to
    that of the argument with highest capacity

func ChainCap[T any, chanT Readable[T]](capacity int, args ...chanT) <-chan T
    ChainCap creates a new channel, with desired capacity, that serves as a
    frontent for the given arguments

func Close[T any](ch chan T)
    Close closes a channel

func Compact[T comparable, chanT Readable[T]](src chanT) <-chan T
    Compact removes all duplicates from a channel in constant time

func CompactFunc[T comparable, chanT Readable[T]](eq func(T, T) bool, src chanT) <-chan T
    CompactFunc removes all duplicates from a channel in linear time

func CompactHash[T any, H comparable, chanT Readable[T]](hash func(T) H, src chanT) <-chan T
    CompactHash removes all duplicates from a channel in constant time

func Context[T any](ctx context.Context, src chan T) chan T
    Context returns a channel that closes as soon as either the context is done
    or the source channel is closed. It does not, otherwise, operate on the
    context object an cannot cancel it

func Count[T any, chanT Readable[T]](c chanT) (out uint64)
    Count returns the number of elements passed through the channel before it is
    closed (externally)

func Discard[T any, int rules.Int, chanT Readable[T]](count int, src chanT)
    Discard ignores "count" items from "src"

func Do[T any, chanT Readable[T]](f func(T), src chanT)
    Do calls a function on every value of a slice

func Drain[T any](src chan T)
    Drain consumes a channel to depletion

func Extend[T any, readable Readable[T]](receiver chan T, args ...readable)
    Extend the first argument with the contents of the successors. non blocking,
    non order-preserving

func Extender[T any, readable Readable[T]](target chan T) func(...readable)
    Extender functions create extensions of the target channel. See Extend for
    more information

func Filter(src chan bool) chan bool
    Filter takes a boolean channel and skips all false values it receives

func FilterPred[T any, channel Readable[T]](pred func(T) bool, src channel) chan T
    FilterPred skips any channel receipts that fail to satisfy the given
    predicate

func Get[T any, chanT Readable[T]](count int, src chanT) []T
    Deprecated: use Pop instead

func Inf[T any](init func() T, capacity int) <-chan T
    Inf produces an infinite channel of given capacity using init to create each
    element/ Inf handles channel closure

func Lazify[T any](arg []T) <-chan T
    Lazify converts a slice into a readonly channel

func Map[I, O any, chanI Readable[I]](f func(I) O, src chanI) <-chan O
    Map calls Cast

func MustUpto[T rules.Real](args ...T) chan T
    MustUpto returns an iterator whose behaviour is equivalent to that of Range

func New[T any]() chan T
    New initializes an unbuffered channel

func Next[T any](src <-chan T) T
    Next extracts one receipt from the channel

func NextSafe[T any, chanT Readable[T]](src chanT) (T, bool)
    NextSafe extracts one receipt from the channel and informs whether it's
    closed

func Pop[T any, chanT Readable[T]](count int, src chanT) []T
    Pop receives "count" items from "src"

func Popper[T any, channel Readable[T]](pred func(T) bool, src channel) func() T
    Popper functions return the next channel receipt that satisfies the given
    predicate. If you don't know when the source channel is closed, use NextSafe
    instead; unless the zero value is known to not satisfy the predicate.

func PopperDefault[T any, chanT Readable[T]](src chanT, defaultVal T) func() T
    PopperDefault returns a popper that yields a default value once the channel
    is closed

func PopperSafe[T any](fn func(T) bool, src chan T) func() (T, bool)
    PopperSafe functions return the next channel receipt that satisfies the
    given predicate and tells the user when the channel is closed

func PredPutter[T any](dst chan T, pred func(T) bool) func(T) error
    PredPutter returns a method of the given channel which sends it the given
    argument if, and only if, the argument satisfies the given predicate.
    It returns a non-nil error when the destination is closed

func Push[T any, int rules.Int, chanT Writable[T]](count int, dst chanT)
    Push writes "count" copies of zero-initialized "T" instances to "dst"

func PushVal[T any, int rules.Int, chanT Writable[T]](count int, dst chanT, val T)
    PushVal writes a given value to a channel. Useful for spawning go routines
    before returning

func Pusher[T any, chanT Writable[T]](dst chanT) func(T) error
    Pusher returns a method of the given channel which sends it a given
    argument. It returns a non nil error when the destination is closed

func Put[T any, int rules.Int, chanT Writable[T]](count int, dst chanT)
    Deprecated: use Push instead

func Putter[T any, chanT Writable[T]](dst chanT) func(T) error
    Deprecated: use Pusher instead

func RO[T any](c chan T) <-chan T
    Deprecated: use ReadOnly instead

func RW[T any](c <-chan T) chan T
    Deprecated: use ReadWrite

func ReadOnly[T any](c chan T) <-chan T
    ReadOnly coerces a read-write channel's type to read-only

func ReadWrite[T any](c <-chan T) chan T
    ReadWrite wraps a read-only channel with a read-write one

func Runes(s string) <-chan rune
    Runes yields each rune in s

func Send[T any, chanT Readable[T]](f func(T), src chanT)
    Deprecated: use Do instead

func Tee[T any](src <-chan T) (one, two chan T)
    Tee returns a pair of mutually independent copies of the source channel.
    Blocks at each iterative step until both channels have been read from,
    use Teen to avoid this behaviour

func Upto[T rules.Real](args ...T) (chan T, error)
    Upto returns an iterator whose content depends on the number of arguments as
    follows

        		# of args 	|| 	behaviour
        	 		 1	 	|| 	stop
        		 	 2	 	|| 	start, stop
        		 	 3	 	|| 	start, stop, step
                   else 	|| 	error

    step defaults to 1 start defaults to 0

func WO[T any](c chan T) chan<- T
    Deprecated: use WriteOnly instead

func Watch[T any](dst chan T, src <-chan T)
    Watch feeds dst with items received from src. It does not close either of
    them

func WriteOnly[T any](c chan T) chan<- T
    WriteOnly coerces a read-write channel's type to write-only


TYPES

type Readable[T any] interface{ chan T | <-chan T }

type Writable[T any] interface{ chan T | chan<- T }

```
