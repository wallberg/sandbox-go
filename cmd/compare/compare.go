package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wallberg/sandbox-go/taocp"
)

// main - compare performance of SAT algorithms/options
func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	langford13, langford13CoverOptions := taocp.SatLangford(13)
	langford13n := len(langford13CoverOptions)

	cases := []struct {
		n       int              // number of strictly distinct literals
		sat     bool             // is satisfiable
		clauses taocp.SatClauses // clauses to satisfy
		desc    string           // description
	}{
		{langford13n, false, langford13, "L3 - langford(13)"},
		{100, false, taocp.SatRand(3, 420, 100, 1), "rand(3,420,100,1)"},
		{97, false, taocp.SatWaerdan(3, 10, 97), "W1 - waerden(3,10,97)"},
	}

	for _, c := range cases {
		fmt.Println(c.desc)
		for _, alg := range []string{"D", "L", "Lcr"} {
			stats := taocp.SatStats{
				// Progress: true,
				// Delta:    100000000,
			}
			options := taocp.SatOptions{}
			var sat bool
			var solution []int

			start := time.Now()

			switch alg {
			case "A":
				sat, solution = taocp.SatAlgorithmA(c.n, c.clauses, &stats, &options)

			case "B":
				sat, solution = taocp.SatAlgorithmB(c.n, c.clauses, &stats, &options)

			case "D":
				sat, solution = taocp.SatAlgorithmD(c.n, c.clauses, &stats, &options)

			case "L":
				optionsL := taocp.NewSatAlgorithmLOptions()
				sat, solution = taocp.SatAlgorithmL(c.n, c.clauses, &stats, &options, optionsL)

			case "Lcr":
				optionsL := taocp.NewSatAlgorithmLOptions()
				optionsL.CompensationResolvants = true
				sat, solution = taocp.SatAlgorithmL(c.n, c.clauses, &stats, &options, optionsL)
			}

			duration := time.Since(start)
			var durationPerNode string
			if stats.Nodes == 0 {
				durationPerNode = "N/A"
			} else {
				durationPerNode = (duration / time.Duration(stats.Nodes)).String()
			}

			valid := true
			if sat {
				valid = taocp.SatTest(c.n, c.clauses, solution)
			}
			// levels := stats.Levels
			// if stats.MaxLevel > -1 {
			// 	levels = levels[0:stats.MaxLevel]
			// } else {
			// 	levels = levels[0:0]
			// }

			fmt.Printf("%4s: duration=%s (%s/node) sat=%t, valid=%t, n=%d, #clauses=%d, #nodes=%d, #levels=%d\n",
				alg, duration.Round(time.Millisecond), durationPerNode, sat, valid, c.n, len(c.clauses), stats.Nodes, stats.MaxLevel)
		}
		fmt.Println()
	}
}
