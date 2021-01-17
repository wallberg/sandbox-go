package graph

import (
	"github.com/yourbasic/graph"
)

// Enhance the yourbasic/graph package.

// Path generates a path (P) of order n.
func Path(n int) (g *graph.Mutable) {
	g = graph.New(n)
	for i := 0; i < n-1; i++ {
		g.AddBoth(i, i+1)
	}
	return
}

// Cycle generates a cycle (C) of order n.
func Cycle(n int) (g *graph.Mutable) {
	g = graph.New(n)
	for i := 0; i < n-1; i++ {
		g.AddBoth(i, i+1)
	}
	g.AddBoth(0, n-1)
	return
}

// Complete generates a complete graph (K) of order n.
func Complete(n int) (g *graph.Mutable) {
	g = graph.New(n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				g.AddBoth(i, j)
			}
		}
	}
	return
}

// CartesianProduct generates the cartesian product of graphs g and h
func CartesianProduct(g graph.Iterator, h graph.Iterator) *graph.Mutable {
	gOrder, hOrder := g.Order(), h.Order()

	gh := graph.New(gOrder * hOrder)

	// Iterate over g vertices
	for gFrom := 0; gFrom < gOrder; gFrom++ {

		// Iterate over g edges
		g.Visit(gFrom, func(gTo int, gCost int64) bool {

			// Iterate over h vertices
			for hFrom := 0; hFrom < hOrder; hFrom++ {

				// Iterate over  h edges
				h.Visit(hFrom, func(hTo int, hCost int64) bool {

					ghFrom := gFrom*hOrder + hFrom
					gh.AddCost(ghFrom, gTo*hOrder+hFrom, gCost)
					gh.AddCost(ghFrom, gFrom*hOrder+hTo, hCost)
					return false

				})
			}
			return false
		})
	}

	return gh
}
