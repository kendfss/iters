// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kendfss/oprs"
	"github.com/kendfss/oprs/math/real"
	"github.com/kendfss/oracle"
)

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1},
		nil,
		false,
	},
	{
		[]int{},
		nil,
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       []float64
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		[]float64{1, 2},
		[]float64{1, 2},
		true,
		true,
	},
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		false,
		true,
	},
}

func TestEqual(t *testing.T) {
	for _, test := range equalIntTests {
		if got := Equal(test.s1, test.s2); got != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := Equal(test.s1, test.s2); got != test.wantEqual {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
	}
}

func TestEqualFunc(t *testing.T) {
	for _, test := range equalIntTests {
		// if got := EqualFunc(test.s1, test.s2, equal[int]); got != test.want {
		if got := EqualFunc(equal[int], test.s1, test.s2); got != test.want {
			t.Errorf("EqualFunc(%v, %v, equal[int]) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		// if got := EqualFunc(test.s1, test.s2, equal[float64]); got != test.wantEqual {
		if got := EqualFunc(equal[float64], test.s1, test.s2); got != test.wantEqual {
			t.Errorf("Equal(%v, %v, equal[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
		// if got := EqualFunc(test.s1, test.s2, equalNaN[float64]); got != test.wantEqualNaN {
		if got := EqualFunc(equalNaN[float64], test.s1, test.s2); got != test.wantEqualNaN {
			t.Errorf("Equal(%v, %v, equalNaN[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqualNaN)
		}
	}

	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	if EqualFunc(offByOne[int], s1, s1) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = true, want false", s1, s1)
	}
	if !EqualFunc(offByOne[int], s1, s2) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = false, want true", s1, s2)
	}

	s3 := []string{"a", "b", "c"}
	s4 := []string{"A", "B", "C"}
	if !EqualFunc(strings.EqualFold, s3, s4) {
		t.Errorf("EqualFunc(%v, %v, strings.EqualFold) = false, want true", s3, s4)
	}

	cmpIntString := func(v1 int, v2 string) bool {
		return string(rune(v1)-1+'a') == v2
	}
	if !EqualFunc(cmpIntString, s1, s3) {
		t.Errorf("EqualFunc(%v, %v, cmpIntString) = false, want true", s1, s3)
	}
}

var compareIntTests = []struct {
	s1, s2 []int
	want   int
}{
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		-1,
	},
	{
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3},
		+1,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 4, 3},
		-1,
	},
	{
		[]int{1, 4, 3},
		[]int{1, 2, 3},
		+1,
	},
}

var compareFloatTests = []struct {
	s1, s2 []float64
	want   int
}{
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, math.NaN(), 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3, 4},
		[]float64{1, 2, math.NaN()},
		+1,
	},
}

func TestCompare(t *testing.T) {
	intWant := func(want bool) string {
		if want {
			return "0"
		}
		return "!= 0"
	}
	for _, test := range equalIntTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}
	for _, test := range equalFloatTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.wantEqualNaN {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqualNaN))
		}
	}

	for _, test := range compareIntTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range compareFloatTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
}

