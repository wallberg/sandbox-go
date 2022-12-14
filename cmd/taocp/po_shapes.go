package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wallberg/sandbox-go/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {

	if poCommand := parser.Find("po"); poCommand != nil {
		var command poShapesCommand
		_, err := poCommand.AddCommand("shapes",
			"Generate Polyominoes Shapes",
			"Generate YAML format Polyomino shapes",
			&command,
		)
		if err != nil {
			log.Fatalf("Error adding po shapes subcommand: %v", err)
		}
	} else {
		log.Fatalf("Error adding xc sub-command: Unable to find parent 'po' command")

	}
}

type poShapesCommand struct {
	N          int    `short:"n" long:"n" description:"Generate pieces of size n <= 62" default:"5"`
	Output     string `short:"o" long:"output" description:"Output YAML file" default:"-"`
	Convex     bool   `short:"c" long:"convex" description:"Limit to convex shapes (default: false)"`
	X          int    `short:"x" long:"x" description:"Limit to shape size on one axis" default:"62"`
	Y          int    `short:"y" long:"y" description:"Limit to shape size on other axis" default:"62"`
	Placements bool   `short:"p" long:"placements" description:"Include all base placements"`
}

func (command poShapesCommand) Execute(args []string) error {
	// Error check the input
	if command.N > 62 {
		return fmt.Errorf("got n=%d; want n <= 62 because we use [0-9a-zA-Z] to represent the coordinates", command.N)
	}

	// Open output file for writing
	var err error
	var output *os.File
	if command.Output == "-" {
		output = os.Stdout
	} else {
		if output, err = os.Create(command.Output); err != nil {
			return err
		}
	}
	defer output.Close()

	// Setup the YAML output structure
	shapes := taocp.NewPolyominoShapes()
	setName := fmt.Sprintf("%d", command.N)
	shapes.PieceSets[setName] = make(map[string]*taocp.PolyominoShape)

	// Generate the shapes
	for i, po := range taocp.GeneratePolyominoShapes(command.N) {
		// Skip if the shape must be convex and it
		if command.Convex && !po.IsConvex() {
			continue
		}

		// Skip if the shape does not fit in the required bounding box
		if command.X < 62 || command.Y < 62 {
			xMin, yMin, xMax, yMax := po.Bounds()
			xSize, ySize := xMax-xMin+1, yMax-yMin+1

			// Try both orientations of the bounding box
			if !(xSize <= command.X && ySize <= command.Y) && !(ySize <= command.X && xSize <= command.Y) {
				continue
			}
		}

		// Add the piece to the YAML output strucuture
		pieceName := fmt.Sprintf("%d", i)
		var shapeString strings.Builder
		for _, point := range po {
			if shapeString.Len() > 0 {
				shapeString.WriteString(" ")
			}
			shapeString.WriteString(point.String())
		}
		shape := taocp.PolyominoShape{Shape: shapeString.String()}
		if command.Placements {
			shape.Placements = taocp.BasePlacements(po, true)
		}
		shapes.PieceSets[setName][pieceName] = &shape
	}

	// Generate the YAML
	var data []byte
	if data, err = yaml.Marshal(shapes); err != nil {
		return err
	}

	// Write the YAML
	if _, err = output.Write(data); err != nil {
		return err
	}

	return nil
}
