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

// A word stair is a cyclic arrangement of words, offset stepwise, that contains
// 2p distinct words across and down. They exist in two varieties, left and
// right (!left). Each word in words must be the same length.

// WordStair finds word stairs of period p and word length n.
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

// WordStairKernel finds word stair kernels for word length n=5, see Exercise
// 7.2.2.1-91
func WordStairKernel(words []string, left bool, stats *ExactCoverStats,
	visit func(string) bool) {

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

	if left {
		log.Fatalf("left=%t, but this function currently only supports left=false", left)
	}

	// Use XCC to build the list of kernels
	var (
		items   []string   // Primary items
		sitems  []string   // Secondary items
		options [][]string // Options
	)

	xccOptions := &XCCOptions{
		Exercise83: true, // Remove transpose solutions
	}

	// Setup the 8 primary items
	items = []string{
		"x3x4x5c2c3",
		"c4c5c6c7c8",
		"c9c10c11c12x6",
		"c13c14x7x8x9",
		"x1x2x5c5c9",
		"c1c2c6c10c13",
		"c3c7c11c14x10",
		"c8c12x7x11x12",
	}

	// Setup the 14 + 2W secondary items, c1-c14 (colored) + words + words'
	for i := 1; i <= 14; i++ {
		sitems = append(sitems, fmt.Sprintf("c%d", i))
	}
	sitems = append(sitems, words...)
	for _, word := range words {
		// prime values to remove tranpose solutions
		sitems = append(sitems, word+"'")
	}

	// Setup the 8W options
	for _, word := range words {
		// x3 x4 x5 c2 c3
		options = append(options, []string{"x3x4x5c2c3",
			"c2:" + word[3:4], "c3:" + word[4:5],
			word})

		// c4 c5 c6 c7 c8
		options = append(options, []string{"c4c5c6c7c8",
			"c4:" + word[0:1], "c5:" + word[1:2], "c6:" + word[2:3], "c7:" + word[3:4], "c8:" + word[4:5],
			word, word + "'"})

		// c9 c10 c11 c12 x6
		options = append(options, []string{"c9c10c11c12x6",
			"c9:" + word[0:1], "c10:" + word[1:2], "c11:" + word[2:3], "c12:" + word[3:4],
			word})

		// c13 c14 x7 x8 x9
		options = append(options, []string{"c13c14x7x8x9",
			"c13:" + word[0:1], "c14:" + word[1:2],
			word})

		// x1 x2 x5 c5 c9
		options = append(options, []string{"x1x2x5c5c9",
			"c5:" + word[3:4], "c9:" + word[4:5],
			word})

		// c1 c2 c6 c10 c13
		options = append(options, []string{"c1c2c6c10c13",
			"c1:" + word[0:1], "c2:" + word[1:2], "c6:" + word[2:3], "c10:" + word[3:4], "c13:" + word[4:5],
			word, word + "'"})

		// c3 c7 c11 c14 x10
		options = append(options, []string{"c3c7c11c14x10",
			"c3:" + word[0:1], "c7:" + word[1:2], "c11:" + word[2:3], "c14:" + word[3:4],
			word})

		// c8 c12 x7 x11 x12
		options = append(options, []string{"c8c12x7x11x12",
			"c8:" + word[0:1], "c12:" + word[1:2],
			word})

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
			kernel := make([]byte, 14)

			for _, option := range solution {
				switch option[0] {
				case "c1c2c6c10c13":
					kernel[0] = option[1][3]
				case "x3x4x5c2c3":
					kernel[1] = option[1][3]
					kernel[2] = option[2][3]
				case "c4c5c6c7c8":
					kernel[3] = option[1][3]
					kernel[4] = option[2][3]
					kernel[5] = option[3][3]
					kernel[6] = option[4][3]
					kernel[7] = option[5][3]
				case "c9c10c11c12x6":
					kernel[8] = option[1][3]
					kernel[9] = option[2][4]
					kernel[10] = option[3][4]
					kernel[11] = option[4][4]
				case "c13c14x7x8x9":
					kernel[12] = option[1][4]
					kernel[13] = option[2][4]
				}
			}
			resume := visit(string(kernel))
			return resume
		})
}
