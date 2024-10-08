package taocp

import (
	"fmt"
	"iter"
	"log"

	"codeberg.org/ac/base62"
	"github.com/wallberg/sandbox-go/slice"
)

// DoubleWordSquare iterates over n x n arrays whose rows and columns contain 2n
// different words, using XCC. Each word in words must be the same length.
func DoubleWordSquare(words []string, stats *ExactCoverStats,
	xccOptions *XCCOptions) iter.Seq[[]string] {

	return func(yield func([]string) bool) {

		// Get value of n and verify all words have length n
		n := len(words[0])
		for _, word := range words {
			if len(word) != n {
				log.Fatalf("n=%d, but '%s' has length of %d", n, word, len(word))
			}

		}

		var (
			items   []string   // Primary items
			sitems  []string   // Secondary items
			options [][]string // Options
		)

		// Setup the 2n items
		for i := 1; i <= n; i++ {
			items = append(items, fmt.Sprintf("a%d", i)) // across
			items = append(items, fmt.Sprintf("d%d", i)) // down
		}

		// Setup the secondary items
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				sitems = append(sitems, fmt.Sprintf("%d%d", i, j))
			}
		}
		sitems = append(sitems, words...)
		if xccOptions.Exercise83 {
			for _, word := range words {
				// prime values to remove tranpose solutions
				sitems = append(sitems, word+"'")
			}
		}

		// Setup the 2Wn options
		for _, word := range words {
			for i := 1; i <= n; i++ {
				aOption := []string{fmt.Sprintf("a%d", i)}
				dOption := []string{fmt.Sprintf("d%d", i)}
				for c, char := range word {
					aOption = append(aOption, fmt.Sprintf("%d%d:%c", i, c+1, char))
					dOption = append(dOption, fmt.Sprintf("%d%d:%c", c+1, i, char))
				}
				aOption = append(aOption, word)
				dOption = append(dOption, word)

				if xccOptions.Exercise83 && i == 1 {
					// prime values to remove tranpose solutions
					aOption = append(aOption, word+"'")
					dOption = append(dOption, word+"'")
				}

				options = append(options, aOption, dOption)
			}
		}

		if stats.Debug && stats.Verbosity > 1 {
			log.Print("items", items)
			log.Print("sitems", sitems)
			log.Print("options")
			for _, option := range options {
				log.Print("  ", option)
			}
		}

		// Get the solutions
		for solution, err := range XCC(items, options, sitems, stats, xccOptions) {

			if err != nil {
				log.Fatal(err)
			}

			// Build the solution, a_1 .. a_n, then d_1 .. d_n
			var x []string
			for i := 1; i <= n; i++ {
				a := fmt.Sprintf("a%d", i)
				for _, option := range solution {
					if option[0] == a {
						x = append(x, option[n+1])
						break
					}
				}
			}
			for i := 1; i <= n; i++ {
				d := fmt.Sprintf("d%d", i)
				for _, option := range solution {
					if option[0] == d {
						x = append(x, option[n+1])
						break
					}
				}
			}
			if !yield(x) {
				return
			}
		}
	}
}

// A word stair is a cyclic arrangement of words, offset stepwise, that contains
// 2p distinct words across and down. They exist in two varieties, left and
// right (!left). Each word in words must be the same length.

