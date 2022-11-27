package math

import (
	"testing"
)

func TestCountDigits(t *testing.T) {
	cases := []struct {
		n        int64 // input number
		expected int   // output count
	}{
		{0, 0},
		{1, 1},
		{11, 2},
		{999, 3},
		{1000, 4},
		{1234567890, 10},
	}

	for i, c := range cases {
		got := CountDigits(c.n)

		if got != c.expected {
			t.Errorf("Got %v for case %d; want %v", got, i, c.expected)
		}
	}
}
