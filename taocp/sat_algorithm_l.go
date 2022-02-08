package taocp

import (
	"fmt"
	"log"
	"strings"
)

const (
	MaxInt = int(^uint(0) >> 1)
	rt     = MaxInt - 2 // RT - real truth
	nt     = MaxInt - 4 // NT - near truth
	pt     = MaxInt - 6 // PT - proto truth
)

// SatAlgorithmL implements Algorithm L (7.2.2.2), satisfiability by DPLL with lookahead.
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
		nOrig      int      // original value of n, before conversion to 3SAT
		m          int      // total number of clauses
		varx       []int    // VAR - permutation of {1,...,n} (VAR[k] = x iff INX[x] = k)
		inx        []int    // INX
		varN       int      // N - number of free variables in VAR
		d          int      // depth of the implicit search tree
		f          int      // number of fixed variables
		timp       []int    // TIMP - ternary clauses
		tsize      []int    // TSIZE - number of clauses for each l in TIMP
		link       []int    // LINK - circular list of the three literals in each clause in TIMP
		bimp       [][]int  // BIMP - instead of the buddy system, trying using built-in slices
		bsize      []int    // BSIZE - number of clauses for each l in BIMP
		p, pp, ppp int      // index into TIMP
		units      int      // U - number of distinct variables in unit clauses
		force      []int    // FORCE - stack of U unit variables which have a forced value
		istamp     int      // ISTAMP - stamp to make downdating BIMP tables easier
		ist        []int    // IST - private stamp for literal l
		istack     [][2]int // ISTACK - stack of previous values of (l, BSIZE[l])
		istackI    int      // I - size of ISTACK
		branch     []int    // BRANCH - record each branching decision
		dec        []int    // DEC - ??
		backf      []int    // BACKF - ??
		backi      []int    // BACKI - ??
		t          int      // T - truth context
		val        []int    // VAL - track if literal l is fixed in context T
		r          []int    // R - record the names of literals that have received values
		e          int      // E - current stack size of R; 0 <= E <= n
		g          int      // G - ??
		h          int      // H - ??
		conflict   int      // CONFLICT - algorithm L step to goto in case of conflict
		l          int      // literal l
		x          int      // X - variable
		ntL        int      // L - "nearly true" literal l
		j, k       int      // indices
		debug      bool     // debugging is enabled
		progress   bool     // progress tracking is enabled
	)

	fmt.Println(m, ist, istack, istackI, istamp, t, nt, rt, pt, r, e, g, val, conflict, h)

	// dump
	dump := func() {

		var b strings.Builder
		b.WriteString("\n")

		// FORCE
		b.WriteString("FORCE\n")
		b.WriteString(fmt.Sprintf("U=%d: ", units))
		for i := 0; i < units; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}", force[i]))
		}
		b.WriteString("\n\n")

		// VAR
		b.WriteString("VAR\n")
		b.WriteString(fmt.Sprintf("N=%d: ", varN))
		for k := 0; k < varN; k++ {
			if k > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}", varx[k]))
		}
		b.WriteString("\n\n")

		// VAL and R
		b.WriteString("VAL and R\n")
		b.WriteString(fmt.Sprintf("E=%d: ", e))
		for k := 0; k < e; k++ {
			if k > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}=%d", r[k], val[k]))
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

	// l2k
	// l2k := func(l int) int {
	// 	if l%2 == 0 {
	// 		return l >> 2
	// 	} else {
	// 		return (l >> 2) * -1
	// 	}
	// }

	// binary_propogation uses a siimsple breadth-first search procedure
	// to propagate the binarary consequences of a literal l inn context T
	// returns. Returns false if no conflict, true if there is conflict.
	// Formula (62)
	binary_propagation := func(l int) bool {
		h = e

		// Take account of l
		if val[l>>1] >= t {
			// l is fixed in context t
			if val[l>>1]&1 == l&1 {
				// l is fixed true, do nothing

			} else if val[l>>1] == (l^1)&1 {
				// l is fixed false, goto CONFLICT
				return true
			}
		} else {
			val[l>>1] = t + (l & 1)
			r[e] = l
			e += 1
		}

		for h < e {
			l = r[h]
			h += 1
			// For each l' in BIMP(l)
			for j := 0; j < bsize[l]; j++ {
				lp := bimp[l][j]

				// Take account of l'
				if val[lp>>1] >= t {
					// l' is fixed in context t
					if val[lp>>1]&1 == lp&1 {
						// l' is fixed true, do nothing

					} else if val[lp>>1] == (lp^1)&1 {
						// l' is fixed false, goto CONFLICT
						return true
					}
				} else {
					val[lp>>1] = t + (lp & 1)
					r[e] = lp
					e += 1
				}
			}
		}
		return false
	}

	// lvisit prepares the solution
	lvisit := func() []int {
		solution := make([]int, nOrig)
		for i := 0; i < n; i++ {
			l := force[i]
			solution[(l>>2)-1] = (l & 1) ^ 1
		}
		if debug {
			log.Printf("visit solution=%v", solution)
		}

		return solution
	}

	// appendBimp adds x to BIMP[l]
	appendBimp := func(l, x int) {

		// Update private stamp IST, if necessary. Formula (63)
		if ist[l] != istamp {
			ist[l] = istamp
			istack[istackI][0] = l
			istack[istackI][1] = bsize[l]
			istackI += 1
		}

		// Append x to l
		if bsize[l] == len(bimp[l]) {
			bimp[l] = append(bimp[l], x)
		} else {
			bimp[l][bsize[l]] = x
		}
		bsize[l] += 1
	}

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

	// Configure initial permutation of the "free variable" list, that is,
	// not fixed in context RT. A variable becomes fixed by swapping it to the
	// end of the free list and decreasing N; then we can free it later by
	// simply increasing N, without swapping.
	varx = make([]int, n)
	inx = make([]int, n+1)
	for k = 0; k < n; k++ {
		varx[k] = k + 1
		inx[k+1] = k
	}
	varN = n

	d = 0
	f = 0

	istamp = 0
	ist = make([]int, 2*n+2)
	istack = make([][2]int, 2*n+2)
	istackI = 0

	dec = make([]int, n)
	backf = make([]int, n)
	backi = make([]int, n)
	branch = make([]int, n)

	val = make([]int, n+1)
	r = make([]int, n+1)

	if debug {
		dump()
	}

	if progress {
		showProgress()
	}

	//
	// L2 [New node.]
	//
