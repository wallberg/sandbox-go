package sortx

import (
	"sort"
)

// Provide sort functionality which extends the core sort package

// InsertInt inserts an int into an already sorted slice of ints, if not in the
// slice already, and returns its index location. Expand the underlying array
// if necessary.
func InsertInt(values *[]int, value int) int {
	// Search for the insertion point
	i := sort.SearchInts(*values, value)

	// Check if this value is already in the slice
	n := len(*values)
	if i == n || (*values)[i] != value {

		// No, insert into the slice
		if n < cap(*values) {

			// Expand the slice
			*values = (*values)[:n+1]
			copy((*values)[i+1:n+1], (*values)[i:n])

		} else {

			// The slice is full, expand the array
			newCap := n
			if newCap == 0 {
				newCap = 8
			} else {
				newCap *= 2
			}
			newValues := make([]int, newCap)[0 : n+1]
			copy(newValues[0:i], (*values)[0:i])
			copy(newValues[i+1:n+1], (*values)[i:n])
			(*values) = newValues
		}

		(*values)[i] = value
	}

	return i
}
