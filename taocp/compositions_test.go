package taocp

import (
	"reflect"
	"testing"
)

func TestCompositions(t *testing.T) {

	cases := []struct {
		n        int     // input n value
		expected [][]int // output values
	}{
		{2, [][]int{
			{2},
			{1, 1},
		}},
		{4, [][]int{
			{4},
			{3, 1},
			{2, 1, 1},
			{2, 2},
			{1, 1, 2},
			{1, 1, 1, 1},
			{1, 2, 1},
			{1, 3},
		}},
	}

	for i, c := range cases {

		got := make([][]int, 0)
		Compositions(c.n, func(x []int) bool {
			y := make([]int, len(x))
			copy(y, x)
			got = append(got, y)
			return true
		})

		if !reflect.DeepEqual(got, c.expected) {
			t.Errorf("Expected %v for case %d; got %v", c.expected, i, got)
		}
	}
}
