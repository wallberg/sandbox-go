package taocp

import (
	"log"
	"reflect"
	"testing"

	"github.com/wallberg/sandbox/sgb"
)

func TestDoubleWordSquare(t *testing.T) {

	cases := []struct {
		words           []string
		removeTranspose bool
		count           int
	}{
		{
			[]string{
				"abcde",
				"wcdef",
				"xdefg",
				"yefgh",
				"zfghi",
				"awxyz",
				"bcdef",
				"cdefg",
				"defgh",
				"efghi",
			},
			false,
			2,
		},
		{
			[]string{
				"abcde",
				"wcdef",
				"xdefg",
				"yefgh",
				"zfghi",
				"awxyz",
				"bcdef",
				"cdefg",
				"defgh",
				"efghi",
			},
			true,
			1,
		},
		// too slow
		// {
		// 	nil,
		//  true,
		// 	323264,
		// },
	}

	// var err error
	// cases[1].words, err = sgb.LoadWords()
	// if err != nil {
	// 	t.Errorf("Error getting words: %v", err)
	// 	return
	// }

	for i, c := range cases {

		stats := &ExactCoverStats{
			// Progress: true,
			// Delta:    50000000,
			// Debug:    true,
			// Verbosity:    2,
			// SuppressDump: true,
		}

		xccOptions := &XCCOptions{Exercise83: c.removeTranspose}

		count := 0
		DoubleWordSquare(c.words, stats, xccOptions, func(s []string) bool {
			count++
			return true
		})

		if count != c.count {
			t.Errorf("Got %d solutions for case #%d; want %d", count, i, c.count)
		}
	}
}

func TestDoubleWordSquareMinimax(t *testing.T) {
	// These tests verify Exercise 7.2.2.1-88
	cases := []struct {
		words    []string
		solution []string
	}{
		{
			[]string{
				"arena",
				"sinks",
				"onset",
				"needs",
				"mason",
				"urine",
				"sense",
				"inked",
				"casts",
				"utero",
				"three",
				"earth",
				"steps",
				"anger",
				"tense",
				"beast",
				"blast",
				"scope",
				"lance",
				"argon",
				"sexed",
				"piing",
				"weald",
				"iodic",
				"abuzz",
				"poxed",
				"hurly",
				"music",
			},
			[]string{
				"beast", "lance", "argon", "steps", "three",
				"blast", "earth", "anger", "scope", "tense",
			},
		},
		// { // too slow (80s)
		// 	nil,
		// 	[]string{
		// 		"beast", "lance", "argon", "steps", "three",
		// 		"blast", "earth", "anger", "scope", "tense",
		// 	},
		// },
	}

	// var err error
	// cases[1].words, err = sgb.LoadWords()
	// if err != nil {
	// 	t.Errorf("Error getting words: %v", err)
	// 	return
	// }

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

		var got []string
		DoubleWordSquare(c.words, stats, xccOptions, func(s []string) bool {
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
			return true
		})

		if !reflect.DeepEqual(got, c.solution) {
			t.Errorf("Got solution %v for case #%d; want %v", got, i, c.solution)
		}
	}
}

func TestExercise_7221_89(t *testing.T) {
	// These tests verify Exercise 7.2.2.1-89
	cases := []struct {
		size     int
		solution []string
	}{
		{
			2,
			[]string{
				"is", "to",
				"it", "so",
			},
		},
		{
			3,
			[]string{
				"may", "age", "not",
				"man", "ago", "yet",
			},
		},
		{
			4,
			[]string{
				"show", "none", "open", "west",
				"snow", "hope", "ones", "went",
			},
		},
		// { // too slow
		// 	5,
		// 	[]string{
		// 		"start", "three", "roofs", "asset", "peers",
		// 		"strap", "those", "arose", "refer", "tests",
		// 	},
		// },
		// { // unverified
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
		DoubleWordSquare(words, stats, xccOptions, func(s []string) bool {
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
			return true
		})

		if !reflect.DeepEqual(got, c.solution) {
			t.Errorf("For case #%d got solution %v; want %v", i, got, c.solution)
		}
	}
}
