package slice

// package slice provides helpful slice utilities

// FindString returns the smallest index i at which x == a[i],
// or -1 if there is no such index.
func FindString(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

// IsCycleString checks a and b contain the same cycle of strings
func IsCycleString(a []string, b []string) bool {

	if len(a) == 0 && len(b) == 0 {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	// Find the first word of got in c.solution
	j := FindString(b, a[0])
	if j == -1 {
		return false
	}

	n := len(a)
	for x := 1; x < n; x++ {
		if a[x] != b[(j+x)%n] {
			return false
		}
	}

	return true

}

// AppendUniqueString inserts x into a, if not already present.  Returns new value of a, a la append()
func AppendUniqueString(a [][]string, x []string) [][]string {

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
