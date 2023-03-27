package main

import (
	"fmt"
	"log"
)

// sccData contains all common data for a single scc search.
type sccData struct {

	// nodes - all nodes we've seen, ordered by their index value
	nodes []sccNode

	// stack - current node stack
	stack []int

	// indexes - map a vertex v to it's index value
	indexes map[int]int
}

// sccNode stores data for a single vertex node in the search process.
type sccNode struct {
	// lowlink - smallest index of any node on the stack known to be reachable from this node.
	// Set it initially to the value of index. After we've looked at all child nodes, if lowlink
	// still equals index then we know it's the root of the SCC on the stack.
	lowlink int

	// onStack - is this node currently on the stack?
	onStack bool
}

// @note scc()
// scc runs Tarjan's algorithm recursivley and outputs a grouping of
// strongly connected vertices.
// Returns: a) *node - the v node currently processed, and b) bool - did we find a contradiction?
//
// TODO: add license and acknowledgement
func (data *sccData) scc(v int, BIMP [][]int, BSIZE []int, visit func([]int) bool) (*sccNode, bool) {

	// Set the depth index for v to the smallest unused index
	vIndex := len(data.nodes)
	data.indexes[v] = vIndex

	// Add v to the list of "seen" nodes
	vNode := &sccNode{lowlink: vIndex, onStack: true}
	data.nodes = append(data.nodes, *vNode)

	// Push v into the stack
	data.stack = append(data.stack, v)

	fmt.Printf("v=%d, S=%v\n", v, data.stack)

	// Consider successors of v
	for i := 0; i < BSIZE[v]; i++ {
		w := BIMP[v][i]

		wIndex, seen := data.indexes[w]
		if !seen {

			// Successor w has not yet been visited; recurse on it
			wNode, contradiction := data.scc(w, BIMP, BSIZE, visit)
			if contradiction {
				// Contradiction found, halt the search
				return vNode, true
			}

			// vNode.lowlink = min(vNode.lowlink, wNode.lowlink)
			if wNode.lowlink < vNode.lowlink {
				vNode.lowlink = wNode.lowlink
			}

		} else if data.nodes[wIndex].onStack {

			// Successor w is in stack S and hence in the current SCC.

			// If successor w is not in stack S, then (v, w) is an edge pointing
			// to an SCC previously found and must be ignored.

			// vNode.lowlink = min(vNode.lowlink, wIndex)
			if wIndex < vNode.lowlink {
				vNode.lowlink = wIndex
			}
		}
	}

	// If v is a root node, pop the stack and generate an SCC
	if vNode.lowlink == vIndex {
		var scc []int
		i := len(data.stack) - 1
		for {
			w := data.stack[i]
			wIndex := data.indexes[w]
			data.nodes[wIndex].onStack = false
			scc = append(scc, w)
			if wIndex == vIndex {
				break
			}
			i--
		}
		data.stack = data.stack[:i]
		fmt.Printf("scc=%v\n", scc)
		if visit(scc) {
			// Contradiction found, halt the search
			return vNode, true
		}
	}

	return vNode, false
}

