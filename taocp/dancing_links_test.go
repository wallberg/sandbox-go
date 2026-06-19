package taocp

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

var (
	xcItems = []string{"a", "b", "c", "d", "e", "f", "g"}

	xcOptions = [][]string{
		{"c", "e"},
		{"a", "d", "g"},
		{"b", "c", "f"},
		{"a", "d", "f"},
		{"b", "g"},
		{"d", "e", "g"},
	}

	xcExpected = [][]string{
		{"a", "d", "f"},
		{"b", "g"},
		{"c", "e"},
	}
)

func TestExactCoverYaml(t *testing.T) {

	// Build YAML struct
	xcYaml := NewExactCoverYaml(xccItems, xccSItems, xccOptions)

	// Serialize to YAML
	data, err := yaml.Marshal(xcYaml)
	if err != nil {
		t.Errorf("Error serializing ExactCoverYaml: %v", err)
		return
	}

	// Deserialize from YAML
	var xcYaml2 ExactCoverYaml
	err = yaml.Unmarshal([]byte(data), &xcYaml2)
	if err != nil {
		t.Errorf("Error deserializing ExactCoverYaml: %v", err)
	}

	// Test the round trip
	if !reflect.DeepEqual(*xcYaml, xcYaml2) {
		t.Errorf("Got back %v; want %v", xcYaml2, *xcYaml)
	}
}

func TestExactCover(t *testing.T) {

	count := 0
	var stats ExactCoverStats

	for solution := range ExactCover(xcItems, xcOptions, []string{}, &stats) {
		if !reflect.DeepEqual(solution, xcExpected) {
			t.Errorf("Expected %v; got %v", xcExpected, solution)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 solution; got %d", count)
	}

	if stats.Solutions != 1 {
		t.Errorf("Expected 1 stats.Solution; got %d", stats.Solutions)
	}
}

func TestLangfordPairs(t *testing.T) {

	var count int

	expected := []int{3, 1, 2, 1, 3, 2}
	count = 0
	for solution := range LangfordPairs(3, nil) {
		if !reflect.DeepEqual(solution, expected) {
			t.Errorf("Expected %v; got %v", expected, solution)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 solution; got %d", count)
	}

	count = 0
	for range LangfordPairs(7, nil) {
		count++
		break // halt after the first solution
	}

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
	for range LangfordPairs(n, nil) {
		count++
	}

	if count != expected {
		t.Errorf("Expected 1 solution; got %d", count)
	}
}

func BenchmarkLangfordPairs(b *testing.B) {
	for _, n := range []int{6, 8, 10, 12} {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				for range LangfordPairs(n, nil) {
				}
			}
		})
	}
}

func TestNQueens(t *testing.T) {

	expected0 := []string{
		"r1", "c2",
		"r2", "c4",
		"r3", "c1",
		"r4", "c3",
	}

	expected1 := []string{
		"r1", "c3",
		"r2", "c1",
		"r3", "c4",
		"r4", "c2",
	}

	count := 0
	for solution := range NQueens(4, nil) {
		if count == 0 && !reflect.DeepEqual(solution, expected0) {
			t.Errorf("Expected %v; got %v", expected0, solution)
		}
		if count == 1 && !reflect.DeepEqual(solution, expected1) {
			t.Errorf("Expected %v; got %v", expected1, solution)
		}
		count++
	}

	if count != 2 {
		t.Errorf("Expected 1 solution; got %d", count)
	}

	testNQueens(t, 8, 92)
	testNQueens(t, 11, 2680)
}

func testNQueens(t *testing.T, n int, expected int) {

	count := 0
	for range NQueens(n, nil) {
		count++
	}

	if count != expected {
		t.Errorf("Expected 1 solution; got %d", count)
	}
}

func BenchmarkNQueens(b *testing.B) {
	for _, n := range []int{8, 11, 13} {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				for range NQueens(n, nil) {
				}
			}
		})
	}
}

