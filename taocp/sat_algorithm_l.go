package taocp

import (
	"fmt"
	"log"
	"strings"
)

// SatAlgorithmL implements Algorithm D (7.2.2.2), satisfiability by DPLL with lookahead.
// The task is to determine if the clause set is satisfiable, and if it is return
// one satisfying assignment of the clauses.
//
// Arguments:
// n       -- number of strictly distinct literals
// clauses -- list of clauses to satisfy
// stats   -- SAT processing statistics
// options -- runtime options
//
func SatAlgorithmL(n int, clauses SatClauses,
	stats *SatStats, options *SatOptions) (bool, []int) {

	// State represents a single cell in the state table
	type State struct {
		L int // literal
	}

	var (
		nOrig       int     // original value of n, before conversion to SAT3
		m           int     // total number of clauses
		stateSize   int     // total size of the state table
		state       []State // search state
		start       []int   // start of each clause in the table
		watch       []int   // list of all clauses that currently watch l
		link        []int   // the number of another clause with the same watch literal
		h           []int   // the literal being watched at depth d
		next        []int   // active ring : not-yet-set variables whose watch lists aren't empty
		head        int     // head pointer into the active ring
		tail        int     // tail pointer into the active ring
		d           int     // depth of the implicit search tree
		x           []int   // selected literal at depth d
		l, lp       int     // literal
		p, q        int     // index into the state table
		b           int     // branch on literal?
		f           int     // flag?
		i, j, jp, k int     // indices
		moves       []int   // store current progress
		debug       bool    // debugging is enabled
		progress    bool    // progress tracking is enabled
	)

	// dump
	dump := func() {

		var b strings.Builder
		b.WriteString("\n")

		// State, p
		b.WriteString("   p = ")
		for p := range state {
			b.WriteString(fmt.Sprintf(" %2d", p))
		}
		b.WriteString("\n")

		// State, L
		b.WriteString("L(p) = ")
		for p := range state {
			if state[p].L == 0 {
				b.WriteString("  -")
			} else {
				b.WriteString(fmt.Sprintf(" %2d", state[p].L))
			}
		}
		b.WriteString("\n")

		// l
		b.WriteString("       l = ")
		for l := range watch {
			b.WriteString(fmt.Sprintf(" %2d", l))
		}
		b.WriteString("\n")

		// WATCH
		b.WriteString("WATCH(l) = ")
		for _, val := range watch {
			b.WriteString(fmt.Sprintf(" %2d", val))
		}
		b.WriteString("\n")

		// NEXT
		b.WriteString(" NEXT(v) = ")
		for _, val := range next {
			b.WriteString(fmt.Sprintf(" %2d", val))
		}
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("head=%d, tail=%d\n", head, tail))
		b.WriteString(fmt.Sprintf("active ring=%s\n", activeRing()))

		// j
		b.WriteString("       j = ")
		for j := range start {
			b.WriteString(fmt.Sprintf(" %2d", j))
		}
		b.WriteString("\n")

		// START
		b.WriteString("START(j) = ")
		for _, val := range start {
			b.WriteString(fmt.Sprintf(" %2d", val))
		}
		b.WriteString("\n")

		// LINK
		b.WriteString(" LINK(j) = ")
		for _, val := range link {
			b.WriteString(fmt.Sprintf(" %2d", val))
		}
		b.WriteString("\n")

		log.Print(b.String())
	}

	// showProgress
	showProgress := func() {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("Nodes=%d, d=%d, l=%d, moves=%v\n", stats.Nodes, d, l, moves[1:d+1]))

		log.Print(b.String())
	}

	// initialize
	initialize := func() {

		if stats != nil {
			stats.Theta = stats.Delta
			stats.MaxLevel = -1
			if stats.Levels == nil {
				stats.Levels = make([]int, n)
			} else {
				for len(stats.Levels) < n {
					stats.Levels = append(stats.Levels, 0)
				}
			}
			debug = stats.Debug
			progress = stats.Progress
		}

		// Initialize the state table
		m = len(clauses)
		start = make([]int, m+1)
		link = make([]int, m+1)
		watch = make([]int, 2*n+2)
		moves = make([]int, n+1)
		h = make([]int, n+1)
		next = make([]int, n+1)
		x = make([]int, n+1)

		// Record all binary clauses in the BIMP array

		// Record all ternary clauses in the TIMP array

		// Let U be the number of distinct variable in unit clauses

		// Terminate unsuccessfully if two unit clauses contradict each other

		// Record all distinct unit literals in FORCE[k] for 0 <= k < U

		for _, clause := range clauses {
			stateSize += len(clause)
		}
		state = make([]State, stateSize)

		start[0] = stateSize

		// index into state
		p = stateSize - 1

		// Iterate over the clauses
		for j := 1; j <= len(clauses); j++ {
			clause := clauses[j-1]
			start[j] = p - len(clause) + 1

			// Iterate over literal values of the clauses
			for _, k := range clause {
				// compute literal l
				var l int
				if k >= 0 {
					l = 2 * k
				} else {
					l = -2*k + 1
				}

				// insert into the state table
				state[p].L = l

				// Check if this is the last literal in the clause
				if p == start[j] {
					// Insert this literal into the watch list of clauses.
					// watch is the head of the list, with link containing
					// the next pointers. The last clause in the list has
					// value of 0.
					if watch[l] == 0 {
						// Insert the first clause into the list
						watch[l] = j
					} else {
						// Insert this clause to the end of the list
						jp := watch[l]
						for link[jp] != 0 {
							jp = link[jp]
						}
						link[jp] = j
					}
				}

				// advance to the next position in the table
				p -= 1
			}
		}
	}

	// lvisit prepares the solution
	lvisit := func() []int {
		solution := make([]int, nOrig)
		// for i := 1; i < n+1; i++ {
		// 	if h[i] > 0 {
		// 		solution[h[i]-1] = (moves[i] % 2) ^ 1
		// 	}
		// }
		if debug {
			log.Printf("visit solution=%v", solution)
		}

		return solution
	}

	//
	// L1 [Initialize.]
	//

	// Convert the input to SAT3, if it isn't already
	nOrig = n
	_, n, clauses = Sat3(n, clauses)

	initialize()

	if debug {
		log.Printf("L1. Initialize")
	}

	if debug {
		dump()
	}

	if progress {
		showProgress()
	}

	// 	//
	// 	// L2. [Success?]
	// 	//
	// L2:
	// 	if tail == 0 {
	// 		// visit the solution
	// 		if debug {
	// 			log.Println("L2. [Success!]")
	// 		}
	// 		if stats != nil {
	// 			stats.Solutions++
	// 		}

	// 		return true, lvisit()
	// 	}

	// 	k = tail

	// 	//
	// 	// L3. [Look for unit clauses.]
	// 	//
	// L3:
	// 	head = next[k]

	// 	if debug {
	// 		log.Printf("L3. [Look for unit clauses.] active=%s, k=%d, head=%d", activeRing(), k, head)
	// 	}

	// 	// Compute f = [2h is a unit] + 2[2h + 1 is a unit]
	// 	f = isUnit(2*head) + 2*isUnit(2*head+1)
	// 	if debug {
	// 		log.Printf("L3. f=%d", f)
	// 	}

	// 	if f == 3 {
	// 		goto L7 // [Backtrack.]
	// 	}

	// 	if f == 1 || f == 2 {
	// 		moves[d+1] = f + 3
	// 		tail = k
	// 		goto L5 // [Move on.]
	// 	}

	// 	if head != tail {
	// 		k = head
	// 		goto L3 // [Look for unit clauses.]
	// 	}

	// 	//
	// 	// L4. [Two-way branch.]
	// 	//

	// 	head = next[tail]
	// 	if watch[2*head] == 0 || watch[2*head+1] != 0 {
	// 		moves[d+1] = 1
	// 	} else {
	// 		moves[d+1] = 0
	// 	}

	// 	if debug {
	// 		log.Printf("L4. [Two-way branch.] d=%d, x=%v, moves=%v", d, x[1:], moves[1:d+1])
	// 	}

	// L5:
	// 	//
	// 	// L5. [Move on.]
	// 	//
	// 	if debug {
	// 		log.Printf("L5. [Move on.]")
	// 	}

	// 	d += 1
	// 	h[d] = head
	// 	k = head

	// 	if stats != nil {
	// 		stats.Levels[d-1]++
	// 		stats.Nodes++

	// 		if progress {
	// 			if d > stats.MaxLevel {
	// 				stats.MaxLevel = d
	// 			}
	// 			if stats.Nodes >= stats.Theta {
	// 				showProgress()
	// 				stats.Theta += stats.Delta
	// 			}
	// 		}
	// 	}

	// 	if tail == k {
	// 		tail = 0
	// 	} else {
	// 		// delete variable k from the ring
	// 		next[tail] = next[k]
	// 		head = next[k]
	// 	}

	// 	//
	// 	// L6. [Update watches.]
	// 	//
	// L6:
	// 	if debug {
	// 		log.Printf("L6. [Update watches.]")
	// 	}

	// 	b = (moves[d] + 1) % 2
	// 	x[k] = b

	// 	// Clear the watch list for ^(x_k)
	// 	l = 2*k + b
	// 	j = watch[l]
	// 	watch[l] = 0

	// 	for j != 0 {

	// 		// step (i)
	// 		jp = link[j]
	// 		i = start[j]
	// 		p = i + 1

	// 		// step (ii) - loop while L(p) is false
	// 		// will end before p == start[j-1]
	// 		for x[state[p].L>>1] == state[p].L&1 {
	// 			p += 1
	// 		}

	// 		// step (iii)
	// 		lp = state[p].L
	// 		state[p].L = l
	// 		state[i].L = lp

	// 		// step (iv)
	// 		p = watch[lp]
	// 		q = watch[lp^1]
	// 		if p != 0 || q != 0 || x[lp>>1] >= 0 {
	// 			goto L6vi
	// 		}

	// 		// step (v) - Insert |l'| into the ring as its new head
	// 		if tail == 0 {
	// 			tail = lp >> 1
	// 			head = lp >> 1
	// 			next[tail] = head
	// 		} else {
	// 			next[lp>>1] = head
	// 			head = lp >> 1
	// 			next[tail] = head
	// 		}

	// 		// step (vi)
	// 	L6vi:
	// 		// Insert j into the watch list of l'
	// 		link[j] = p
	// 		watch[lp] = j

	// 		// step (vii)
	// 		j = jp
	// 	}

	// 	if debug {
	// 		log.Printf("L6. d=%d, l=%d, active=%s, x=%v, moves=%v", d, l, activeRing(), x[1:], moves[1:d+1])
	// 	}

	// 	goto L2

	// 	//
	// 	// L7. [Backtrack.]
	// 	//
	// L7:
	// 	if debug {
	// 		log.Printf("L7. [Backtrack.]")
	// 	}

	// 	tail = k

	// 	for moves[d] >= 2 {
	// 		k = h[d]
	// 		x[k] = -1
	// 		if watch[2*k] != 0 || watch[2*k+1] != 0 {
	// 			next[k] = head
	// 			head = k
	// 			next[tail] = head
	// 		}
	// 		d -= 1
	// 	}

	// 	//
	// 	// L8. [Failure?]
	// 	//
	// 	if debug {
	// 		log.Printf("L8. [Failure?]")
	// 	}

	// 	if d > 0 {
	// 		moves[d] = 3 - moves[d]
	// 		k = h[d]
	// 		goto L6
	// 	}

	// 	// Terminate, the clauses aren't satisfiable
	// 	return false, nil
}
