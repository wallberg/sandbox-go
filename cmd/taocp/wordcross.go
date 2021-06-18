package main

import (
	"log"
)

// initialize this command by adding it to the parser
func init() {
	var command wcCommand

	_, err := parser.AddCommand("wc",
		"WordCross",
		`Solve WordCross puzzles`,
		&command,
	)
	if err != nil {
		log.Fatalf("Error adding wc command: %v", err)
	}
}

type wcCommand struct {
}
