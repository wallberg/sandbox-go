package taocp

// Explore Boolean Basics from The Art of Computer Programming, Volume 4a,
// Combinatorial Algorithms, Part 1, 2011
//
// §7.1.1 Boolean Basics

// BitPairs returns pairs of indexes of subcubes in v which have all the same
// bit values except at bit postion j.  Pairs are returned sequentially on the
// out channel.
//
// Exercise 29.
func BitPairs(v []int, j int, out chan int) {

	defer close(out)

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
			out <- k
			out <- kp
		}

		// B7. [Advance k.]
		k++
	}
}
