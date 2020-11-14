package taocp

import (
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

	results := make(chan int)
	go BitPairs(v, j, results)

	i := 0
	for result := range results {
		if result != expected[i] {
			t.Errorf("For case v=%d and j=%d, expected %d for i=%d; got %d",
				v, j, expected[i], i, result)
		}
		i++
	}

	if i != len(expected) {
		t.Errorf("For case v=%d and j=%d, expected %d results; got %d",
			v, j, len(expected), i)
	}

}