// WordStair iterates over word stairs of period p and word length n.
func WordStair(words []string, p int, left bool, stats *ExactCoverStats,
	xccOptions *XCCOptions) iter.Seq[[]string] {

	return func(yield func([]string) bool) {

		// Get value of n and verify all words have length n
		n := len(words[0])
		for _, word := range words {
			if len(word) != n {
				log.Fatalf("n=%d, but '%s' has length of %d", n, word, len(word))
			}

		}

		var (
			items   []string   // Primary items
			sitems  []string   // Secondary items
			options [][]string // Options
		)

		// Setup the 2p primary items
		for i := 0; i < p; i++ {
			items = append(items, fmt.Sprintf("a%d", i)) // across
			items = append(items, fmt.Sprintf("d%d", i)) // down
		}

		// Setup the pn + W secondary items
		for i := 0; i < p; i++ {
			for j := 1; j <= n; j++ {
				sitems = append(sitems, fmt.Sprintf("%d%d", i, j))
			}
		}
		sitems = append(sitems, words...)
		if xccOptions.Exercise83 {
			for _, word := range words {
				// prime values to remove tranpose solutions
				sitems = append(sitems, word+"'")
			}
		}

		// Setup the 2Wp options
		for _, word := range words {
			// across
			for i := 0; i < p; i++ {
				aOption := []string{fmt.Sprintf("a%d", i)}
				for c, char := range word {
					aOption = append(aOption, fmt.Sprintf("%d%d:%c", i, c+1, char))
				}
				aOption = append(aOption, word)

				if xccOptions.Exercise83 && i == 0 {
					// prime values to remove tranpose solutions
					aOption = append(aOption, word+"'")
				}

				options = append(options, aOption)
			}

			// down
			for i := 0; i < p; i++ {
				dOption := []string{fmt.Sprintf("d%d", i)}
				for c, char := range word {
					if left {
						// left stair
						dOption = append(dOption, fmt.Sprintf("%d%d:%c", (i+c)%p, c+1, char))
					} else {
						// right stair
						dOption = append(dOption, fmt.Sprintf("%d%d:%c", (i+c)%p, c+1, word[n-c-1]))
					}
				}
				dOption = append(dOption, word)

				if xccOptions.Exercise83 && i == 0 {
					// prime values to remove tranpose solutions
					dOption = append(dOption, word+"'")
				}

				options = append(options, dOption)
			}
		}

		if stats.Debug && stats.Verbosity > 1 {
			log.Print("items", items)
			log.Print("sitems", sitems)
			log.Print("options")
			for _, option := range options {
				log.Print("  ", option)
			}
		}

		// Get the solutions
		for solution, err := range XCC(items, options, sitems, stats, xccOptions) {

			if err != nil {
				log.Fatal(err)
			}

			// Build the solution, a_0 .. a_(p-1), then d_0 .. d_(p-1)
			var x []string
			for i := 0; i < p; i++ {
				a := fmt.Sprintf("a%d", i)
				for _, option := range solution {
					if option[0] == a {
						x = append(x, option[n+1])
						break
					}
				}
			}
			for i := 0; i < p; i++ {
				d := fmt.Sprintf("d%d", i)
				for _, option := range solution {
					if option[0] == d {
						x = append(x, option[n+1])
						break
					}
				}
			}
			if !yield(x) {
				return
			}
		}
	}
}

