// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var (
	ints             = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
	float64s         = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8, 74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3}
	float64sWithNaNs = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.NaN(), math.NaN(), math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8}
	strs             = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}
)

func TestSortIntSlice(t *testing.T) {
	data := ints[:]
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortFuncIntSlice(t *testing.T) {
	data := ints[:]
	SortFunc(func(a, b int) bool { return a < b }, data)
	if !IsSorted(data) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortFloat64Slice(t *testing.T) {
	data := float64s[:]
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestSortFloat64SliceWithNaNs(t *testing.T) {
	data := float64sWithNaNs[:]
	input := make([]float64, len(float64sWithNaNs))
	for i := range input {
		input[i] = float64sWithNaNs[i]
	}
	// Make sure Sort doesn't panic when the slice contains NaNs.
	Sort(data)
	// Check whether the result is a permutation of the input.
	sort.Float64s(data)
	sort.Float64s(input)
	for i, v := range input {
		if data[i] != v && !(math.IsNaN(data[i]) && math.IsNaN(v)) {
			t.Fatalf("the result is not a permutation of the input\ngot %v\nwant %v", data, input)
		}
	}
}

func TestSortStringSlice(t *testing.T) {
	data := strs[:]
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sorted %v", strs)
		t.Errorf("   got %v", data)
	}
}

func TestSortLarge_Random(t *testing.T) {
	n := 1000000
	if testing.Short() {
		n /= 100
	}
	data := make([]int, n)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(100)
	}
	if IsSorted(data) {
		t.Fatalf("terrible rand.rand")
	}
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sort didn't sort - 1M ints")
	}
}

type intPair struct {
	a, b int
}

type intPairs []intPair

// Pairs compare on a only.
func intPairLess(x, y intPair) bool {
	return x.a < y.a
}

// Record initial order in B.
func (d intPairs) initB() {
	for i := range d {
		d[i].b = i
	}
}

// InOrder checks if a-equal elements were not reordered.
func (d intPairs) inOrder() bool {
	lastA, lastB := -1, 0
	for i := 0; i < len(d); i++ {
		if lastA != d[i].a {
			lastA = d[i].a
			lastB = d[i].b
			continue
		}
		if d[i].b <= lastB {
			return false
		}
		lastB = d[i].b
	}
	return true
}

func TestStability(t *testing.T) {
	n, m := 100000, 1000
	if testing.Short() {
		n, m = 1000, 100
	}
	data := make(intPairs, n)

	// random distribution
	for i := 0; i < len(data); i++ {
		data[i].a = rand.Intn(m)
	}
	if IsSortedFunc(intPairLess, data) {
		t.Fatalf("terrible rand.rand")
	}
	data.initB()
	SortStableFunc(intPairLess, data)
	if !IsSortedFunc(intPairLess, data) {
		t.Errorf("Stable didn't sort %d ints", n)
	}
	if !data.inOrder() {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}

	// already sorted
	data.initB()
	SortStableFunc(intPairLess, data)
	if !IsSortedFunc(intPairLess, data) {
		t.Errorf("Stable shuffled sorted %d ints (order)", n)
	}
	if !data.inOrder() {
		t.Errorf("Stable shuffled sorted %d ints (stability)", n)
	}

	// sorted reversed
	for i := 0; i < len(data); i++ {
		data[i].a = len(data) - i
	}
	data.initB()
	SortStableFunc(intPairLess, data)
	if !IsSortedFunc(intPairLess, data) {
		t.Errorf("Stable didn't sort %d ints", n)
	}
	if !data.inOrder() {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}
}

func TestBinarySearch(t *testing.T) {
	str1 := []string{"foo"}
	str2 := []string{"ab", "ca"}
	str3 := []string{"mo", "qo", "vo"}
	str4 := []string{"ab", "ad", "ca", "xy"}

	// slice with repeating elements
	strRepeats := []string{"ba", "ca", "da", "da", "da", "ka", "ma", "ma", "ta"}

	// slice with all element equal
	strSame := []string{"xx", "xx", "xx"}

	tests := []struct {
		target    string
		data      []string
		wantPos   int
		wantFound bool
	}{
		{data: []string{}, target: "foo", wantPos: 0, wantFound: false},
		{data: []string{}, target: "", wantPos: 0, wantFound: false},

		{data: str1, target: "foo", wantPos: 0, wantFound: true},
		{data: str1, target: "bar", wantPos: 0, wantFound: false},
		{data: str1, target: "zx", wantPos: 1, wantFound: false},

		{data: str2, target: "aa", wantPos: 0, wantFound: false},
		{data: str2, target: "ab", wantPos: 0, wantFound: true},
		{data: str2, target: "ad", wantPos: 1, wantFound: false},
		{data: str2, target: "ca", wantPos: 1, wantFound: true},
		{data: str2, target: "ra", wantPos: 2, wantFound: false},

		{data: str3, target: "bb", wantPos: 0, wantFound: false},
		{data: str3, target: "mo", wantPos: 0, wantFound: true},
		{data: str3, target: "nb", wantPos: 1, wantFound: false},
		{data: str3, target: "qo", wantPos: 1, wantFound: true},
		{data: str3, target: "tr", wantPos: 2, wantFound: false},
		{data: str3, target: "vo", wantPos: 2, wantFound: true},
		{data: str3, target: "xr", wantPos: 3, wantFound: false},

		{data: str4, target: "aa", wantPos: 0, wantFound: false},
		{data: str4, target: "ab", wantPos: 0, wantFound: true},
		{data: str4, target: "ac", wantPos: 1, wantFound: false},
		{data: str4, target: "ad", wantPos: 1, wantFound: true},
		{data: str4, target: "ax", wantPos: 2, wantFound: false},
		{data: str4, target: "ca", wantPos: 2, wantFound: true},
		{data: str4, target: "cc", wantPos: 3, wantFound: false},
		{data: str4, target: "dd", wantPos: 3, wantFound: false},
		{data: str4, target: "xy", wantPos: 3, wantFound: true},
		{data: str4, target: "zz", wantPos: 4, wantFound: false},

		{data: strRepeats, target: "da", wantPos: 2, wantFound: true},
		{data: strRepeats, target: "db", wantPos: 5, wantFound: false},
		{data: strRepeats, target: "ma", wantPos: 6, wantFound: true},
		{data: strRepeats, target: "mb", wantPos: 8, wantFound: false},

		{data: strSame, target: "xx", wantPos: 0, wantFound: true},
		{data: strSame, target: "ab", wantPos: 0, wantFound: false},
		{data: strSame, target: "zz", wantPos: 3, wantFound: false},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			{
				pos, found := BinarySearch(tt.target, tt.data)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearch got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}

			{
				pos, found := BinarySearchFunc(strings.Compare, tt.target, tt.data)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearchFunc got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}
		})
	}
}

func TestBinarySearchInts(t *testing.T) {
	data := []int{20, 30, 40, 50, 60, 70, 80, 90}
	tests := []struct {
		target    int
		wantPos   int
		wantFound bool
	}{
		{20, 0, true},
		{23, 1, false},
		{43, 3, false},
		{80, 6, true},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.target), func(t *testing.T) {
			{
				pos, found := BinarySearch(tt.target, data)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearch got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}

			{
				cmp := func(a, b int) int {
					return a - b
				}
				pos, found := BinarySearchFunc(cmp, tt.target, data)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearchFunc got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}
		})
	}
}
