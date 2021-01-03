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
				Point{x: 1, y: 2}, Point{x: 4, y: 2}, Point{x: 5, y: 0},
				Point{x: 5, y: 1}, Point{x: 5, y: 2}, Point{x: 5, y: 3},
				Point{x: 6, y: 2}, Point{x: 7, y: 2},
			},
			false,
		},
		{
			"[0-2][a-c]",
			Polyomino{
				Point{x: 0, y: 10}, Point{x: 0, y: 11}, Point{x: 0, y: 12},
				Point{x: 1, y: 10}, Point{x: 1, y: 11}, Point{x: 1, y: 12},
				Point{x: 2, y: 10}, Point{x: 2, y: 11}, Point{x: 2, y: 12},
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
			Polyomino{Point{x: 1, y: 0}},
			[]Polyomino{{Point{x: 0, y: 0}}},
			true,
		},
		{
			Polyomino{Point{x: 0, y: 1}, Point{x: 0, y: 2}, Point{x: 0, y: 3}},
			[]Polyomino{
				{Point{x: 0, y: 0}, Point{x: 0, y: 1}, Point{x: 0, y: 2}},
				{Point{x: 0, y: 0}, Point{x: 1, y: 0}, Point{x: 2, y: 0}},
			},
			true,
		},
		{
			Polyomino{Point{x: 0, y: 1}, Point{x: 0, y: 2}, Point{x: 0, y: 3}},
			[]Polyomino{
				{Point{x: 0, y: 0}, Point{x: 0, y: 1}, Point{x: 0, y: 2}},
			},
			false,
		},
		{
			Polyomino{Point{x: 0, y: 0}, Point{x: 0, y: 1}, Point{x: 0, y: 2}, Point{x: 1, y: 0}},
			[]Polyomino{
				{Point{x: 0, y: 0}, Point{x: 0, y: 1}, Point{x: 0, y: 2}, Point{x: 1, y: 0}},
				{Point{x: 0, y: 0}, Point{x: 0, y: 1}, Point{x: 0, y: 2}, Point{x: 1, y: 2}},
				{Point{x: 0, y: 0}, Point{x: 0, y: 1}, Point{x: 1, y: 0}, Point{x: 2, y: 0}},
				{Point{x: 0, y: 0}, Point{x: 0, y: 1}, Point{x: 1, y: 1}, Point{x: 2, y: 1}},
				{Point{x: 0, y: 0}, Point{x: 1, y: 0}, Point{x: 1, y: 1}, Point{x: 1, y: 2}},
				{Point{x: 0, y: 0}, Point{x: 1, y: 0}, Point{x: 2, y: 0}, Point{x: 2, y: 1}},
				{Point{x: 0, y: 1}, Point{x: 1, y: 1}, Point{x: 2, y: 0}, Point{x: 2, y: 1}},
				{Point{x: 0, y: 2}, Point{x: 1, y: 0}, Point{x: 1, y: 1}, Point{x: 1, y: 2}},
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
