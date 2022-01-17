package taocp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Explore Satisfiability from The Art of Computer Programming, Volume 4,
// Fascicle 6, Satisfiability, 2015
//
// §7.2.2.2 Satisfiability (SAT)

// SatStats is a struct for tracking SAT statistics and reporting
// runtime progress
type SatStats struct {
	// Input parameters
	Progress     bool // Display runtime progress
	Debug        bool // Enable debug logging
	Verbosity    int  // Debug verbosity level (0 or 1)
	Delta        int  // Display progress every Delta number of Nodes
	SuppressDump bool // Don't display the dump()

	// Statistics collectors
	MaxLevel  int   // Maximum level reached
	Theta     int   // Display progress at next Theta number of Nodes
	Levels    []int // Count of times each level is entered
	Nodes     int   // Count of nodes processed
	Solutions int   // Count of solutions returned
}

// SatOptions provides SAT runtime options
type SatOptions struct {
}

// String returns a String representation of type SATStats struct
func (s SatStats) String() string {
	// Find first non-zero level count
	i := len(s.Levels)
	for s.Levels[i-1] == 0 && i > 1 {
		i--
	}

	return fmt.Sprintf("nodes=%d, solutions=%d, levels=%v", s.Nodes,
		s.Solutions, s.Levels[:i])
}

// SatClause represents a single clause
type SatClause []int

// SatClauses represents a list of clauses
type SatClauses []SatClause

// AppendUniqueSatClause inserts x into a, if not already present.  Returns new value of a, a la append()
func AppendUniqueSatClause(a SatClauses, x SatClause) SatClauses {

	// Iterate over the existing values of a
	for _, y := range a {
		// Check if their lengths are identical
		if len(x) == len(y) {
			// Check if their values are identical
			for i, v := range x {
				if v != y[i] {
					// Not identical
					goto END
				}
			}
			// lengths are the identical and all values are identical
			// return a unchanged
			return a
		}
	END:
		// Not identical, continue to next value of y
	}

	return append(a, x)
}

// SatRead reads a SAT file in Knuth format and returns
// a list of clauses along with the mapping of variables
// (numeric to string name)
func SatRead(filename string) (SatClauses, map[int]string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening %s for reading: %v", filename, err)
	}

	var clauses SatClauses
	variable2name := make(map[int]string)
	name2variable := make(map[string]int)
	nextVariable := 1 // next literal to use

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if a comment line
		if !strings.HasPrefix(line, "~ ") {
			// Create a new clause
			var clause SatClause

			// Iterate over the literals of the clause
			for _, name := range strings.Fields(line) {

				// Determine if the literal is negated
				sign := 1
				if strings.HasPrefix(name, "~") {
					name = name[1:]
					sign = -1
				}

				// Determine the variable number for this name
				var found bool
				var variable int
				if variable, found = name2variable[name]; !found {
					variable = nextVariable
					name2variable[name] = variable
					variable2name[variable] = name
					nextVariable += 1
				}

				literal := sign * variable

				clause = append(clause, literal)
			}
			clauses = append(clauses, clause)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error scanning %s: %v", filename, err)
	}

	return clauses, variable2name, nil
}

// SatWaerdan returns the SAT clauses for waerden(j,k;n)
func SatWaerdan(j, k, n int) SatClauses {
	var clauses SatClauses

	for d := 1; n > (j-1)*d; d++ {
		for i := 1; i <= n-(j-1)*d; i++ {
			var clause SatClause
			for v := i; v <= i+(j-1)*d; v += d {
				clause = append(clause, v)
			}
			clauses = append(clauses, clause)
		}
	}

	for d := 1; n > (k-1)*d; d++ {
		for i := 1; i <= n-(k-1)*d; i++ {
			var clause SatClause
			for v := i; v <= i+(k-1)*d; v += d {
				clause = append(clause, -1*v)
			}
			clauses = append(clauses, clause)
		}
	}

	return clauses
}

type LangfordOption struct {
	D  int // digit (1..n)
	S1 int // slot1 in the sequence (1..2n)
	S2 int // slot2 in the sequence (1..2n)
}

