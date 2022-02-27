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
		varx       []int    // VAR - permutation of {1,...,n} (VAR[k] = x iff INX[x] = k)
		inx        []int    // INX
		varN       int      // N - number of free variables in VAR
		varX       int      // X - variable of L promoted to real truth
		d          int      // d - depth of the implicit search tree
		f          int      // F - number of fixed variables
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
		branch     []int    // BRANCH - where are we in decision making (-1, 0, 1)
		dec        []int    // DEC - decision on l at each branch
		backf      []int    // BACKF - ??
		backi      []int    // BACKI - ??
		backl      []int    // BACKL - ?? (for showProgress(), Exercise 142)
		t          int      // T - truth context
		val        []int    // VAL - track if variable x is fixed in context T
		r          []int    // R - record the names of literals that have received values
		e          int      // E - current stack size of R; 0 <= E <= n
		g          int      // G - number of really true literals in R (starting from 0)		// and "nearly true" for G <= k < E
		h          int      // H - ??
		conflict   int      // CONFLICT - algorithm L step to goto in case of conflict
		l          int      // literal l
		x          int      // variable x
		ntL        int      // L - "nearly true" literal l
		k, j       int      // indices
		debug      bool     // debugging is enabled
		progress   bool     // progress tracking is enabled
	)

	// truth returns a string description of truth values
	truth := func(t int) string {
		switch t {
		case rt + 1:
			return "RF"
		case rt:
			return "RT"
		case nt + 1:
			return "NF"
		case nt:
			return "NT"
		case pt + 1:
			return "PF"
		case pt:
			return "PT"
		default:
			return fmt.Sprintf("%d", t)
		}
	}

	// showProgress - Exercise 142
	showProgress := func() {
		var b strings.Builder
		b.WriteString("  Progress: ")
		backl[d] = f
		localr := 0
		k := 0

		for k < d {
			for localr < backf[k] {
				b.WriteString(fmt.Sprintf("%d ", 6+(r[localr]&1)))
				localr += 1
			}
			if branch[k] < 0 {
				b.WriteString("| ")
			} else {
				b.WriteString(fmt.Sprintf("%d ", (2*branch[k])+r[localr]&1))
				localr += 1
			}
			for localr < backl[k+1] {
				b.WriteString(fmt.Sprintf("%d ", 4+(r[localr]&1)))
				localr += 1
			}
			k += 1
		}

		log.Print(b.String())
	}

	// dump
	dump := func() {

		var b strings.Builder
		b.WriteString("\n")

		showProgress()

		b.WriteString(fmt.Sprintf("n=%d, d=%d, f=%d, h=%d\n", n, d, f, h))
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

		// R
		b.WriteString("R\n")
		b.WriteString(fmt.Sprintf("E=%d, G=%d: ", e, g))
		for k := 0; k < e; k++ {
			if k > 0 {
				b.WriteString(", ")
			}
			l := r[k]
			x := l >> 1
			b.WriteString(fmt.Sprintf("{%d}=%s", l, truth(val[x])))
		}
		b.WriteString("\n\n")

		// VAL
		b.WriteString("VAL\n")
		for x := 1; x <= n; x++ {
			if x > 1 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}=%s", x, truth(val[x])))
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

	// binary_propogation uses a simple breadth-first search procedure
	// to propagate the binarary consequences of a literal l inn context T
	// returns. Returns false if no conflict, true if there is conflict.
	// Formula (62)
	binary_propagation := func(l int) bool {

		if debug {
			log.Printf("  binary_propagation l=%d", l)
		}
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
		solution := make([]int, n)

		for i := 0; i < n; i++ {
			l := r[i]
			solution[(l>>1)-1] = (l & 1) ^ 1
		}
		if debug {
			log.Printf("visit solution=%v (%v)", solution[:nOrig], solution[nOrig:])
		}

		return solution[:nOrig]
	}

	// appendBimp adds x to BIMP[l]
	appendBimp := func(l, x int) {

		// Update private stamp IST, if necessary. Formula (63)
		if ist[l] != istamp {
			ist[l] = istamp
			if istackI == len(ist) {
				istack = append(istack, [2]int{l, bsize[l]})
			} else {
				istack[istackI][0] = l
				istack[istackI][1] = bsize[l]
			}
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
	// @note L1 [Initialize.]
	//

	// Convert the input to 3SAT, if it isn't already
	nOrig = n
	_, n, clauses = Sat3(n, clauses)

	initialize()

	if debug {
		log.Printf("L1. Initialize")
	}

	//
	// Record all unit clauses with forced variable values
	//
	// TODO: Determine why L4 and L5 seem to wipe out what we've
	// done here without any restoration. Maybe instead of adding
	// to the force stack, we should add to the fixed stack or R stack.
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
	for k = 1; k <= n; k++ {
		varx[k-1] = k
		inx[k] = k - 1
	}
	varN = n

	d = 0
	f = 0

	istamp = 0
	ist = make([]int, 2*n+2)
	istack = make([][2]int, 1024) // Grow dynamically, when needed
	istackI = 0

	dec = make([]int, n+1)
	backf = make([]int, n+1)
	backi = make([]int, n+1)
	backl = make([]int, n+1)
	branch = make([]int, n+1)

	val = make([]int, n+1)
	r = make([]int, n+1)

	if debug && stats.Verbosity > 0 {
		dump()
	}

	//
	// @note L2 [New node.]
	//
L2:
	if debug {
		log.Printf("L2. New node")
	}

	branch[d] = -1 // No decision yet

	if debug || progress {
		showProgress()
	}

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
		// Choose whatever literal happens to be first in the current list
		// of free variables.
		x = varx[0]
		l = 2 * x

		if debug {
			log.Printf("  Trying d=%d, branch=%v, x=%d, l=%d from free variable list", d, branch[0:d+1], x, l)
		}

		contradiction := false

		// Record the forced values, looking for a contradiction
		for i := 0; !contradiction && i < bsize[l]; i++ {
			lp := bimp[l][i]

			// Look for a contradiction
			for k := 0; !contradiction && k < units; k++ {
				if lp^1 == force[k] || lp^1 == l {
					// A contradiction
					contradiction = true
				}
			}

			if !contradiction {
				// No contradiction, add it to the force stack
				force[units] = lp
				units += 1

				val[lp>>1] = rt
			}
		}

		if contradiction {
			// Try again with l^1
			l = l ^ 1
			branch[d] = 1
			units = 0

			if debug {
				log.Printf("  Trying d=%d, branch=%v, x=%d, l=%d from free variable list", d, branch[0:d+1], x, l)
			}

			// Record the forced values, looking for a contradiction
			for i := 0; i < bsize[l]; i++ {
				lp := bimp[l][i]

				// Look for a contradiction
				for k := 0; k < units; k++ {
					if lp^1 == force[k] || lp^1 == l {
						// A contradiction
						if debug && stats.Verbosity > 0 {
							dump()
							log.Printf("L2. Found unit clause contradictions; neither %d nor %d will work", l^1, l)
						}
						goto L15
					}
				}

				// No contradiction, add it to the force stack
				force[units] = lp
				units += 1

				val[lp>>1] = rt + lp&1
			}
		}

		if debug {
			log.Printf("  Selected d=%d, branch=%v, l=%d from free variable list", d, branch[0:d], l)
			dump()
		}

	} else { // units > 0
		goto L5
	}

	//
	// @note L3 [Choose l.]
	//
L3:
	if debug {
		log.Printf("L3. Choose l")
	}

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
	branch[d] = 0 // We are trying l

	//
	// @note L4 [Try l.]
	//
L4:

	if debug {
		log.Printf("L4. Try l")
	}

	units = 1
	force[0] = l

	//
	// @note L5 [Accept near truths.]
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

	if debug && stats.Verbosity > 0 {
		dump()
	}
	//
	// @note L6 [Choose a nearly true L.]
	//
L6:
	if debug {
		log.Printf("L6. Choose a nearly true L")
	}

	// At this point the stacked literals R_k are "really true" for 0 <= k < G,
	// and "nearly true" for G <= k < E. We want them all to be really true.

	if debug {
		// assertion
		for k := 0; k < g; k++ {
			l := r[k]
			x := l >> 1
			if val[x] != rt {
				log.Panicf("assertion failed: variable {%d}=%s", x, truth(val[x]))
			}
		}
	}

	if g == e {
		// No nearly true literals
		goto L10
	}

	ntL = r[g]
	g += 1

	//
	// @note L7 [Promote L to real truth.]
	//

	if debug {
		log.Printf("L7. Promote L=%d to real truth", ntL)
	}

	varX = ntL >> 1
	val[varX] = rt + ntL&1

	// Remove variable X from the free list (Exercise 137)
	varN = n - g
	x = varx[varN]
	j = inx[varX]
	varx[j] = x
	inx[x] = j
	varx[varN] = varX
	inx[varX] = varN

	// Remove variable X from all TIMP pairs (Exercise 137)
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

	if debug && stats.Verbosity > 0 {
		dump()
	}

	for i := 0; i < tsize[ntL]; i++ {
		p := timp[ntL] + 2*i
		u, v := timp[p], timp[p+1]

		//
		// @note L8 [Consider u or v.]
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
			// @note L9 [Exploit u or v.]
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
	// @note L10 [Accept real truths.]
	//
L10:
	if debug {
		log.Printf("L10. Accept real truths")
	}

	f = e

	if branch[d] >= 0 {
		d += 1
		if debug {
			log.Printf("  branch[%d]=%d, incremented d to %d, going to L2", d-1, branch[d-1], d)
		}
		goto L2
	} else if d > 0 {
		if debug {
			log.Printf("  branch[%d]=%d and d=%d > 0, going to L3", d, branch[d], d)
		}
		goto L3
	} else { // d == 0
		if debug {
			log.Printf("  branch[%d]=%d and d=0, going to L2", d, branch[d])
		}
		goto L2
	}

	//
	// @note L11 [Unfix near truths.]
	//
L11:

	if debug {
		log.Printf("L11. Unfix near truths")
	}

	for e > g {
		e -= 1
		val[r[e]>>1] = 0
	}

	//
	// @note L12 [Unfix real truths.]
	//
L12:
	if debug {
		log.Printf("L12. Unfix real truths")
	}

	for e > f {
		// Implicitly restore X to the free list because N + E = n
		// (Exercise 137)
		e -= 1
		varX = r[e] >> 1

		// Reactivate the TIMP pairs that involve X
		// (Exercise 137)
		for _, l := range []int{2 * varX, 2*varX + 1} {
			for i := tsize[l] - 1; i >= 0; i-- {
				p := timp[l] + 2*i
				u, v := timp[p], timp[p+1]

				tsize[v^1] += 1
				tsize[u^1] += 1
			}
		}
		val[x] = 0
	}

	//
	// @note L13 [Downdate BIMPs.]
	//

	if debug {
		log.Printf("L13. Downdate BIMPs")
	}

	if branch[d] >= 0 {
		for istackI > backi[d] {
			istackI -= 1
			l, s := istack[istackI][0], istack[istackI][1]
			bsize[l] = s
		}
	}

	//
	// @note L14 [Try again?]
	//

	if debug {
		log.Printf("L14. Try again?")
	}

	// We've discovered that DEC[d] doesn't work
	if branch[d] == 0 {
		l = dec[d]
		dec[d] = l ^ 1
		l = l ^ 1
		branch[d] = 1 // l didn't work out, so try ^l

		if debug {
			log.Printf("  Trying again, d=%d, branch=%v, l=%d", d, branch[0:d], l)
		}

		goto L4
	}

	//
	// @note L15 [Backtrack.]
	//
L15:
	if debug {
		log.Printf("L15. Backtrack")
	}

	if d == 0 {
		// Terminate unsuccessfully
		return false, nil
	}

	d -= 1
	e = f
	f = backf[d]
	goto L12
}
