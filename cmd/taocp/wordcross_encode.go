package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/wallberg/sandbox-go/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {

	if wcCommand := parser.Find("wc"); wcCommand != nil {
		var command wcEncodeCommand
		_, err := wcCommand.AddCommand("encode",
			"Encode WordCross puzzle as XCC problem",
			`Encode WordCross puzzle as XCC problem. Words are read from stdin`,
			&command,
		)
		if err != nil {
			log.Fatalf("Error adding wc encode subcommand: %v", err)
		}
	} else {
		log.Fatalf("Error adding decode subcommand: Unable to find parent 'wc' command")
	}
}

type wcEncodeCommand struct {
	M int `short:"m" long:"m" description:"number of rows" default:"8"`
	N int `short:"n" long:"n" description:"number of columngs" default:"8"`
}

func (command wcEncodeCommand) Execute(args []string) error {
	var words []string

	// Read words from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Generate XCC input
	items, options, sitems := taocp.WordCross(words, command.M, command.N)

	// Build YAML struct
	xcYaml := taocp.NewExactCoverYaml(items, sitems, options)

	// Serialize to YAML on stdout
	data, err := yaml.Marshal(xcYaml)
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	return nil
}
