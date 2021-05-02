package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/wallberg/sandbox/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {
	var command xcCommand

	_, err := parser.AddCommand("xc",
		"Exact Cover w/ Colors (XCC)",
		`Solve Exact Cover w/ Colors (XCC) problems using taocp.ExactCoverColors
Uses YAML for input and output`,
		&command,
	)
	if err != nil {
		log.Fatalf("Error adding xc command: %v", err)
	}
}

type xcCommand struct {
	Input         string `short:"i" long:"input" description:"Input YAML" default:"-"`
	Output        string `short:"o" long:"output" description:"Output YAML" default:"-"`
	Verbosity     int    `short:"v" long:"verbosity" description:"Verbosity level" default:"1"`
	Delta         int    `short:"d" long:"delta" description:"Display progress ~Delta nodes (Verbosity > 0)" default:"10000000"`
	Compact       bool   `short:"c" long:"compact" description:"Output solutions in compact format, one per line"`
	Minimax       bool   `short:"m" long:"minimax" description:"Return minimax solutions (multiple)"`
	MinimaxSingle bool   `short:"s" long:"minimax-single" description:"Return minimax solutions (single)"`
	Exercise83    bool   `short:"e" long:"exercise83" description:"Use the curious extension of Exercise 7.2.2.1-83"`
}

func (command xcCommand) Execute(args []string) error {
	var err error

	// Validate Minimax options
	if command.Minimax && command.MinimaxSingle {
		return fmt.Errorf("please select only one of --minimax, --minimax-single")
	}

	// Open input file for reading
	var input *os.File
	if command.Input == "-" {
		input = os.Stdin
	} else {
		if input, err = os.Open(command.Input); err != nil {
			return err
		}
	}
	defer input.Close()

	// Read the input
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, input); err != nil {
		return err
	}

	// Deserialize from YAML
	var xcYaml taocp.ExactCoverYaml
	err = yaml.Unmarshal(buf.Bytes(), &xcYaml)
	if err != nil {
		return err
	}
	options := make([][]string, len(xcYaml.Options))
	for i, option := range xcYaml.Options {
		options[i] = strings.Split(option, " ")
	}

	// Open output file for writing
	var output *os.File
	if command.Output == "-" {
		output = os.Stdout
	} else {
		if output, err = os.Create(command.Output); err != nil {
			return err
		}
	}
	defer output.Close()

	// Solve
	stats := &taocp.ExactCoverStats{
		Debug:     command.Verbosity > 1,
		Progress:  command.Verbosity > 0,
		Verbosity: command.Verbosity - 2,
		Delta:     command.Delta,
	}

	// XCC processing options
	xccOptions := &taocp.XCCOptions{
		Minimax:       command.Minimax || command.MinimaxSingle,
		MinimaxSingle: command.MinimaxSingle,
		Exercise83:    command.Exercise83,
	}

	if !command.Compact {
		output.WriteString("solutions:\n")
	}
	err = taocp.XCC(xcYaml.Items, options, xcYaml.SItems, stats, xccOptions,
		func(solution [][]string) bool {
			if !command.Compact {
				output.WriteString("  -\n")
				for _, option := range solution {
					output.WriteString("    - \"")
					output.WriteString(strings.Join(option, " "))
					output.WriteString("\"\n")
				}
			} else {
				var s strings.Builder
				for _, option := range solution {
					if s.Len() > 0 {
						s.WriteString(", ")
					}
					s.WriteString("\"")
					s.WriteString(strings.Join(option, " "))
					s.WriteString("\"")
				}
				s.WriteString("\n")
				output.WriteString(s.String())
			}
			return true
		})
	if err != nil {
		return err
	}

	log.Println("Stats:", *stats)

	return nil
}