var (
	input1 = [9][9]int{
		{0, 8, 3, 9, 2, 1, 6, 5, 7},
		{9, 6, 7, 3, 4, 5, 8, 2, 1},
		{2, 5, 1, 8, 7, 6, 4, 9, 3},
		{5, 4, 8, 1, 3, 2, 9, 7, 0},
		{7, 2, 9, 5, 6, 4, 1, 3, 8},
		{1, 3, 6, 7, 9, 8, 2, 4, 5},
		{3, 7, 2, 6, 8, 9, 5, 1, 4},
		{8, 1, 4, 2, 5, 3, 7, 6, 9},
		{6, 9, 5, 4, 1, 7, 3, 8, 0},
	}

	expected1 = [9][9]int{
		{4, 8, 3, 9, 2, 1, 6, 5, 7},
		{9, 6, 7, 3, 4, 5, 8, 2, 1},
		{2, 5, 1, 8, 7, 6, 4, 9, 3},
		{5, 4, 8, 1, 3, 2, 9, 7, 6},
		{7, 2, 9, 5, 6, 4, 1, 3, 8},
		{1, 3, 6, 7, 9, 8, 2, 4, 5},
		{3, 7, 2, 6, 8, 9, 5, 1, 4},
		{8, 1, 4, 2, 5, 3, 7, 6, 9},
		{6, 9, 5, 4, 1, 7, 3, 8, 2},
	}

	input2 = [9][9]int{
		{3, 0, 0, 2, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 7, 0, 0, 0},
		{7, 0, 6, 0, 3, 0, 5, 0, 0},
		{0, 7, 0, 0, 0, 9, 0, 8, 0},
		{9, 0, 0, 0, 2, 0, 0, 0, 4},
		{0, 1, 0, 8, 0, 0, 0, 5, 0},
		{0, 0, 9, 0, 4, 0, 3, 0, 1},
		{0, 0, 0, 7, 0, 2, 0, 0, 0},
		{0, 0, 0, 0, 0, 8, 0, 0, 6},
	}

	expected2 = [9][9]int{
		{3, 5, 1, 2, 8, 6, 4, 9, 7},
		{4, 9, 2, 1, 5, 7, 6, 3, 8},
		{7, 8, 6, 9, 3, 4, 5, 1, 2},
		{2, 7, 5, 4, 6, 9, 1, 8, 3},
		{9, 3, 8, 5, 2, 1, 7, 6, 4},
		{6, 1, 4, 8, 7, 3, 2, 5, 9},
		{8, 2, 9, 6, 4, 5, 3, 7, 1},
		{1, 6, 3, 7, 9, 2, 8, 4, 5},
		{5, 4, 7, 3, 1, 8, 9, 2, 6},
	}

	// 29a
	input3 = [9][9]int{
		{0, 0, 3, 0, 1, 0, 0, 0, 0},
		{4, 1, 5, 0, 0, 0, 0, 9, 0},
		{2, 0, 6, 5, 0, 0, 3, 0, 0},
		{5, 0, 0, 0, 8, 0, 0, 0, 9},
		{0, 7, 0, 9, 0, 0, 0, 3, 2},
		{0, 3, 8, 0, 0, 4, 0, 6, 0},
		{0, 0, 0, 2, 6, 0, 4, 0, 3},
		{0, 0, 0, 3, 0, 0, 0, 0, 8},
		{3, 2, 0, 0, 0, 7, 9, 5, 0},
	}

	expected3 = [9][9]int{
		{7, 9, 3, 4, 1, 2, 6, 8, 5},
		{4, 1, 5, 6, 3, 8, 2, 9, 7},
		{2, 8, 6, 5, 7, 9, 3, 1, 4},
		{5, 6, 2, 1, 8, 3, 7, 4, 9},
		{1, 7, 4, 9, 5, 6, 8, 3, 2},
		{9, 3, 8, 7, 2, 4, 5, 6, 1},
		{8, 5, 9, 2, 6, 1, 4, 7, 3},
		{6, 4, 7, 3, 9, 5, 1, 2, 8},
		{3, 2, 1, 8, 4, 7, 9, 5, 6},
	}

	// 29b
	input4 = [9][9]int{
		{0, 0, 0, 0, 0, 0, 3, 0, 0},
		{1, 0, 0, 4, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 0, 5},
		{9, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 2, 6, 0, 0},
		{0, 0, 0, 0, 5, 3, 0, 0, 0},
		{0, 5, 0, 8, 0, 0, 0, 0, 0},
		{0, 0, 0, 9, 0, 0, 0, 7, 0},
		{0, 8, 3, 0, 0, 0, 0, 4, 0},
	}

	expected4 = [9][9]int{
		{5, 9, 7, 2, 1, 8, 3, 6, 4},
		{1, 3, 2, 4, 6, 5, 8, 9, 7},
		{8, 6, 4, 3, 7, 9, 1, 2, 5},
		{9, 1, 5, 6, 8, 4, 7, 3, 2},
		{3, 4, 8, 7, 9, 2, 6, 5, 1},
		{2, 7, 6, 1, 5, 3, 4, 8, 9},
		{6, 5, 9, 8, 4, 7, 2, 1, 3},
		{4, 2, 1, 9, 3, 6, 5, 7, 8},
		{7, 8, 3, 5, 2, 1, 9, 4, 6},
	}
	// Euler Problem 96 - Grid 49
	input5 = [9][9]int{
		{0, 0, 0, 0, 0, 3, 0, 1, 7},
		{0, 1, 5, 0, 0, 9, 0, 0, 8},
		{0, 6, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 7, 0, 0, 0},
		{0, 0, 9, 0, 0, 0, 2, 0, 0},
		{0, 0, 0, 5, 0, 0, 0, 0, 4},
		{0, 0, 0, 0, 0, 0, 0, 2, 0},
		{5, 0, 0, 6, 0, 0, 3, 4, 0},
		{3, 4, 0, 2, 0, 0, 0, 0, 0},
	}
	// Euler Problem 96 - Grid 50
	input6 = [9][9]int{
		{3, 0, 0, 2, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 7, 0, 0, 0},
		{7, 0, 6, 0, 3, 0, 5, 0, 0},
		{0, 7, 0, 0, 0, 9, 0, 8, 0},
		{9, 0, 0, 0, 2, 0, 0, 0, 4},
		{0, 1, 0, 8, 0, 0, 0, 5, 0},
		{0, 0, 9, 0, 4, 0, 3, 0, 1},
		{0, 0, 0, 7, 0, 2, 0, 0, 0},
		{0, 0, 0, 0, 0, 8, 0, 0, 6},
	}
)

