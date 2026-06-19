package taocp

import "sort"

// Explore Permutations from The Art of Computer Programming, Volume 4A,
// Combinatorial Algorithms, Part 1, 2011
//
// §7.2.1.2 Generating All Permutations

// Permutations calls the visit function once for every permutation >= x, until
// the permutations are exhausted or visit returns false. The contents of x are
// changed in place.
func Permutations(x []int, visit func() bool) {
	xSort := sort.IntSlice(x)

	another := visit()
	for another && nextPermutation(xSort) {
		another = visit()
	}
}

// nextPermutation generates the next permutation of the sortable collection x
// in lexical order.  It returns false if the permutations are exhausted.
// Algorithm L, p. 319
//
// Take from https://play.golang.org/p/ljft9xhOEn
func nextPermutation(x sort.Interface) bool {
	n := x.Len() - 1
	if n < 1 {
		return false
	}
	j := n - 1
	for ; !x.Less(j, j+1); j-- {
		if j == 0 {
			return false
		}
	}
	l := n
	for !x.Less(j, l) {
		l--
	}
	x.Swap(j, l)
	for k, l := j+1, n; k < l; {
		x.Swap(k, l)
		k++
		l--
	}
	return true
}
