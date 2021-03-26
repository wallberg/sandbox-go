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
		// too long
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

	cases := []struct {
		words    []string
		solution []string
	}{
		{
			[]string{
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
			},
			[]string{
				"beast", "lance", "argon", "steps", "three",
				"blast", "earth", "anger", "scope", "tense",
			},
		},
		{
			nil,
			[]string{
				"beast", "lance", "argon", "steps", "three",
				"blast", "earth", "anger", "scope", "tense",
			},
		},
	}

	var err error
	cases[1].words, err = sgb.LoadWords()
	if err != nil {
		t.Errorf("Error getting words: %v", err)
		return
	}

	for i, c := range cases {

		stats := &ExactCoverStats{
			Progress:     true,
			Delta:        5000000,
			Debug:        false,
			Verbosity:    2,
			SuppressDump: true,
		}

		xccOptions := &XCCOptions{
			Minimax:       true,
			MinimaxSingle: false,
			Exercise83:    false,
		}

		var got []string
		DoubleWordSquare(c.words, stats, xccOptions, func(s []string) bool {
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
			got = s
			log.Printf("m=%d, %v", m, s)
			return true
		})

		if !reflect.DeepEqual(got, c.solution) {
			t.Errorf("Got solution %v for case #%d; want %v", got, i, c.solution)
		}
	}
}
