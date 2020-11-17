package taocp

import (
	"reflect"
	"testing"
)

func TestExactCover(t *testing.T) {

	items := []string{"a", "b", "c", "d", "e", "f", "g"}

	options := [][]string{
		{"c", "e"},
		{"a", "d", "g"},
		{"b", "c", "f"},
		{"a", "d", "f"},
		{"b", "g"},
		{"d", "e", "g"},
	}

	expected := [][]string{
		{"a", "d", "f"},
		{"b", "g"},
		{"c", "e"},
	}

	count := 0
	var stats Stats
	// stats.Progress = true
	// stats.Debug = true

	ExactCover(items, options, []string{}, &stats,
		func(solution [][]string) {
			if !reflect.DeepEqual(solution, expected) {
				t.Errorf("Expected %v; got %v", expected, solution)
			}
			count++
		})

	if count != 1 {
		t.Errorf("Expected 1 solution; got %d", count)
	}
}