func TestCompareFunc(t *testing.T) {
	intWant := func(want bool) string {
		if want {
			return "0"
		}
		return "!= 0"
	}
	for _, test := range equalIntTests {
		if got := CompareFunc(equalToCmp(equal[int]), test.s1, test.s2); (got == 0) != test.want {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[int])) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}
	for _, test := range equalFloatTests {
		if got := CompareFunc(equalToCmp(equal[float64]), test.s1, test.s2); (got == 0) != test.wantEqual {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[float64])) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqual))
		}
	}

	for _, test := range compareIntTests {
		if got := CompareFunc(cmp[int], test.s1, test.s2); got != test.want {
			t.Errorf("CompareFunc(%v, %v, cmp[int]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range compareFloatTests {
		if got := CompareFunc(cmp[float64], test.s1, test.s2); got != test.want {
			t.Errorf("CompareFunc(%v, %v, cmp[float64]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}

	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	if got := CompareFunc(equalToCmp(offByOne[int]), s1, s2); got != 0 {
		t.Errorf("CompareFunc(%v, %v, offByOne) = %d, want 0", s1, s2, got)
	}

	s3 := []string{"a", "b", "c"}
	s4 := []string{"A", "B", "C"}
	if got := CompareFunc(strings.Compare, s3, s4); got != 1 {
		t.Errorf("CompareFunc(%v, %v, strings.Compare) = %d, want 1", s3, s4, got)
	}

	compareLower := func(v1, v2 string) int {
		return strings.Compare(strings.ToLower(v1), strings.ToLower(v2))
	}
	if got := CompareFunc(compareLower, s3, s4); got != 0 {
		t.Errorf("CompareFunc(%v, %v, compareLower) = %d, want 0", s3, s4, got)
	}

	cmpIntString := func(v1 int, v2 string) int {
		return strings.Compare(string(rune(v1)-1+'a'), v2)
	}
	if got := CompareFunc(cmpIntString, s1, s3); got != 0 {
		t.Errorf("CompareFunc(%v, %v, cmpIntString) = %d, want 0", s1, s3, got)
	}
}

var indexTests = []struct {
	s    []int
	v    int
	want int
}{
	{
		nil,
		0,
		-1,
	},
	{
		[]int{},
		0,
		-1,
	},
	{
		[]int{1, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 3, 2},
		2,
		1,
	},
}

func TestIndex(t *testing.T) {
	for _, test := range indexTests {
		if got := Index(test.v, test.s); got != test.want {
			t.Errorf("Index(%v, %v) = %d, want %d", test.s, test.v, got, test.want)
		}
	}
}

func TestIndexFunc(t *testing.T) {
	for _, test := range indexTests {
		if got := IndexFunc(equal[int], test.v, test.s); got != test.want {
			t.Errorf("IndexFunc(%v, equalToIndex(equal[int], %v)) = %d, want %d", test.s, test.v, got, test.want)
		}
		// if got := IndexFunc(test.s, equalToIndex(equal[int], test.v)); got != test.want {
		// 	t.Errorf("IndexFunc(%v, equalToIndex(equal[int], %v)) = %d, want %d", test.s, test.v, got, test.want)
		// }
	}

	s1 := []string{"hi", "HI"}
	first, second := s1[0], s1[1]
	if got := IndexFunc(equal[string], second, s1); got != 1 {
		// if got := IndexFunc(s1, equalToIndex(equal[string], "HI")); got != 1 {
		t.Errorf("IndexFunc(%v, equalToIndex(equal[string], %q)) = %d, want %d", s1, second, got, 1)
	}
	if got := IndexFunc(equal[string], first, s1); got != 0 {
		// if got := IndexFunc(s1, equalToIndex(strings.EqualFold, "HI")); got != 0 {
		t.Errorf("IndexFunc(%v, equalToIndex(strings.EqualFold, %q)) = %d, want %d", s1, first, got, 0)
	}
}

func TestContains(t *testing.T) {
	for _, test := range indexTests {
		if got := Contains(test.s, test.v); got != (test.want != -1) {
			t.Errorf("Contains(%v, %v) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}
}

var insertTests = []struct {
	s    []int
	add  []int
	want []int
	i    int
}{
	{
		s:    []int{1, 2, 3},
		i:    0,
		add:  []int{4},
		want: []int{4, 1, 2, 3},
	},
	{
		s:    []int{1, 2, 3},
		i:    1,
		add:  []int{4},
		want: []int{1, 4, 2, 3},
	},
	{
		s:    []int{1, 2, 3},
		i:    3,
		add:  []int{4},
		want: []int{1, 2, 3, 4},
	},
	{
		s:    []int{1, 2, 3},
		i:    2,
		add:  []int{4, 5},
		want: []int{1, 2, 4, 5, 3},
	},
}

func TestInsert(t *testing.T) {
	s := []int{1, 2, 3}
	if got := Insert(s, 0); !Equal(got, s) {
		t.Errorf("Insert(%v, 0) = %v, want %v", s, got, s)
	}
	for _, test := range insertTests {
		copy := Clone(test.s)
		if got := Insert(copy, test.i, test.add...); !Equal(got, test.want) {
			t.Errorf("Insert(%v, %d, %v...) = %v, want %v", test.s, test.i, test.add, got, test.want)
		}
	}
}

var deleteTests = []struct {
	s    []int
	want []int
	i    int
	j    int
}{
	{
		s:    []int{1, 2, 3},
		i:    0,
		j:    0,
		want: []int{1, 2, 3},
	},
	{
		s:    []int{1, 2, 3},
		i:    0,
		j:    1,
		want: []int{2, 3},
	},
	{
		s:    []int{1, 2, 3},
		i:    3,
		j:    3,
		want: []int{1, 2, 3},
	},
	{
		s:    []int{1, 2, 3},
		i:    0,
		j:    2,
		want: []int{3},
	},
	{
		s:    []int{1, 2, 3},
		i:    0,
		j:    3,
		want: []int{},
	},
}

func TestDelete(t *testing.T) {
	for _, test := range deleteTests {
		copy := Clone(test.s)
		if got := Delete(copy, test.i, test.j); !Equal(got, test.want) {
			t.Errorf("Delete(%v, %d, %d) = %v, want %v", test.s, test.i, test.j, got, test.want)
		}
	}
}

func TestClone(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := Clone(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clone(%v) = %v, want %v", s1, s2, s1)
	}
	s1[0] = 4
	want := []int{1, 2, 3}
	if !Equal(s2, want) {
		t.Errorf("Clone(%v) changed unexpectedly to %v", want, s2)
	}
	if got := Clone([]int(nil)); got != nil {
		t.Errorf("Clone(nil) = %#v, want nil", got)
	}
	if got := Clone(s1[:0]); got == nil || len(got) != 0 {
		t.Errorf("Clone(%v) = %#v, want %#v", s1[:0], got, s1[:0])
	}
}

var compactTests = []struct {
	s    []int
	want []int
}{
	{
		nil,
		nil,
	},
	{
		[]int{1},
		[]int{1},
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	},
	{
		[]int{1, 1, 2},
		[]int{1, 2},
	},
	{
		[]int{1, 2, 1},
		[]int{1, 2, 1},
	},
	{
		[]int{1, 2, 2, 3, 3, 4},
		[]int{1, 2, 3, 4},
	},
}

func TestCompact(t *testing.T) {
	for _, test := range compactTests {
		copy := Clone(test.s)
		if got := Compact(copy); !Equal(got, test.want) {
			t.Errorf("Compact(%v) = %v, want %v", test.s, got, test.want)
		}
	}
}

func TestCompactFunc(t *testing.T) {
	for _, test := range compactTests {
		copy := Clone(test.s)
		if got := CompactFunc(equal[int], copy); !Equal(got, test.want) {
			t.Errorf("CompactFunc(%v, equal[int]) = %v, want %v", test.s, got, test.want)
		}
	}

	s1 := []string{"a", "a", "A", "B", "b"}
	copy := Clone(s1)
	want := []string{"a", "B"}
	if got := CompactFunc(strings.EqualFold, copy); !Equal(got, want) {
		t.Errorf("CompactFunc(%v, strings.EqualFold) = %v, want %v", s1, got, want)
	}
}

func TestGrow(t *testing.T) {
	s1 := []int{1, 2, 3}
	copy := Clone(s1)
	s2 := Grow(copy, 1000)
	if !Equal(s1, s2) {
		t.Errorf("Grow(%v) = %v, want %v", s1, s2, s1)
	}
	if cap(s2) < 1000+len(s1) {
		t.Errorf("after Grow(%v) cap = %d, want >= %d", s1, cap(s2), 1000+len(s1))
	}
}

func TestClip(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5, 6}[:3]
	orig := Clone(s1)
	if len(s1) != 3 {
		t.Errorf("len(%v) = %d, want 3", s1, len(s1))
	}
	if cap(s1) < 6 {
		t.Errorf("cap(%v[:3]) = %d, want >= 6", orig, cap(s1))
	}
	s2 := Clip(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clip(%v) = %v, want %v", s1, s2, s1)
	}
	if cap(s2) != 3 {
		t.Errorf("cap(Clip(%v)) = %d, want 3", orig, cap(s2))
	}
}

const (
	nTests = 10
	nItems = 10
	nMax   = 10
)

func TestShortest(t *testing.T) {
	type test struct {
		short, long []int
	}

	const (
		nItems = 100
	)

	for i := 0; i < nTests; i++ {
		size := rand.Intn(nItems)
		data := test{}
		for len(data.short) == len(data.long) {
			data = test{short: oracle.Mkr(size, nMax), long: oracle.Mkr(size+1, nMax)}
		}
		result := Shortest(data.short, data.long)
		if result != 0 {
			oracle.Quitf(t, "#%d (%#v): result is %d but should be 0", i, data, result)
		}
	}
}

func TestLongest(t *testing.T) {
	type test struct {
		short, long []int
	}

	const (
		nItems = 100
	)

	for i := 0; i < nTests; i++ {
		size := rand.Intn(nItems)
		data := test{}
		for len(data.short) == len(data.long) {
			data = test{short: oracle.Mkr(size, nMax), long: oracle.Mkr(size+1, nMax)}
			data.short = Compact(data.short)
			data.long = Compact(data.long)
		}
		result := Longest(data.short, data.long)
		if result != 1 {
			oracle.Quitf(t, "#%d (%#v): result is %d but should be 1; %d !< %v", i, data, result, len(data.short), len(data.long))
		}
	}
}

func TestFilter(t *testing.T) {
	for i := 0; i < nTests; i++ {
		data := Upton[int](nItems)
		result := FilterFunc(oprs.IsEven[int], data)
		for j, e := range result {
			if e%2 == 1 {
				oracle.Quitf(t, "#%d.%d (%v -> %v): %d is odd, not even", i, j, data, result, e)
			}
		}
	}
}

func TestZip(t *testing.T) {
	for i := range Upton[int](nTests) {
		left := oracle.Mkr(nItems, nMax)
		right := oracle.Mkr(nItems, nMax)
		result := Zip2(left, right)
		for j, pair := range result {
			if want, have := left[j], pair.Left; want != have {
				oracle.Quitf(t, "#%d.%d (%v -> %v): want %d, have %d", i, j, left, pair, want, have)
			}
			if want, have := right[j], pair.Right; want != have {
				oracle.Quitf(t, "#%d.%d (%v -> %v): want %d, have %d", i, j, right, pair, want, have)
			}
		}
	}
}

func TestFlatter(t *testing.T) {
	const (
		nItems = 4
		nMax   = 1000
	)
	for i := range Upton[int](nTests) {
		data := [][]int{
			oracle.Mkr(nItems, nMax),
			oracle.Mkr(nItems, nMax),
		}
		result := Flatter(data)

		ptr := -1
		for j, have := range result {
			if j%nItems == 0 {
				ptr++
			}
			want := data[ptr][j%nItems]
			if want != have {
				oracle.Quitf(t, "#%d.%d (%v -> %v): want %d, have %d", i, j, data, result, want, have)
			}
		}
	}
}

func TestChain(t *testing.T) {
	const (
		nMax = 100
	)
	for i := range Upton[int](nTests) {
		first, second := oracle.Mkr(nItems, nMax), oracle.Mkr(nItems, nMax)
		both := [][]int{first, second}
		expected := append(Clone(first), second...)
		result := Chain(first, second)
		if len(expected) != len(result) {
			oracle.Quitf(t, "#%d (%v -> %v): result should have %d items but has %d", i, expected, result, len(expected), len(result))
		}
		for j, ptr := 0, -1; j < 2*nItems; j++ {
			if j%nItems == 0 {
				ptr++
			}
			have := result[j]
			want := both[ptr][j%nItems]
			if have != want {
				oracle.Quitf(t, "#%d.%d (%v -> %v): result#%d should be %d but is %d", i, j, expected, result, j, want, have)
			}
		}
	}
}

// func TestChained(t *testing.T) {
// 	const (
// 		nMax = 100
// 	)
// 	for i := range Upton[int](nTests) {
// 		first, second := oracle.Mkr(nItems, nMax), oracle.Mkr(nItems, nMax)
// 		both := [][]int{first, second}
// 		expected := append(Clone(first), second...)
// 		result := Chain(first, second)
// 		if len(expected) != len(result) {
// 			oracle.Quitf(t, "#%d (%v -> %v): result should have %d items but has %d", i, expected, result, len(expected), len(result))
// 		}
// 		for j, ptr := 0, -1; j < 2*nItems; j++ {
// 			if j%nItems == 0 {
// 				ptr++
// 			}
// 			have := result[j]
// 			want := both[ptr][j%nItems]
// 			if have != want {
// 				oracle.Quitf(t, "#%d.%d (%v -> %v): result#%d should be %d but is %d", i, j, expected, result, j, want, have)
// 			}
// 		}
// 	}
// }

func TestReversed(t *testing.T) {
	for i := range Upton[int](nTests) {
		data := oracle.Mkr(nItems, nMax)
		first := Reversed(data)
		second := Reversed(first)
		if want, have := len(data), len(first); want != have {
			oracle.Quitf(t, "#%d: first call has %d items but should have %d", i, want, have)
		}
		if want, have := len(data), len(second); want != have {
			oracle.Quitf(t, "#%d: first call has %d items but should have %d", i, want, have)
		}
		for j, have := range second {
			if want := data[j]; have != want {
				oracle.Quitf(t, "#%d.%d: have %d, want %d", i, j, have, want)
			}
		}
	}
}

func TestSwap(t *testing.T) {
	for i := range Upton[int](nTests) {
		orig := oracle.Mkr(nItems, nMax)
		clon := append([]int{}, orig...)

		indices := oracle.Mkr(2, nItems)
		for orig[indices[0]] == orig[indices[1]] {
			indices = oracle.Mkr(2, nItems)
		}
		j, k := indices[0], indices[1]
		Swap(clon, j, k)
		have1, have2 := clon[k], clon[j]
		want1, want2 := orig[j], orig[k]
		if have1 != want1 || have2 != want2 {
			fmt.Println(i, j, k)
			fmt.Println(orig)
			fmt.Println(clon)
			oracle.Quitf(t, "#%d @(%d, %d): have (%d, %d), want (%d, %d)", i, j, k, have1, have2, want1, want2)
		}
	}
}

func TestLen(t *testing.T) {
	for i := range Upton[int](nTests) {
		data := oracle.Mkr(rand.Intn(nMax), nMax)
		if want, have := len(data), Len[int](data); want != have {
			oracle.Inequiv(t, i, have, want)
		}
	}
}

func TestSelect(t *testing.T) {
	for i := range Upton[int](nTests) {
		data := oracle.Mkr(nItems, nMax)
		indices := oracle.Mkr(rand.Intn(nItems), nMax)
		results := Select(data, indices)
		if len(indices) != len(results) {
			t.Fatal("lengths do not agree")
		}
		for j, e := range indices {
			if want, have := data[e], results[j]; have != want {
				t.Log(i, j)
				t.Log(data)
				t.Log(indices)
				t.Log(results)
				oracle.Inequiv(t, j, have, want)
			}
		}

	}
}

func TestAll(t *testing.T) {
	data := Ones(nItems)
	pred := oprs.Is(1)
	t.Run("true test", func(t *testing.T) {
		if !All(Cast(pred, data)) {
			oracle.Inequiv(t, 0, false, true)
		}
	})
	t.Run("false test", func(t *testing.T) {
		data = append(data, 2)
		if All(Cast(pred, data)) {
			oracle.Inequiv(t, 0, true, false)
		}
	})
}

func TestAny(t *testing.T) {
	data := Ones(nItems)
	pred := oprs.Is(2)
	t.Run("false test", func(t *testing.T) {
		if Any(Cast(pred, data)...) {
			oracle.Inequiv(t, 0, false, true)
		}
	})
	t.Run("true test", func(t *testing.T) {
		data = append(data, 2)
		if !Any(Cast(pred, data)...) {
			oracle.Inequiv(t, 0, true, false)
		}
	})
}

func TestMax(t *testing.T) {
	data := Upton[int](2)
	t.Run("basic", func(t *testing.T) {
		want := 1
		if have := Max(data...); have != want {
			oracle.Inequiv(t, 0, have, want)
		}
	})
	t.Run("reversed", func(t *testing.T) {
		Reverse(data)
		want := 0
		if have := Max(data...); have != want {
			t.Log(data)
			oracle.Inequiv(t, 1, have, want)
		}
	})
	t.Run("take first", func(t *testing.T) {
		data = Ones(2)
		want := 0
		if have := Max(data...); have != want {
			oracle.Inequiv(t, 2, have, want)
		}
	})
}

func TestMin(t *testing.T) {
	data := Upton[int](2)
	t.Run("basic", func(t *testing.T) {
		want := 0
		if have := Min(data...); have != want {
			oracle.Inequiv(t, 0, have, want)
		}
	})
	t.Run("revesed", func(t *testing.T) {
		Reverse(data)
		want := 1
		if have := Min(data...); have != want {
			t.Log(data)
			oracle.Inequiv(t, 1, have, want)
		}
	})
	t.Run("take first", func(t *testing.T) {
		data = Ones(2)
		want := 0
		if have := Min(data...); have != want {
			oracle.Inequiv(t, 2, have, want)
		}
	})
}

func TestExtremal(t *testing.T) {
	t.Run("lt ==> Min", func(t *testing.T) {
		for i := range Upton[int](nTests) {
			data := oracle.Mkr(nItems, nMax)
			want := Min(data...)
			have := Extremal(oprs.Lt[int], data...)
			if want != have {
				t.Log(data)
				oracle.Inequiv(t, i, have, want)
			}
		}
	})
	t.Run("gt ==> Max", func(t *testing.T) {
		for i := range Upton[int](nTests) {
			data := oracle.Mkr(nItems, nMax)
			want := Max(data...)
			have := Extremal(oprs.Gt[int], data...)
			if want != have {
				t.Log(data, Max(data...))
				oracle.Inequiv(t, i, have, want)
			}
		}
	})
}

func TestOnes(t *testing.T) {
	for i := 0; i < 10; i++ {
		result := Ones(i)
		if len(result) != i {
			oracle.Quitf(t, "result has length %d, expected %d", len(result), i)
		}
		for j, e := range result {
			if e != 1 {
				oracle.Quitf(t, "element #%d (%d) != 1", j, e)
			}
		}
	}
}

func TestReduce(t *testing.T) {
	type test[T any] struct {
		ans T
		op  func(T, T) T
		seq []T
	}

	intTests := []test[int]{
		{seq: []int{1, 1, 2}, op: real.Mul[int], ans: 2},
		{seq: []int{99, 11, 3}, op: real.Div[int], ans: 3},
		{seq: []int{8, 2, 2}, op: real.Sub[int], ans: 4},
	}
	for i, test := range intTests {
		result := Reduce(test.op, test.seq)
		if result != test.ans {
			oracle.Quitf(t, "#%d (%v): expected %q, got %q", i, test.seq, test.ans, result)
		}
	}

	strTests := []test[string]{
		{seq: strings.Split("quick brown dog", " "), op: real.Add[string], ans: "quickbrowndog"},
		{seq: strings.Split("jumps over the lazy fox", " "), op: real.Add[string], ans: "jumpsoverthelazyfox"},
		{seq: strings.Split("quick brown dog jumps over the lazy fox", " "), op: real.Add[string], ans: "quickbrowndogjumpsoverthelazyfox"},
	}
	for i, test := range strTests {
		result := Reduce(test.op, test.seq)
		if result != test.ans {
			oracle.Quitf(t, "#%d (%v): expected %q, got %q", i, test.seq, test.ans, result)
		}
	}
}

func TestUpto(t *testing.T) {
	type argSet struct {
		start, stop, step int
	}

	const (
		nMax = 100
	)
	t.Run("continuity", func(t *testing.T) {
		const (
			nMax = 50
		)
		t.Run("positive", func(t *testing.T) {
			for i := 0; i < nTests; i++ {
				a := oracle.Mkr(3, nMax)
				Sort(a)
				start, stop, step := a[0], a[2], a[1]
				args := argSet{start: start, stop: stop, step: step}
				result := Upto[int](start, stop, step)
				for j, e := range result[1:] {
					diff := e - step
					if result[j] != diff {
						oracle.Quitf(t, "#%d.%d (%v.%v): difference is %d, should be %d", i, j, args, result, diff, args.step)
					}
				}
			}
		})
		t.Run("negative", func(t *testing.T) {
			for i := 0; i < nTests; i++ {
				a := oracle.Mkr(3, nMax)
				Sort(a)
				start, stop, step := a[2], a[1], -a[1]
				args := argSet{start: start, stop: stop, step: step}
				result := Upto[int](start, stop, step)
				for j, e := range result[1:] {
					diff := e - step
					if result[j] != diff {
						oracle.Quitf(t, "#%d.%d (%v.%v): difference is %d, should be %d", i, j, args, result, diff, args.step)
					}
				}
			}
		})
	})

	t.Run("series", func(t *testing.T) {
		t.Run("triangular", func(t *testing.T) {
			t.Run("positive", func(t *testing.T) {
				for i := 0; i < nTests; i++ {
					val := rand.Intn(nMax)
					result := Upto[int](0, val, 1)
					sum, trg := Reduce(real.Add[int], result), oracle.Triangular(val)
					if sum != trg {
						oracle.Quitf(t, "#%d (%d -> %v): Range(0, %d, 1) should add up to %d but is %d", i, val, result, val, trg, sum)
					}
				}
			})
			t.Run("negative", func(t *testing.T) {
				for i := 0; i < nTests; i++ {
					val := rand.Intn(nMax)
					result := Upto[int](val, 0, -1)
					sum, trg := Reduce(real.Add[int], result), oracle.TriangularR(val)
					if sum != trg {
						oracle.Quitf(t, "#%d (%d -> %v): Range(%d, 0, -1) should add up to %d but is %d", i, val, result, val, trg, sum)
					}
				}
			})
		})
	})
}

func TestSnap(t *testing.T) {
	type tst struct {
		len int
		wid int
	}
	tests := []tst{
		{wid: 1, len: 3},
		{wid: 1, len: 5},
		{wid: 1, len: 4},
		{wid: 2, len: 4},
		{wid: 3, len: 4},
		{wid: 8, len: 21},
		{wid: 5, len: 21},
		{wid: 9, len: 21},
		{wid: 9, len: 30},
	}

	for i, test := range tests {
		result := Snap(test.wid, Anify(Ones(test.len)))
		chunks := real.Subtractions(test.len, test.wid)

		if len(result) != len(chunks) {
			oracle.Quitf(t, "result #%d has %d elements but expected %d elements", i, len(result), len(chunks))
		}
		for j, wedge := range result {
			if len(wedge) != chunks[j] {
				oracle.Quitf(t, "wedge #%d (%v) has %d elements but expected %d elements", j, wedge, len(wedge), chunks[j])
			}
		}
	}
}

func TestMap(t *testing.T) {
	for i := 0; i < nTests; i++ {
		data := oracle.Mkr(nItems, nMax)
		result := Cast(real.Succ[int], data)
		if len(result) != len(data) {
			oracle.Quitf(t, "#%d (%v -> %v): result has length %d, data has length %d", i, data, result, len(result), len(data))
		}
		for j, e := range result {
			if e-1 != data[j] {
				oracle.Quitf(t, "#%d.%d (%v -> %v): got %d, should have %d", i, j, data, result, e, result[j])
			}
		}
	}
}

func TestSplit(t *testing.T) {
	type tst struct {
		str string
		brk rune
	}

	tests := []tst{
		{str: "papyrus", brk: 'p'},
		{str: "cabbage", brk: 'b'},
		{str: "cabbage", brk: 'a'},
		{str: "balance", brk: 'a'},
		{str: "calculus", brk: 'c'},
		{str: "calculus", brk: 'l'},
		{str: "calculus", brk: 'u'},
	}
	toStr := func(arg []rune) string {
		return string(arg)
	}

	for i, test := range tests {
		result := Split(oracle.Runes(test.str), test.brk)
		expected := oracle.Runes2(strings.Split(test.str, string(test.brk)))
		if len(result) != len(expected) {
			t.Log(expected, result)
			oracle.Quitf(t, "result #%d has %d elements but expected %d elements", i, len(result), len(expected))
		}
		for j, wedge := range result {
			if len(wedge) != len(expected[j]) {
				t.Log(expected, result)
				t.Log(Cast(toStr, expected))
				t.Log(Cast(toStr, result))
				oracle.Quitf(t, "wedge #%d.%d (%v) has %d elements but expected %d elements", i, j, wedge, len(wedge), expected[j])
			}
			if !Equal(wedge, expected[j]) {
				t.Log(expected, result)
				t.Log(Cast(toStr, expected))
				t.Log(Cast(toStr, result))
				oracle.Quitf(t, "wedge #%d is %v but should be %v", j, wedge, expected[j])
			}
		}
	}
}

func TestSplitAfter(t *testing.T) {
	type tst struct {
		str string
		brk rune
	}

	tests := []tst{
		{str: "papyrus", brk: 'p'},
		{str: "cabbage", brk: 'b'},
		{str: "cabbage", brk: 'a'},
		{str: "balance", brk: 'a'},
		{str: "calculus", brk: 'c'},
		{str: "calculus", brk: 'l'},
		{str: "calculus", brk: 'u'},
	}
	toStr := func(arg []rune) string {
		return string(arg)
	}

	for i, test := range tests {
		result := SplitAfter(oracle.Runes(test.str), test.brk)
		expected := oracle.Runes2(strings.SplitAfter(test.str, string(test.brk)))
		if len(result) != len(expected) {
			t.Log(expected, result)
			oracle.Quitf(t, "result #%d has %d elements but expected %d elements", i, len(result), len(expected))
		}
		for j, wedge := range result {
			if len(wedge) != len(expected[j]) {
				t.Log(expected, result)
				t.Log(Cast(toStr, expected))
				t.Log(Cast(toStr, result))
				oracle.Quitf(t, "wedge #%d.%d (%v) has %d elements but expected %d elements", i, j, wedge, len(wedge), expected[j])
			}
			if !Equal(wedge, expected[j]) {
				t.Log(expected, result)
				t.Log(Cast(toStr, expected))
				t.Log(Cast(toStr, result))
				oracle.Quitf(t, "wedge #%d is %v but should be %v", j, wedge, expected[j])
			}
		}
	}
}

func TestRotated(t *testing.T) {
	type test struct {
		slice []int
		want  []int
		steps int
	}
	tests := []test{
		{slice: []int{0, 1}, steps: 0, want: []int{0, 1}},
		{slice: []int{0, 1}, steps: 1, want: []int{1, 0}},
		{slice: []int{0, 1}, steps: 2, want: []int{0, 1}},
		{slice: []int{0, 1}, steps: 3, want: []int{1, 0}},
		{slice: []int{0, 1}, steps: -1, want: []int{1, 0}},
		{slice: []int{0, 1}, steps: -2, want: []int{0, 1}},
		{slice: []int{0, 1}, steps: -3, want: []int{1, 0}},
		{slice: Upton[int](rand.Int() % rand.Intn(917)), steps: randSign(rand.Int())},
		{slice: Upton[int](rand.Int() % rand.Intn(917)), steps: randSign(rand.Int())},
		{slice: Upton[int](rand.Int() % rand.Intn(917)), steps: randSign(rand.Int())},
		{slice: Upton[int](rand.Int() % rand.Intn(917)), steps: randSign(rand.Int())},
		{slice: Upton[int](rand.Int() % rand.Intn(917)), steps: randSign(rand.Int())},
		{slice: Upton[int](rand.Int() % rand.Intn(917)), steps: randSign(rand.Int())},
	}
	for i, test := range tests {
		have := Rotated(test.slice, test.steps)
		if len(test.want) > 0 {
			assert.Equal(t, len(have), len(test.want), "#%d: length failure", i)
			assert.Equal(t, have, test.want, "#%d: value failure", i)

		}
		assert.Equal(t, Rotated(have, -test.steps), test.slice, "#%d: value failure", i)
	}
}

func TestRepeat(t *testing.T) {
	for range Upton[int](nTests) {
		count, seed := rand.Intn(nItems), rand.Intn(nMax)
		slice := Repeat(seed, count)
		assert.Equal(t, count, len(slice), "len(slice) != count")
		for _, e := range slice {
			assert.Equal(t, seed, e, "element != seed")
		}
	}
	t.Run("equals Tee", func(l *testing.T) {
		for i := range Upton[int](nTests) {
			count, seed := rand.Intn(nItems)+1, make([]int, rand.Intn(nItems)+1)
			// fmt.Printf("count:\t%d\nseed:\t%d\n", count, seed)
			repeat := Repeat(seed, count)
			tee := Tee(seed, count)
			assert.Equal(l, count, len(repeat), "#%d:\tlen(slice) != count", i)
			assert.Equal(l, count, len(tee), "#%d:\tlen(tee) != count", i)
			for j, f := range Zip3(tee, repeat) {
				r, t := f()
				if !assert.Equal(l, seed, t, "#%d:\ttee[%d] != seed", i, j) {
					return
				}
				if !assert.Equal(l, seed, r, "#%d:\trepeat[%d] != seed", i, j) && r != nil {
					return
				}
			}
		}
	})
}

func TestTee(l *testing.T) {
	// tested under TestRepeat
	l.Run("chain zero != self", func(l *testing.T) {
		t := oracle.RandNums[int](rand.Intn(20))
		assert.NotEqual(l, t, Chain(Tee(t, 0)...))
	})
	l.Run("chain one == self", func(l *testing.T) {
		t := oracle.RandNums[int](rand.Intn(20))
		assert.Equal(l, t, Chain(Tee(t, 1)...))
	})
}

// func TestPermutations(l *testing.T) {
// 	// # permutations('ABCD', 2) --> AB AC AD BA BC BD CA CB CD DA DB DC
// 	// # permutations(range(3)) --> 012 021 102 120 201 210
// 	type test struct {
// 		want [][]int
// 		arg  []int
// 		r    int
// 	}
// 	tests := []test{
// 		{r: 2, arg: []int{'A', 'B', 'C', 'D'}, want: [][]int{{'A', 'B'}, {'A', 'C'}, {'A', 'D'}, {'B', 'A'}, {'B', 'C'}, {'B', 'D'}, {'C', 'A'}, {'C', 'B'}, {'C', 'D'}, {'D', 'A'}, {'D', 'B'}, {'D', 'C'}}},
// 		{r: 0, arg: Upton[int](3), want: [][]int{{'0', '1', '2'}, {'0', '2', '1'}, {'1', '0', '2'}, {'1', '2', '0'}, {'2', '0', '1'}, {'2', '1', '0'}}},
// 		// {arg: Upton[int](0), want: [][]int{{}}},
// 		// {arg: Upton[int](1), want: [][]int{{0}}},
// 		// {arg: Upton[int](2), want: [][]int{{0, 1}, {1, 0}}},
// 		// {arg: Upton[int](3), want: [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}},
// 		// {arg: Upton[int](4), want: [][]int{{0, 1, 2, 3}, {0, 1, 3, 2}, {0, 2, 1, 3}, {0, 2, 3, 1}, {0, 3, 1, 2}, {0, 3, 2, 1}, {1, 0, 2, 3}, {1, 0, 3, 2}, {1, 2, 0, 3}, {1, 2, 3, 0}, {1, 3, 0, 2}, {1, 3, 2, 0}, {2, 0, 1, 3}, {2, 0, 3, 1}, {2, 1, 0, 3}, {2, 1, 3, 0}, {2, 3, 0, 1}, {2, 3, 1, 0}, {3, 0, 1, 2}, {3, 0, 2, 1}, {3, 1, 0, 2}, {3, 1, 2, 0}, {3, 2, 0, 1}, {3, 2, 1, 0}}},
// 		// {arg: Upton[int](5), want: [][]int{{0, 1, 2, 3, 4}, {0, 1, 2, 4, 3}, {0, 1, 3, 2, 4}, {0, 1, 3, 4, 2}, {0, 1, 4, 2, 3}, {0, 1, 4, 3, 2}, {0, 2, 1, 3, 4}, {0, 2, 1, 4, 3}, {0, 2, 3, 1, 4}, {0, 2, 3, 4, 1}, {0, 2, 4, 1, 3}, {0, 2, 4, 3, 1}, {0, 3, 1, 2, 4}, {0, 3, 1, 4, 2}, {0, 3, 2, 1, 4}, {0, 3, 2, 4, 1}, {0, 3, 4, 1, 2}, {0, 3, 4, 2, 1}, {0, 4, 1, 2, 3}, {0, 4, 1, 3, 2}, {0, 4, 2, 1, 3}, {0, 4, 2, 3, 1}, {0, 4, 3, 1, 2}, {0, 4, 3, 2, 1}, {1, 0, 2, 3, 4}, {1, 0, 2, 4, 3}, {1, 0, 3, 2, 4}, {1, 0, 3, 4, 2}, {1, 0, 4, 2, 3}, {1, 0, 4, 3, 2}, {1, 2, 0, 3, 4}, {1, 2, 0, 4, 3}, {1, 2, 3, 0, 4}, {1, 2, 3, 4, 0}, {1, 2, 4, 0, 3}, {1, 2, 4, 3, 0}, {1, 3, 0, 2, 4}, {1, 3, 0, 4, 2}, {1, 3, 2, 0, 4}, {1, 3, 2, 4, 0}, {1, 3, 4, 0, 2}, {1, 3, 4, 2, 0}, {1, 4, 0, 2, 3}, {1, 4, 0, 3, 2}, {1, 4, 2, 0, 3}, {1, 4, 2, 3, 0}, {1, 4, 3, 0, 2}, {1, 4, 3, 2, 0}, {2, 0, 1, 3, 4}, {2, 0, 1, 4, 3}, {2, 0, 3, 1, 4}, {2, 0, 3, 4, 1}, {2, 0, 4, 1, 3}, {2, 0, 4, 3, 1}, {2, 1, 0, 3, 4}, {2, 1, 0, 4, 3}, {2, 1, 3, 0, 4}, {2, 1, 3, 4, 0}, {2, 1, 4, 0, 3}, {2, 1, 4, 3, 0}, {2, 3, 0, 1, 4}, {2, 3, 0, 4, 1}, {2, 3, 1, 0, 4}, {2, 3, 1, 4, 0}, {2, 3, 4, 0, 1}, {2, 3, 4, 1, 0}, {2, 4, 0, 1, 3}, {2, 4, 0, 3, 1}, {2, 4, 1, 0, 3}, {2, 4, 1, 3, 0}, {2, 4, 3, 0, 1}, {2, 4, 3, 1, 0}, {3, 0, 1, 2, 4}, {3, 0, 1, 4, 2}, {3, 0, 2, 1, 4}, {3, 0, 2, 4, 1}, {3, 0, 4, 1, 2}, {3, 0, 4, 2, 1}, {3, 1, 0, 2, 4}, {3, 1, 0, 4, 2}, {3, 1, 2, 0, 4}, {3, 1, 2, 4, 0}, {3, 1, 4, 0, 2}, {3, 1, 4, 2, 0}, {3, 2, 0, 1, 4}, {3, 2, 0, 4, 1}, {3, 2, 1, 0, 4}, {3, 2, 1, 4, 0}, {3, 2, 4, 0, 1}, {3, 2, 4, 1, 0}, {3, 4, 0, 1, 2}, {3, 4, 0, 2, 1}, {3, 4, 1, 0, 2}, {3, 4, 1, 2, 0}, {3, 4, 2, 0, 1}, {3, 4, 2, 1, 0}, {4, 0, 1, 2, 3}, {4, 0, 1, 3, 2}, {4, 0, 2, 1, 3}, {4, 0, 2, 3, 1}, {4, 0, 3, 1, 2}, {4, 0, 3, 2, 1}, {4, 1, 0, 2, 3}, {4, 1, 0, 3, 2}, {4, 1, 2, 0, 3}, {4, 1, 2, 3, 0}, {4, 1, 3, 0, 2}, {4, 1, 3, 2, 0}, {4, 2, 0, 1, 3}, {4, 2, 0, 3, 1}, {4, 2, 1, 0, 3}, {4, 2, 1, 3, 0}, {4, 2, 3, 0, 1}, {4, 2, 3, 1, 0}, {4, 3, 0, 1, 2}, {4, 3, 0, 2, 1}, {4, 3, 1, 0, 2}, {4, 3, 1, 2, 0}, {4, 3, 2, 0, 1}, {4, 3, 2, 1, 0}}},
// 	}
// 	for i, test := range tests {
// 		have := Permutations(test.r, test.arg)
// 		// have := Permutations(test.arg)
// 		assert.Equal(l, test.want, have, "#%d:\n\targ:%d\n\trep:%d", i, test.arg, test.r)
// 	}
// }

// func TestProduct(l *testing.T) {
// 	type test struct {
// 		want [][]int
// 		arg  []int
// 		r    int
// 	}
// 	tests := []test{
// 		{r: 2, arg: []int{'A', 'B', 'C', 'D'}, want: [][]int{{'A', 'B'}, {'A', 'C'}, {'A', 'D'}, {'B', 'A'}, {'B', 'C'}, {'B', 'D'}, {'C', 'A'}, {'C', 'B'}, {'C', 'D'}, {'D', 'A'}, {'D', 'B'}, {'D', 'C'}}},
// 		{r: 0, arg: Upton[int](3), want: [][]int{{'0', '1', '2'}, {'0', '2', '1'}, {'1', '0', '2'}, {'1', '2', '0'}, {'2', '0', '1'}, {'2', '1', '0'}}},
// 	}
// 	for _, test := range tests {
// 		have := Product(test.r, test.arg)
// 		assert.Equal(l, test.want, have)
// 	}
// }

// func TestCombinations(t *testing.T) {
// 	type test struct {
// 		want [][]int
// 		arg  []int
// 		r    int
// 	}
// 	tests := []test{
// 		// Combinations('ABCD', 2) --> AB AC AD BC BD CD
// 		{r: 2, arg: []int{'a', 'b', 'c', 'd'}, want: [][]int{{'a', 'b'}, {'a', 'c'}, {'a', 'd'}, {'b', 'c'}, {'b', 'd'}, {'c', 'd'}}},
// 		// Combinations(range(4), 3) --> 012 013 023 123
// 		{r: 3, arg: Upton[int](4), want: [][]int{{0, 1, 2}, {0, 1, 3}, {0, 2, 3}, {1, 2, 3}}},
// 	}

// 	for i, test := range tests {
// 		have := Combinations(test.arg, test.r)
// 		assert.Equal(t, test.want, have, "#%d:\n\targ:%d\n\trep:%d", i, test.arg, test.r)
// 	}
// }

func TestGetxy(l *testing.T) {
	img := [][]int8{
		{0, 1, 2, 3},
		{4, 5, 6, 7},
		{8, 9, 10, 11},
		{12, 13, 14, 15},
	}
	pix := Chain(img...)
	pairs := [][]int{rand.Perm(len(img[0])), rand.Perm(len(img))}

	fmt.Println(pairs[0])
	fmt.Println(pairs[1])

	pairs = Zip(pairs...)
	fmt.Println(pairs)
	// for i, pair := range pairs {
	for _, pair := range pairs {
		x, y := pair[0], pair[1]
		fmt.Println(x, y)
		assert.Equal(l, img[y][x], Getxy(pix, 4, x, y), "x, y := %d, %d")
	}
}

func TestWindows(t *testing.T) {
	arg := Upton[int](10)
	wants := map[int][][]int{
		0:  nil,
		1:  {{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}},
		2:  {{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}},
		3:  {{0, 1, 2}, {1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}, {5, 6, 7}, {6, 7, 8}, {7, 8, 9}},
		4:  {{0, 1, 2, 3}, {1, 2, 3, 4}, {2, 3, 4, 5}, {3, 4, 5, 6}, {4, 5, 6, 7}, {5, 6, 7, 8}, {6, 7, 8, 9}},
		5:  {{0, 1, 2, 3, 4}, {1, 2, 3, 4, 5}, {2, 3, 4, 5, 6}, {3, 4, 5, 6, 7}, {4, 5, 6, 7, 8}, {5, 6, 7, 8, 9}},
		6:  {{0, 1, 2, 3, 4, 5}, {1, 2, 3, 4, 5, 6}, {2, 3, 4, 5, 6, 7}, {3, 4, 5, 6, 7, 8}, {4, 5, 6, 7, 8, 9}},
		7:  {{0, 1, 2, 3, 4, 5, 6}, {1, 2, 3, 4, 5, 6, 7}, {2, 3, 4, 5, 6, 7, 8}, {3, 4, 5, 6, 7, 8, 9}},
		8:  {{0, 1, 2, 3, 4, 5, 6, 7}, {1, 2, 3, 4, 5, 6, 7, 8}, {2, 3, 4, 5, 6, 7, 8, 9}},
		9:  {{0, 1, 2, 3, 4, 5, 6, 7, 8}, {1, 2, 3, 4, 5, 6, 7, 8, 9}},
		10: {{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for size, want := range wants {
		have := Windows(arg, size)
		assert.Equal(t, want, have)
	}
}

func TestEnumerate(t *testing.T) {
	type check struct {
		arg []int
	}
	checks := []check{}
	for i, check := range checks {
		have := Enumerate[int](check.arg)
		if len(have) > 0 {
			k0, e0 := have[0]()
			for j, pair := range have[1:] {
				k, e := pair()
				assert.Equal(t, k0, k-1, "#%d.%d\n\t index violation: have %v, want %v", i, j, k, k0)
				assert.Equal(t, e, check.arg[k0-1], "#%d.%d\n\telem violation: have %v, want %v", i, j, e, check.arg[k0+1])
				assert.Equal(t, e0, check.arg[k-1], "#%d.%d\n\telem violation: have %v, want %v", i, j, e0, check.arg[k-1])
				k0, e0 = k, e
			}
		}
	}
}

func TestWalks(t *testing.T) {
	type check struct {
		slice  []int
		length int
		want   [][]int
	}

	checks := []check{
		{slice: []int{0, 1, 2, 3}, length: 2, want: [][]int{{0, 1}, {1, 2}, {2, 3}}},
	}

	for _, check := range checks {
		have := Walks(check.length, check.slice)
		assert.Equal(t, check.want, have)
	}
}

func TestPairwise(t *testing.T) {
	require.Equal(t, [][]byte{{'A', 'B'}, {'B', 'C'}, {'C', 'D'}, {'D', 'E'}, {'E', 'F'}, {'F', 'G'}}, Pairwise([]byte("ABCDEFG")...))
	require.Equal(t, [][]rune{{'A', 'B'}, {'B', 'C'}, {'C', 'D'}, {'D', 'E'}, {'E', 'F'}, {'F', 'G'}}, Pairwise([]rune("ABCDEFG")...))
}
