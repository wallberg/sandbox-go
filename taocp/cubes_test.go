package taocp

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestRotations(t *testing.T) {

	cases := []struct {
		cube Cube
	}{
		// Test cases go here
		{
			"abcdef",
		},
	}

	for i, c := range cases {

		// Test goes here
		cubes := c.cube.Rotations()
		got := mapset.NewSet[Cube]()
		for _, cube := range cubes {
			got.Add(cube)
			fmt.Println(cube)
		}

		if got.Cardinality() != 24 {
			t.Errorf("Got %d rotations for case #%d, cube=%v; want 24", got.Cardinality(), i, c.cube)
		}

		if cubes[23] == (Cube)("abdcfe") {
			t.Errorf("Got lost rotation %v for case #%d, cube=%v; want abdcfe", cubes[23], i, c.cube)
		}
	}
}

func TestCubes(t *testing.T) {

	cases := []struct {
		cubes []string
	}{
		// Test cases go here
		{
			[]string{
				"abcdef", "abcdfe", "acdefb", "acdebf", "adefbc", "adefcb", "aefbcd", "aefbdc", "afbcde", "afbced",
				"abcefd", "abcedf", "acdfbe", "acdfeb", "adebcf", "adebfc", "aefcdb", "aefcbd", "afbdec", "afbdce",
				"abcfde", "abcfed", "acdbef", "acdbfe", "adecfb", "adecbf", "aefdbc", "aefdcb", "afbecd", "afbedc",
			},
		},
	}

	for i, c := range cases {

		// Test goes here
		cubes := Cubes()
		got := mapset.NewSet[Cube]()
		for _, cube := range cubes {
			got.Add(cube)
		}

		expected := mapset.NewSet[Cube]()
		for _, cube := range c.cubes {
			expected.Add((Cube)(cube))
		}

		if got.Cardinality() != expected.Cardinality() {
			t.Errorf("Got %d solutions for case #%d; want %d", got.Cardinality(), i, expected.Cardinality())
		}

		// Check that each cube from Cubes() contains exactly one of its rotations
		// in the expected set
		for _, cube := range cubes {
			matches := 0
			for _, rotation := range cube.Rotations() {
				if got.Contains(rotation) {
					matches += 1
				}
			}
			if matches != 1 {
				t.Errorf("Got %d matched rotations for case #%d, cube=%v; want 1", matches, i, cube)
			}
		}
	}
}

func TestBrick(t *testing.T) {

	cases := []struct {
		l, m, n int
		count   int
	}{
		{1, 1, 1, 720},
	}

	for i, c := range cases {

		stats := &ExactCoverStats{
			// Progress: true,
			// Delta:    50000000,
			// Debug:    true,
			// Verbosity:    2,
			// SuppressDump: true,
		}

		xccOptions := &XCCOptions{}

		items, options, sitems := Brick(c.l, c.m, c.n)

		count := 0
		XCC(items, options, sitems, stats, xccOptions,
			func(solution [][]string) bool {
				count++
				return true
			})

		if count != c.count {
			t.Errorf("Got %d solutions for case #%d, l=%d, m=%d, n=%d; want %d",
				count, i, c.l, c.m, c.n, c.count)
		}
	}
}
