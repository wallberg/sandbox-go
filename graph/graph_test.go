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

func TestArcs(t *testing.T) {
	cases := []struct {
		g     *graph.Mutable // graph
		v     int            // edges from vertex v
		count int            // number of edges
	}{
		{
			Cycle(5),
			0,
			2,
		},
		{
			Complete(5),
			0,
			4,
		},
	}

	for _, c := range cases {
		count := 0
		for a := Arcs(c.g, c.v); a != nil; a = a.next {
			count++
		}

		if count != c.count {
			t.Errorf("Got %d; want %d", count, c.count)
		}
	}
}

func TestConnectedSubsetsVertex(t *testing.T) {

	cases := []struct {
		g    *graph.Mutable // input graph
		n    int            // input size of subsets
		v    int            // included vertex
		want [][]int        // expected solutions
	}{
		{
			CartesianProduct(Path(3), Path(3)),
			5,
			0,
			[][]int{
				{0, 1, 3, 2, 4},
				{0, 1, 3, 2, 6},
				{0, 1, 3, 2, 5},
				{0, 1, 3, 4, 6},
				{0, 1, 3, 4, 5},
				{0, 1, 3, 4, 7},
				{0, 1, 3, 6, 7},
				{0, 1, 2, 4, 5},
				{0, 1, 2, 4, 7},
				{0, 1, 2, 5, 8},
				{0, 1, 4, 5, 7},
				{0, 1, 4, 5, 8},
				{0, 1, 4, 7, 6},
				{0, 1, 4, 7, 8},
				{0, 3, 4, 6, 5},
				{0, 3, 4, 6, 7},
				{0, 3, 4, 5, 7},
				{0, 3, 4, 5, 2},
				{0, 3, 4, 5, 8},
				{0, 3, 4, 7, 8},
				{0, 3, 6, 7, 8},
			},
		},
		{
			Path(1),
			1,
			0,
			[][]int{
				{0},
			},
		},
	}

	for _, c := range cases {
		var got [][]int

		ConnectedSubsetsVertex(c.g, c.n, c.v, func(solution []int) bool {
			cp := make([]int, c.n)
			copy(cp, solution)
			got = append(got, cp)
			return false
		})

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Got %v; want %v", got, c.want)
		}
	}
}

func TestConnectedSubsets(t *testing.T) {

	cases := []struct {
		n     int
		count int // number of connected subsets
	}{
		{1, 1},
		{2, 4},
		{3, 22},
		{4, 113},
		{5, 571},
		{6, 2816},
		{7, 13616},
		{8, 64678},
		{9, 302574},
	}

	for _, c := range cases {

		g := CartesianProduct(Path(c.n), Path(c.n))
		count := 0

		ConnectedSubsets(g, c.n, func(solution []int) bool {
			count++
			return false
		})

		if !reflect.DeepEqual(count, c.count) {
			t.Errorf("Got %v; want %v", count, c.count)
		}
	}
}

func TestRemoveIsolated(t *testing.T) {
	g := graph.New(5)
	g.AddBoth(0, 1)
	g.AddCost(1, 4, 10)
	g.AddCost(4, 1, 12)
	g.Add(1, 2)
	g.AddCost(1, 3, 11)
	g.Add(3, 2)

	var mapping map[int]int
	var h *graph.Mutable

	// First iteration
	g, mapping = RemoveIsolated(g)

	if g.Order() != 4 {
		t.Errorf("Got order=%d after first iteration; want 4", g.Order())
	}
	if len(mapping) != 4 {
		t.Errorf("Got len(mapping)=%d after first iteration; want 4", len(mapping))
	}

	// Second iteration
	g, mapping = RemoveIsolated(g)

	if g.Order() != 3 {
		t.Errorf("Got order=%d after second iteration; want 3", g.Order())
	}
	if len(mapping) != 3 {
		t.Errorf("Got len(mapping)=%d after second iteration; want 3", len(mapping))
	}

	// Third iteration
	h, mapping = RemoveIsolated(g)

	if h != g {
		t.Errorf("g itself not returned after third iteration")
	}
	if mapping != nil {
		t.Errorf("nil mapping not returned after third iteration")
	}

}
