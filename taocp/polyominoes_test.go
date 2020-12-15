package taocp

import (
	"fmt"
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

func TestBasePlacements(t *testing.T) {

	cases := []struct {
		first      []int
		placements [][]int
	}{
		{
			[]int{65536},
			[][]int{{0}},
		},
		{
			[]int{1, 2, 3},
			[][]int{
				{0, 1, 2},
				{0, 65536, 131072},
			},
		},
		{
			[]int{0, 1, 2, 65536},
			[][]int{
				{0, 1, 2, 65536},
				{0, 1, 2, 65538},
				{0, 1, 65536, 131072},
				{0, 1, 65537, 131073},
				{0, 65536, 65537, 65538},
				{0, 65536, 131072, 131073},
				{1, 65537, 131072, 131073},
				{2, 65536, 65537, 65538},
			},
		},
	}

	for _, c := range cases {
		placements := BasePlacements(c.first)

		if !reflect.DeepEqual(placements, c.placements) {
			fmt.Println(placements)
			fmt.Println(c.placements)
			t.Errorf("placements = %v; want %v", placements, c.placements)
		}
	}
}

func TestLoadPolyominoes(t *testing.T) {

	sets, err := LoadPolyominoes()
	if err != nil {
		t.Error(err)
		return
	}

	cases := []struct {
		name  string // name of the set
		count int    // number of shapes in the set
	}{
		{"1", 1},
		{"2", 1},
		{"3", 2},
		{"4", 5},
		{"5", 12},
	}

	for _, c := range cases {
		if set, ok := sets[c.name]; !ok {
			t.Errorf("Did not find set name='%s'", c.name)
		} else {
			fmt.Println(set)
			if len(set.shapes) != c.count {
				t.Errorf("Set '%s' has %d shapes; want %d",
					set.name, len(set.shapes), c.count)
			}
		}
	}
}
