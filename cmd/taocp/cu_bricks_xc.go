package main

import (
	"log"
)

// initialize this command by adding it to the parser
func init() {

	_, err = cuBricksCommand.AddCommand("xc",
		"Generate Bricks XCC",
		"Generate YAML format input to XCC solver for Bricks",
		&cuBricksXcCommandDataType{},
	)
	if err != nil {
		log.Fatalf("Error adding cu bricks xc subcommand: %v", err)
	}
}

type cuBricksXcCommandDataType struct {
}

func (command cuBricksXcCommandDataType) Execute(args []string) error {

	return nil
}
