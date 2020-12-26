package main

import (
	"fmt"
	"log"

	"github.com/wallberg/sandbox/taocp"
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
	N int `short:"n" long:"n" description:"generate pieces of size n" default:"5"`
}

func (command poShapesCommand) Execute(args []string) error {
	for _, poly := range taocp.PolyominoShapes(command.N) {
		fmt.Println(poly)
	}

	return nil
}
