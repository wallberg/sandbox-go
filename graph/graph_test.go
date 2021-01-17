package graph

import (
	"reflect"
	"testing"

	"github.com/yourbasic/graph"
)

func TestPath(t *testing.T) {
	want := graph.New(3)
	want.AddBoth(0, 1)
	want.AddBoth(1, 2)

	got := Path(3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v; want %v", got, want)
	}
}

func TestCycle(t *testing.T) {
	want := graph.New(3)
	want.AddBoth(0, 1)
	want.AddBoth(1, 2)
	want.AddBoth(2, 0)

	got := Cycle(3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v; want %v", got, want)
	}
}

func TestComplete(t *testing.T) {
	want := graph.New(4)
	want.AddBoth(0, 1)
	want.AddBoth(0, 2)
	want.AddBoth(0, 3)
	want.AddBoth(1, 2)
	want.AddBoth(1, 3)
	want.AddBoth(2, 3)

	got := Complete(4)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v; want %v", got, want)
	}
}

func TestCartesianProduce(t *testing.T) {
	var want, got *graph.Mutable

	want = graph.New(12)
	want.AddBoth(0, 3)
	want.AddBoth(0, 1)
	want.AddBoth(1, 4)
	want.AddBoth(1, 2)
	want.AddBoth(2, 5)
	want.AddBoth(3, 6)
	want.AddBoth(3, 4)
	want.AddBoth(4, 7)
	want.AddBoth(4, 5)
	want.AddBoth(5, 8)
	want.AddBoth(6, 9)
	want.AddBoth(6, 7)
	want.AddBoth(7, 10)
	want.AddBoth(7, 8)
	want.AddBoth(8, 11)
	want.AddBoth(9, 10)
	want.AddBoth(10, 11)

	got = CartesianProduct(Path(4), Path(3))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v; want %v", got, want)
	}

}
