package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/wallberg/sandbox/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {
	var command mccCommand

	_, err := parser.AddCommand("mcc",
		"Exact Cover w/ Multiplicities and Colors (MCC)",
		`Solve Exact Cover w/ Multiplicities and Colors (MCC) problems using taocp.MCC
Uses YAML for input and output`,
		&command,
	)
	if err != nil {
		log.Fatalf("Error adding xc command: %v", err)
	}
}

type mccCommand struct {
	Input     string `short:"i" long:"input" description:"Input YAML" default:"-"`
	Output    string `short:"o" long:"output" description:"Output YAML" default:"-"`
	Verbosity int    `short:"v" long:"verbosity" description:"Verbosity level" default:"1"`
	Delta     int    `short:"d" long:"delta" description:"Display progress ~Delta nodes (Verbosity > 0)" default:"100000000"`
	Compact   bool   `short:"c" long:"compact" description:"Output solutions in compact format, one per line"`
}

func (command mccCommand) Execute(args []string) error {
	var err error

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

	// Setup multiplicities.
	// If the item ends with {u,v}, where 0 <= u <= v, these are taken as the
	// multiplicity range for that item. Default values for u and v are 1.
	multiplicities := make([][2]int, len(xcYaml.Items))
	reMultiplicities := regexp.MustCompile(`^(.*)\{([0-9]+),([0-9]+)\}$`)
	for i, item := range xcYaml.Items {
		// Default values
		multiplicities[i][0] = 1
		multiplicities[i][1] = 1

		if m := reMultiplicities.FindStringSubmatch(item); m != nil {
			// Remove the multiplicities from the item name
			xcYaml.Items[i] = m[1]

			// Set the u value
			if multiplicities[i][0], err = strconv.Atoi(m[2]); err != nil {
				return err
			}

			// Set the v value
			if multiplicities[i][1], err = strconv.Atoi(m[3]); err != nil {
				return err
			}
		}
	}

	// Solve
	stats := &taocp.ExactCoverStats{
		Debug:     command.Verbosity > 1,
		Progress:  command.Verbosity > 0,
		Verbosity: command.Verbosity - 2,
		Delta:     command.Delta,
	}

	if !command.Compact {
		output.WriteString("solutions:\n")
	}
	err = taocp.MCC(xcYaml.Items, multiplicities, options, xcYaml.SItems, stats,
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

	return nil
}
