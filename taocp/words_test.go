package taocp

import (
	"testing"
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
