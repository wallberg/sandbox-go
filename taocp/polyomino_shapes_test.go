package taocp

import (
	"reflect"
	"testing"
)

func TestPolyominoShapes(t *testing.T) {

	cases := []struct {
		n      int     // size
		count  int     // number of shapes generated
		shapes [][]int // generated shapes
	}{
		{
			1,
			1,
			[][]int{{0}},
		},
		{
			2,
			1,
			[][]int{{0, 1}},
		},
		{
			3,
			2,
			[][]int{{0, 1, 2}, {0, 1, 65536}},
		},
		{
			4,
			5,
			nil,
		},
		{
			5,
			12,
			nil,
		},
		{
			6,
			35,
			nil,
		},
		{
			7,
			108,
			nil,
		},
		{
			8,
			369,
			nil,
		},
		{
			9,
			1285,
			nil,
		},
		{
			10,
			4655,
			nil,
		},
		// { // too slow
		// 	11,
		// 	17073,
		// 	nil,
		// },
	}

	for _, c := range cases {
		shapes := PolyominoShapes(c.n)

		if count := len(shapes); count != c.count {
			t.Errorf("for n=%d, got number of shapes %d; want %d", c.n, count, c.count)
		}

		if c.shapes != nil && !reflect.DeepEqual(shapes, c.shapes) {
			t.Errorf("for n=%d, got shapes %v; want %v", c.n, shapes, c.shapes)
		}
	}
}
