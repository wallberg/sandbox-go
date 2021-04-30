package main

import (
	"log"
)

var wsCommand wsCommandType

// initialize this command by adding it to the parser
func init() {

	_, err := parser.AddCommand("ws",
		"Word Stairs",
		"A word stair is a cyclic arrangement of words, offset stepwise, that contains 2p distinct words across and down. They exist in two varieties, left and right.",
		&wsCommand,
	)
	if err != nil {
		log.Fatalf("Error adding po command: %v", err)
	}
}

type wsCommandType struct {
}
