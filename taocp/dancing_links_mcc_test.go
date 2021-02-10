package taocp

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func sortSolutions(solutions [][][]string) {
	for _, solution := range solutions {
		for _, option := range solution {
			// Sort the items in the option
			sort.Strings(option)
		}

		// Sort the options in the solution
		sort.Slice(solution, func(i, j int) bool {
			return strings.Join(solution[i], " ") < strings.Join(solution[j], " ")
		})
	}

	// Sort the solutions
	sort.Slice(solutions, func(i, j int) bool {
		k := 0
		for ; k < len(solutions[i]) && k < len(solutions[j]); k++ {
			iOption := strings.Join(solutions[i][k], " ")
			jOption := strings.Join(solutions[j][k], " ")

			if iOption < jOption {
				return true
			} else if iOption > jOption {
				return false
			}
		}

		if k < len(solutions[j]) {
			// We ran out of i options first
			return true
		}

		// We ran out of j options first or at the same time
		return false

	})
}

func TestMCC(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		items          []string
		multiplicities [][2]int
		options        [][]string
		secondary      []string
		solutions      [][][]string
	}{
		{
			[]string{"a", "b"},
			[][2]int{{0, 1}, {3, 3}},
			[][]string{
				{"a", "b"},
				{"a"},
				{"b"},
			},
			[]string{},
			[][][]string{},
		},

		{
			[]string{"a", "b"},
			[][2]int{{1, 1}, {2, 3}},
			[][]string{
				{"a", "b"},
				{"a"},
				{"b"},
			},
			[]string{},
			[][][]string{
				{{"a", "b"}, {"b"}},
			},
		},

		{
			[]string{"a", "b"},
			[][2]int{{1, 1}, {1, 1}},
			[][]string{
				{"a", "b"},
				{"a"},
				{"b"},
			},
			[]string{},
			[][][]string{
				{{"a", "b"}},
				{{"a"}, {"b"}},
			},
		},

		{
			[]string{"a", "b"},
			[][2]int{{0, 1}, {1, 1}},
			[][]string{
				{"a", "b"},
				{"a"},
				{"b"},
			},
			[]string{},
			[][][]string{
				{{"a", "b"}},
				{{"b"}, {"a"}},
				{{"b"}},
			},
		},

		{
			[]string{"a", "b"},
			[][2]int{{0, 1}, {1, 2}},
			[][]string{
				{"a", "b"},
				{"a"},
				{"b"},
			},
			[]string{},
			[][][]string{
				{{"a", "b"}, {"b"}},
				{{"a", "b"}},
				{{"b"}, {"a"}},
				{{"b"}},
			},
		},

		{
			[]string{"a", "b"},
			[][2]int{{0, 2}, {0, 2}},
			[][]string{
				{"a", "b"},
				{"a"},
				{"b"},
			},
			[]string{},
			[][][]string{
				{{"a", "b"}, {"a"}, {"b"}},
				{{"a", "b"}, {"a"}},
				{{"a", "b"}, {"b"}},
				{{"a", "b"}},
				{{"a"}, {"b"}},
				{{"a"}},
				{{"b"}},
				{},
			},
		},

		{
			[]string{"#1", "#2", "00", "01", "10", "11"},
			[][2]int{{2, 2}, {0, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}},
			[][]string{
				{"#1", "00"},
				{"#1", "01"},
				{"#1", "10"},
				{"#1", "11"},
				{"#2", "00", "01"},
				{"#2", "10", "11"},
				{"#2", "01", "11"},
				{"#2", "00", "10"},
			},
			[]string{},
			[][][]string{
				{{"#1", "00"}, {"#1", "01"}, {"#2", "10", "11"}},
				{{"#1", "00"}, {"#1", "10"}, {"#2", "01", "11"}},
				{{"#1", "01"}, {"#1", "11"}, {"#2", "00", "10"}},
				{{"#1", "10"}, {"#1", "11"}, {"#2", "00", "01"}},
			},
		},

		{
			[]string{"#1", "#2", "00", "01", "10", "11"},
			[][2]int{{2, 3}, {0, 2}, {1, 1}, {1, 1}, {1, 1}, {1, 1}},
			[][]string{
				{"#1", "00"},
				{"#1", "01"},
				{"#1", "10"},
				{"#1", "11"},
				{"#2", "00", "10"},
				{"#2", "10", "11"},
				{"#2", "01", "11"},
				{"#2", "00", "01"},
			},
			[]string{},
			[][][]string{
				{{"#1", "00"}, {"#1", "01"}, {"#2", "10", "11"}},
				{{"#1", "00"}, {"#1", "10"}, {"#2", "01", "11"}},
				{{"#1", "01"}, {"#1", "11"}, {"#2", "00", "10"}},
				{{"#1", "10"}, {"#1", "11"}, {"#2", "00", "01"}},
			},
		},

		{
			[]string{"#1", "#2", "00", "01", "10", "11"},
			[][2]int{{0, 4}, {0, 2}, {1, 1}, {1, 1}, {1, 1}, {1, 1}},
			[][]string{
				{"#1", "00"},
				{"#1", "01"},
				{"#1", "10"},
				{"#1", "11"},
				{"#2", "00", "10"},
				{"#2", "10", "11"},
				{"#2", "01", "11"},
				{"#2", "00", "01"},
			},
			[]string{},
			[][][]string{
				{{"#1", "00"}, {"#1", "01"}, {"#1", "10"}, {"#1", "11"}},
				{{"#1", "00"}, {"#1", "01"}, {"#2", "10", "11"}},
				{{"#1", "00"}, {"#1", "10"}, {"#2", "01", "11"}},
				{{"#1", "01"}, {"#1", "11"}, {"#2", "00", "10"}},
				{{"#1", "10"}, {"#1", "11"}, {"#2", "00", "01"}},
				{{"#2", "00", "01"}, {"#2", "10", "11"}},
				{{"#2", "00", "10"}, {"#2", "01", "11"}},
			},
		},
	}

	for _, c := range cases {
		got := make([][][]string, 0)
		stats := &ExactCoverStats{
			Progress:  false,
			Delta:     0,
			Debug:     false,
			Verbosity: 2,
		}
		err := MCC(c.items, c.multiplicities, c.options, c.secondary, stats,
			func(solution [][]string) bool {
				got = append(got, solution)
				return true
			})

		if err != nil {
			t.Error(err)
		}

		sortSolutions(got)
		sortSolutions(c.solutions)

		if !reflect.DeepEqual(got, c.solutions) {
			t.Errorf("Got solutions %v; want %v", got, c.solutions)
		}
	}
}

