package taocp

import (
	"log"
	"testing"
)

func TestSatAlgorithmAAll(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		n       int        // number of strictly distinct literals
		sat     bool       // is satisfiable
		clauses SatClauses // clauses to satisfy
	}{
		{3, true, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}}},
		{3, false, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}, {1, 2, -3}}},
		{4, true, ClausesRPrime},
		{4, false, ClausesR},
	}

	for _, c := range cases {

		stats := SatStats{
			// Debug:    true,
			// Progress: true,
		}
		options := SatOptions{}

		got := false
		SatAlgorithmAAll(c.n, c.clauses, &stats, &options,
			func(solution []int) bool {
				got = true
				// log.Printf("solution=%v", solution)
				return true
			})

		if got != c.sat {
			t.Errorf("expected satisfiable=%t for clauses %v; got %t", c.sat, c.clauses, got)
		}
	}
}
