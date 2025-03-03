package taocp

import (
	"bytes"
	"iter"
	"sync"
)

// Visit returns a string representation of the word rectangle
func toString(x []byte, m int, l int) string {
	b := bytes.Buffer{}
	for i := 0; i < l; i += m {
		if i > 0 {
			b.WriteString(":")
		}
		for j := 0; j < m && i+j < l; j++ {
			b.WriteByte(x[i+j] + 97)
		}

	}
	return b.String()
}

// WordRectangles returns m x n word rectangles
//
// # Example 5 x 6 word rectangle
//
// -       n, i ⟶
// - m, j  a b l a t e
// -    ↓  a l a t e d
// - 	   r a p i n g
// - 	   g r i l s e
// - 	   h e n t e d
func WordRectangles(mTrie *CPrefixTrie, nTrie *PrefixTrie,
	max int, initials []byte) iter.Seq[string] {

	return func(yield func(string) bool) {

		// B1 [Initialize.]
		count := 0 // count of returned results

		m := mTrie.Size
		n := nTrie.Size
		mn := m * n

		// Level of the backtrack tree (ie, index into x)
		l := 0

		// a is an m x n lookup table for nTrie nodes corresponding to the prefixes
		// of the first i letters of partial solution for n-length words
		a := make([][]int, m)
		for i := range a {
			a[i] = make([]int, n+1)
		}

		// b is an m x n lookup table for links of letters for m-length words in
		// mTrie
		b := make([][]*Link, m)
		for i := range b {
			b[i] = make([]*Link, n)
		}
		b[0][0] = &mTrie.Nodes[0]

		// Solution tracker
		x := make([]byte, mn)
		x[0] = b[0][0].Letter

		var step byte = 2
		for {
			i, j := l%m, l/m

			switch step {

			case 2: // B2 [Enter level l.]

				if l == mn {
					// Visit x
					if !yield(toString(x, m, l)) {
						return
					}
					count++
					if max > 0 && count == max {
						return
					}
					step = 5
				} else {
					// Set x_l = min D_l
					x[l] = b[i][j].Letter
					step = 3
				}

			case 3: // B3 [Try x_l.]

				skipLetter := false
				if l == 0 && initials != nil {
					// Check if this is an initial letter we should process
					skipLetter = true
					for _, initial := range initials {
						if x[0] == initial {
							skipLetter = false
							break
						}
					}
				}

				// Test if P_l holds
				// ie, Does this possible next letter for m match a prefix for n
				if node := nTrie.Nodes[a[i][j]][b[i][j].Letter]; node != 0 && !skipLetter {
					// Update data structures to facilitate testing P_(l+1)
					a[i][j+1] = node
					if i == m-1 {
						if j < n-1 {
							b[0][j+1] = &mTrie.Nodes[0]
						}
					} else {
						b[i+1][j] = &mTrie.Nodes[b[i][j].Node]
					}

					l++
					step = 2
				} else {
					step = 4
				}

			case 4: // B4 [Try again.]

				// Check if x_l == max D_l
				if link := b[i][j].Right; link.Right != nil {
					// No, set x_l to next larger element of D_l
					b[i][j] = link
					x[l] = link.Letter
					step = 3
				} else {
					step = 5
				}

			case 5: // B5 [Backtrack.]

				l--
				if l < 0 {
					return // all done
				}
				step = 4

			}
		}
	}
}

// MultiWordRectangles iterates over m x n word rectangles, running in n parallel
// threads
func MultiWordRectangles(mTrie *CPrefixTrie, nTrie *PrefixTrie,
	max int, n int, i int) iter.Seq[string] {

	return func(yield func(string) bool) {
		if n < 1 || n > 26 {
			return
		}

		if i < 0 || i > n {
			return
		}

		// Wait for all WordRectangles() threads to complete
		var wg sync.WaitGroup

		// Fan in the output of each WordRectangles() thread
		fanin := func(in iter.Seq[string]) {
			defer wg.Done()
			for s := range in {
				if !yield(s) {
					return
				}
			}
		}

		if i == 0 {
			wg.Add(n)
		} else {
			wg.Add(1)
		}

		// Initial letters for each WordRectangles() thread to process
		// Results are certainly not evenly distributed across initial letters
		initials := [26]byte{}
		for i := 0; i < 26; i++ {
			initials[i] = byte(i)
		}

		chunkSize := (25 + n) / n

		j := 1
		for start := 0; start < 26; start += chunkSize {
			end := start + chunkSize
			if end > 26 {
				end = 26
			}

			if i == 0 || i == j {
				go fanin(WordRectangles(mTrie, nTrie, max, initials[start:end]))
			}

			j++
		}

		wg.Wait()
	}
}
