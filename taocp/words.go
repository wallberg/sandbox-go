package taocp

import (
	"fmt"
	"log"
)

// DoubleWordSquare finds n x n arrays whose rows and columns contain 2n
// different words, using XCC. Each word in words must be the same length.
func DoubleWordSquare(words []string, stats *ExactCoverStats,
	xccOptions *XCCOptions, visit func([]string) bool) {

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
	XCC(items, options, sitems, stats, xccOptions,
		func(solution [][]string) bool {

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
			visit(x)
			return true
		})
}

// WordStair finds word stairs of period p and word length n, a cyclic
// arrangement of words, offset stepwise, that contains 2p distinct words across
// and down. They exist in two varieties, left and right (!left). Each word in words
// must be the same length.
func WordStair(words []string, p int, left bool, stats *ExactCoverStats,
	xccOptions *XCCOptions, visit func([]string) bool) {

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
	XCC(items, options, sitems, stats, xccOptions,
		func(solution [][]string) bool {

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
			visit(x)
			return true
		})
}
