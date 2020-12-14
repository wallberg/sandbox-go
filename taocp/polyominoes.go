package taocp

import (
	"fmt"
	"regexp"
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
					value := (x << 16) + y
					sortx.InsertInt(&pairs, value)
				}
			}
		} else {
			return nil, fmt.Errorf("Unable to parse pair: '%s'", pairString)
		}
	}

	return pairs, nil
}
