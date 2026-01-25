package taocp

import (
	"iter"
)

// StrictPartitions returns all integer partitions of n, with distinct values,
// with partition size between min and max (inclusive).
// Not a TAOCP implementation.
func StrictPartitions(n int, min int, max int) iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		var partition []int
		var generate func(int, int)

		if n <= 0 {
			return
		}

		generate = func(remaining int, lastValue int) {
			if remaining == 0 {
				if len(partition) >= min && len(partition) <= max {
					p := make([]int, len(partition))
					copy(p, partition)
					if !yield(p) {
						return
					}
				}
				return
			}

			for value := lastValue - 1; value >= 1; value-- {
				if value > remaining {
					continue
				}
				partition = append(partition, value)
				generate(remaining-value, value)
				partition = partition[:len(partition)-1]
			}
		}

		generate(n, n+1)
	}
}