L2:
	if debug {
		log.Printf("L2. New node")
	}

	branch[d] = -1

	if units == 0 {
		//
		// Algorithm X
		//
		if f == n {
			// No variables are free, visit the solution

			if debug {
				log.Println("L2. [Success!]")
			}

			if stats != nil {
				stats.Solutions++
			}

			return true, lvisit()
		}

		// TODO: Goto L15 if Algorithm X discovers a conflict
	}

	if units > 0 {
		goto L5
	}

	//
	// L3 [Choose l.]
	//

	if debug {
		log.Printf("L3. Choose l")
	}

	// Choose whatever literal happens to be first in the current list
	// of free variables.
	l = varx[0]

	if l == 0 {
		d += 1
		goto L2
	}

	if debug {
		log.Printf("  d=%d, l=%d", d, l)
	}

	dec[d] = l
	backf[d] = f
	backi[d] = istackI
	branch[d] = 0

	//
	// L4 [Try l.]
	//

	if debug {
		log.Printf("L4. Try l")
	}

	units = 1
	force[0] = l

	//
	// L5 [Accept near truths.]
	//
L5:
	if debug {
		log.Printf("L5. Accept near truths")
	}

	t = nt
	g, e = f, f
	istamp += 1
	conflict = 11 // L11

	// Iterate over each l in the FORCE stack
	for i := 0; i < units; i++ {
		l := force[i]

		// Perform the binary propogation routine
		if binary_propagation(l) {
			// There was a conflict
			switch conflict {
			case 11:
				goto L11
			default:
				log.Panicf("Unknown value of CONFLICT: %d", conflict)
			}

		}
	}

	units = 0

	if debug {
		dump()
	}

	//
	// L6 [Choose a nearly true L.]
	//
