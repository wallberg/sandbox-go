package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/wallberg/sandbox/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {
	var command poCommand

	_, err := parser.AddCommand("po",
		"Polyominoes",
		"Generate YAML format input to XCC solver for Polyominoes",
		&command,
	)
	if err != nil {
		log.Fatalf("Error adding po command: %v", err)
	}
}

type poCommand struct {
	List   bool     `short:"l" long:"list" description:"list available piece sets and board shapes"`
	Pieces []string `short:"p" long:"pieces" description:"comma separated list of piece sets" default:"5"`
	Board  string   `short:"b" long:"board" description:"board name"`
}

func (command poCommand) Execute(args []string) error {

	if command.List {
		// List piece sets and boards

		fmt.Println("Piece Sets")

		// Get sorted list of piece set names
		setNames := make([]string, 0)
		for setName := range taocp.PolyominoSets {
			setNames = append(setNames, setName)
		}
		sort.Strings(setNames)

		// Display piece sets
		for _, setName := range setNames {
			set := taocp.PolyominoSets[setName]
			if setName != "Boards" {
				fmt.Printf("  %s (", setName)
				for i, shape := range set.Shapes {
					if i > 0 {
						fmt.Print(", ")
					}
					fmt.Print(shape.Name)
				}
				fmt.Println(")")
			}
		}

		fmt.Println("\nBoards")

		// Get sorted list of board names
		boardNames := make([]string, 0)
		for _, board := range taocp.PolyominoSets["Boards"].Shapes {
			boardNames = append(boardNames, board.Name)
		}
		sort.Strings(boardNames)

		// Display board names
		for _, boardName := range boardNames {
			fmt.Printf("  %s\n", boardName)
		}
	} else {
		// Validate command line parameters
		if command.Board == "" {
			return fmt.Errorf("the required flag `-b, --board' was not specified")
		}

		// Generate XCC input
		items, options, sitems := taocp.Polyominoes(command.Pieces, command.Board)

		// Build YAML struct
		xcYaml := taocp.NewExactCoverYaml(items, sitems, options)

		// Serialize to YAML
		data, err := yaml.Marshal(xcYaml)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}

	return nil
}
