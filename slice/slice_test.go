package slice

import (
	"reflect"
	"testing"
)

func TestFindString(t *testing.T) {

	cases := []struct {
		a   []string
		x   string
		pos int
	}{
		{[]string{}, "x", -1},
		{[]string{"w", "y", "z"}, "x", -1},
		{[]string{"y", "x"}, "x", 1},
		{[]string{"x", "y"}, "x", 0},
	}

	for i, c := range cases {

		pos := FindString(c.a, c.x)

		if pos != c.pos {
			t.Errorf("For case #%d (a=%v, x=%s) got %d; want %d", i, c.a, c.x, pos, c.pos)
		}
	}
}

func TestIsCycleString(t *testing.T) {

	cases := []struct {
		a      []string
		b      []string
		result bool
	}{
		{[]string{}, []string{}, true},
		{[]string{"w", "y", "z"}, []string{}, false},
		{[]string{"y", "x"}, []string{"x", "y"}, true},
		{[]string{"y", "x", "w"}, []string{"w", "y", "x"}, true},
		{[]string{"w", "y", "x"}, []string{"w", "y", "x"}, true},
		{[]string{"y", "x", "w"}, []string{"w", "x", "y"}, false},
		{[]string{"y", "x", "w"}, []string{"y", "w", "x"}, false},
	}

	for i, c := range cases {

		result := IsCycleString(c.a, c.b)

		if result != c.result {
			t.Errorf("For case #%d (a=%v, b=%v) got %t; want %t", i, c.a, c.b, result, c.result)
		}
	}
}

func TestReverseString(t *testing.T) {

	cases := []struct {
		a      []string
		result []string
	}{
		{[]string{}, []string{}},
		{[]string{"x"}, []string{"x"}},
		{[]string{"w", "y", "z"}, []string{"z", "y", "w"}},
	}

	for i, c := range cases {

		result := ReverseString(c.a)

		if !reflect.DeepEqual(result, c.result) {
			t.Errorf("For case #%d (a=%v) got %v; want %v", i, c.a, result, c.result)
		}
	}
}
