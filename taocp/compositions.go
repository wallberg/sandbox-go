package taocp

import (
	"fmt"
	"iter"
)

// Explore Compositions from The Art of Computer Programming, Volume 4A,
// Combinatorial Algorithms, Part 1, 2011
//
// §7.2.1.1 Generating All n-Tuples, Exercise 12

// Compositions iterates over each composition of n, until
// the permutations are exhausted or visit returns false.
func Compositions(n int) iter.Seq[[]int] {

	return func(yield func([]int) bool) {

		if n <= 1 {
			panic(fmt.Sprintf("Expected n > 1; got %d", n))
		}

		// C1. [Initialize.]
		t := 1
		s := make([]int, n+1)
		s[1] = n

	C2:
		// C2. [Visit.]
		if !yield(s[1 : t+1]) {
			return
		}
		if t%2 == 0 {
			goto C4
		}

		// C3. [Odd step.]
		if s[t] > 1 {
			s[t] -= 1
			s[t+1] = 1
			t += 1
		} else {
			t -= 1
			s[t] += 1
		}
		goto C2

	C4:
		// C4. [Even step.]
		if t == 1 {
			return
		}
		if s[t-1] > 1 {
			s[t-1] -= 1
			s[t+1] = s[t]
			s[t] = 1
			t += 1
		} else {
			t -= 1
			if t == 1 {
				return
			}
			s[t] = s[t+1]
			s[t-1] += 1
		}
		goto C2
	}
}
