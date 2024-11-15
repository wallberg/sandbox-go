package taocp

import (
	"fmt"
	"log"
	"testing"
)

func TestSatAlgorithmL(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		n          int        // number of strictly distinct literals
		sat        bool       // is satisfiable
		bigClauses bool       // use BigClauses option?
		clauses    SatClauses // clauses to satisfy
	}{
		{1, true, false, SatClauses{{1}}},
		{1, true, false, SatClauses{{-1}}},
		{1, false, false, SatClauses{{1}, {-1}}},
		{2, true, false, SatClauses{{1}, {2}}},
		{3, true, false, SatClauses{{1}, {2}, {-3}}},
		{3, true, false, SatClauses{{-1}, {2}, {3}}},
		{2, true, false, SatClauses{{1, 2}}},
		{2, true, false, SatClauses{{1, 2}, {1, -2}}},
		{2, true, false, SatClauses{{-1, -2}}},
		{2, false, false, SatClauses{{1, 2}, {-1, -2}, {1, -2}, {-1, 2}}},
		{2, true, false, SatClauses{{-1, 2}, {1, -2}}},
		{2, true, false, SatClauses{{1, -2}, {-1, 2}}},
		{5, true, false, SatClauses{{1, -2}, {2, 2}, {-1, 3}, {2, 4}, {-4, 5}}},
		{5, true, false, SatClauses{
			{1, 2}, {2, 3}, {3, 4}, {4, 5},
			{-1, -2}, {-1, -3}, {-1, -4}, {-1, -5}}},
		{100, true, false, SatRand(2, 80, 100, 0)},
		{100, true, false, SatRand(2, 100, 100, 0)},
		{100, false, false, SatRand(2, 400, 100, 0)},
		{1000, true, false, SatRand(2, 1000, 1000, 0)},
		{1000, true, false, SatRand(2, 1100, 1000, 0)},
		{1000, false, false, SatRand(2, 2000, 1000, 0)},
		{100, true, false, SatRand(3, 420, 100, 0)},
		{100, false, false, SatRand(3, 500, 100, 0)},
		{3, true, false, SatClauses{{1, 2, 3}}},
		{3, true, false, SatClauses{{-1, -2, 3}}},
		{3, true, false, SatClauses{{1, -2}, {2, 3}, {-1, -2, 3}}},
		{3, true, false, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}}},
		{3, false, false, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}, {1, 2, -3}}},
		{4, true, false, ClausesRPrime},
		{4, false, false, ClausesR},
		{8, true, false, SatWaerdan(3, 3, 8)},
		{9, false, false, SatWaerdan(3, 3, 9)},
		{10, false, false, SatWaerdan(3, 3, 10)},
		{3, false, true, SatComplete(3)},
		{4, false, true, SatComplete(4)},
		{5, false, true, SatComplete(5)},
		{8, true, true, SatWaerdan(3, 3, 8)},
		{9, false, true, SatWaerdan(3, 3, 9)},
		{10, false, true, SatWaerdan(3, 3, 10)},
	}

	for i, c := range cases {

		stats := SatStats{
			// Debug:     true,
			// Verbosity: 1,
			// Progress:  true,
		}
		options := SatOptions{}
		optionsL := NewSatAlgorithmLOptions()
		optionsL.CompensationResolvants = false
		optionsL.SuppressBigClauses = !c.bigClauses

		// var clausesStr string
		// if len(c.clauses) < 10 {
		// 	clausesStr = fmt.Sprintf("%v", c.clauses)
		// } else {
		// 	clausesStr = fmt.Sprintf("#%d", len(c.clauses))
		// }

		// t.Logf("Executing test case #%d, n=%d, sat=%t, bigc=%t, clauses=%s", i, c.n, c.sat, c.bigClauses, clausesStr)

		sat, solution := SatAlgorithmL(c.n, c.clauses, &stats, &options, optionsL)

		if sat != c.sat {
			t.Errorf("expected satisfiable=%t for case %d, clauses %v; got %t", c.sat, i, c.clauses, sat)
		}
		if sat {
			validSolution := SatTest(c.n, c.clauses, solution)
			if !validSolution {
				t.Errorf("expected a valid solution for case %d, n=%d, clauses=%v; did not get one (solution=%v)", i, c.n, c.clauses, solution)
			}
		}

	}
}