// WordStairKernel finds word stair kernels for word length n=5, see Exercise
// 7.2.2.1-91
func WordStairKernel(words []string, left bool) ([]string, [][]string, []string) {

	// Get value of n and verify all words have length n
	n := len(words[0])

	if n != 5 {
		log.Fatalf("n=%d, but this function currently only supports n=5", n)
	}

	for _, word := range words {
		if len(word) != n {
			log.Fatalf("n=%d, but '%s' has length of %d", n, word, len(word))
		}
	}

	// Use XCC to build the list of kernels
	var (
		items   []string   // Primary items
		sitems  []string   // Secondary items
		options [][]string // Options
	)

	// Setup the 8 primary items
	items = []string{
		"c4a",
		"c13a",
		"x3a",
		"c9a",
		"x1d",
		"c1d",
		"c3d",
		"c8d",
	}

	// Setup the 14+2+W secondary items, c1-c14,x5,x7 (colored) and every word
	for i := 1; i <= 14; i++ {
		sitems = append(sitems, fmt.Sprintf("c%d", i))
	}
	sitems = append(sitems, "x5", "x7")
	sitems = append(sitems, words...)

	// Setup the options
	for _, word := range words {
		if left {
			// x3 x4 x5 c2 c3
			options = slice.AppendUniqueString(options, []string{"x3a",
				"x5:" + word[2:3], "c2:" + word[1:2], "c3:" + word[0:1]})

			// c4 c5 c6 c7 c8
			options = slice.AppendUniqueString(options, []string{"c4a",
				"c4:" + word[4:5], "c5:" + word[3:4], "c6:" + word[2:3], "c7:" + word[1:2], "c8:" + word[0:1], word})

			// c9 c10 c11 c12 x6
			options = slice.AppendUniqueString(options, []string{"c9a",
				"c9:" + word[4:5], "c10:" + word[3:4], "c11:" + word[2:3], "c12:" + word[1:2]})

			// c13 c14 x7 x8 x9
			options = slice.AppendUniqueString(options, []string{"c13a",
				"c13:" + word[4:5], "c14:" + word[3:4], "x7:" + word[2:3]})
		} else {
			// x3 x4 x5 c2 c3
			options = slice.AppendUniqueString(options, []string{"x3a",
				"x5:" + word[2:3], "c2:" + word[3:4], "c3:" + word[4:5]})

			// c4 c5 c6 c7 c8
			options = slice.AppendUniqueString(options, []string{"c4a",
				"c4:" + word[0:1], "c5:" + word[1:2], "c6:" + word[2:3], "c7:" + word[3:4], "c8:" + word[4:5], word})

			// c9 c10 c11 c12 x6
			options = slice.AppendUniqueString(options, []string{"c9a",
				"c9:" + word[0:1], "c10:" + word[1:2], "c11:" + word[2:3], "c12:" + word[3:4]})

			// c13 c14 x7 x8 x9
			options = slice.AppendUniqueString(options, []string{"c13a",
				"c13:" + word[0:1], "c14:" + word[1:2], "x7:" + word[2:3]})
		}

		// x1 x2 x5 c5 c9
		options = slice.AppendUniqueString(options, []string{"x1d",
			"x5:" + word[2:3], "c5:" + word[3:4], "c9:" + word[4:5]})

		// c1 c2 c6 c10 c13
		options = slice.AppendUniqueString(options, []string{"c1d",
			"c1:" + word[0:1], "c2:" + word[1:2], "c6:" + word[2:3], "c10:" + word[3:4], "c13:" + word[4:5], word})

		// c3 c7 c11 c14 x10
		options = slice.AppendUniqueString(options, []string{"c3d",
			"c3:" + word[0:1], "c7:" + word[1:2], "c11:" + word[2:3], "c14:" + word[3:4]})

		// c8 c12 x7 x11 x12
		options = slice.AppendUniqueString(options, []string{"c8d",
			"c8:" + word[0:1], "c12:" + word[1:2], "x7:" + word[2:3]})
	}

	return items, options, sitems
}

// EncodeCell converts (i, j) int coordinates to base62 string
func EncodeCell(i int, j int) string {
	return base62.Encode(uint32(i)) + base62.Encode(uint32(j))
}

// DecodeCell converts base62 string coordinates to int (i, j)
func DecodeCell(ij string) (int, int, error) {
	var i uint32
	var j uint32
	var err error

	if i, err = base62.Decode(ij[0:1]); err != nil {
		return -1, -1, err
	}

	if j, err = base62.Decode(ij[1:2]); err != nil {
		return -1, -1, err
	}

	return int(i), int(j), nil
}

