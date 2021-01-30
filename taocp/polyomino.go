package taocp

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/wallberg/sandbox/graph"
	"gopkg.in/yaml.v2"
)

// Explore Dancing Links from The Art of Computer Programming, Volume 4,
// Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
// Dancing Links, 2020
//
// §7.2.2.1 Dancing Links - Polyominoes

var (
	rePair = regexp.MustCompile(`[0-9a-zA-Z]|\[[0-9a-zA-Z-]+?\]`)

	// PolyominoSets contains sets of common shapes
	PolyominoSets = LoadPolyominoes()
)

// LoadPolyominoes loads the standard sets of shapes
func LoadPolyominoes() *PolyominoShapes {
	// Load in ./assets/polyominoes.yaml
	box := packr.NewBox("./assets")

	data, err := box.FindString("polyominoes.yaml")
	if err != nil {
		log.Fatalf("Error reading assets/polyominoes.yaml: %v\n", err)
	}

	// Read the yaml file
	shapes := NewPolyominoShapes()
	err = yaml.Unmarshal([]byte(data), &shapes)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Enrich the pieces
	for _, pieceset := range shapes.PieceSets {
		for _, shape := range pieceset {

			points, err := ParsePlacementPairs(shape.Shape)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			shape.Placements = BasePlacements(points, true)
			shape.Points = shape.Placements[0]
		}

	}

	// Enrich the boards
	for _, shape := range shapes.Boards {

		points, err := ParsePlacementPairs(shape.Shape)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		shape.Placements = BasePlacements(points, false)
		shape.Points = shape.Placements[0]

	}

	return shapes
}

// ParsePlacementPairs parses a single placement pair specification string
// format.
func ParsePlacementPairs(s string) (Polyomino, error) {

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

	var po Polyomino
	pset := make(pointset)

	// Split on single space
	for _, pairString := range strings.Split(s, " ") {

		// Find 2 values in each pair
		m := rePair.FindAllStringSubmatch(pairString, -1)
		if m != nil && len(m) == 2 {
			xValues := getValues(m[0][0])
			yValues := getValues(m[1][0])

			for _, x := range xValues {
				for _, y := range yValues {
					point := Point{X: x, Y: y}
					if !pset[point] {
						po = append(po, point)
						pset[point] = true
					}
				}
			}
		} else {
			return nil, fmt.Errorf("Unable to parse pair: '%s'", pairString)
		}
	}

	// Sort the points
	sortPoints(po)

	return po, nil
}

