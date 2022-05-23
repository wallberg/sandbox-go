package main

import (
	"log"

	flags "github.com/jessevdk/go-flags"
)

var (
	cuCommand *flags.Command
)

// initialize this command by adding it to the parser
func init() {

	cuCommand, err = parser.AddCommand("cu",
		"Cubes",
		"Operations on Cubes, with a different color {a,b,c,d,e,f} on each face",
		&cuCommandDataType{},
	)
	if err != nil {
		log.Fatalf("Error adding cu command: %v", err)
	}
}

type cuCommandDataType struct {
}
