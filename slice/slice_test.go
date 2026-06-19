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

func TestAppendUniqueString(t *testing.T) {

	cases := []struct {
		xValues [][]string // Values of x to insert
		a       [][]string // Expected resulting [][]string
	}{
		{
			[][]string{},
			nil,
		},
		{
			[][]string{{}},
			[][]string{{}},
		},
		{
			[][]string{{"a"}},
			[][]string{{"a"}},
		},
		{
			[][]string{{"a"}, {"b"}},
			[][]string{{"a"}, {"b"}},
		},
		{
			[][]string{{"a"}, {"b"}, {"a"}},
			[][]string{{"a"}, {"b"}},
		},
		{
			[][]string{{"a"}, {"b"}, {"a", "b"}},
			[][]string{{"a"}, {"b"}, {"a", "b"}},
		},
		{
			[][]string{{"a"}, {"b"}, {"a", "b"}, {"a", "b"}},
			[][]string{{"a"}, {"b"}, {"a", "b"}},
		},
		{
			[][]string{{"a"}, {"b"}, {"a", "b"}, {"a", "b"}, {"b"}, {"a", "b", "c"}},
			[][]string{{"a"}, {"b"}, {"a", "b"}, {"a", "b", "c"}},
		},
	}

	for i, c := range cases {

		// Insert all values of x
		var a [][]string
		for _, x := range c.xValues {
			a = AppendUniqueString(a, x)
		}

		if !reflect.DeepEqual(a, c.a) {
			t.Errorf("For case #%d (xValues=%v) got %v; want %v", i, c.xValues, a, c.a)
		}
	}
}
