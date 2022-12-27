package taocp

import (
	"fmt"
	"log"
	"testing"
)

func TestSatAlgorithmL(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		n       int        // number of strictly distinct literals
		sat     bool       // is satisfiable
		clauses SatClauses // clauses to satisfy
	}{
		{1, true, SatClauses{{1}}},
		{1, true, SatClauses{{-1}}},
		{1, false, SatClauses{{1}, {-1}}},
		{2, true, SatClauses{{1}, {2}}},
		{3, true, SatClauses{{1}, {2}, {-3}}},
		{3, true, SatClauses{{-1}, {2}, {3}}},
		{2, true, SatClauses{{1, 2}}},
		{2, true, SatClauses{{1, 2}, {1, -2}}},
		{2, true, SatClauses{{-1, -2}}},
		{2, false, SatClauses{{1, 2}, {-1, -2}, {1, -2}, {-1, 2}}},
		{2, true, SatClauses{{-1, 2}, {1, -2}}},
		{2, true, SatClauses{{1, -2}, {-1, 2}}},
		{5, true, SatClauses{{1, -2}, {2, 2}, {-1, 3}, {2, 4}, {-4, 5}}},
		{5, true, SatClauses{
			{1, 2}, {2, 3}, {3, 4}, {4, 5},
			{-1, -2}, {-1, -3}, {-1, -4}, {-1, -5}}},
		// {3, true, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}}},
		// {3, false, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}, {1, 2, -3}}},
		// {4, true, ClausesRPrime},
		// {4, false, ClausesR},
		// {9, false, ClausesWaerden339},
	}

	for i, c := range cases {

		stats := SatStats{
			Debug:     true,
			Verbosity: 1,
			// Progress:  true,
		}
		options := SatOptions{}

		t.Logf("Executing test case #%d, c=%v", i, c)

		sat, solution := SatAlgorithmL(c.n, c.clauses, &stats, &options)

		if sat != c.sat {
			t.Errorf("expected satisfiable=%t for case %d, clauses %v; got %t", c.sat, i, c.clauses, sat)
		}
		if sat {
			validSolution := SatTest(c.n, c.clauses, solution)
			if !validSolution {
				t.Errorf("expected a valid solution for n=%d, clauses=%v; did not get one (solution=%v)", c.n, c.clauses, solution)
			}
		}

	}
}

func TestSatAlgorithmLFromFile(t *testing.T) {

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

			sat, solution := SatAlgorithmL(len(variables), clauses, &stats, &options)

			if sat != c.sat {
				t.Errorf("expected satisfiable=%t for filename %s; got %t", c.sat, c.filename, sat)
			} else if sat {
				validSolution := SatTest(c.numVariables, clauses, solution)
				if !validSolution {
					t.Errorf("expected a valid solution for filename %s; did not get one", c.filename)
				}
			}
		})
	}
}

func TestSatAlgorithmLLangford(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	for n := 2; n <= 12; n++ {

		t.Run(fmt.Sprintf("langford(%d)", n), func(t *testing.T) {
			stats := SatStats{
				// Debug: true,
				// Progress: true,
			}
			options := SatOptions{}

			expected := false
			if n%4 == 0 || n%4 == 3 {
				expected = true
			}

			clauses, coverOptions := SatLangford(n)

			sat, solution := SatAlgorithmL(len(coverOptions), clauses, &stats, &options)

			if sat != expected {
				t.Errorf("expected langford(%d) satisfiable=%t; got %t", n, expected, sat)
			} else if sat {
				validSolution := SatTest(len(coverOptions), clauses, solution)
				if !validSolution {
					t.Errorf("expected a valid solution for langford(%d); did not get one", n)
				}
			}
		})
	}
}

func BenchmarkSatAlgorithmLFromFile(b *testing.B) {

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

		firstExecution := true

		clauses, variables, _ := SatRead(c.filename)

		b.Run(c.filename, func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				stats := SatStats{}
				options := SatOptions{}

				sat, _ := SatAlgorithmL(len(variables), clauses, &stats, &options)

				if firstExecution {
					b.Logf("SAT=%t, n=%d, m=%d, nodes=%d", sat, len(variables), len(clauses), stats.Nodes)
					firstExecution = false
				}
			}

		})
	}
}

func BenchmarkSatAlgorithmLLangford(b *testing.B) {

	for _, n := range []int{5, 9, 13} {

		firstExecution := true

		clauses, coverOptions := SatLangford(n)

		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				stats := SatStats{}
				options := SatOptions{}

				sat, _ := SatAlgorithmL(len(coverOptions), clauses, &stats, &options)

				if firstExecution {
					b.Logf("SAT=%t, n=%d, m=%d, nodes=%d", sat, len(coverOptions), len(clauses), stats.Nodes)
					firstExecution = false
				}
			}

		})
	}
}
