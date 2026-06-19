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

	if wsCommand := parser.Find("ws"); wsCommand != nil {
		var command wsKernelsCommand
		_, err := wsCommand.AddCommand("kernels",
			"Generate Word Stair Kernels XCC",
			"Generate YAML format input to XCC solver for Word Stair Kernels (Exercise 7.2.2.1-91). Words are read from stdin (n=5).",
			&command,
		)
		if err != nil {
			log.Fatalf("Error adding ws xc subcommand: %v", err)
		}
	} else {
		log.Fatalf("Error adding kernels sub-command: Unable to find parent 'ws' command")

	}
}

type wsKernelsCommand struct {
	Right bool `short:"r" long:"right" description:"generate a right word stair; (default: a left word stair)"`
}

func (command wsKernelsCommand) Execute(args []string) error {

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
	items, options, sitems := taocp.WordStairKernel(words, !command.Right)

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
