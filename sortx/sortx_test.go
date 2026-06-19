package sortx

import (
	"math/rand"
	"testing"
)

func TestInsertInt(t *testing.T) {

	n := 1000

	values := make([]int, 0)

	// Insert new values
	for _, value := range rand.Perm(n) {
		InsertInt(&values, value)
	}

	// Insert existing values
	for _, value := range []int{0, 1, n - 1, n} {
		InsertInt(&values, value)
	}

	for i := 0; i < n; i++ {
		if values[i] != i {
			t.Errorf("values[%d] = %d; want %d", i, i, i)
		}
	}
}

func BenchmarkInsertInt(b *testing.B) {
	cases := []struct {
		name string
		size int
	}{
		{"10^3", 1000},
		{"10^4", 10000},
		{"10^5", 100000},
		// {"10^6", 1000000}, # too slow
	}

	for _, c := range cases {

		b.Run(c.name, func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				values := make([]int, 0)

				// Insert new values
				for _, value := range rand.Perm(c.size) {
					InsertInt(&values, value)
				}
			}
		})
	}
}
