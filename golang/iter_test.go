package golang

import (
	"iter"
	"maps"
	"reflect"
	"slices"
	"testing"
)

// Tests to support understanding of the new Iterator feature in Go 1.23
// https://pkg.go.dev/iter

// TestIter tests standard library support in maps and slices packages
func TestIter(t *testing.T) {

	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3

	var got []string

	got = append(got, slices.Sorted(maps.Keys(m))...)

	want := []string{"a", "b", "c"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v; want %v", got, want)
	}
}

type Things struct{}

func (*Things) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		items := []string{"a", "b", "c"}
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	}
}

// TestIterPush tests the default Push model
func TestIterPush(t *testing.T) {

	var thing Things

	var got []string
	for v := range thing.All() {
		got = append(got, v)
	}

	want := []string{"a", "b", "c"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v; want %v", got, want)
	}
}

func (things *Things) PullAll() iter.Seq[string] {
	return func(yield func(string) bool) {
		next, stop := iter.Pull(things.All())
		defer stop()
		for {
			v, ok := next()
			if !ok {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// TestIterPull tests the Pull model
func TestIterPull(t *testing.T) {

	var thing Things

	var got []string
	for v := range thing.PullAll() {
		got = append(got, v)
	}

	want := []string{"a", "b", "c"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v; want %v", got, want)
	}
}
