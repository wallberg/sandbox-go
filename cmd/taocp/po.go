package main

import (
	"log"
)

var poCommand poCommandType

// initialize this command by adding it to the parser
func init() {

	_, err := parser.AddCommand("po",
		"Polyominoes",
		"Operations on Polyominoes, plane geometric figures formed by joining one or more equal squares edge to edge (aka n-ominoes)",
		&poCommand,
	)
	if err != nil {
		log.Fatalf("Error adding po command: %v", err)
	}
}

type poCommandType struct {
}
