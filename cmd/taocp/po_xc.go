package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/wallberg/sandbox-go/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {

	if poCommand := parser.Find("po"); poCommand != nil {
		var command poXcCommand
		_, err := poCommand.AddCommand("xc",
			"Generate Polyominoes XCC",
			"Generate YAML format input to XCC solver for Polyominoes",
			&command,
		)
		if err != nil {
			log.Fatalf("Error adding po xc subcommand: %v", err)
		}
	} else {
		log.Fatalf("Error adding xc sub-command: Unable to find parent 'po' command")

	}
}

type poXcCommand struct {
	List   bool   `short:"l" long:"list" description:"list available piece sets and board shapes"`
	Pieces string `short:"p" long:"pieces" description:"comma separated list of piece sets" default:"5"`
	Board  string `short:"b" long:"board" description:"board name"`
}

func (command poXcCommand) Execute(args []string) error {

	if command.List {
		// List piece sets and boards

		fmt.Println("Piece Sets")

		// Get sorted list of piece set names
		setNames := make([]string, 0)
		for setName := range taocp.PolyominoSets.PieceSets {
			setNames = append(setNames, setName)
		}
		sort.Strings(setNames)

		// Display piece sets
		for _, setName := range setNames {
			set := taocp.PolyominoSets.PieceSets[setName]
			fmt.Printf("  %s (", setName)
			i := 0
			for shapeName := range set {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(shapeName)
				i++
			}
			fmt.Println(")")
		}

		fmt.Println("\nBoards")

		// Get sorted list of board names
		boardNames := make([]string, 0)
		for boardName := range taocp.PolyominoSets.Boards {
			boardNames = append(boardNames, boardName)
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
		pieces := strings.Split(command.Pieces, ",")

		// Generate XCC input
		items, options, sitems := taocp.Polyominoes(pieces, command.Board)

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
