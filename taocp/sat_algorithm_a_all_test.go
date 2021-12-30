package taocp

import (
	"log"
	"testing"
)

func TestSATAlgorithmAAll(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		n       int        // number of strictly distinct literals
		sat     bool       // is satisfiable
		clauses SATClauses // clauses to satisfy
	}{
		{3, true, SATClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}}},
		{3, false, SATClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}, {1, 2, -3}}},
		{4, true, ClausesRPrime},
		{4, false, ClausesR},
	}

	for _, c := range cases {

		stats := SATStats{
			// Debug:    true,
			// Progress: true,
		}
		options := SATOptions{}

		got := false
		SATAlgorithmAAll(c.n, c.clauses, &stats, &options,
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