L6:
	if debug {
		log.Printf("L6. Choose a nearly true L")
	}

	// At this point the stacked literals R_k are "really true" for 0 <= k < G,
	// and "nearly true" for G <= k < E. We want them all to be really true.
	if g == e {
		goto L10
	}

	ntL = r[g]
	g += 1

	//
	// L7 [Promote L to real truth.]
	//

	if debug {
		log.Printf("L7. Promote L=%d to real truth", ntL)
	}

	x = ntL >> 1
	val[x] = rt + ntL&1

	// Remove variable X from the free list and from all TIMP pairs (Exercise 137)
	varN = n - g
	x = varx[varN]
	j = inx[x]
	varx[j] = x
	inx[x] = j
	varx[varN] = x
	inx[x] = varN

	for _, l := range []int{2 * x, 2*x + 1} {
		for i := 0; i < tsize[l]; i++ {
			p := timp[l] + 2*i
			u, v := timp[p], timp[p+1]

			pp = link[p]
			ppp = link[pp]

			s := tsize[u^1] - 1
			tsize[u^1] = s
			t := timp[u^1] + 2*s // local t, not T

			if pp != t {
				// Swap pairs
				up, vp := timp[t], timp[t+1]
				q := link[t]
				qp := link[q]
				link[qp], link[p] = pp, t
				timp[pp], timp[pp+1] = up, vp
				link[pp] = q
				timp[t], timp[t+1] = v, l^1
				link[t] = ppp
				pp = t
			}

			s = tsize[v^1] - 1
			tsize[v^1] = s
			t = timp[v^1] + 2*s // local t, not T

			if ppp != t {
				// Swap pairs
				up, vp := timp[t], timp[t+1]
				q := link[t]
				qp := link[q]
				link[qp], link[pp] = ppp, t
				timp[ppp], timp[ppp+1] = up, vp
				link[ppp] = q
				timp[t], timp[t+1] = l^1, u
				link[t] = p
				pp = t
			}
		}
	}

	for i := 0; i < tsize[ntL]; i++ {
		p := timp[ntL] + 2*i
		u, v := timp[p], timp[p+1]

		//
		// L8 [Consider u or v.]
		//

		if debug {
			log.Printf("L8. Consider u or v")
		}

		// We have deduced that u or v must be true; five cases arise.
		// TODO: don't calculate these values until necessary

		uFixed := val[u>>1] >= t
		uFixedTrue := uFixed && val[u>>1]&1 == u&1
		uFixedFalse := uFixed && val[(u^1)>>1]&1 == (u^1)&1

		vFixed := val[v>>1] >= t
		vFixedTrue := vFixed && val[v>>1]&1 == v&1
		vFixedFalse := vFixed && val[(v^1)>>1]&1 == (v^1)&1

		if uFixedTrue || vFixedTrue {

			// Case 1. u or v is fixed true, do nothing

		} else if uFixedFalse && vFixedFalse {

			// Case 2. u and v are fixed false
			switch conflict {
			case 11:
				goto L11
			default:
				log.Panicf("Unknown value of CONFLICT: %d", conflict)
			}

		} else if uFixedFalse && !vFixed {

			// Case 3. u is fixed false but v isn't fixed
			if binary_propagation(v) {
				switch conflict {
				case 11:
					goto L11
				default:
					log.Panicf("Unknown value of CONFLICT: %d", conflict)
				}
			}

		} else if vFixedFalse && !uFixed {

			// Case 4. v is fixed false but u isn't fixed
			if binary_propagation(u) {
				switch conflict {
				case 11:
					goto L11
				default:
					log.Panicf("Unknown value of CONFLICT: %d", conflict)
				}
			}

		} else {

			// Case 5. Neither u nor v is fixed

			//
			// L9 [Exploit u or v.]
			//

			// TODO: Use Exercise 139 to improve this step by deducing
			// further implications called "compensation resolvents".

			if debug {
				log.Printf("L9. Exploit u or v")
			}

			var vInBimp, notvInBimp bool
			for i := 0; i < bsize[u^1]; i++ {
				if bimp[u^1][i] == v {
					vInBimp = true
				}
				if bimp[u^1][i] == v^1 {
					notvInBimp = true
				}
			}

			if notvInBimp {
				if binary_propagation(u) {
					switch conflict {
					case 11:
						goto L11
					default:
						log.Panicf("Unknown value of CONFLICT: %d", conflict)
					}
				}
			} else if vInBimp {
				// do nothing, we already have the clause u or v
			} else {

				var notuInBimp bool
				for i := 0; i < bsize[v^1]; i++ {
					if bimp[v^1][i] == u^1 {
						notuInBimp = true
					}
				}

				if notuInBimp {
					if binary_propagation(v) {
						switch conflict {
						case 11:
							goto L11
						default:
							log.Panicf("Unknown value of CONFLICT: %d", conflict)
						}
					}
				} else {
					// append v to BIMP[^u] and u to BIMP[^v]
					appendBimp(u^1, v)
					appendBimp(v^1, u)
				}

			}

		}

	}

	goto L6

	//
	// L10 [Accept real truths.]
	//
L10:
	if debug {
		log.Printf("L10. Accept real truths")
	}

	//
	// L11 [Unfix near truths.]
	//
L11:

	if debug {
		log.Printf("L11. Unfix near truths")
	}

	//
	// L12 [Unfix real truths.]
	//

	if debug {
		log.Printf("L12. Unfix real truths")
	}

	//
	// L13 [Downdate BIMPs.]
	//

	if debug {
		log.Printf("L13. Downdate BIMPs")
	}

	//
	// L14 [Try again?]
	//

	if debug {
		log.Printf("L14. Try again?")
	}

	//
	// L15 [Backtrack.]
	//

	if debug {
		log.Printf("L15. Backtrack")
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
