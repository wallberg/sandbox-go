package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	graphx "github.com/wallberg/sandbox-go/graph"
	"github.com/wallberg/sandbox-go/taocp"
	"github.com/yourbasic/graph"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {

	if wcCommand := parser.Find("wc"); wcCommand != nil {
		var command wcDecodeCommand
		_, err := wcCommand.AddCommand("decode",
			"Decode XCC solutions to WordCross solutions",
			"Decode XCC solutions to WordCross solutions",
			&command,
		)
		if err != nil {
			log.Fatalf("Error adding wc decode subcommand: %v", err)
		}
	} else {
		log.Fatalf("Error adding decode subcommand: Unable to find parent 'wc' command")

	}
}

type wcDecodeCommand struct {
	M         int    `short:"m" long:"m" description:"number of rows" default:"8"`
	N         int    `short:"n" long:"n" description:"number of columngs" default:"8"`
	Input     string `short:"i" long:"input" description:"Input YAML" default:"-"`
	Distinct  bool   `short:"d" long:"distinct" description:"Limit to distinct solutions"`
	Connected bool   `short:"c" long:"connected" description:"Limit to connected solutions"`
}

func (command wcDecodeCommand) Execute(args []string) error {
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
	var solutions taocp.ExactCoverSolutions
	err = yaml.Unmarshal(buf.Bytes(), &solutions)
	if err != nil {
		return err
	}

	if len(solutions.Solutions) == 0 {
		return nil
	}

	// Process the solutions
	m, n := command.M, command.N
	var i, j int
	reCell := regexp.MustCompile(`^([0-9A-Za-z]{2}):([A-Za-z])$`)

	// getKey creates a key to uniquely identify a grid
	getKey := func(grid [][]string) string {
		var key strings.Builder
		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				key.WriteString(grid[i][j])
			}
		}

		return key.String()
	}

	grids := make(map[string]bool)
	count := 0

	for _, solution := range solutions.Solutions {

		// Setup an empty grid
		grid := make([][]string, m)
		for i = 0; i < m; i++ {
			grid[i] = make([]string, n)
			for j = 0; j < n; j++ {
				grid[i][j] = " "
			}
		}

		// Fill out the grid
		for _, option := range solution {
			for _, item := range strings.Fields(string(option)) {
				if match := reCell.FindStringSubmatch(item); match != nil {
					if i, j, err = taocp.DecodeCell(match[1]); err != nil {
						return err
					}

					grid[i-1][j-1] = match[2]
				}
			}
		}

		// Determine if we've already seen this grid
		key := getKey(grid)
		distinct := !grids[key]
		if distinct {
			grids[key] = true
		}

		// Determine if connected
		g := graph.New(m * n)
		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				if grid[i][j] != " " {
					if i+1 < m && grid[i+1][j] != " " {
						g.AddBoth(i*n+j, (i+1)*n+j)
					}
					if j+1 < n && grid[i][j+1] != " " {
						g.AddBoth(i*n+j, i*n+j+1)
					}
				}
			}
		}
		g, _ = graphx.RemoveIsolated(g)
		connected := graph.Connected(g)

		if (!command.Distinct || distinct) && (!command.Connected || connected) {

			// Print out the grid as a valid solution
			for i = 0; i < m; i++ {
				fmt.Println(grid[i])
			}
			fmt.Println("---")

			count++
		}

	}

	fmt.Println("Count: ", count)

	return nil
}