// WordCross packs a given set of words into an m x n rectangle with conditions.
// Exercise 7.2.2.1-109
func WordCross(words []string, m int, n int) ([]string, [][]string, []string) {

	// Generate the problem as XCC
	var (
		items   []string   // Primary items
		sitems  []string   // Secondary items
		options [][]string // Options
	)

	// Setup the w + m(n-1) + (m-1)n primary items
	for w := range words {
		// words
		items = append(items, fmt.Sprintf("#%d", w))
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if j != n {
				// horizontal edges
				items = append(items, "H"+EncodeCell(i, j))
			}
			if i != m {
				// vertical edges
				items = append(items, "V"+EncodeCell(i, j))
			}
		}
	}

	// Setup the 2mn secondary items
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			// letter or blank in the cell
			sitems = append(sitems, EncodeCell(i, j))
			// 0 - cell is empty, 1 - cell has a letter
			sitems = append(sitems, EncodeCell(i, j)+"'")
		}
	}

	// Setup the options
	for w, word := range words {
		l := len(word)

		// Horizontal word options
		for i := 1; i <= m; i++ {
			for j := 0; j+l <= n; j++ {
				option := []string{fmt.Sprintf("#%d", w)}
				if j > 0 {
					// Leading blank
					option = append(option, EncodeCell(i, j)+":.")
					option = append(option, EncodeCell(i, j)+"':0")

					// Horizontal edge
					option = append(option, "H"+EncodeCell(i, j))
				}

				// Cell contents
				for c := 1; c <= l; c++ {
					option = append(option, EncodeCell(i, j+c)+":"+word[c-1:c])
					option = append(option, EncodeCell(i, j+c)+"':1")

					if c < l {
						// Horizontal edge
						option = append(option, "H"+EncodeCell(i, j+c))
					}
				}

				if j+l < n {
					// Horizontal edge
					option = append(option, "H"+EncodeCell(i, j+l))

					// Trailing blank
					option = append(option, EncodeCell(i, j+l+1)+":.")
					option = append(option, EncodeCell(i, j+l+1)+"':0")
				}
				options = append(options, option)
			}
		}

		// Vertical word options
		// (omit word #0 when m = n)
		if w > 0 || m != n {
			for j := 1; j <= n; j++ {
				for i := 0; i+l <= m; i++ {
					option := []string{fmt.Sprintf("#%d", w)}
					if i > 0 {
						// Leading blank
						option = append(option, EncodeCell(i, j)+":.")
						option = append(option, EncodeCell(i, j)+"':0")

						// Vertical edge
						option = append(option, "V"+EncodeCell(i, j))
					}

					// Cell contents
					for c := 1; c <= l; c++ {
						option = append(option, EncodeCell(i+c, j)+":"+word[c-1:c])
						option = append(option, EncodeCell(i+c, j)+"':1")

						if c < l {
							// Vertical edge
							option = append(option, "V"+EncodeCell(i+c, j))
						}
					}

					if i+l < m {
						// Vertical edge
						option = append(option, "V"+EncodeCell(i+l, j))

						// Trailing blank
						option = append(option, EncodeCell(i+l+1, j)+":.")
						option = append(option, EncodeCell(i+l+1, j)+"':0")
					}
					options = append(options, option)
				}
			}
		}
	}

	// Edge options
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if j != n {
				// horizontal edges
				options = append(options, []string{
					"H" + EncodeCell(i, j),
					EncodeCell(i, j) + "':0",
					EncodeCell(i, j+1) + "':1",
					EncodeCell(i, j) + ":.",
				})
				options = append(options, []string{
					"H" + EncodeCell(i, j),
					EncodeCell(i, j) + "':1",
					EncodeCell(i, j+1) + "':0",
					EncodeCell(i, j+1) + ":.",
				})
				options = append(options, []string{
					"H" + EncodeCell(i, j),
					EncodeCell(i, j) + "':0",
					EncodeCell(i, j+1) + "':0",
					EncodeCell(i, j) + ":.",
					EncodeCell(i, j+1) + ":.",
				})
			}
			if i != m {
				// vertical edges
				options = append(options, []string{
					"V" + EncodeCell(i, j),
					EncodeCell(i, j) + "':0",
					EncodeCell(i+1, j) + "':1",
					EncodeCell(i, j) + ":.",
				})
				options = append(options, []string{
					"V" + EncodeCell(i, j),
					EncodeCell(i, j) + "':1",
					EncodeCell(i+1, j) + "':0",
					EncodeCell(i+1, j) + ":.",
				})
				options = append(options, []string{
					"V" + EncodeCell(i, j),
					EncodeCell(i, j) + "':0",
					EncodeCell(i+1, j) + "':0",
					EncodeCell(i, j) + ":.",
					EncodeCell(i+1, j) + ":.",
				})
			}
		}
	}

	return items, options, sitems
}
