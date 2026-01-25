package taocp

import (
	"log"
	"reflect"
	"testing"
)

func TestStrictPartitions(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		n    int
		min  int
		max  int
		want [][]int
	}{
		{
			0,
			0,
			0,
			[][]int{},
		},
		{
			1,
			0,
			0,
			[][]int{},
		},
		{
			1,
			1,
			1,
			[][]int{{1}},
		},
		{
			2,
			1,
			2,
			[][]int{
				{2},
			},
		},
		{
			3,
			1,
			3,
			[][]int{
				{3},
				{2, 1},
			},
		},
		{
			8,
			1,
			3,
			[][]int{
				{8},
				{7, 1},
				{6, 2},
				{5, 3},
				{5, 2, 1},
				{4, 3, 1},
			},
		},
		{
			8,
			2,
			3,
			[][]int{
				{7, 1},
				{6, 2},
				{5, 3},
				{5, 2, 1},
				{4, 3, 1},
			},
		},
		{
			8,
			1,
			2,
			[][]int{
				{8},
				{7, 1},
				{6, 2},
				{5, 3},
			},
		},
		{
			12,
			1,
			10,
			[][]int{
				{12},
				{11, 1},
				{10, 2},
				{9, 3},
				{9, 2, 1},
				{8, 4},
				{8, 3, 1},
				{7, 5},
				{7, 4, 1},
				{7, 3, 2},
				{6, 5, 1},
				{6, 4, 2},
				{6, 3, 2, 1},
				{5, 4, 3},
				{5, 4, 2, 1},
			},
		},
		{
			12,
			1,
			1,
			[][]int{
				{12},
			},
		},
		{
			12,
			2,
			3,
			[][]int{
				{11, 1},
				{10, 2},
				{9, 3},
				{9, 2, 1},
				{8, 4},
				{8, 3, 1},
				{7, 5},
				{7, 4, 1},
				{7, 3, 2},
				{6, 5, 1},
				{6, 4, 2},
				{5, 4, 3},
			},
		},
		{
			12,
			4,
			4,
			[][]int{
				{6, 3, 2, 1},
				{5, 4, 2, 1},
			},
		},
	}

	for i, c := range cases {

		got := make([][]int, 0)
		for partition := range StrictPartitions(c.n, c.min, c.max) {
			got = append(got, partition)
		}

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("For case #%d, n=%d, min=%d, max=%d, got solutions %v; want %v", i, c.n, c.min, c.max, got, c.want)
		}
	}

	for partition := range StrictPartitions(38, 3, 5) {
		log.Printf("%v", partition)
	}
	t.Errorf("Done")
}
