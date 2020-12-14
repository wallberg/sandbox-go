package taocp

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/wallberg/sandbox/sortx"
)

// Explore Dancing Links from The Art of Computer Programming, Volume 4,
// Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
// Dancing Links, 2020
//
// §7.2.2.1 Dancing Links - Polyominoes

var (
	valueMap = []byte{'0', '1', '2', '3', '4', '5',
		'6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
		'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

	rePair = regexp.MustCompile(`[0-9a-zA-Z]|\[[0-9a-zA-Z-]+?\]`)
)

// pack stores a placement pair in an int
func pack(x int, y int) int { return (x << 16) + y }

// unpack pulls a placement pair out of an int
func unpack(pair int) (int, int) { return pair >> 16, pair & 65535 }

// ParsePlacementPairs parses a single placement pair specification string
// format.
func ParsePlacementPairs(s string) ([]int, error) {

	// find gets the index of value in valueMap
	find := func(value byte) int {
		for i, v := range valueMap {
			if value == v {
				return i
			}
		}
		return -1
	}

	// getValues parses the string format for lists of values
	getValues := func(valuesString string) []int {
		if valuesString[0] == '[' {
			valuesString = valuesString[1 : len(valuesString)-1]
		}
		values := make([]int, 0)
		for i := 0; i < len(valuesString); {
			start, stop := find(valuesString[i]), 0
			if i+2 < len(valuesString) && valuesString[i+1] == '-' {
				stop = find(valuesString[i+2])
				i += 3
			} else {
				stop = start
				i++
			}
			for j := start; j <= stop; j++ {
				values = append(values, j)
			}
		}

		return values
	}

	pairs := make([]int, 0)

	// Split on single space
	for _, pairString := range strings.Split(s, " ") {

		// Find 2 values in each pair
		m := rePair.FindAllStringSubmatch(pairString, -1)
		if m != nil && len(m) == 2 {
			xValues := getValues(m[0][0])
			yValues := getValues(m[1][0])

			for _, x := range xValues {
				for _, y := range yValues {
					sortx.InsertInt(&pairs, pack(x, y))
				}
			}
		} else {
			return nil, fmt.Errorf("Unable to parse pair: '%s'", pairString)
		}
	}

	return pairs, nil
}

// BasePlacements takes one placement pair as input, shifted to minimum
// coordinates, and generates every possible transformation using rotate and
// reflect.
func BasePlacements(first []int) [][]int {

	// minmax finds minimum and maximum (x, y) values
	minmax := func(placement []int) (int, int, int) {
		// Get xMin, yMin, xMax
		xMin, yMin, xMax := -1, -1, -1
		for _, pair := range placement {
			x, y := unpack(pair)
			if xMin == -1 || x < xMin {
				xMin = x
			}
			if yMin == -1 || y < yMin {
				yMin = y
			}
			if xMax == -1 || x > xMax {
				xMax = x
			}
		}
		return xMin, yMin, xMax
	}

	xMin, yMin, _ := minmax(first)

	// Shift, if necessary
	if xMin > 0 || yMin > 0 {
		firstNew := make([]int, len(first))
		for i, pair := range first {
			x, y := unpack(pair)
			firstNew[i] = pack(x-xMin, y-yMin)
		}
		first = firstNew
	}

	n := len(first)

	// Generate placements
	placements := make([][]int, 1)
	placements[0] = first

	for i := 0; i < len(placements); i++ {
		// Generate the rotation and reflection
		rotate := make([]int, len(placements[i]))
		reflect := make([]int, len(placements[i]))

		_, _, xMax := minmax(placements[i])
		for j, pair := range placements[i] {
			x, y := unpack(pair)
			rotate[j] = pack(y, xMax-x)
			reflect[j] = pack(y, x)
		}
		sort.Ints(rotate)
		sort.Ints(reflect)

		// Add each to the list of placements, if not already there
		for _, placement := range [][]int{rotate, reflect} {
			// See if this placement already exists
			exists := false
			// Iterate over each existing placement
			for j := range placements {
				same := true
				for k := 0; k < n; k++ {
					if placement[k] != placements[j][k] {
						same = false
						break
					}
				}
				if same {
					exists = true
					break
				}
			}

			if !exists {
				placements = append(placements, placement)
			}
		}

	}

	// Sort the list of placements
	sort.Slice(placements, func(i, j int) bool {
		for k := 0; k < n; k++ {
			if placements[i][k] < placements[j][k] {
				return true
			} else if placements[i][k] > placements[j][k] {
				return false
			}
		}
		return false
	})

	return placements
}
