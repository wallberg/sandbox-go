package taocp

import (
	"fmt"
	"log"
)

// DoubleWordSquare finds n x n arrays whose rows and columns contain 2n
// different words, using XCC. Each word in words must be the same length.
func DoubleWordSquare(words []string, stats *ExactCoverStats, visit func([]string) bool) {

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
	for _, word := range words {
		sitems = append(sitems, word+"'")
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

			if i == 1 {
				aOption = append(aOption, word+"'")
				dOption = append(dOption, word+"'")
			}

			options = append(options, aOption, dOption)
		}
	}

	// Get the solutions
	XCC(items, options, sitems, stats, false, false,
		func(solution [][]string) bool {
			// fmt.Println(solution)
			visit([]string{"solution goes here"})
			return true
		})
}
