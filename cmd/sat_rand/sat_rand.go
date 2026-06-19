package main

import (
	"fmt"

	"github.com/wallberg/sandbox-go/taocp"
)

func main() {
	cases := []struct {
		k, m, n int
	}{
		{2, 80, 100},
		{2, 100, 100},
		{2, 400, 100},
		{2, 1000, 1000},
		{2, 1100, 1000},
		{2, 2000, 1000},
		{2, 10000, 10000},
		// {2, 11000, 10000}, // probably not sat
	}

	for i, c := range cases {
		clauses := taocp.SatRand(c.k, c.m, c.n, 0)

		stats := taocp.SatStats{}
		options := taocp.SatOptions{}

		sat, solution := taocp.SatAlgorithmD(c.n, clauses, &stats, &options)

		// Verify a valid solution
		if sat {
			validSolution := taocp.SatTest(c.n, clauses, solution)
			if !validSolution {
				fmt.Printf("For case #%d, k=%d, m=%d, n=%d, satisfiable=true, but received an invalid solution\n", i, c.k, c.m, c.n)
			}
		}

		fmt.Printf("For case #%d, k=%d, m=%d, n=%d, satisfiable=%t\n", i, c.k, c.m, c.n, sat)
	}
}
