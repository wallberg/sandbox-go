package golang

import (
	"testing"
)

// Tests to support understanding of Generics
// https://go.dev/doc/tutorial/generics

// SumInts adds together the values of m.
func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

// SumFloats adds together the values of m.
func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func Sum[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func TestSumInts(t *testing.T) {

	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	got := SumInts(ints)

	var want int64 = 46

	if got != want {
		t.Errorf("Got %v from SumInts; want %v", got, want)
	}

	got = Sum(ints)

	if got != want {
		t.Errorf("Got %v from Sum; want %v", got, want)
	}
}

func TestSumFloats(t *testing.T) {

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	got := SumFloats(floats)

	var want float64 = 62.97

	if got != want {
		t.Errorf("Got %v from SumFloats; want %v", got, want)
	}

	got = Sum[string, float64](floats)

	if got != want {
		t.Errorf("Got %v from Sum; want %v", got, want)
	}
}
