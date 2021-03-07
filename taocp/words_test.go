package taocp

import (
	"testing"

	"github.com/wallberg/sandbox/sgb"
)

func TestDoubleWordSquare(t *testing.T) {

	var err error

	cases := []struct {
		words []string
		count int
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
			1,
		},
		{
			nil,
			323264,
		},
	}

	cases[1].words, err = sgb.LoadWords()
	if err != nil {
		t.Errorf("Error getting words: %v", err)
		return
	}

	for i, c := range cases {

		stats := &ExactCoverStats{
			Progress: true,
			Delta:    50000000,
			// Debug:        true,
			// Verbosity:    2,
			// SuppressDump: true,
		}

		count := 0
		DoubleWordSquare(c.words, stats, func(s []string) bool {
			count++
			return true
		})

		if count != c.count {
			t.Errorf("Got %d solutions for case #%d; want %d", count, i, c.count)
		}
	}
}
