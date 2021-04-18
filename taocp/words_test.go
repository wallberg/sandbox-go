package taocp

import (
	"log"
	"math"
	"reflect"
	"testing"

	"github.com/wallberg/sandbox/sgb"
	"github.com/wallberg/sandbox/slice"
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
		// { // too slow
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

func TestExercise_7221_90(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// These tests verify Exercise 7.2.2.1-89
	cases := []struct {
		p        int
		left     bool
		solution []string
		mMin     int
	}{
		{
			1,
			false,
			[]string{
				"spots",
			},
			1300,
		},
		{
			2,
			true,
			[]string{
				"write", "whole",
			},
			500,
		},
		{
			2,
			false,
			[]string{
				"stall", "spies",
			},
			1633,
		},
		{
			3,
			true,
			[]string{
				"makes", "lived", "waxes",
			},
			500,
		},
		{
			3,
			false,
			[]string{
				"stood", "holes", "leaps",
			},
			1561,
		},
		// too slow
		// {
		// 	4,
		// 	true,
		// 	[]string{
		// 		"there", "share", "whole", "whose",
		// 	},
		// 	500,
		// },
		// {
		// 	4,
		// 	false,
		// 	[]string{
		// 		"mixed", "tears", "slept", "salad",
		// 	},
		// 	1300,
		// },
		// {
		// 	5,
		// 	true,
		// 	[]string{
		// 		"stood", "thank", "share", "ships", "store",
		// 	},
		// 	500,
		// },
		// {
		// 	5,
		// 	false,
		// 	[]string{
		// 		"years", "steam", "sales", "marks", "dried",
		// 	},
		// 	1300,
		// },
		// {
		// 	6,
		// 	true,
		// 	[]string{
		// 		"where", "sheep", "small", "still", "whole", "share",
		// 	},
		// 	500,
		// },
		// {
		// 	6,
		// 	false,
		// 	[]string{
		// 		"steps", "seals", "draws", "knots", "traps", "drops",
		// 	},
		// 	1300,
		// },
		// {
		// 	7,
		// 	true,
		// 	[]string{
		// 		"makes", "based", "tired", "works", "lands", "lives", "gives",
		// 	},
		// 	500,
		// },
		// {
		// 	7,
		// 	false,
		// 	[]string{
		// 		"tried", "fears", "slips", "seams", "draws", "erect", "tears",
		// 	},
		// 	1300,
		// },
		// {
		// 	8,
		// 	true,
		// 	[]string{
		// 		"water", "makes", "loved", "gives", "lakes", "based", "notes", "tones",
		// 	},
		// 	504,
		// },
		// {
		// 	8,
		// 	false,
		// 	[]string{
		// 		"years", "stops", "hooks", "fried", "tears", "slant", "sword", "sweep",
		// 	},
		// 	1300,
		// },
		// {
		// 	9,
		// 	true,
		// 	[]string{
		// 		"where", "sheet", "still", "shall", "white", "shape", "stars", "whole", "shore",
		// 	},
		// 	500,
		// },
		// {
		// 	9,
		// 	false,
		// 	[]string{
		// 		"start", "spear", "sales", "tests", "steer", "speak", "skies", "slept", "sport",
		// 	},
		// 	1300,
		// },
		// {
		// 	10,
		// 	true,
		// 	[]string{
		// 		"there", "shoes", "shirt", "stone", "shook", "start", "while", "shell", "steel", "sharp",
		// 	},
		// 	500,
		// },
		// {
		// 	10,
		// 	false,
		// 	[]string{
		// 		"years", "stock", "horns", "fuels", "beets", "speed", "tears", "plant", "sword", "sweep",
		// 	},
		// 	1300,
		// },
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

		stats := &ExactCoverStats{
			// Progress: true,
			// Delta:    20000000,
			// Debug:        true,
			// Verbosity:    2,
			// SuppressDump: true,
		}

		xccOptions := &XCCOptions{
			Minimax:       true,
			MinimaxSingle: false,
			Exercise83:    true,
		}

		mMin := math.MaxInt64
		var got [][]string // list of solutions with the minimum m value

		WordStair(words, c.p, c.left, stats, xccOptions, func(s []string) bool {
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
			return true
		})

		if mMin >= c.mMin {
			t.Errorf("For case #%d, p=%d, left=%t, got m=%d; want m < %d",
				i, c.p, c.left, mMin, c.mMin)
		}

		// Check that we got a matching cycle of words
		isCycle := false
		for _, s := range got {
			if slice.IsCycleString(s[:c.p], c.solution) ||
				slice.IsCycleString(s[c.p:], c.solution) ||
				slice.IsCycleString(slice.ReverseString(s[:c.p]), c.solution) ||
				slice.IsCycleString(slice.ReverseString(s[c.p:]), c.solution) {
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
	}
}

func TestWordStairKernel(t *testing.T) {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// words, err := sgb.LoadWords()
	// if err != nil {
	// 	t.Errorf("Error getting words: %v", err)
	// 	return
	// }

	// These tests verify Exercise 7.2.2.1-89
	cases := []struct {
		words []string
		left  bool
		count int
	}{
		{
			[]string{"dried", "years", "steam", "sales", "skies", "seats", "dream", "salad"},
			false,
			1,
		},
	}

	for i, c := range cases {

		stats := &ExactCoverStats{
			// Progress: true,
			// Delta:    20000000,
			// Debug:        true,
			// Verbosity:    2,
			// SuppressDump: true,
		}

		kernels := make(map[string]bool)

		WordStairKernel(c.words, c.left, stats, func(s string) bool {
			kernels[s] = true
			return true
		})

		if len(kernels) != c.count {
			t.Errorf("For case #%d, left=%t, got %d kernels; want %d",
				i, c.left, len(kernels), c.count)
		}
	}
}