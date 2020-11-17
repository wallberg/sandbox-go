package taocp

import (
	"fmt"
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

func TestLangfordPairs(t *testing.T) {

	var count int

	expected := []int{3, 1, 2, 1, 3, 2}
	count = 0
	LangfordPairs(3, nil,
		func(solution []int) {
			if !reflect.DeepEqual(solution, expected) {
				t.Errorf("Expected %v; got %v", expected, solution)
			}
			count++
		})

	if count != 1 {
		t.Errorf("Expected 1 solution; got %d", count)
	}

	testLangfordPairs(t, 7, 26)
	testLangfordPairs(t, 8, 150)
	testLangfordPairs(t, 10, 0)
	testLangfordPairs(t, 11, 17792)
}

func testLangfordPairs(t *testing.T, n int, expected int) {

	count := 0
	LangfordPairs(n, nil, func(solution []int) { count++ })

	if count != expected {
		t.Errorf("Expected 1 solution; got %d", count)
	}
}

func BenchmarkLangfordPairs(b *testing.B) {
	for _, n := range []int{6, 8, 10, 12} {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				LangfordPairs(n, nil, func([]int) {})
			}
		})
	}
}
