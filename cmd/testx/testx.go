package main

import (
	"fmt"
	"log"
)

// data contains all common data for a single operation.
type data struct {

	// nodes - all nodes we've seen, ordered by their index value
	nodes []node

	// S - current node stack
	S []int

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
// Returns the v node currently processed, or nil if the visit function requests
// a halt.
// TODO: add license and acknowledgement
func (data *data) strongConnect(v int, BIMP [][]int, BSIZE []int, visit func([]int) bool) *node {

	// Set the depth index for v to the smallest unused index
	vIndex := len(data.nodes)
	data.indexes[v] = vIndex

	// Add v to the list of "seen" nodes
	vNode := &node{lowlink: vIndex, onStack: true}
	data.nodes = append(data.nodes, *vNode)

	// Push v into the stack
	data.S = append(data.S, v)

	// fmt.Printf("v=%d, nodes=%v, S=%v\n", v, data.nodes, data.S)

	// Consider successors of v
	for i := 0; i < BSIZE[v]; i++ {
		w := BIMP[v][i]

		wIndex, seen := data.indexes[w]
		if !seen {

			// Successor w has not yet been visited; recurse on it
			wNode := data.strongConnect(w, BIMP, BSIZE, visit)
			if wNode == nil {
				// Halt the search
				return nil
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
		i := len(data.S) - 1
		for {
			w := data.S[i]
			wIndex := data.indexes[w]
			data.nodes[wIndex].onStack = false
			scc = append(scc, w)
			if wIndex == vIndex {
				break
			}
			i--
		}
		data.S = data.S[:i]
		fmt.Printf("scc=%v\n", scc)
		if visit(scc) {
			// Halt the search
			return nil
		}
	}

	return vNode
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

		var reps []int

		data := &data{
			nodes:   make([]node, 0, len(vertices)),
			indexes: make(map[int]int, len(vertices)),
		}

		for v := range vertices {
			if _, seen := data.indexes[v]; !seen {
				result := data.strongConnect(v, BIMP, BSIZE, func(scc []int) bool {

					var rep int // representative from each SCC
					if len(scc) == 1 {
						rep = scc[0]
					} else {
						// Check if this SCC contains both l and ¬l, and if so
						// terminate with a contradiction.
						// TODO: Determine if this would be faster with a map
						for i := 0; i < len(scc); i++ {
							for j := i + 1; j < len(scc); j++ {
								if scc[i] == scc[j]^1 {
									// Contradiction
									return true
								}
							}
						}

						// TODO: Choose a representatve l with maximum h(l)
						rep = scc[0]
					}

					reps = append(reps, rep)
					return false
				})
				if result == nil {
					fmt.Println("Contradiction!")
					return
				}
			}
		}
		fmt.Printf("Representatives=%v", reps)
		fmt.Println()
	}

	n = 4
	BIMP = make([][]int, 2*n+2)
	BSIZE = make([]int, 2*n+2)
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
	for l := 2; l <= 2*n+1; l++ {
		switch l {
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
			BIMP[l] = []int{6}
		case 9:
			BIMP[l] = []int{3}
		default:
			BIMP[l] = []int{}
		}
		BSIZE[l] = len(BIMP[l])
	}

	test()

}
