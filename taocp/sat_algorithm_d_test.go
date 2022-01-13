package taocp

import (
	"log"
	"reflect"
	"testing"
)

func TestSatAlgorithmD(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		n        int        // number of strictly distinct literals
		sat      bool       // is satisfiable
		solution []int      // solution
		clauses  SatClauses // clauses to satisfy
	}{
		{3, true, []int{0, 0, 1}, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}}},
		{3, false, nil, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}, {1, 2, -3}}},
		{4, true, []int{0, 1, 0, 1}, ClausesRPrime},
		{4, false, nil, ClausesR},
		{9, false, nil, ClausesWaerden339},
	}

	for _, c := range cases {

		stats := SatStats{
			// Debug: true,
			// Progress: true,
		}
		options := SatOptions{}

		sat, solution := SatAlgorithmD(c.n, c.clauses, &stats, &options)

		if sat != c.sat {
			t.Errorf("expected satisfiable=%t for clauses %v; got %t", c.sat, c.clauses, sat)
			continue
		}
		if sat && !reflect.DeepEqual(solution, c.solution) {
			t.Errorf("expected solution=%v for clauses %v; got %v", c.solution, c.clauses, solution)
			continue
		}
	}
}

func TestSatAlgorithmDFromFile(t *testing.T) {

	// log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		filename     string // file name of the SAT data file
		numVariables int    // number of variable
		numClauses   int    // number of clauses to satisfy
		sat          bool   // is satisfiable
	}{
		{"testdata/SATExamples/L1.sat", 130, 2437, false},
		{"testdata/SATExamples/L2.sat", 273, 1020, false},
		{"testdata/SATExamples/L5.sat", 1472, 102922, true},
		{"testdata/SATExamples/X2.sat", 129, 354, false},
		{"testdata/SATExamples/P3.sat", 144, 529, true},
		{"testdata/SATExamples/P4.sat", 400, 2509, true},
	}

	for _, c := range cases {

		t.Run(c.filename, func(t *testing.T) {

			clauses, variables, err := SatRead(c.filename)

			if err != nil {
				t.Errorf("expected to read file %s; got error %v", c.filename, err)
				return
			}
			if len(variables) != c.numVariables {
				t.Errorf("expected %d variables; got %d", c.numVariables, len(variables))
				return
			}
			if len(clauses) != c.numClauses {
				t.Errorf("expected %d clauses; got %d", c.numClauses, len(clauses))
				return
			}

			stats := SatStats{
				// Debug:    true,
				// Progress: true,
				// Delta:    1000000000,
			}
			options := SatOptions{}

			got, _ := SatAlgorithmD(len(variables), clauses, &stats, &options)

			if got != c.sat {
				t.Errorf("expected satisfiable=%t for filename %s; got %t", c.sat, c.filename, got)
			}
		})
	}
}
