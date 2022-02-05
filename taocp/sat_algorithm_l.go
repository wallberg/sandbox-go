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
	stats *SatStats, options *SatOptions) (sat bool, solution []int) {

	var (
		nOrig      int     // original value of n, before conversion to 3SAT
		m          int     // total number of clauses
		varx       []int   // VAR - permutation of {1,...,n} (VAR[k] = x iff INX[x] = k)
		inx        []int   // INX
		d          int     // depth of the implicit search tree
		f          int     // number of fixed variables
		timp       []int   // TIMP - ternary clauses
		tsize      []int   // TSIZE - number of clauses for each l in TIMP
		link       []int   // LINK - circular list of the three literals in each clause in TIMP
		bimp       [][]int // BIMP - instead of the buddy system, trying using built-in slices
		bsize      []int   // BSIZE - number of clauses for each l in BIMP
		p, pp, ppp int     // index into TIMP
		units      int     // U - number of distinct variables in unit clauses
		force      []int   // FORCE - stack of U unit variables which have a forced value
		istackSize int     // size of istack
		istamp     int     // stamp to make downdating BIMP tables easier
		k          int     // indices
		debug      bool    // debugging is enabled
		progress   bool    // progress tracking is enabled
	)

	fmt.Println(nOrig, d, m, f, istackSize, istamp)

	// dump
	dump := func() {

		var b strings.Builder
		b.WriteString("\n")

		// FORCE
		b.WriteString("FORCE\n")
		b.WriteString(fmt.Sprintf("units=%d: ", units))
		for i := 0; i < units; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}", force[i]))
		}
		b.WriteString("\n\n")

		// BIMP
		b.WriteString("BIMP\n")
		for l := 2; l <= 2*n+1; l++ {
			b.WriteString(fmt.Sprintf("%d: ", l))
			for i := 0; i < bsize[l]; i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(fmt.Sprintf("%d", bimp[l][i]))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")

		// TIMP
		b.WriteString("TIMP\n")
		for l := 2; l <= 2*n+1; l++ {
			b.WriteString(fmt.Sprintf("%d: ", l))
			for i := 0; i < tsize[l]; i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				p := timp[l] + 2*i
				b.WriteString(fmt.Sprintf("{%d,%d}->", timp[p], timp[p+1]))
				p = link[p]
				b.WriteString(fmt.Sprintf("{%d,%d}->", timp[p], timp[p+1]))
				p = link[p]
				b.WriteString(fmt.Sprintf("{%d,%d}", timp[p], timp[p+1]))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")

		log.Print(b.String())
	}

	// showProgress
	showProgress := func() {
		var b strings.Builder
		// b.WriteString(fmt.Sprintf("Nodes=%d, d=%d, l=%d, moves=%v\n", stats.Nodes, d, l, moves[1:d+1]))

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
	}

	// k2l converts variable k to literal 2k if positive, 2k+1 if negative
	k2l := func(k int) int {
		if k < 0 {
			return -2*k + 1
		} else {
			return 2 * k
		}
	}

	// // l2k
	// l2k := func(l int) int {
	// 	if l%2 == 0 {
	// 		return l >> 2
	// 	} else {
	// 		return (l >> 2) * -1
	// 	}
	// }

	// // lvisit prepares the solution
	// lvisit := func() []int {
	// 	solution := make([]int, nOrig)
	// 	// for i := 1; i < n+1; i++ {
	// 	// 	if h[i] > 0 {
	// 	// 		solution[h[i]-1] = (moves[i] % 2) ^ 1
	// 	// 	}
	// 	// }
	// 	if debug {
	// 		log.Printf("visit solution=%v", solution)
	// 	}

	// 	return solution
	// }

	//
	// L1 [Initialize.]
	//

	// Convert the input to 3SAT, if it isn't already
	nOrig = n
	_, n, clauses = Sat3(n, clauses)

	initialize()

	if debug {
		log.Printf("L1. Initialize")
	}

	m = len(clauses)

	//
	// Record all unit clauses with forced variable values
	//
	force = make([]int, 2*n+2)
	units = 0
	for _, clause := range clauses {
		if len(clause) == 1 {
			l := k2l(clause[0])

			// Look for a contradiction
			for k := 0; k < units; k++ {
				if l^1 == force[k] {
					// A contradiction
					if debug {
						log.Printf("L1. Found a unit clause contradiction")
					}
					return false, nil
				}
			}

			// Add l to the stack of distinct unit clauses
			force[units] = l
			units += 1

		}
	}

	//
	// Record all binary clauses in the BIMP array
	//
	bimp = make([][]int, 2*n+2)
	for l := 2; l <= 2*n+1; l++ {
		bimp[l] = make([]int, 4)
	}
	bsize = make([]int, 2*n+2)

	// Insert binary clauses into BIMP
	for _, clause := range clauses {
		// Check for clause of length 2
		if len(clause) == 2 {
			u, v := k2l(clause[0]), k2l(clause[1])

			if bsize[u^1] == len(bimp[u^1]) {
				bimp[u^1] = append(bimp[u^1], v)
			} else {
				bimp[u^1][bsize[u^1]] = v
			}
			bsize[u^1] += 1

			if bsize[v^1] == len(bimp[v^1]) {
				bimp[v^1] = append(bimp[v^1], u)
			} else {
				bimp[v^1][bsize[v^1]] = u
			}
			bsize[v^1] += 1
		}
	}

	//
	// Record all ternary clauses in the TIMP array
	//
	timp = make([]int, 2*n+2)
	tsize = make([]int, 2*n+2)

	// Get the values of TIMP[l] and TSIZE[l] for each l
	for l := 2; l <= 2*n+1; l++ {
		// Look for clauses containing this literal
		for _, clause := range clauses {
			// Check for clause of length 3
			if len(clause) == 3 {
				if l == k2l(-1*clause[0]) || l == k2l(-1*clause[1]) || l == k2l(-1*clause[2]) {
					// Found l in this clause
					if timp[l] == 0 {
						// This is the first clause in the list for l
						timp[l] = len(timp)
					}
					timp = append(timp, 0, 0)
					tsize[l] += 1
				}
			}
		}
	}

	// Add each clause to TIMP and set their LINK values
	link = make([]int, len(timp))
	tindex := make([]int, 2*n+2) // tindex[l] is the index for next insertion point in TIMP[l]

	for _, clause := range clauses {
		// Check for clause of length 3
		if len(clause) == 3 {
			u, v, w := k2l(clause[0]), k2l(clause[1]), k2l(clause[2])

			p = timp[u^1] + tindex[u^1]
			timp[p] = v
			timp[p+1] = w
			tindex[u^1] += 2

			pp = timp[v^1] + tindex[v^1]
			timp[pp] = u
			timp[pp+1] = w
			tindex[v^1] += 2

			ppp = timp[w^1] + tindex[w^1]
			timp[ppp] = u
			timp[ppp+1] = v
			tindex[w^1] += 2

			link[p] = pp
			link[pp] = ppp
			link[ppp] = p
		}
	}

	// Let U be the number of distinct variable in unit clauses

	// Terminate unsuccessfully if two unit clauses contradict each other

	// Record all distinct unit literals in FORCE[k] for 0 <= k < U

	// Configure initial permutation
	varx = make([]int, n)
	inx = make([]int, n+1)
	for k = 0; k < n; k++ {
		varx[k] = k + 1
		inx[k+1] = k
	}

	d = 0
	f = 0
	istackSize = 0
	istamp = 0

	if debug {
		dump()
	}

	if progress {
		showProgress()
	}

	return false, nil

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