func TestSudoku(t *testing.T) {

	testSudoku(t, input1, expected1)
	testSudoku(t, input2, expected2)
	testSudoku(t, input3, expected3)
	testSudoku(t, input4, expected4)
}

func testSudoku(t *testing.T, input [9][9]int, expected [9][9]int) {
	count := 0
	for solution := range Sudoku(input, nil) {
		if !reflect.DeepEqual(solution, expected) {
			t.Errorf("Expected %v; got %v", expected, solution)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 solution; got %d", count)
	}
}

func BenchmarkSudoku(b *testing.B) {
	b.Run("input1", func(b *testing.B) {
		for repeat := 0; repeat < b.N; repeat++ {
			for range Sudoku(input1, nil) {
			}
		}
	})
	b.Run("input2", func(b *testing.B) {
		for repeat := 0; repeat < b.N; repeat++ {
			for range Sudoku(input2, nil) {
			}
		}
	})
	b.Run("input3", func(b *testing.B) {
		for repeat := 0; repeat < b.N; repeat++ {
			for range Sudoku(input3, nil) {
			}
		}
	})
	b.Run("input4", func(b *testing.B) {
		for repeat := 0; repeat < b.N; repeat++ {
			for range Sudoku(input4, nil) {
			}
		}
	})
	b.Run("input5", func(b *testing.B) {
		for repeat := 0; repeat < b.N; repeat++ {
			for range Sudoku(input5, nil) {
			}
		}
	})
	b.Run("input6", func(b *testing.B) {
		for repeat := 0; repeat < b.N; repeat++ {
			for range Sudoku(input6, nil) {
			}
		}
	})
}
