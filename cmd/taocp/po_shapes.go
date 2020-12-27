package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wallberg/sandbox/taocp"
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
	N      int    `short:"n" long:"n" description:"generate pieces of size n <= 62" default:"5"`
	Output string `short:"o" long:"output" description:"Output YAML file" default:"-"`
}

func (command poShapesCommand) Execute(args []string) error {
	// Error check the input
	if command.N > 62 {
		return fmt.Errorf("Got n=%d; want n <= 62 because we use [0-9a-zA-Z] to represent the coordinates", command.N)
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
	shapes.PieceSets[setName] = make(map[string]string)

	// Generate the shapes
	for i, shape := range taocp.GeneratePolyominoShapes(command.N) {
		// Add the piece to the YAML output strucuture
		pieceName := fmt.Sprintf("%d", i)
		var shapeString strings.Builder
		for _, point := range shape {
			if shapeString.Len() > 0 {
				shapeString.WriteString(" ")
			}
			shapeString.WriteString(point.String())
		}
		shapes.PieceSets[setName][pieceName] = shapeString.String()
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
