package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

var options Options

var parser = flags.NewParser(&options, flags.Default)

// Options provides the top-level usage for this program
type Options struct {
	// None
}

func main() {

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