// @note build_lookahead()
// build_lookahead builds the LL and LO lookahead tables,
// returning new values of i and degree
func build_lookahead(LL []int, LO []int, CHILDREN [][]int, l int, i int, degree int) (int, int) {

	// Build LL in preorder
	LL[i] = l

	// Visit all children
	nexti := i + 1
	for _, lp := range CHILDREN[l] {
		nexti, degree = build_lookahead(LL, LO, CHILDREN, lp, nexti, degree)
	}

	// Build LO in postorder
	degree += 2
	LO[i] = degree
	i++

	return nexti, degree
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var (
		// Number of variables
		n int

		// BIMP - binary clauses; instead of the buddy system, we are using built-in slices
		BIMP [][]int

		// BSIZE - number of clauses for each l in BIMP (literal indexed)
		BSIZE []int

		// sccs - list of all Strongly Connected Components (SCCs) (Algorithm X)
		sccs [][]int

		// S - number of SCCs found in the dependency digraph (Algorithm X)
		S int

		// reps - list of all representatives of each SCC in the dependency subforest (Algorithm X)
		reps []int

		// PARENT - parent node in the dependency subforest (Algorithm X, literal indexed)
		PARENT []int

		// CHILDREN - child nodes in the dependency subforest (Algorithm X, literal indexed)
		CHILDREN [][]int

		// LL - lookahead literal in the Algorithm X dependency subforest
		LL []int

		// LO - lookahead offset in the Algrorithm X dependency subforest
		LO []int
	)

	// @note test()
	test := func() {
		fmt.Println("---")

		//
		// @note X4 - Dependency Digraph
		//
		// Construct the dependency graph on the 2C candidate literals, by
		// extracting a subset of arcs from the BIMP tables. (This computation
		// needn't be exact, because we're only calculating heuristics; an upper
		// bound can be placed on the number of arcs considered, so that we
		// don't spend too much time here. However, it is important to have the
		// arc u --> v iff ¬v --> ¬u is also present.)
		// TODO: Extract a subset, instead of everything

		vertices := make(map[int]bool, 2*n)

		for l := 2; l <= 2*n+1; l++ {
			if BSIZE[l] > 0 {
				vertices[l] = true
				CHILDREN[l] = CHILDREN[l][:0]
				for i := 0; i < BSIZE[l]; i++ {
					lp := BIMP[l][i]
					vertices[lp] = true
					CHILDREN[lp] = CHILDREN[lp][:0]
				}
			}
		}

		fmt.Printf("vertices: %v\n", vertices)

		// Use Tarjan's algorithm to get strongly connected components and a
		// subforest at the same time. The subforest must include all 2C
		// candidate literals, have no cycles, and have no nodes with > 1
		// outbound edge.
		// TODO: split these comments to correct positions

		sccs = sccs[:0]

		// TODO: determine if we can be faster without these maps
		data := &sccData{
			nodes:   make([]sccNode, 0, len(vertices)),
			indexes: make(map[int]int, len(vertices)),
		}

		for v := range vertices {
			if _, seen := data.indexes[v]; !seen {
				_, contradiction := data.scc(v, BIMP, BSIZE, func(scc []int) bool {

					if len(scc) > 1 {

						// Check if this SCC contains both l and ¬l, and if so terminate with a contradiction.
						// TODO: Determine if this would be faster with a map
						for i := 0; i < len(scc); i++ {
							for j := i + 1; j < len(scc); j++ {
								if scc[i] == scc[j]^1 {
									// Contradiction found, halt the search
									return true
								}
							}
						}
					}

					sccs = append(sccs, scc)

					return false
				})
				if contradiction {
					fmt.Println("Contradiction!")
					return
				}
			}
		}

		//
		// @note X4 - Dependency Subforest
		//
		// Select a representative from each SCC
		S = len(sccs)
		reps = reps[:S]

		// Select representatives from each SCC of size = 1
		for r, scc := range sccs {

			if len(scc) == 1 {
				l := scc[0]

				if BSIZE[l] == 0 {
					PARENT[l] = 0
				} else {
					lp := BIMP[l][0] // must select one parent for the subforest
					PARENT[l] = lp
					CHILDREN[lp] = append(CHILDREN[lp], l)
				}
				reps[r] = l
			}
		}

		// Select representatives from each SCC of size > 1
		for r, scc := range sccs {

			if len(scc) > 1 {

				// Choose a representatve l with maximum h(l), ensuring that if l is a
				// representative then ¬l is also a representative

				var l int    // l with maximum h(l) value
				var maxi int // index of l in reps
				var maxh int // maximum  value of in reps

				// Search for l ∈ SCC which has maximum h(l) and ¬l is a representative
				i := 0
				for ; i < len(scc); i++ {

					// Search for next value of l ∈ SCC which has maximum h(l)
					maxi = i
					maxh = scc[i] // TODO: replace simulated h(l) with real value

					for j := i + 1; j < len(scc); j++ {
						if scc[j] > maxh {
							maxi = j
							maxh = scc[j]
						}
					}

					l = scc[maxi]

					// Determine if ¬l is a representative
					isRep := false
					for _, rep := range reps {
						if l == rep^1 {
							isRep = true
							break
						}
					}

					if isRep {
						// Found the representative we are looking for
						break
					} else {
						// ¬l is not a representative so swap l to the beginning of the list and try again with the
						// remaining members of the SCC.
						if maxi != i {
							scc[i], scc[maxi] = scc[maxi], scc[i]
						}
					}
				}

				if i == len(scc) {
					// We did not find a representative l with matching ¬l
					// Let's pick the l with maximum h(l) and assume that the matching
					// ¬l will arrive later
					l = scc[0]
				}

				reps[r] = l

				// Set the PARENT for all l ∈ SCC
				for _, lp := range scc {
					if l == lp {
						PARENT[l] = 0
					} else {
						PARENT[lp] = l
						CHILDREN[l] = append(CHILDREN[l], lp)
					}
				}
			}
		}

		// Assert that every variable x has 0 or 2 literal representatives
		// TODO: wrap this assertion in a debug
		xCount := make([]int, n+1)
		for _, l := range reps {
			xCount[l>>1] += 1
		}
		for _, x := range xCount {
			if xCount[x] != 0 && xCount[x] != 2 {
				log.Panicf("assertion failed: x=%d has %d representatives of l or ¬l", x, xCount[x])
			}
		}

		fmt.Printf("S=%d, representatives=%v\n", S, reps)
		fmt.Print("PARENT: ")
		for l := 2; l <= 2*n+1; l++ {
			if l > 2 {
				fmt.Print(", ")
			}
			fmt.Printf("%d=%d", l, PARENT[l])
		}
		fmt.Println()

		//
		// @note X4 - Lookahead Tables
		//
		// Construct lookahead tables LL and LO, for the S candidate literals,
		// with LL containing the candidate literals in preorder, and LO containing
		// each literal's truth degree (2 * literal's postorder position)
		LL = LL[:len(vertices)]
		LO = LO[:len(vertices)]

		i := 0
		degree := 0

		// Iterate over the representatives
		for _, rep := range reps {
			if PARENT[rep] == 0 {
				i, degree = build_lookahead(LL, LO, CHILDREN, rep, i, degree)
			}
		}

		fmt.Print("\nLL: ")
		for i := 0; i < len(LL); i++ {
			fmt.Printf("%2d ", LL[i])
		}
		fmt.Print("\nLO: ")
		for i := 0; i < len(LO); i++ {
			fmt.Printf("%2d ", LO[i])
		}
		fmt.Println()
	}

	n = 4
	BIMP = make([][]int, 2*n+2)
	BSIZE = make([]int, 2*n+2)
	PARENT = make([]int, 2*n+2)
	CHILDREN = make([][]int, 2*n+2)
	sccs = make([][]int, 0, 2*n)
	reps = make([]int, 0, 2*n)
	LL = make([]int, 0, 2*n)
	LO = make([]int, 0, 2*n)

	for l := 2; l <= 2*n+1; l++ {
		switch l {
		case 3:
			BIMP[l] = []int{5, 7}
		case 4:
			BIMP[l] = []int{2}
		case 6:
			BIMP[l] = []int{2}
		case 7:
			BIMP[l] = []int{9}
		case 8:
			BIMP[l] = []int{6}
		default:
			BIMP[l] = []int{}
		}
		BSIZE[l] = len(BIMP[l])
		CHILDREN[l] = make([]int, 0)
	}

	test()

	n = 9
	BIMP = make([][]int, 2*n+2)
	BSIZE = make([]int, 2*n+2)
	PARENT = make([]int, 2*n+2)
	CHILDREN = make([][]int, 2*n+2)
	sccs = make([][]int, 0, 2*n)
	reps = make([]int, 0, 2*n)
	LL = make([]int, 0, 2*n)
	LO = make([]int, 0, 2*n)

	for l := 2; l <= 2*n+1; l++ {
		switch l {
		case 3:
			BIMP[l] = []int{6, 18}
		case 5:
			BIMP[l] = []int{16}
		case 7:
			BIMP[l] = []int{8, 2, 14}
		case 9:
			BIMP[l] = []int{6, 12}
		case 13:
			BIMP[l] = []int{8, 14}
		case 15:
			BIMP[l] = []int{12, 6, 18}
		case 17:
			BIMP[l] = []int{4}
		case 19:
			BIMP[l] = []int{14, 2}
		default:
			BIMP[l] = []int{}
		}
		BSIZE[l] = len(BIMP[l])
		CHILDREN[l] = make([]int, 0)
	}

	test()

	n = 4
	BIMP = make([][]int, 2*n+2)
	BSIZE = make([]int, 2*n+2)
	PARENT = make([]int, 2*n+2)
	CHILDREN = make([][]int, 2*n+2)
	sccs = make([][]int, 0, 2*n)
	reps = make([]int, 0, 2*n)
	LL = make([]int, 0, 2*n)
	LO = make([]int, 0, 2*n)

	for l := 2; l <= 2*n+1; l++ {
		switch l {
		case 2:
			BIMP[l] = []int{8}
		case 3:
			BIMP[l] = []int{5, 7}
		case 4:
			BIMP[l] = []int{2}
		case 5:
			BIMP[l] = []int{9}
		case 6:
			BIMP[l] = []int{2}
		case 7:
			BIMP[l] = []int{9}
		case 8:
			BIMP[l] = []int{4, 6}
		case 9:
			BIMP[l] = []int{3}
		default:
			BIMP[l] = []int{}
		}
		BSIZE[l] = len(BIMP[l])
		CHILDREN[l] = make([]int, 0)
	}

	test()

}
