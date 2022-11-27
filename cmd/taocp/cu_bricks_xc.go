package main

import (
	"fmt"
	"log"

	"github.com/wallberg/sandbox-go/taocp"
	"gopkg.in/yaml.v2"
)

// initialize this command by adding it to the parser
func init() {

	_, err = cuBricksCommand.AddCommand("xc",
		"Generate Bricks (l x m x n) XCC",
		"Generate YAML format input to XCC solver for Bricks (l x m x n)",
		&cuBricksXcCommandDataType{},
	)
	if err != nil {
		log.Fatalf("Error adding cu bricks xc subcommand: %v", err)
	}
}

type cuBricksXcCommandDataType struct {
	L        int  `short:"l" description:"l dimension size of the brick" default:"1"`
	M        int  `short:"m" description:"m dimension size of the brick" default:"1"`
	N        int  `short:"n" description:"n dimension size of the brick" default:"1"`
	FixFirst bool `short:"f" long:"fix" description:"fix first cube, reduces solutions by factor of 720"`
}

func (command cuBricksXcCommandDataType) Execute(args []string) error {

	// Generate XCC input
	items, options, sitems := taocp.Bricks(command.L, command.M, command.N, command.FixFirst)

	// Build YAML struct
	xcYaml := taocp.NewExactCoverYaml(items, sitems, options)

	// Serialize to YAML
	data, err := yaml.Marshal(xcYaml)
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	return nil
}
