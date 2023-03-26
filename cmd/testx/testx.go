package main

import (
	"fmt"
	"log"
)

// data contains all common data for a single operation.
type data struct {

	// nodes - all nodes we've seen, ordered by their index value
	nodes []node

	// stack - current node stack
	stack []int

	// indexes - map a vertex v to it's index value
	indexes map[int]int
}

// node stores data for a single node in the connection process.
type node struct {
	// lowlink - smallest index of any node on the stack known to be reachable from this node.
	// Set it initially to the value of index. After we've looked at all child nodes, if lowlink
	// still equals index then we know it's the root of the SCC on the stack.
	lowlink int

	// onStack - is this node currently on the stack?
	onStack bool
}

// strongConnect runs Tarjan's algorithm recursivley and outputs a grouping of
// strongly connected vertices.
// Returns: a) *node - the v node currently processed, and b) bool - did we find a contradiction?
//
// TODO: add license and acknowledgement
func (data *data) strongConnect(v int, BIMP [][]int, BSIZE []int, visit func([]int) bool) (*node, bool) {

	// Set the depth index for v to the smallest unused index
	vIndex := len(data.nodes)
	data.indexes[v] = vIndex

	// Add v to the list of "seen" nodes
	vNode := &node{lowlink: vIndex, onStack: true}
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
			wNode, contradiction := data.strongConnect(w, BIMP, BSIZE, visit)
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

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var (
		// Number of variables
		n int

		// BIMP - binary clauses; instead of the buddy system, we are using built-in slices
		BIMP [][]int

		// BSIZE - number of clauses for each l in BIMP (literal indexed)
		BSIZE []int

		// S - number of SCCs found in the dependency graph
		S int

		// PARENT - parent node in the dependency subforest (literal indexed)
		PARENT []int
	)

	test := func() {
		fmt.Println("---")

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
				for i := 0; i < BSIZE[l]; i++ {
					vertices[BIMP[l][i]] = true
				}
			}
		}

		fmt.Printf("vertices: %v\n", vertices)

		// Use Tarjan's algorithm to get strongly connected components and a
		// subforest at the same time. The subforest must include all 2C
		// candidate literals, have no cycles, and have no nodes with > 1
		// outbound edge.

		// sccs - list of all SCCs
		var sccs [][]int

		data := &data{
			nodes:   make([]node, 0, len(vertices)),
			indexes: make(map[int]int, len(vertices)),
		}

		for v := range vertices {
			if _, seen := data.indexes[v]; !seen {
				_, contradiction := data.strongConnect(v, BIMP, BSIZE, func(scc []int) bool {

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

		// Select a representative from each SCC
		S = len(sccs)
		reps := make([]int, S)

		// Select representatives from each SCC of size = 1
		for r, scc := range sccs {

			if len(scc) == 1 {
				l := scc[0]

				if BSIZE[l] == 0 {
					PARENT[l] = 0
				} else {
					PARENT[l] = BIMP[l][0] // must select one parent for the subforest
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

				if i < len(scc) {
					// Found a representative l
					reps[r] = l

					// Set the PARENT for all l ∈ SCC
					for _, lp := range scc {
						if l == lp {
							if BSIZE[l] == 0 {
								PARENT[l] = 0
							} else {
								PARENT[l] = BIMP[l][0] // must select one parent for the subforest
							}
						} else {
							PARENT[lp] = l
						}
					}
				} else {
					log.Panicf("assertion failed: we did not find any members of l ∈ SCC=%v with a matching ¬l among the other representatives", scc)
				}
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
	}

	n = 4
	BIMP = make([][]int, 2*n+2)
	BSIZE = make([]int, 2*n+2)
	PARENT = make([]int, 2*n+2)
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
	}

	test()

	n = 9
	BIMP = make([][]int, 2*n+2)
	BSIZE = make([]int, 2*n+2)
	PARENT = make([]int, 2*n+2)
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
	}

	test()

	n = 4
	BIMP = make([][]int, 2*n+2)
	BSIZE = make([]int, 2*n+2)
	PARENT = make([]int, 2*n+2)
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
	}

	test()

}
