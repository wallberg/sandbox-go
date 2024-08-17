package taocp

import "iter"

// Explore Boolean Basics from The Art of Computer Programming, Volume 4a,
// Combinatorial Algorithms, Part 1, 2011
//
// §7.1.1 Boolean Basics

// BitPairs iterates over pairs of indexes of subcubes in v which have all the same
// bit values except at bit postion j.
//
// Exercise 29.
func BitPairs(v []int, j int) iter.Seq2[int, int] {

	return func(yield func(k, kp int) bool) {
		// B1. [Initialize.]
		m := len(v)

		k, kp := 0, 0

		for {
			// B2. [Find a zero.]
			for {
				if k == m {
					return
				}
				if v[k]&(1<<j) == 0 {
					break
				}
				k++
			}

			// B3. [Make k-prime > k.]
			if kp <= k {
				kp = k + 1
			}

			// B4. [Advance k-prime.]
			for {
				if kp == m {
					return
				}
				if v[kp] >= v[k]+(1<<j) {
					break
				}
				kp++
			}

			// B5. [Skip past a big mismatch.]
			if v[k]^v[kp] >= 1<<(j+1) {
				k = kp
				continue // Goto B2
			}

			// B6. [Record a match.]
			if v[kp] == v[k]+(1<<j) {
				if !yield(k, kp) {
					return
				}
			}

			// B7. [Advance k.]
			k++
		}
	}
}

// MaximalSubcubes iterates over all maximal subcubes (aka prime implicant) of v.
// Subcubes are represented as tuples (a, b) where a records the position of
// the asterisks and b records the bits in non-* positions. n is the size in
// bits of bitstrings in v. v is the list of implicants (sorted in ascending
// order) represented as integer bitstrings.
//
// Exercise 30.
func MaximalSubcubes(n int, v []int) iter.Seq2[int, int] {

	return func(yield func(a, b int) bool) {
		// P1. [Initialize.]
		m := len(v)

		// The current value of a being processed
		A := 0

		// Stack S contains |A| + 1 lists of subcubes; each list contains subcubes
		// with the same a value, in increasing order of a. This includes all lists
		// with wildcards plus the first list with a=0, equivalent to the input v
		// list of bitstrings
		S := make([]int, 2*m+n)

		// Tag bits indicating a matching j-buddy, for each corresponding subcube
		// in S
		T := make([]int, 2*m+n)

		// Determine the j-buddy pairs for the initial subcube list
		for j := 0; j < n; j++ {
			for k, kp := range BitPairs(v, j) {
				T[k] |= (1 << j)
				T[kp] |= (1 << j)
			}
		}

		// For each subcube, either output it as maximal or advance it
		// (with j-buddy) to the next subcube list, with additional wildcard
		r, s, t := 0, 0, 0
		for s < m {
			if T[s] == 0 {
				if !yield(0, v[s]) {
					return
				}
			} else {
				S[t] = v[s]
				T[t] = T[s]
				t++
			}
			s++
		}
		S[t] = 0

		for {
			// P2. [Advance A.]
			j := 0
			if S[t] == t { // the topmost list is empty
				for j < n && A&(1<<j) == 0 {
					j++
				}
			}
			for j < n && A&(1<<j) != 0 {
				t = S[t] - 1
				A -= (1 << j)
				j++
			}
			if j >= n {
				return
			}
			A += (1 << j)

			// P3. [Generate list A.]
			r, s = t, S[t]
			for k, kp := range BitPairs(S[s:r], j) {
				x := (T[s+k] & T[s+kp]) - (1 << j)
				if x == 0 {
					if !yield(A, S[s+k]) {
						return
					}
				} else {
					t++
					S[t] = S[s+k]
					T[t] = x
				}
			}
			t++
			S[t] = r + 1
		}
	}
}
