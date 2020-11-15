package taocp

// Explore Boolean Basics from The Art of Computer Programming, Volume 4a,
// Combinatorial Algorithms, Part 1, 2011
//
// §7.1.1 Boolean Basics

// BitPairs returns pairs of indexes of subcubes in v which have all the same
// bit values except at bit postion j.  Pairs are returned are by calling
// the visit function once for each pair.
//
// Exercise 29.
func BitPairs(v []int, j int, visit func(k int, kp int)) {

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
			visit(k, kp)
		}

		// B7. [Advance k.]
		k++
	}
}

// MaximalSubcubes returns all maximal subcubes (aka prime implicant) of v.
// Subcubes are represented as tuples (a, b) where a records the position of
// the asterisks and b records the bits in non-* positions. n is the size in
// bits of bitstrings in v. v is the list of implicants (sorted in ascending
// order) represented as integer bitstrings.
//
// Exercise 30.
func MaximalSubcubes(n int, v []int, visit func(a int, b int)) {

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
		BitPairs(v, j, func(k int, kp int) {
			T[k] |= (1 << j)
			T[kp] |= (1 << j)
		})
	}
	// For each subcube, either output it as maximal or advance it
	// (with j-buddy) to the next subcube list, with additional wildcard
	r, s, t := 0, 0, 0
	for s < m {
		if T[s] == 0 {
			visit(0, v[s])
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
		BitPairs(S[s:r], j, func(k int, kp int) {
			x := (T[s+k] & T[s+kp]) - (1 << j)
			if x == 0 {
				visit(A, S[s+k])
			} else {
				t++
				S[t] = S[s+k]
				T[t] = x
			}
		})
		t++
		S[t] = r + 1
	}
}
