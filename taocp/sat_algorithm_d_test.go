package taocp

import (
	"fmt"
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

			sat, solution := SatAlgorithmD(len(variables), clauses, &stats, &options)

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

func TestSatAlgorithmDLangford(t *testing.T) {

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

			sat, solution := SatAlgorithmD(len(coverOptions), clauses, &stats, &options)

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

// TestSatAlgorithmDSat3 tests Sat3() using Algorithm D.
func TestSat3(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	cases := []struct {
		n        int        // number of strictly distinct literals
		sat      bool       // is satisfiable
		solution []int      // solution
		clauses  SatClauses // clauses to satisfy
	}{
		{4, true, []int{0, 0, 1, 0}, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3, 4}, {1, -4}}},
		{4, false, nil, SatClauses{{1, -2}, {2, 3}, {-1, -3}, {-1, -2, 3}, {1, 2, -3, 4}, {-4}}},
		{4, true, []int{0, 1, 0, 1}, ClausesRPrime},
		{4, false, nil, ClausesR},
		{9, false, nil, ClausesWaerden339},
	}

	for _, c := range cases {

		sat3, n3, clauses3 := Sat3(c.n, c.clauses)

		if !sat3 {
			if n3 <= c.n {
				t.Errorf("expected number of SAT3 variables for filename to be greater than %d; got %d", c.n, n3)
			}
			if len(clauses3) <= len(c.clauses) {
				t.Errorf("expected number of SAT3 clauses for filename to be greater than %d; got %d", len(c.clauses), len(clauses3))
			}
		}

		stats := SatStats{}
		options := SatOptions{}

		sat, solution := SatAlgorithmD(n3, clauses3, &stats, &options)
		if solution != nil {
			solution = solution[0:c.n]
		}

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

func BenchmarkSatAlgorithmDFromFile(b *testing.B) {

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

				sat, _ := SatAlgorithmD(len(variables), clauses, &stats, &options)

				if firstExecution {
					b.Logf("SAT=%t, n=%d, m=%d, nodes=%d", sat, len(variables), len(clauses), stats.Nodes)
					firstExecution = false
				}
			}

		})
	}
}

func BenchmarkSatAlgorithmDLangford(b *testing.B) {

	for _, n := range []int{5, 9, 13} {

		firstExecution := true

		clauses, coverOptions := SatLangford(n)

		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				stats := SatStats{}
				options := SatOptions{}

				sat, _ := SatAlgorithmD(len(coverOptions), clauses, &stats, &options)

				if firstExecution {
					b.Logf("SAT=%t, n=%d, m=%d, nodes=%d", sat, len(coverOptions), len(clauses), stats.Nodes)
					firstExecution = false
				}
			}

		})
	}
}
