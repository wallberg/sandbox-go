//go:build longtests

package taocp

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/wallberg/sandbox-go/sgb"
	"github.com/wallberg/sandbox-go/slice"
)

// Run all long tests, using maximum parallelization

func TestLong(t *testing.T) {
	t.Run("taocp long tests", func(t *testing.T) {
		t.Run("DoubleWordSquare", doubleWordSquare)
		t.Run("DoubleWordSquareMinimax", doubleWordSquareMinimax)
		t.Run("Exercise 7.2.2.1-89", exercise_7221_89)
		t.Run("Exercise 7.2.2.1-90", exercise_7221_90)
		t.Run("MagicHexagon", magicHexagon)
	})
}

func doubleWordSquare(t *testing.T) {
	t.Parallel()

	cases := []struct {
		words           []string
		removeTranspose bool
		count           int
	}{
		{ // Passed in 3.3 hours
			nil,
			true,
			323264,
		},
	}

	all_words, err := sgb.LoadWords()
	if err != nil {
		t.Errorf("Error getting words: %v", err)
		return
	}

	for i, c := range cases {

		stats := &ExactCoverStats{
			// Progress: true,
			// Delta:    50000000,
			// Debug:    true,
			// Verbosity:    2,
			// SuppressDump: true,
		}

		xccOptions := &XCCOptions{Exercise83: c.removeTranspose}

		if c.words == nil {
			c.words = all_words
		}

		count := 0
		for range DoubleWordSquare(c.words, stats, xccOptions) {
			count++
		}

		if count != c.count {
			t.Errorf("Got %d solutions for case #%d; want %d", count, i, c.count)
		}
	}
}

func doubleWordSquareMinimax(t *testing.T) {
	t.Parallel()

	// These tests verify Exercise 7.2.2.1-88
	cases := []struct {
		words    []string
		solution []string
	}{
		{
			nil,
			[]string{
				"beast", "lance", "argon", "steps", "three",
				"blast", "earth", "anger", "scope", "tense",
			},
		},
	}

	all_words, err := sgb.LoadWords()
	if err != nil {
		t.Errorf("Error getting words: %v", err)
		return
	}

	for i, c := range cases {

		stats := &ExactCoverStats{
			// Progress:     true,
			// Delta:        1000000, //5000000,
			// Debug:        false,
			// Verbosity:    2,
			// SuppressDump: true,
		}

		xccOptions := &XCCOptions{
			Minimax:       true,
			MinimaxSingle: true,
			Exercise83:    true,
		}

		if c.words == nil {
			c.words = all_words
		}

		var got []string
		for s := range DoubleWordSquare(c.words, stats, xccOptions) {
			got = s
			if stats.Debug || stats.Progress {
				// Determine max word position
				m := 0
				for _, word1 := range s {
					for j, word2 := range c.words {
						if word1 == word2 {
							if j > m {
								m = j
							}
							break
						}
					}
				}
				log.Printf("m=%d, %v", m, s)
			}
		}

		if !reflect.DeepEqual(got, c.solution) {
			t.Errorf("Got solution %v for case #%d; want %v", got, i, c.solution)
		}
	}
}

func exercise_7221_89(t *testing.T) {
	t.Parallel()

	// These tests verify Exercise 7.2.2.1-89
	cases := []struct {
		size     int
		solution []string
	}{
		{
			5,
			[]string{
				"start", "three", "roofs", "asset", "peers",
				"strap", "those", "arose", "refer", "tests",
			},
		},
		// { // Passed in 5.7 hours
		// 	6,
		// 	[]string{
		// 		"chests", "lustre", "obtain", "arenas", "circle", "assess",
		// 		"cloaca", "hubris", "esters", "stance", "trials", "senses",
		// 	},
		// },
		// { // unverified
		// 	7,
		// 	[]string{
		// 		"hertzes", "operate", "mimical", "acerate", "genetic", "endmost", "resents",
		// 		"homager", "epicene", "remends", "trireme", "zacaton", "etatist", "selects",
		// 	},
		// },
	}

	for i, c := range cases {
		name := fmt.Sprintf("size=%v", c.size)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			words, err := sgb.LoadOSPD4(c.size)
			if err != nil {
				t.Errorf("Error getting words: %v", err)
				return
			}

			stats := &ExactCoverStats{
				// Progress: true,
				// Delta:    1000000, //5000000,
				// Debug:        false,
				// Verbosity:    2,
				// SuppressDump: true,
			}

			xccOptions := &XCCOptions{
				Minimax:       true,
				MinimaxSingle: true,
				Exercise83:    true,
			}

			var got []string
			for s := range DoubleWordSquare(words, stats, xccOptions) {
				got = s
				if stats.Debug || stats.Progress {
					// Determine max word position
					m := 0
					for _, word1 := range s {
						for j, word2 := range words {
							if word1 == word2 {
								if j > m {
									m = j
								}
								break
							}
						}
					}
					log.Printf("m=%d, %v", m, s)
				}
			}

			if !reflect.DeepEqual(got, c.solution) {
				t.Errorf("For case #%d got solution %v; want %v", i, got, c.solution)
			}
		})
	}
}

