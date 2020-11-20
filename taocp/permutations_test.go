package taocp

import (
	"reflect"
	"testing"
)

func TestPermutations(t *testing.T) {

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
	Permutations(x, func() bool {
		permCopy := make([]int, len(x))
		copy(permCopy, x)
		solution = append(solution, permCopy)
		return true
	})

	if !reflect.DeepEqual(solution, expected) {
		t.Errorf("Expected %v; got %v", expected, solution)
	}

	x = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	count := 0
	Permutations(x, func() bool {
		count++
		return true
	})

	if count != 362880 {
		t.Errorf("Expected 362880 permutations; got %d", count)
	}

}