// BasePlacements takes one placement pair as input and shifts it to minimum
// coordinates, and optionally generates every possible transformation using
// rotate and reflect.
func BasePlacements(first Polyomino, transform bool) []Polyomino {

	xMin, yMin, _, _ := minmax(first)

	// Shift, if necessary
	if xMin > 0 || yMin > 0 {
		firstNew := make(Polyomino, len(first))
		for i, point := range first {
			firstNew[i] = Point{X: point.X - xMin, Y: point.Y - yMin}
		}
		first = firstNew
	}

	n := len(first)

	// Generate placements
	placements := make([]Polyomino, 1)
	placements[0] = first

	if transform {
		for i := 0; i < len(placements); i++ {
			// Generate the rotation and reflection
			rotate := make(Polyomino, len(placements[i]))
			reflect := make(Polyomino, len(placements[i]))

			_, _, xMax, _ := minmax(placements[i])
			for j, po := range placements[i] {
				rotate[j] = Point{X: po.Y, Y: xMax - po.X}
				reflect[j] = Point{X: po.Y, Y: po.X}
			}
			sortPoints(rotate)
			sortPoints(reflect)

			// Add each to the list of placements, if not already there
			for _, placement := range []Polyomino{rotate, reflect} {
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
	sortPolyominoes(placements)

	return placements
}

// Polyominoes uses the list of piece shape set names and the board shape name
// found in PolyominoSets to generate items, options, and secondary items
// to find solutions using ExactCover().
func Polyominoes(shapeSetNames []string, boardName string) ([]string, [][]string, []string) {

	// Build the list of items
	items := make([]string, 0)
	cells := make(map[Point]bool)

	// Add the piece items
	for _, piecesetName := range shapeSetNames {
		pieceset := PolyominoSets.PieceSets[piecesetName]
		for shapeName := range pieceset {
			name := fmt.Sprintf("s%sp%s", piecesetName, shapeName)
			items = append(items, name)
		}
	}

	// Add the board items
	board := PolyominoSets.Boards[boardName]
	if board == nil {
		log.Fatalf("Can't find board shape named '%s'", boardName)
	}
	_, _, xMaxBoard, yMaxBoard := minmax(board.Placements[0])

	for _, point := range board.Placements[0] {
		cellItem := fmt.Sprintf("%c%c", valueMap[point.X], valueMap[point.Y])
		items = append(items, cellItem)
		cells[point] = true
	}

	// Build the list of options
	options := make([][]string, 0)

	// Iterate over each shape
	for _, piecesetName := range shapeSetNames {
		pieceset := PolyominoSets.PieceSets[piecesetName]
		for shapeName, shape := range pieceset {

			// Iterate over each shape base placement
			for _, placement := range shape.Placements {

				// Get the bounds of this placement
				_, _, xMax, yMax := minmax(placement)

				// Iterate over delta placements
				for xDelta := 0; xDelta+xMax <= xMaxBoard; xDelta++ {
					for yDelta := 0; yDelta+yMax <= yMaxBoard; yDelta++ {

						// Add the option, if all cells are in the board
						option := make([]string, len(placement)+1)
						name := fmt.Sprintf("s%sp%s", piecesetName, shapeName)
						option[0] = name
						addOption := true
						for i, point := range placement {
							x, y := point.X, point.Y
							x += xDelta
							y += yDelta
							if !cells[Point{X: x, Y: y}] {
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

// PolyominoPacking generates polyominoes of size n which fit into an x by y
// bounding box. Optionally exclude straight pieces and non-convex pieces.
// Returns a list of polyomino shapes.
func PolyominoPacking(x int, y int, n int, includeStraight bool,
	includeNonConvex bool) (pos []Polyomino) {

	// Generate the bounding grid
	g := graph.CartesianProduct(graph.Path(x), graph.Path(y))

	// Generate shapes which fit into the grid
	graph.ConnectedSubsets(g, n, func(solution []int) bool {

		// Translate the subgraph into a polyomino shape
		var po Polyomino
		for _, v := range solution {
			po = append(po, Point{v % y, v / y})
		}

		// Determine if this shape should be included
		if includeStraight || !po.IsStraight() {
			if includeNonConvex || po.IsConvex() {
				pos = append(pos, po)
			}
		}

		return false
	})

	return pos
}

// PolyominoXC generates items and options for XC solving given an input
// Polyomino board and list of Polyomino shapes. This function assumes that the
// shapes already cover all the intended board positions.
func PolyominoXC(board Polyomino, shapes []Polyomino) (items []string, options [][]string) {

	// Create a Poinset for the board points
	boardSet := board.toPointset()

	// pointItem generates an item name for a Point
	pointItem := func(point Point) string {
		return fmt.Sprintf("%c%c", valueMap[point.X], valueMap[point.Y])
	}

	// Add the items, one for each board position
	for _, point := range board {
		items = append(items, pointItem(point))
	}

	// Add the options
	// Iterate over each shape
	for _, shape := range shapes {

		var option []string

		for _, point := range shape {
			if !boardSet[point] {
				log.Fatalf("Shape %v contains point %v which is not in the board",
					shape, point)
			}

			option = append(option, pointItem(point))
		}

		options = append(options, option)
	}

	return items, options
}

// PolyominoFill translates the board to origin and for each shape translates
// the shape to origin, determines if this shape has not been handled yet,
// then generates all placements which fit on the board.
func PolyominoFill(board Polyomino, shapes []Polyomino) (boardOut Polyomino, shapesOut []Polyomino) {

	// Shift the board to the origin
	boardOut = board.translateToOrigin()

	// Get board bounds
	_, _, xBoardMax, yBoardMax := boardOut.Bounds()

	// Create a Poinset for the board points
	boardSet := boardOut.toPointset()

	seenSet := make(map[string]bool)

	// seen checks if a shape has already been handled; if yes it returns true,
	// if no it marks as handled and returns false
	seen := func(shape Polyomino) bool {
		// Generate the string key
		keyBuilder := strings.Builder{}
		for _, point := range shape {
			keyBuilder.WriteString(fmt.Sprintf("%c%c", valueMap[point.X], valueMap[point.Y]))
		}
		key := keyBuilder.String()

		// Check if seen
		if seenSet[key] {
			return true
		}

		seenSet[key] = true
		return false
	}

	// Iterate over each shape
	for _, shape := range shapes {

		// Shift the shape to the origin and sort
		shapeOut := shape.translateToOrigin()

		// Check if we've handled this shape already
		if !seen(shapeOut) {

			// Get shape bounds
			_, _, xMax, yMax := shapeOut.Bounds()

			// Iterate over shape placements
			for xDelta := 0; xDelta+xMax <= xBoardMax; xDelta++ {
				for yDelta := 0; yDelta+yMax <= yBoardMax; yDelta++ {

					// Add the placement if all the points are in the board
					placement := make(Polyomino, len(shapeOut))
					addPlacement := true
					for i, point := range shapeOut {
						x := point.X + xDelta
						y := point.Y + yDelta
						placementPoint := Point{X: x, Y: y}
						if !boardSet[placementPoint] {
							addPlacement = false
							break
						}
						placement[i] = placementPoint
					}
					if addPlacement {
						shapesOut = append(shapesOut, placement)
					}
				}
			}
		}
	}

	return boardOut, shapesOut
}
