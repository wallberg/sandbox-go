package taocp

import (
	"slices"
	"testing"
)

// Explore Constraint Satisfication Problems from The Art of Computer Programming, Volume 4,
// Fascicle 7, Constraint Satisfaction, 2024?
//
// §7.2.2.3 Constraint Satisfaction (CSP)

// Exercise_7222_3 expresses CSP of (1) and (2) as a SAT problem.
func TestExercise_7222_3(t *testing.T) {

	const (

		// D1
		x1B = iota + 1
		x1S

		// D2
		x2C
		x2L

		// D3
		x3A
		x3I
		x3U

		// D4
		x4E
		x4O

		// D5
		x5D
		x5N

		// R1
		r11
		r12
		r13

		// R2
		r21
		r22
		r23

		// R3
		r31
		r32
		r33

		last
	)

	stats := SatStats{
		// Debug:     true,
		// Verbosity: 1,
		// Progress:  true,
	}
	options := SatOptions{}

	n := last - 1
	sat := true
	clauses := SatClauses{
		// D1
		{x1B, x1S},
		{-1 * x1B, -1 * x1S},

		// D2
		{x2C, x2L},
		{-1 * x2C, -1 * x2L},

		// D3
		{x3A, x3I, x3U},
		{-1 * x3A, -1 * x3I},
		{-1 * x3A, -1 * x3U},
		{-1 * x3I, -1 * x3U},

		// D4
		{x4E, x4O},
		{-1 * x4E, -1 * x4O},

		// D5
		{x5D, x5N},
		{-1 * x5D, -1 * x5N},

		// R1
		{r11, r12, r13},
		{-1 * r11, x1B},
		{-1 * r11, x3A},
		{-1 * r11, x5N},
		{-1 * r12, x1B},
		{-1 * r12, x3U},
		{-1 * r12, x5D},
		{-1 * r13, x1S},
		{-1 * r13, x3I},
		{-1 * r13, x5N},

		// R2
		{r21, r22},
		{-1 * r21, x1B},
		{-1 * r21, x4E},
		{-1 * r22, x1S},
		{-1 * r22, x4E},
		{-1 * r23, x1S},
		{-1 * r23, x4O},

		// R3
		{r31, r32, r33},
		{-1 * r31, x2C},
		{-1 * r31, x4O},
		{-1 * r31, x5D},
		{-1 * r32, x2C},
		{-1 * r32, x4O},
		{-1 * r32, x5N},
		{-1 * r33, x2L},
		{-1 * r33, x4E},
		{-1 * r33, x5D},
	}

	got, solution := SatAlgorithmA(n, clauses, &stats, &options)

	if got != sat {
		t.Errorf("expected satisfiable=%t, clauses %v; got %t", sat, clauses, got)
	}

	if sat {
		validSolution := SatTest(n, clauses, solution)
		if !validSolution {
			t.Errorf("expected a valid solution, n=%d, clauses=%v; did not get one (solution=%v)", n, clauses, solution)
		}

		BLUED := []int{1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0}
		SCION := []int{0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1}

		if !slices.Equal(solution[0:11], BLUED) && !slices.Equal(solution[0:11], SCION) {

			t.Logf("n=%v, incorrect solution: D1=%v, D2=%v, D3=%v, D4=%v, D5=%v, R1=%v, R2=%v, R3=%v",
				n,
				solution[0:2],   // D1
				solution[2:4],   // D2
				solution[4:7],   // D3
				solution[7:9],   // D4
				solution[9:11],  // D5
				solution[11:14], // R1
				solution[14:17], // R2
				solution[17:20], // R3
			)
		}
	}
}
