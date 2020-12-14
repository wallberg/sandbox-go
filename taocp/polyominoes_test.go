package taocp

import (
	"reflect"
	"testing"
)

func TestParsePlacementPairs(t *testing.T) {

	cases := []struct {
		s     string // string to parse
		pairs []int  // sorted pairs
		err   bool   // true if error is expected
	}{
		{
			"[14-7]2 5[0-3]",
			[]int{65538, 262146, 327680, 327681, 327682, 327683, 393218,
				458754},
			false,
		},
		{
			"[0-2][a-c]",
			[]int{10, 11, 12,
				65546, 65547, 65548,
				131082, 131083, 131084,
			},
			false,
		},
		{
			"",
			nil,
			true,
		},
		{
			"x",
			nil,
			true,
		},
	}

	for _, c := range cases {
		pairs, err := ParsePlacementPairs(c.s)

		if (err != nil) != c.err {
			t.Errorf("(err != nil) = %v; want %v", err != nil, c.err)
		}

		if !reflect.DeepEqual(c.pairs, pairs) {
			t.Errorf("pairs = %v; want %v", pairs, c.pairs)
		}
	}
}
