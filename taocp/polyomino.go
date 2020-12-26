package taocp

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/wallberg/sandbox/sortx"
	"gopkg.in/yaml.v2"
)

// Explore Dancing Links from The Art of Computer Programming, Volume 4,
// Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
// Dancing Links, 2020
//
// §7.2.2.1 Dancing Links - Polyominoes

// Polyomino represents a single polyomino shape
type Polyomino struct {
	Name       string  // name of the shape
	Shape      string  // string specification of the shape
	Placements [][]int // base placements of the shape (rotation, reflection)
}

// PolyominoSet represents a set of polyomino shapes
type PolyominoSet struct {
	Name   string      // name of the set
	Shapes []Polyomino // list of Polyomino shapes
}

var (
	valueMap = []byte{'0', '1', '2', '3', '4', '5',
		'6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
		'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

	rePair = regexp.MustCompile(`[0-9a-zA-Z]|\[[0-9a-zA-Z-]+?\]`)

	// PolyominoSets contains sets of common shapes
	PolyominoSets = LoadPolyominoes()
)

// LoadPolyominoes loads the standard sets of shapes
func LoadPolyominoes() map[string]PolyominoSet {
	// Load in ./assets/polyominoes.yaml
	box := packr.NewBox("./assets")

	data, err := box.FindString("polyominoes.yaml")
	if err != nil {
		log.Fatalf("Error reading assets/polyominoes.yaml: %v\n", err)
	}

	// Read the yaml file into map
	yamlSets := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &yamlSets)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Read the sets
	sets := make(map[string]PolyominoSet)

	for yamlSetName, yamlSet := range yamlSets {

		set := PolyominoSet{Name: yamlSetName.(string)}

		// Iterate over the shapes
		for yamlShapeName, yamlShape := range yamlSet.(map[interface{}]interface{}) {
			// Add a shape to the set
			shape := Polyomino{Name: yamlShapeName.(string)}
			pairs, err := ParsePlacementPairs(yamlShape.(string))
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			shape.Placements = BasePlacements(pairs, set.Name != "Boards")
			set.Shapes = append(set.Shapes, shape)
		}

		sets[set.Name] = set
	}

	return sets
}

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

// minmax finds minimum and maximum (x, y) values
func minmax(placement []int) (int, int, int, int) {
	xMin, yMin, xMax, yMax := -1, -1, -1, -1
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
		if yMax == -1 || y > yMax {
			yMax = y
		}
	}
	return xMin, yMin, xMax, yMax
}

// BasePlacements takes one placement pair as input and shifts it to minimum
// coordinates, and optionally generates every possible transformation using
// rotate and reflect.
func BasePlacements(first []int, transform bool) [][]int {

	xMin, yMin, _, _ := minmax(first)

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

	if transform {
		for i := 0; i < len(placements); i++ {
			// Generate the rotation and reflection
			rotate := make([]int, len(placements[i]))
			reflect := make([]int, len(placements[i]))

			_, _, xMax, _ := minmax(placements[i])
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

// Polyominoes uses the list of piece shape set names and the board shape name
// found in PolyominoSets to generate items, options, and secondary items
// to find solutions using ExactCover().
func Polyominoes(shapeSetNames []string, boardName string) ([]string, [][]string, []string) {

	// Get the board shape
	var board *Polyomino
	for _, x := range PolyominoSets["Boards"].Shapes {
		if x.Name == boardName {
			board = &x
			break
		}
	}
	if board == nil {
		log.Fatalf("Can't find board shape named '%s'", boardName)
	}
	_, _, xMaxBoard, yMaxBoard := minmax(board.Placements[0])

	// Build the list of items and map of available cells
	items := make([]string, 0)
	cells := make(map[int]bool)

	for _, shapeSetName := range shapeSetNames {
		for _, shape := range PolyominoSets[shapeSetName].Shapes {
			items = append(items, shape.Name)
		}
	}

	for _, cell := range board.Placements[0] {
		x, y := unpack(cell)
		cellItem := fmt.Sprintf("%c%c", valueMap[x], valueMap[y])
		items = append(items, cellItem)
		cells[cell] = true
	}

	// Build the list of options
	options := make([][]string, 0)

	// Iterate over each shape
	for _, shapeSetName := range shapeSetNames {
		for _, shape := range PolyominoSets[shapeSetName].Shapes {

			// Iterate over each shape base placement
			for _, placement := range shape.Placements {

				// Get the bounds of this placement
				_, _, xMax, yMax := minmax(placement)

				// Iterate over delta placements
				for xDelta := 0; xDelta+xMax <= xMaxBoard; xDelta++ {
					for yDelta := 0; yDelta+yMax <= yMaxBoard; yDelta++ {

						// Add the option, if all cells are in the board
						option := make([]string, len(placement)+1)
						option[0] = shape.Name
						addOption := true
						for i, cell := range placement {
							x, y := unpack(cell)
							x += xDelta
							y += yDelta
							if !cells[pack(x, y)] {
								addOption = false
								break
							}
							cellItem := fmt.Sprintf("%c%c",
								valueMap[x], valueMap[y])
							option[i+1] = cellItem
						}
						if addOption {
							options = append(options, option)
						}
					}
				}
			}
		}
	}

	return items, options, []string{}

}
