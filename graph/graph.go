package graph

import (
	"fmt"
	"log"
	"sort"
	"strings"

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

// Arc represents an arc variable for use with ARCS, NAME, and TIP,
// as in TAOCP Section 7.
type Arc struct {
	v    int   // an edge desination vertex
	cost int64 // edge cost
	next *Arc  // next Arc node in the list
}

// String describes an Arc as a string
func (a *Arc) String() string {
	return fmt.Sprintf("(arc=%d:%d)", a.v, a.cost)
}

// Arcs generates a singly linked list of Arc nodes, with one node for each
// edge that emanates from vertex v in graph g.  Arcs are sorted in
// vertex order
func Arcs(g graph.Iterator, v int) (arcs *Arc) {
	costs := make(map[int]int64) // Map distination vertices to costs
	var edges []int

	g.Visit(v, func(w int, c int64) bool {
		edges = append(edges, w)
		costs[w] = c
		return false
	})

	sort.Ints(edges)

	a := arcs
	for _, w := range edges {
		if a == nil {
			arcs = &Arc{v: w, cost: costs[w]}
			a = arcs
		} else {
			a.next = &Arc{v: w, cost: costs[w]}
			a = a.next
		}
	}

	return arcs
}

var debug bool

func init() {
	debug = false
}

// ConnectedSubsetsVertex generates all connected subsets in g of size n which
// contain vertex v. Implements TAOCP Algorithm R from 7.2.2 Exercise 75.
func ConnectedSubsetsVertex(g graph.Iterator, n int, v int,
	visit func([]int) (halt bool)) {

	var (
		l   int    // backtrack level
		i   int    // an index
		a   *Arc   // a linked list of edges to neighbor
		u   int    // a vertex
		vs  []int  // list of vertices
		is  []int  // list of indices
		as  []*Arc // list of linked list of edges to neightbors
		tag []int  // number of times a vertex has been tagged
	)

	dump := func() {
		log.Printf("| i=%d, a=%v, v=%d, u=%d", i, a, v, u)
		log.Printf("| is=%v, as=%v, vs=%v",
			is[0:l], as[0:l], vs[0:l])
		log.Printf("| tag=%v", tag)
	}

	if debug {
		log.Printf("R1. input graph g")
		for v := 0; v < g.Order(); v++ {
			var l strings.Builder
			l.WriteString(fmt.Sprintf("  v=%d ->", v))
			for a := Arcs(g, v); a != nil; a = a.next {
				l.WriteString((fmt.Sprintf(" %d", a.v)))
			}
			log.Print(l.String())
		}
	}

	// R1. [Initialize.]

	vs = make([]int, n)
	is = make([]int, n)
	as = make([]*Arc, n)

	tag = make([]int, g.Order())
	vs[0] = v
	i = 0
	a = Arcs(g, v)
	as[0] = a
	tag[v] = 1
	l = 1

	if debug {
		log.Printf("R1. Initialized at level %d", l)
		dump()
	}

	goto R4

R2:
	// R2. [Enter level l.]

	if debug {
		log.Printf("R2. Enter level %d", l)
		dump()
	}

	if l == n {
		visit(vs)
		l = n - 1
	}

R3:
	// R3. [Advance a.]
	if debug {
		if a.next != nil {
			log.Printf("R3. Advance a from %v to %v", a, a.next)
			dump()
		}
	}

	a = a.next

R4:
	// R4. [Done with level?]
	if a != nil {
		goto R5
	}

	if i == l-1 {
		goto R6
	}

	i++
	v = vs[i]
	a = Arcs(g, v)

	if debug {
		log.Printf("R4. Advance i=%d, v=%d, a=%v", i, v, a)
		dump()
	}

R5:
	// R5. [Try a.]
	u = a.v
	tag[u]++

	if debug {
		log.Printf("R5. Try a=%v", a)
		dump()
	}

	if tag[u] > 1 {
		if debug {
			log.Printf("R5. tag[%d]=%d > 1", u, tag[u])
		}
		goto R3
	}

	is[l] = i
	as[l] = a
	vs[l] = u
	l++

	goto R2

R6:
	// R6. [Backtrack.]
	if debug {
		log.Printf("R6. Backtrack")
		dump()
	}

	l--
	if l == 0 {
		return
	}

	i = is[l]
	v = vs[i]

	// untag all neighbors of v_k, for l >= k > i
	for k := i + 1; k <= l; k++ {
		if debug {
			log.Printf("|  Untagging neighbors of vs[%d]=%d", k, vs[k])
		}
		g.Visit(vs[k], func(w int, c int64) bool {
			tag[w]--
			return false
		})
	}

	a = as[l].next
	for a != nil {
		tag[a.v]--
		a = a.next
	}

	a = as[l]

	if debug {
		log.Printf("R6. Untagging complete")
		dump()
	}

	goto R3
}
