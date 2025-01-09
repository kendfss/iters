package chans

import (
	"math/rand"
	"testing"
	"unicode"

	"github.com/kendfss/oracle"
)

const (
	nTests = 10
	nItems = 10
	nMax   = 314159
)

func TestClose(t *testing.T) {
	ch := make(chan struct{})
	Close(ch)
	_, open := <-ch
	if open {
		t.Error("channel still open")
	}
}

func TestDrain(t *testing.T) {
	ctr := 0
	for range [nTests]struct{}{} {
		buf := oracle.RandNums[int8](nItems)
		var ch chan int8
		done := New[struct{}]()
		ch = RW(Lazify(buf))
		go Drain(ch)
		go func() {
			defer close(done)
			defer Push(1, done)
			for {
				_, open := <-ch
				if open {
					ctr++
					continue
				}
				break
			}
		}()
		<-done
	}
	if ctr == 0 {
		t.Error("wasn't able to increment counter")
	}
}

func TestPut(t *testing.T) {
	for range [nTests]struct{}{} {
		count := rand.Intn(nMax)
		ch := New[int]()
		go func() { defer close(ch); Push(count, ch) }()
		for i := 0; i < count; i++ {
			<-ch
		}
		_, open := <-ch
		if open {
			t.Fatal("chan still open")
		}
	}
}

func TestLazify(t *testing.T) {
	for range [nTests]struct{}{} {
		buf := oracle.RandNums[int](nItems)
		src := Lazify(buf)
		for i, want := range buf {
			have := <-src
			if have != want {
				t.Log("buf", buf)
				t.Errorf("%d: have %d, want %d", i, have, want)
			}
		}
		_, ok := <-src
		if ok {
			t.Error("channel still open")
		}
	}
}

// func TestExtend(t *testing.T) {
// 	nItems := 314
// 	for range [nTests]struct{}{} {
// 		n := rand.Intn(nItems)
// 		src := From(0, oracle.RandBytes(n)...)
// 		chans := Teen(src, nItems)
// 		ch := RW(chans[0])
// 		go Extend(ch, chans[1:]...)
// 	}
// }

// func TestTeen(t *testing.T) {
// 	for range [1]struct{}{} {
// 		buf := oracle.RandNums[int](nItems)
// 		src := Lazify(buf)
// 		chans := Teen(src, nItems)
// 		if len(chans) != nItems {
// 			t.Fatalf("have %d chans, want %d", len(chans), nItems)
// 		}
// 		counter := map[int]int{}
// 		for _, e := range buf {
// 			counter[e] += nItems
// 		}
// 		if len(counter) != len(buf) {
// 			t.Fatalf("misconfigured: have %d counters, want %d", len(counter), len(buf))
// 		}
// 		wg := sync.WaitGroup{}
// 		lock := sync.Mutex{}
// 		for range buf {
// 			wg.Add(len(chans))
// 			for i, channel := range chans {
// 				go func(i int, channel <-chan int) {
// 					defer wg.Done()
// 					lock.Lock()
// 					defer lock.Unlock()
// 					counter[<-channel]--
// 				}(i, channel)
// 			}
// 			wg.Wait()
// 		}
// 		for key, have := range counter {
// 			if have != 0 {
// 				t.Errorf("%d: have %d, want 0", key, have)
// 			}
// 		}
// 		for i, ch := range chans {
// 			_, open := <-ch
// 			if open {
// 				t.Errorf("chan %d still open", i)
// 			}
// 		}
// 	}
// }

func TestTee(t *testing.T) {
	for range [1]struct{}{} {
		buf := oracle.RandNums[int](nItems)
		src := Lazify(buf)
		l, r := Tee(src)
		for i, want := range buf {
			have1 := <-l
			have2 := <-r
			if have1 != have2 {
				t.Fatalf("%d: receipts disagree: left %d, right %d", i, have1, have2)
			}
			if have1 != want {
				t.Log(buf)
				t.Fatalf("%dl: have %d, want %d", i, have1, want)
			}
			if have2 != want {
				t.Log(buf)
				t.Fatalf("%dr: have %d, want %d", i, have2, want)
			}
		}
		_, open1 := <-l
		_, open2 := <-r
		if open1 || open2 {
			t.Fatalf("chan still open")
		}
	}
}

func TestGet(t *testing.T) {
	for range [nTests]struct{}{} {
		buf := oracle.RandNums[int](nItems)
		ch := Lazify(buf)
		want := rand.Intn(len(buf))
		got := Pop(want, ch)
		have := len(got)
		if have != want {
			t.Errorf("have %d, want %d: %d -> %d", have, want, buf[:want], want)
		}
	}
}

func TestCount(t *testing.T) {
	for range [nTests]struct{}{} {
		want := rand.Intn(nItems)
		buf := oracle.RandNums[int](want)
		ch := Lazify(buf)
		have := Count(ch)
		if have != uint64(want) {
			t.Errorf("have %d, want %d", have, want)
		}
	}
}

func TestChars(t *testing.T) {
	pred := unicode.IsPrint
	t.Run("ascii", func(t *testing.T) {
		gen := oracle.RandBytes
		for range [nTests]struct{}{} {
			word := chars(pred, gen, rand.Intn(nItems))
			println(string(word))
			haves := Bytes(string(word))
			wants := Lazify(word)
			i := 0
			for {
				have, moreHaves := <-haves
				want, moreWants := <-wants
				if !(moreHaves && moreWants) {
					break
				}
				if moreHaves && !moreWants {
					t.Fatal("wants closed earlier")
				}
				if !moreHaves && moreHaves {
					t.Fatal("haves closed early")
				}
				if have != want {
					t.Errorf("%d: have %q, want %q", i, have, want)
				}
				i++
			}
		}
	})
	t.Run("unicode", func(t *testing.T) {
		gen := oracle.RandRunes
		for range [nTests]struct{}{} {
			word := chars(pred, gen, rand.Intn(nItems))
			println(string(word))
			haves := Runes(string(word))
			wants := Lazify(word)
			i := 0
			for {
				have, moreHaves := <-haves
				want, moreWants := <-wants
				if !(moreHaves && moreWants) {
					break
				}
				if moreHaves && !moreWants {
					t.Fatal("wants closed earlier")
				}
				if !moreHaves && moreHaves {
					t.Fatal("haves closed early")
				}
				if have != want {
					t.Errorf("%d: have %q, want %q", i, have, want)
				}
				i++
			}
		}
	})
}

func TestChain(t *testing.T) {
	t.Run("contents", func(t *testing.T) {
		for range nTests {
			nChans := rand.Intn(nTests)
			bufs := make([][]int, nChans)
			chans := make([]<-chan int, nChans)
			table := map[int]int{}
			for i := range nChans {
				buf := oracle.RandNums[int](nItems)
				bufs[i] = buf
				chans[i] = Lazify(buf)
				for _, e := range buf {
					table[e]++
				}
			}
			t.Log(table)
			for e := range Chain(chans...) {
				table[e]--
			}
			t.Log(table)
			for k, v := range table {
				if v != 0 {
					t.Errorf("%d: have %d, want 0", k, v)
				}
			}
		}
	})
	t.Run("capacity", func(t *testing.T) {
		for range nTests {
			want := rand.Intn(nMax)
			have := cap(Chain(make(chan struct{}), make(chan struct{}, want)))
			if have != want {
				t.Errorf("have %d, want %d", have, want)
			}
		}
	})
}