func TestExercise_7221_69(t *testing.T) {
	// This verifies Exercise 7.2.2.1-69, Gerrymandering in Bitland

	// the board
	board := make(Polyomino, 0)
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			board = append(board, Point{X: x, Y: y})
		}
	}

	// the communities - 0 represents B or Big-Endian
	communities := [9][9]int{
		{0, 0, 1, 0, 1, 1, 1, 1, 0},
		{1, 1, 1, 0, 1, 1, 1, 0, 1},
		{0, 0, 1, 0, 1, 0, 0, 1, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1},
		{0, 0, 0, 1, 1, 0, 1, 1, 0},
		{1, 0, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 0, 1, 1},
		{1, 1, 0, 1, 1, 0, 0, 1, 1},
	}

	// Generate the shapes
	shapes := PolyominoPacking(4, 4, 9, true, false)

	// Fill the board with all possible shape placements
	board, shapes = PolyominoFill(board, shapes)

	// Generate items and options for XC solving
	_, shapeOptions := PolyominoXC(board, shapes)

	// Count the number of solutions for the two cases which give seven wins
	// for B
	cases := []struct {
		multiplicities []int // number of occurrences of N_k
		count          int   // number of solutions
	}{
		{
			[]int{2, 0, 0, 0, 0, 6, 1, 0, 0, 0}, // 2 N_0, 6 N_5, 1 N_6
			0,
		},
		{
			[]int{1, 1, 0, 0, 0, 7, 0, 0, 0, 0}, // 1 N_0, 1 N_1, 7 N_5
			60,
		},
	}

	for _, c := range cases {

		var (
			items          []string
			multiplicities [][2]int
			options        [][]string
		)

		// Generate the N_k items with multiplicity > 0
		for k := 0; k < 10; k++ {
			if c.multiplicities[k] > 0 {
				items = append(items, fmt.Sprintf("N%d", k))
				multiplicities = append(multiplicities,
					[2]int{
						c.multiplicities[k],
						c.multiplicities[k],
					})
			}
		}

		// Add the board cell items
		for x := 0; x < 9; x++ {
			for y := 0; y < 9; y++ {
				items = append(items, fmt.Sprintf("%d%d", x, y))
				multiplicities = append(multiplicities, [2]int{1, 1})
			}
		}

		// Add the N_k value to the options
		for _, shape := range shapeOptions {
			k := 0
			for _, point := range shape {
				x, _ := strconv.Atoi(point[0:1])
				y, _ := strconv.Atoi(point[1:2])
				if communities[x][y] == 0 {
					k++
				}
			}
			// Add the option only for N_k with multiplicity > 0
			if c.multiplicities[k] > 0 {
				option := []string{fmt.Sprintf("N%d", k)}
				option = append(option, shape...)
				options = append(options, option)
			}
		}

		// Solve using MCC
		count := 0
		stats := &ExactCoverStats{
			// Progress: true,
			// Delta:    100000000,
			// Debug:    true,
		}
		err := MCC(items, multiplicities, options, []string{}, stats,
			func(solution [][]string) bool {
				count++
				return true
			})

		if err != nil {
			t.Errorf("Got error %v for %v", err, c.multiplicities)
		} else if count != c.count {
			t.Errorf("Got %d solutions for %v; want %d", count, c.multiplicities, c.count)
		}
	}
}
