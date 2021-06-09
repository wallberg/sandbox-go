package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/wallberg/sandbox/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {
	var command wcCommand

	_, err := parser.AddCommand("wc",
		"WordCross",
		`Generate XCC to solve WordCross puzzles. Words are read from stdin`,
		&command,
	)
	if err != nil {
		log.Fatalf("Error adding xcc command: %v", err)
	}
}

type wcCommand struct {
	M int `short:"m" long:"m" description:"number of rows" default:"8"`
	N int `short:"n" long:"n" description:"number of columngs" default:"8"`
}

func (command wcCommand) Execute(args []string) error {
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