// SatLangford returns the SAT clauses for Langford pairs, langford(n), along with the
// set of exact covering options corresponding to each SAT variable
func SatLangford(n int) (clauses SatClauses, options []LangfordOption) {

	symmetric2sat := func(symmetric []int) (clauses [][]int) {
		clauses = append(clauses, symmetric)
		for j := 0; j < len(symmetric); j++ {
			for k := j + 1; k < len(symmetric); k++ {
				clauses = append(clauses, []int{-1 * symmetric[j], -1 * symmetric[k]})
			}
		}
		return clauses
	}

	// Generate the exact covering options
	for d := 1; d <= n; d++ {
		s1 := 1
		s2 := s1 + d + 1
		for s2 <= 2*n {
			// Add the option, but skip some options to prevent symmetric results
			// See Exercise 7.2.2.1-15
			x := 0
			if n%2 == 0 {
				x = 1
			}
			if d != n-x || s1 <= n/2 {
				options = append(options, LangfordOption{d, s1, s2})
			}
			s1++
			s2++
		}
	}

	// Generate the symmetric function for the n digits and 2*n slots
	var symmetrics [][]int
	for d := 1; d <= n; d++ {
		var symmetric []int
		for i, option := range options {
			// Check if this digit is in this option
			if d == option.D {
				symmetric = append(symmetric, i+1)
			}
		}
		symmetrics = append(symmetrics, symmetric)
	}
	for s := 1; s <= 2*n; s++ {
		var symmetric []int
		for i, option := range options {
			// Check if this digit is in this option
			if s == option.S1 || s == option.S2 {
				symmetric = append(symmetric, i+1)
			}
		}
		symmetrics = append(symmetrics, symmetric)
	}

	// Express the symmetric functions as AND of OR SAT clauses
	for _, symmetric := range symmetrics {
		for _, clause := range symmetric2sat(symmetric) {
			// Append the clause to clauses, if not already present
			clauses = AppendUniqueSatClause(clauses, clause)
		}
	}

	return clauses, options
}

// SatMaxR generates new variables and clauses to ensure that x_1 + ... + x_n is at most r,
// which is S_(<= r). (n - r)r new variables will be created, beginning at startV.
// Returns the list of new clauses and the number of variables created,
// (startV,...,startV+numV-1). As a special case, if n >= 4 and r = n-1, then
// n-2 clauses of length 3 will be created with n-3 variables.
func SatMaxR(r int, clause SatClause, startV int) (newclauses SatClauses, numV int) {
	n := len(clause)

	// Special case, n >= 4 and r = n - 1
	if n >= 4 && r == n-1 {
		numV = n - 3

		for j := 1; j < n-1; j++ {
			var v1, v2, v3 int
			if j == 1 {
				// the first clause
				v1 = clause[j-1] * -1
				v2 = clause[j] * -1
				v3 = startV + j - 1
			} else if j == n-2 {
				// the last clause
				v1 = clause[j] * -1
				v2 = clause[j+1] * -1
				v3 = (startV + j - 2) * -1
			} else {
				// a middle clause
				v1 = clause[j] * -1
				v2 = (startV + j - 2) * -1
				v3 = startV + j - 1
			}
			newclauses = append(newclauses, SatClause{v1, v2, v3})
		}

		return newclauses, numV
	}

	// General case
	numV = (n - r) * r

	for j := 1; j < n-r; j++ {
		for k := 1; k <= r; k++ {
			v1 := (startV - 1 + (j-1)*r + k) * -1
			v2 := startV - 1 + j*r + k
			// fmt.Printf("%d [%d %d], %d [%d %d]\n", v1, k, j, v2, k, j+1)
			newclauses = append(newclauses, SatClause{v1, v2})
		}
	}

	for j := 1; j <= n-r; j++ {
		for k := 0; k <= r; k++ {
			v1 := clause[j+k-1] * -1
			v2 := (startV - 1 + (j-1)*r + k) * -1
			v3 := startV - 1 + (j-1)*r + k + 1

			if k == 0 {
				// fmt.Printf("%d, %d [%d %d]\n", v1, v3, k+1, j)
				newclauses = append(newclauses, SatClause{v1, v3})
			} else if k == r {
				// fmt.Printf("%d, %d [%d %d]\n", v1, v2, k, j)
				newclauses = append(newclauses, SatClause{v1, v2})
			} else {
				// fmt.Printf("%d, %d [%d %d], %d [%d %d]\n", v1, v2, k, j, v3, k+1, j)
				newclauses = append(newclauses, SatClause{v1, v2, v3})
			}
		}
	}

	return newclauses, numV
}
