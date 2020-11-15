package taocp

import (
	"fmt"
	"testing"
)

// "random" function 7.1.1-(22)
var F22 = []int{0, 1, 4, 7, 12, 13, 14, 15}
var F22N = 4
var F22Subcubes = []int{
	1, 0, //  000*
	3, 12, // 11**
	4, 0, //  0*00
	8, 4, //  *100
	8, 7} //  *111

func TestBitPairs(t *testing.T) {

	testBitPairs(t, F22, 0, []int{0, 1, 4, 5, 6, 7})
	testBitPairs(t, F22, 1, []int{4, 6, 5, 7})
	testBitPairs(t, F22, 2, []int{0, 2})
	testBitPairs(t, F22, 3, []int{2, 4, 3, 7})
}

func testBitPairs(t *testing.T, v []int, j int, expected []int) {

	i := 0
	BitPairs(v, j, func(k int, kp int) {
		for _, result := range []int{k, kp} {
			if result != expected[i] {
				t.Errorf("For case v=%d and j=%d, expected %d for i=%d; got %d",
					v, j, expected[i], i, result)
			}
			i++
		}
	})

	if i != len(expected) {
		t.Errorf("For case v=%d and j=%d, expected %d results; got %d",
			v, j, len(expected), i)
	}
}

func TestMaximalSubcubes(t *testing.T) {

	var v []int
	var expected []int

	v = []int{1, 2, 4, 8}
	expected = []int{0, 1, 0, 2, 0, 4, 0, 8}
	testMaximalSubcubes(t, 4, v, expected)

	testMaximalSubcubes(t, F22N, F22, F22Subcubes)

	v = make([]int, 32)
	for i := range v {
		v[i] = i
	}
	expected = []int{31, 0}
	testMaximalSubcubes(t, 5, v, expected)
}

func testMaximalSubcubes(t *testing.T, n int, v []int, expected []int) {

	i := 0
	MaximalSubcubes(n, v, func(a int, b int) {
		for _, result := range []int{a, b} {
			if result != expected[i] {
				t.Errorf("For case v=%d and n=%d, expected %d for i=%d; got %d",
					v, n, expected[i], i, result)
			}
			i++
		}
	})

	if i != len(expected) {
		t.Errorf("For case v=%d and n=%d, expected %d results; got %d",
			v, n, len(expected), i)
	}
}

func BenchmarkMaximalSubcubes(b *testing.B) {
	for n := 4; n < 15; n += 2 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				max := (1 << n)
				v := make([]int, max)
				for i := range v {
					v[i] = i
				}
				MaximalSubcubes(n, v, func(a int, b int) {})
			}
		})
	}
}