func exercise_7221_90(t *testing.T) {

	t.Parallel()

	// log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// These tests verify Exercise 7.2.2.1-90
	cases := []struct {
		p        int
		left     bool
		solution []string
		mMin     int
	}{
		{
			5,
			false,
			[]string{
				"years", "steam", "sales", "marks", "dried",
			},
			1300,
		},
		{
			6,
			true,
			[]string{
				"where", "sheep", "small", "still", "whole", "share",
			},
			500,
		},
		{
			6,
			false,
			[]string{
				"steps", "seals", "draws", "knots", "traps", "drops",
			},
			1300,
		},
		{
			7,
			true,
			[]string{
				"makes", "based", "tired", "works", "lands", "lives", "gives",
			},
			500,
		},
		{
			7,
			false,
			[]string{
				"tried", "fears", "slips", "seams", "draws", "erect", "tears",
			},
			1300,
		},
		{
			8,
			true,
			[]string{
				"water", "makes", "loved", "given", "lakes", "based", "notes", "tones",
			},
			504,
		},
		{
			8,
			false,
			[]string{
				"years", "stops", "hooks", "fried", "tears", "slant", "sword", "sweep",
			},
			1300,
		},
		{
			9,
			true,
			[]string{
				"where", "sheet", "still", "shall", "white", "shape", "stars", "whole", "shore",
			},
			500,
		},
		{
			9,
			false,
			[]string{
				"start", "spear", "sales", "tests", "steer", "speak", "skies", "slept", "sport",
			},
			1300,
		},
		{
			10,
			true,
			[]string{
				"there", "shoes", "shirt", "stone", "shook", "start", "while", "shell", "steel", "sharp",
			},
			500,
		},
		{
			10,
			false,
			[]string{
				"years", "stock", "horns", "fuels", "beets", "speed", "tears", "plant", "sword", "sweep",
			},
			1300,
		},
	}

	words, err := sgb.LoadWords()
	if err != nil {
		t.Errorf("Error getting words: %v", err)
		return
	}

	getM := func(a []string) (int, []int) {
		m := 0
		mWords := make([]int, len(a))
		for i, word := range a {
			mWord := slice.FindString(words, word)
			mWords[i] = mWord
			if mWord > m {
				m = mWord
			}
		}
		return m, mWords
	}

	for i, c := range cases {

		name := fmt.Sprintf("p=%v:left=%v:mMin=%v", c.p, c.left, c.mMin)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stats := &ExactCoverStats{
				// Progress: true,
				// Delta:    20000000,
				// Debug:        true,
				// Verbosity:    2,
				// SuppressDump: true,
			}

			xccOptions := &XCCOptions{
				Minimax:       true,
				MinimaxSingle: true,
				Exercise83:    true,
			}

			mMin := math.MaxInt64
			var got [][]string // list of solutions with the minimum m value

			for s := range WordStair(words, c.p, c.left, stats, xccOptions) {
				// Determine max word position
				m, mWords := getM(s)
				if stats.Debug || stats.Progress {
					log.Printf("m=%d, %v, %v", m, s, mWords)
				}

				if m < mMin {
					mMin = m
					got = nil
				}

				got = append(got, s)
			}

			if mMin >= c.mMin {
				t.Errorf("For case #%d, p=%d, left=%t, got m=%d; want m < %d",
					i, c.p, c.left, mMin, c.mMin)
			}

			// Check that we got a matching cycle of words
			isCycle := false
			for _, s := range got {
				sReverse := make([]string, len(s))
				copy(sReverse, s)
				slices.Reverse(sReverse[:c.p])
				slices.Reverse(sReverse[c.p:])
				if slice.IsCycleString(s[:c.p], c.solution) ||
					slice.IsCycleString(s[c.p:], c.solution) ||
					slice.IsCycleString(sReverse[:c.p], c.solution) ||
					slice.IsCycleString(sReverse[c.p:], c.solution) {
					isCycle = true
					break
				}
			}

			if !isCycle {
				m, mWords := getM(c.solution)
				if stats.Debug || stats.Progress {
					log.Printf("Expected: m=%d, %v, %v", m, c.solution, mWords)
				}
				t.Errorf("For case #%d, p=%d, left=%t, got solutions %v; want %v",
					i, c.p, c.left, got, c.solution)
			}
		})
	}
}

func magicHexagon(t *testing.T) {

	stats := &ExactCoverStats{
		// Progress: true,
		// Delta:    10000,
		// Debug:    true,
		// Verbosity: 2,
	}

	for got := range MagicHexagon(stats) {

		want := [][]string{
			{"abc", "16:a", "19:b", "3:c"},
			{"defg", "12:d", "2:e", "7:f", "17:g"},
			{"hijkl", "10:h", "4:i", "5:j", "1:k", "18:l"},
			{"mnop", "13:m", "8:n", "6:o", "11:p"},
			{"qrs", "15:q", "14:r", "9:s"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected solution %v; got %v", want, got)
		}

		break
	}
}
