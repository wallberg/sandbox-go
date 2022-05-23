package main

import (
	"log"

	flags "github.com/jessevdk/go-flags"
)

var (
	cuBricksCommand *flags.Command
)

// initialize this command by adding it to the parser
func init() {

	cuBricksCommand, err = cuCommand.AddCommand("bricks",
		"Bricks formed from cubes",
		"Bricks formed from cubes (Exercise 7.2.2.1-147)",
		&cuBricksCommandDataType{},
	)
	if err != nil {
		log.Fatalf("Error adding cu bricks subcommand: %v", err)
	}
}

type cuBricksCommandDataType struct {
}
