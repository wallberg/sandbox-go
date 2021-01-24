package taocp

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParsePlacementPairs(t *testing.T) {

	cases := []struct {
		s     string    // string to parse
		pairs Polyomino // sorted pairs
		err   bool      // true if error is expected
	}{
		{
			"[14-7]2 5[0-3]",
			Polyomino{
				Point{X: 1, Y: 2}, Point{X: 4, Y: 2}, Point{X: 5, Y: 0},
				Point{X: 5, Y: 1}, Point{X: 5, Y: 2}, Point{X: 5, Y: 3},
				Point{X: 6, Y: 2}, Point{X: 7, Y: 2},
			},
			false,
		},
		{
			"[0-2][a-c]",
			Polyomino{
				Point{X: 0, Y: 10}, Point{X: 0, Y: 11}, Point{X: 0, Y: 12},
				Point{X: 1, Y: 10}, Point{X: 1, Y: 11}, Point{X: 1, Y: 12},
				Point{X: 2, Y: 10}, Point{X: 2, Y: 11}, Point{X: 2, Y: 12},
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
		first      Polyomino
		placements []Polyomino
		transform  bool
	}{
		{
			Polyomino{Point{X: 1, Y: 0}},
			[]Polyomino{{Point{X: 0, Y: 0}}},
			true,
		},
		{
			Polyomino{Point{X: 0, Y: 1}, Point{X: 0, Y: 2}, Point{X: 0, Y: 3}},
			[]Polyomino{
				{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, Point{X: 0, Y: 2}},
				{Point{X: 0, Y: 0}, Point{X: 1, Y: 0}, Point{X: 2, Y: 0}},
			},
			true,
		},
		{
			Polyomino{Point{X: 0, Y: 1}, Point{X: 0, Y: 2}, Point{X: 0, Y: 3}},
			[]Polyomino{
				{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, Point{X: 0, Y: 2}},
			},
			false,
		},
		{
			Polyomino{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, Point{X: 0, Y: 2}, Point{X: 1, Y: 0}},
			[]Polyomino{
				{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, Point{X: 0, Y: 2}, Point{X: 1, Y: 0}},
				{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, Point{X: 0, Y: 2}, Point{X: 1, Y: 2}},
				{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, Point{X: 1, Y: 0}, Point{X: 2, Y: 0}},
				{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, Point{X: 1, Y: 1}, Point{X: 2, Y: 1}},
				{Point{X: 0, Y: 0}, Point{X: 1, Y: 0}, Point{X: 1, Y: 1}, Point{X: 1, Y: 2}},
				{Point{X: 0, Y: 0}, Point{X: 1, Y: 0}, Point{X: 2, Y: 0}, Point{X: 2, Y: 1}},
				{Point{X: 0, Y: 1}, Point{X: 1, Y: 1}, Point{X: 2, Y: 0}, Point{X: 2, Y: 1}},
				{Point{X: 0, Y: 2}, Point{X: 1, Y: 0}, Point{X: 1, Y: 1}, Point{X: 1, Y: 2}},
			},
			true,
		},
	}

	for _, c := range cases {
		placements := BasePlacements(c.first, c.transform)

		if !reflect.DeepEqual(placements, c.placements) {
			fmt.Println(placements)
			fmt.Println(c.placements)
			t.Errorf("placements = %v; want %v", placements, c.placements)
		}
	}
}

func TestLoadPolyominoes(t *testing.T) {

	sets := LoadPolyominoes()

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
		if set, ok := sets.PieceSets[c.name]; !ok {
			t.Errorf("Did not find set name='%s'", c.name)
		} else {
			if len(set) != c.count {
				t.Errorf("Set '%s' has %d shapes; want %d",
					c.name, len(set), c.count)
			}
		}
	}
}

func TestPolyominoes(t *testing.T) {
	cases := []struct {
		shapes []string // names of the piece shapes
		board  string   // name of the board shape
		count  int      // number of expected results
	}{
		{[]string{"5"}, "3x20", 8},
		{[]string{"1"}, "1x1", 1},
		{[]string{"2"}, "1x1", 0},
		{[]string{"1", "2"}, "2x2", 0},
		{[]string{"1", "2"}, "2x2-1", 2},
		{[]string{"2", "3"}, "2x3", 0},
		{[]string{"1", "2", "3"}, "3x3", 48},
		{[]string{"2", "3"}, "3x3-1", 4},
		{[]string{"1", "2", "3", "4"}, "5x6-1", 100593}, // tautology
		// {[]string{"4", "5"}, "8x8", 0},                   // too slow
		// {[]string{"1", "2", "3", "4", "5"}, "5x18-1", 0}, // too slow
	}

	for _, c := range cases {
		items, options, sitems := Polyominoes(c.shapes, c.board)

		// if true /*c.board == "8x8"*/ {
		// 	fmt.Println(items)
		// 	for _, option := range options {
		// 		fmt.Println(option)
		// 	}
		// }

		// Generate solutions
		count := 0
		if len(options) > 0 {
			stats := &ExactCoverStats{Debug: false, Progress: true, Delta: 10000000}
			ExactCover(items, options, sitems, stats,
				func(solution [][]string) bool {
					count++
					return true
				})
		}

		if count != c.count {
			t.Errorf("Found %d solutions for shape sets=%v, board=%s; want %d", count, c.shapes, c.board, c.count)
		}
	}
}

func TestPolyominoPacking(t *testing.T) {
	cases := []struct {
		x                int
		y                int
		n                int
		includeStraight  bool
		includeNonConvex bool
		count            int
	}{
		{2, 3, 2, true, true, 7},
		{2, 3, 2, false, true, 0},
		{3, 3, 5, true, false, 41},
		{3, 3, 3, false, true, 16},
		{4, 4, 4, false, true, 105},
		{5, 5, 5, false, true, 561},
		{6, 6, 6, false, true, 2804},
		{7, 7, 7, false, true, 13602},
	}

	for _, c := range cases {
		pos := PolyominoPacking(c.x, c.y, c.n, c.includeStraight,
			c.includeNonConvex)

		count := len(pos)
		if count != c.count {
			t.Errorf("Got %d shapes for PolyominoPacking(%d,%d,%d,%t,%t); want %d",
				count, c.x, c.y, c.n, c.includeStraight, c.includeNonConvex,
				c.count)
		}
	}
}
