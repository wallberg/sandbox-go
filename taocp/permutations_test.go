package taocp

import (
	"reflect"
	"sort"
	"testing"
)

func TestNextPermutation(t *testing.T) {

	expected := [][]int{
		{1, 2, 2, 3},
		{1, 2, 3, 2},
		{1, 3, 2, 2},
		{2, 1, 2, 3},
		{2, 1, 3, 2},
		{2, 2, 1, 3},
		{2, 2, 3, 1},
		{2, 3, 1, 2},
		{2, 3, 2, 1},
		{3, 1, 2, 2},
		{3, 2, 1, 2},
		{3, 2, 2, 1},
	}

	x := []int{1, 2, 2, 3}

	solution := make([][]int, 0)
	cp := make([]int, len(x))
	copy(cp, x)
	solution = append(solution, cp)

	for i := 1; NextPermutation(sort.IntSlice(x)); i++ {
		cp := make([]int, len(x))
		copy(cp, x)
		solution = append(solution, cp)
	}

	if !reflect.DeepEqual(solution, expected) {
		t.Errorf("Expected %v; got %v", expected, solution)
	}
}