func TestSatAlgorithmLFromFile(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

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
		// {"testdata/SATExamples/P4.sat", 400, 2509, true}, // (long, see 7.2.2.2 p. 304)
	}

	for _, c := range cases {

		t.Run(c.filename, func(t *testing.T) {
			t.Parallel()

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
				// Debug: true,
				// Verbosity: 1,
				// Progress: true,
				// Delta:    100000,
			}
			options := SatOptions{}
			optionsL := NewSatAlgorithmLOptions()

			sat, solution := SatAlgorithmL(len(variables), clauses, &stats, &options, optionsL)

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
			t.Parallel()

			stats := SatStats{
				// Debug: true,
				// Progress: true,
			}
			options := SatOptions{}
			optionsL := NewSatAlgorithmLOptions()
			optionsL.CompensationResolvants = true

			expected := false
			if n%4 == 0 || n%4 == 3 {
				expected = true
			}

			clauses, coverOptions := SatLangford(n)

			sat, solution := SatAlgorithmL(len(coverOptions), clauses, &stats, &options, optionsL)

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
		// {"testdata/SATExamples/P4.sat", 400, 2509, true}, // (long, see 7.2.2.2 p. 304)
	}

	for _, c := range cases {

		firstExecution := true

		clauses, variables, _ := SatRead(c.filename)

		b.Run(c.filename, func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				stats := SatStats{}
				options := SatOptions{}
				optionsL := NewSatAlgorithmLOptions()
				optionsL.CompensationResolvants = true

				sat, _ := SatAlgorithmL(len(variables), clauses, &stats, &options, optionsL)

				if firstExecution {
					b.Logf("SAT=%t, n=%d, m=%d, nodes=%d", sat, len(variables), len(clauses), stats.Nodes)
					firstExecution = false
				}
			}

		})
	}
}

func BenchmarkSatAlgorithmLLangford(b *testing.B) {

	for _, n := range []int{5, 9, 10} {

		firstExecution := true

		clauses, coverOptions := SatLangford(n)

		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				stats := SatStats{}
				options := SatOptions{}
				optionsL := NewSatAlgorithmLOptions()

				sat, _ := SatAlgorithmL(len(coverOptions), clauses, &stats, &options, optionsL)

				if firstExecution {
					b.Logf("SAT=%t, n=%d, m=%d, nodes=%d", sat, len(coverOptions), len(clauses), stats.Nodes)
					firstExecution = false
				}
			}

		})
	}
}

func BenchmarkSatAlgorithmLSatRandom(b *testing.B) {

	cases := []struct {
		k int // clause length (k-SAT)
		m int // number of clauses
		n int // number of strictly distinct literals
	}{
		{2, 80, 100},
		{2, 100, 100},
		{2, 400, 100},
		{2, 1000, 1000},
		{2, 1100, 1000},
		{2, 2000, 1000},
		{3, 420, 100},
		{3, 400, 50},
		{3, 500, 50},
		{3, 1000, 50},
	}

	for _, c := range cases {

		firstExecution := true

		clauses := SatRand(c.k, c.m, c.n, 0)

		b.Run(fmt.Sprintf("k=%d,m=%d,n=%d", c.k, c.m, c.n), func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				stats := SatStats{}
				options := SatOptions{}
				optionsL := NewSatAlgorithmLOptions()

				sat, _ := SatAlgorithmL(c.n, clauses, &stats, &options, optionsL)

				if firstExecution {
					b.Logf("SAT=%t, m=%d, n=%d, nodes=%d", sat, c.m, c.n, stats.Nodes)
					firstExecution = false
				}
			}

		})
	}
}
