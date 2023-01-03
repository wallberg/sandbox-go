package taocp

import (
	"fmt"
	"log"
	"strings"
)

const (
	MaxInt = int(^uint(0) >> 1)
	rt     = MaxInt - 1 // RT - real truth
	nt     = MaxInt - 3 // NT - near truth
	pt     = MaxInt - 5 // PT - proto truth
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
func SatAlgorithmL(n int, clauses SatClauses,
	stats *SatStats, options *SatOptions) (sat bool, solution []int) {

	// @note Global variables
	var (
		// original value of n, before conversion to 3SAT
		nOrig int

		// N - number of free variables in VAR (F + N = n)
		N int

		// F - number of fixed variables (F + N = n)
		F int

		// VAR - free list; permutation of {1,...,n} (VAR[k] = x iff INX[x] = k)
		VAR []int

		// INX - index partner of VAR (free list)
		INX []int

		// X - variable of L promoted to real truth
		X int

		// d - depth of the implicit search tree
		d int

		// TIMP - ternary clauses
		TIMP []int

		// TSIZE - number of clauses for each l in TIMP
		TSIZE []int

		// LINK - circular list of the three literals in each clause in TIMP
		LINK []int

		// BIMP - binary clauses; instead of the buddy system, we are using built-in slices
		BIMP [][]int

		// BSIZE - number of clauses for each l in BIMP
		BSIZE []int

		// index into TIMP
		p, pp, ppp int

		// U - number of distinct variables in unit clauses (at the current depth)
		U int

		// FORCE - stack of U unit variables which have a forced value at the current depth
		FORCE []int

		// ISTAMP - stamp to make downdating BIMP tables easier
		ISTAMP int

		// IST - private stamp for literal l
		IST []int

		// ISTACK - stack of previous values of (l, BSIZE[l])
		ISTACK [][2]int

		// I - size of ISTACK
		I int

		// BRANCH - decision making at depth d; {-1: no decision yet, 0: trying l, 1: trying ^l)}
		BRANCH []int

		// DEC - decision on l at each branch
		DEC []int

		// BACKI - ??
		BACKI []int

		// BACKF - ?? (for showProgress(), Exercise 142)
		BACKF []int

		// BACKL - ?? (for showProgress(), Exercise 142)
		BACKL []int

		// T - truth context
		T int

		// L - nearly true literal l
		L int

		// VAL - track if variable x is fixed in context T
		VAL []int

		// R - record the names of literals that have received values
		R []int

		// E - current stack size of R; 0 <= E <= n
		E int

		// G - number of really true literals in R (starting from 0), and nearly true for G <= k < E
		G int

		// CONFLICT - algorithm L step to goto in case of CONFLICT
		CONFLICT int

		// literal l
		l int

		// variable x
		x int

		// index
		k, j int

		// is debugging enabled?
		debug bool

		// is progress tracking enabled?
		progress bool
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
		BACKL[d] = F
		r := 0
		k := 0

		for k < d {
			for r < BACKF[k] {
				b.WriteString(fmt.Sprintf("%d ", 6+(R[r]&1)))
				r += 1
			}
			if BRANCH[k] < 0 {
				b.WriteString("| ")
			} else {
				b.WriteString(fmt.Sprintf("%d ", (2*BRANCH[k])+R[r]&1))
				r += 1
			}
			for r < BACKL[k+1] {
				b.WriteString(fmt.Sprintf("%d ", 4+(R[r]&1)))
				r += 1
			}
			k += 1
		}

		// misc variables and the R stack
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("            n=%d, d=%d, F=%d, E=%d, G=%d\n", n, d, F, E, G))
		b.WriteString("            ")
		for k := 0; k < E; k++ {
			if k > 0 {
				b.WriteString(", ")
			}
			l := R[k]
			x := l >> 1
			b.WriteString(fmt.Sprintf("{%d}=%s", l, truth(VAL[x])))
		}
		b.WriteString("\n")

		// Statistics
		b.WriteString(fmt.Sprintf("            Nodes=%d, Levels=%v\n", stats.Nodes, stats.Levels))

		b.WriteString("\n")
		log.Print(b.String())
	}

	// dump
	dump := func() {

		var b strings.Builder
		b.WriteString("\n")

		showProgress()

		b.WriteString(fmt.Sprintf("n=%d, d=%d, f=%d\n", n, d, F))
		b.WriteString("\n")

		// FORCE
		b.WriteString("FORCE\n")
		b.WriteString(fmt.Sprintf("U=%d: ", U))
		for i := 0; i < U; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}", FORCE[i]))
		}
		b.WriteString("\n\n")

		// VAR
		b.WriteString("VAR\n")
		b.WriteString(fmt.Sprintf("N=%d: ", N))
		for k := 0; k < N; k++ {
			if k > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}", VAR[k]))
		}
		b.WriteString("\n\n")

		// R
		b.WriteString("R\n")
		b.WriteString(fmt.Sprintf("E=%d, G=%d: ", E, G))
		for k := 0; k < E; k++ {
			if k > 0 {
				b.WriteString(", ")
			}
			l := R[k]
			x := l >> 1
			b.WriteString(fmt.Sprintf("{%d}=%s", l, truth(VAL[x])))
		}
		b.WriteString("\n\n")

		// VAL
		b.WriteString("VAL\n")
		for x := 1; x <= n; x++ {
			if x > 1 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("{%d}=%s", x, truth(VAL[x])))
		}
		b.WriteString("\n\n")

		// BIMP
		b.WriteString("BIMP\n")
		for l := 2; l <= 2*n+1; l++ {
			b.WriteString(fmt.Sprintf("%d: ", l))
			for i := 0; i < BSIZE[l]; i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(fmt.Sprintf("%d", BIMP[l][i]))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")

		// TIMP
		b.WriteString("TIMP\n")
		for l := 2; l <= 2*n+1; l++ {
			b.WriteString(fmt.Sprintf("%d: ", l))
			for i := 0; i < TSIZE[l]; i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				p := TIMP[l] + 2*i
				b.WriteString(fmt.Sprintf("{%d,%d}->", TIMP[p], TIMP[p+1]))
				p = LINK[p]
				b.WriteString(fmt.Sprintf("{%d,%d}->", TIMP[p], TIMP[p+1]))
				p = LINK[p]
				b.WriteString(fmt.Sprintf("{%d,%d}", TIMP[p], TIMP[p+1]))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")

		log.Print(b.String())
	}

	// @note initialize()
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

	// assertRStackInvariant checks for the R stack invariant, that truth degrees never increase
	// as we move from the bottom to the top, using Formula (71), p. 227
	assertRStackInvariant := func() {
		for j := 1; j < E; j++ {
			if VAL[R[j-1]>>1]|1 < VAL[R[j]>>1] {
				dump()
				log.Fatal("assertion failed: violation of the R stack invariant")
			}
		}
	}

	// binary_propogation uses a simple breadth-first search procedure
	// to propagate the binarary consequences of a literal l inn context T
	// returns. Returns false if no conflict, true if there is conflict.
	// Formula (62), p. 221
	// @note binary_propogation()
	binary_propagation := func(l int) bool {

		if debug {
			log.Printf("  binary_propagation l=%d, t=%s", l, truth(T))
			assertRStackInvariant()
		}

		h := E

		// Take account of l
		if VAL[l>>1] >= T {
			// l is fixed in context t
			if VAL[l>>1]&1 == l&1 {
				// l is fixed true, do nothing

			} else if VAL[(l^1)>>1]&1 == (l^1)&1 { // TODO: change to 'else' ?
				// l is fixed false, goto CONFLICT
				return true
			}
		} else {
			VAL[l>>1] = T + (l & 1)
			R[E] = l
			E += 1
		}

		for h < E {
			l = R[h]
			h += 1
			// For each l' in BIMP(l)
			for j := 0; j < BSIZE[l]; j++ {
				lp := BIMP[l][j]

				// Take account of l'
				if VAL[lp>>1] >= T {
					// l' is fixed in context t
					if VAL[lp>>1]&1 == lp&1 {
						// l' is fixed true, do nothing

					} else if VAL[(lp^1)>>1] == (lp^1)&1 { // TODO: change to 'else' ?
						// l' is fixed false, goto CONFLICT
						return true
					}
				} else {
					VAL[lp>>1] = T + (lp & 1)
					R[E] = lp
					E += 1
				}
			}
		}

		if debug {
			assertRStackInvariant()
		}

		return false
	}

	// lvisit prepares the solution
	lvisit := func() []int {
		solution := make([]int, n)

		// Convert the literals from internal back to external format
		for i := 0; i < n; i++ {
			l := R[i]
			solution[(l>>1)-1] = (l & 1) ^ 1
		}
		if debug {
			log.Printf("visit solution=%v (%v)", solution[:nOrig], solution)
		}

		return solution[:nOrig]
	}

	// appendBimp adds x to BIMP[l]
	appendBimp := func(l, x int) {

		// Update private stamp IST, if necessary. Formula (63)
		if IST[l] != ISTAMP {
			IST[l] = ISTAMP
			if I == len(IST) {
				ISTACK = append(ISTACK, [2]int{l, BSIZE[l]})
			} else {
				ISTACK[I][0] = l
				ISTACK[I][1] = BSIZE[l]
			}
			I += 1
		}

		// Append x to l
		if BSIZE[l] == len(BIMP[l]) {
			BIMP[l] = append(BIMP[l], x)
		} else {
			BIMP[l][BSIZE[l]] = x
		}
		BSIZE[l] += 1
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

	// Convert the literals in each clause from external to internal format
	clausesInternal := make(SatClauses, len(clauses))
	for i, clause := range clauses {
		clausesInternal[i] = make(SatClause, len(clause))
		for j, k := range clause {
			if k < 0 {
				clausesInternal[i][j] = -2*k + 1
			} else {
				clausesInternal[i][j] = 2 * k
			}
		}
	}
	clauses = clausesInternal

	//
	// Record all unit clauses as forced variable values at depth 0
	//
	FORCE = make([]int, 2*n+2)
	U = 0
	for _, clause := range clauses {
		if len(clause) == 1 {
			l := clause[0]

			// Look for a contradiction
			for k := 0; k < U; k++ {
				if l^1 == FORCE[k] {
					// A contradiction
					if debug {
						log.Printf("L1. Found a unit clause contradiction")
					}
					return false, nil
				}
			}

			FORCE[U] = l
			U += 1

		}
	}

	//
	// Record all binary clauses in the BIMP array
	//
	BIMP = make([][]int, 2*n+2)
	for l := 2; l <= 2*n+1; l++ {
		BIMP[l] = make([]int, 4)
	}
	BSIZE = make([]int, 2*n+2)

	// Insert binary clauses into BIMP
	for _, clause := range clauses {
		// Check for clause of length 2
		if len(clause) == 2 {
			u, v := clause[0], clause[1]

			if BSIZE[u^1] == len(BIMP[u^1]) {
				BIMP[u^1] = append(BIMP[u^1], v)
			} else {
				BIMP[u^1][BSIZE[u^1]] = v
			}
			BSIZE[u^1] += 1

			if BSIZE[v^1] == len(BIMP[v^1]) {
				BIMP[v^1] = append(BIMP[v^1], u)
			} else {
				BIMP[v^1][BSIZE[v^1]] = u
			}
			BSIZE[v^1] += 1
		}
	}

	//
	// Record all ternary clauses in the TIMP array
	//
	TIMP = make([]int, 2*n+2)
	TSIZE = make([]int, 2*n+2)

	// Get the values of TIMP[l] and TSIZE[l] for each l
	for l := 2; l <= 2*n+1; l++ {
		// Look for clauses containing this literal
		for _, clause := range clauses {
			// Check for clause of length 3
			if len(clause) == 3 {
				u, v, w := clause[0], clause[1], clause[2]

				if l == u^1 || l == v^1 || l == w^1 {
					// Found l in this clause
					if TIMP[l] == 0 {
						// This is the first clause in the list for l
						TIMP[l] = len(TIMP)
					}
					TIMP = append(TIMP, 0, 0)
					TSIZE[l] += 1
				}
			}
		}
	}

	// Add each clause to TIMP and set their LINK values
	LINK = make([]int, len(TIMP))
	tindex := make([]int, 2*n+2) // tindex[l] is the index for next insertion point in TIMP[l]

	for _, clause := range clauses {
		// Check for clause of length 3
		if len(clause) == 3 {
			u, v, w := clause[0], clause[1], clause[2]

			p = TIMP[u^1] + tindex[u^1]
			TIMP[p] = v
			TIMP[p+1] = w
			tindex[u^1] += 2

			pp = TIMP[v^1] + tindex[v^1]
			TIMP[pp] = u
			TIMP[pp+1] = w
			tindex[v^1] += 2

			ppp = TIMP[w^1] + tindex[w^1]
			TIMP[ppp] = u
			TIMP[ppp+1] = v
			tindex[w^1] += 2

			LINK[p] = pp
			LINK[pp] = ppp
			LINK[ppp] = p
		}
	}

	// Configure initial permutation of the "free variable" list, that is,
	// not fixed in context RT. A variable becomes fixed by swapping it to the
	// end of the free list and decreasing N; then we can free it later by
	// simply increasing N, without swapping.
	VAR = make([]int, n)
	INX = make([]int, n+1)
	for k = 1; k <= n; k++ {
		VAR[k-1] = k
		INX[k] = k - 1
	}
	N = n

	d = 0
	F = 0

	ISTAMP = 0
	IST = make([]int, 2*n+2)
	ISTACK = make([][2]int, 1024) // Grow dynamically, when needed
	I = 0

	DEC = make([]int, n+1)
	BACKF = make([]int, n+1)
	BACKI = make([]int, n+1)
	BACKL = make([]int, n+1)
	BRANCH = make([]int, n+1)

	VAL = make([]int, n+1)
	R = make([]int, n+1)

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

	BRANCH[d] = -1 // No decision yet

	if debug || (progress && stats.Delta != 0 && stats.Nodes%stats.Delta == 0) {
		showProgress()
	}

	stats.Levels[d]++
	stats.Nodes++

	if U > 0 {
		goto L5
	}

	//
	// Algorithm X
	//

	// Iterate over each R stack entry, checking for contradictions
	for i := 0; i < E; i++ {
		l := R[i]

		// Iterate over l's BIMP table
		for j := 0; j < BSIZE[l]; j++ {
			lp := BIMP[l][j]

			// Look for a conflict between the BIMP table entry and an R stack entry
			for k := 0; k < E; k++ {
				lpp := R[k]

				if lp^1 == lpp {
					// A contradiction
					if debug && stats.Verbosity > 0 {
						dump()
						log.Printf("L2. BIMP table for %d in R contains %d, which contradicts %d in R", l, lp, lpp)
					}
					goto L15
				}
			}
		}
	}

	if F == n {
		// All variables are fixed, visit the solution

		if debug {
			log.Println("L2. [Success!]")
		}

		if stats != nil {
			stats.Solutions++
		}

		return true, lvisit()
	}

	// Choose whatever literal happens to be first in the current list
	// of free variables.
	x = VAR[0]
	l = 2 * x

	if debug {
		log.Printf("  Selected d=%d, branch=%v, l=%d from free variable list", d, BRANCH[0:d], l)
		dump()
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

	DEC[d] = l
	BACKF[d] = F
	BACKI[d] = I
	BRANCH[d] = 0 // We are trying l

	//
	// @note L4 [Try l.]
	//
L4:

	if debug {
		log.Printf("L4. Try l")
	}

	U = 1
	FORCE[0] = l

	//
	// @note L5 [Accept near truths.]
	//
L5:
	if debug {
		log.Printf("L5. Accept near truths")
	}

	T = nt
	G, E = F, F
	ISTAMP += 1
	CONFLICT = 11 // L11

	// Iterate over each l in the FORCE stack
	for i := 0; i < U; i++ {
		if debug && stats.Verbosity > 0 && i == 0 {
			log.Printf("State before beginning binary_propagation")
			dump()
		}

		l := FORCE[i]

		// Perform the binary propogation routine
		if binary_propagation(l) {
			// There was a conflict
			switch CONFLICT {
			case 11:
				goto L11
			default:
				log.Panicf("Unknown value of CONFLICT: %d", CONFLICT)
			}

		}
	}

	U = 0

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
		for k := 0; k < E; k++ {
			l := R[k]
			x := l >> 1

			if k < G && VAL[x]&(^1) != rt {
				log.Panicf("assertion failed: variable {%d}=%s is not RT at L6", x, truth(VAL[x]))

			} else if k >= G && VAL[x]&(^1) != nt {
				log.Panicf("assertion failed: variable {%d}=%s is not NT at L6", x, truth(VAL[x]))
			}
		}
	}

	if G == E {
		// No nearly true literals
		goto L10
	}

	L = R[G]
	G += 1

	//
	// @note L7 [Promote L to real truth.]
	//

	if debug {
		log.Printf("L7. Promote L=%d to real truth", L)
	}

	X = L >> 1
	VAL[X] = rt + L&1

	// Remove variable X from the free list (Exercise 137 (a))
	N = n - G
	x = VAR[N]
	j = INX[X]
	VAR[j] = x
	INX[x] = j
	VAR[N] = X
	INX[X] = N

	// Remove variable X from all TIMP pairs (Exercise 137 (a))
	for _, l := range []int{2 * X, 2*X + 1} {

		// For each pair in TIMP[l]
		for i := 0; i < TSIZE[l]; i++ {
			p = TIMP[l] + 2*i
			u, v := TIMP[p], TIMP[p+1]

			pp = LINK[p]
			ppp = LINK[pp]

			s := TSIZE[u^1] - 1
			TSIZE[u^1] = s
			t := TIMP[u^1] + 2*s // local t, not T

			if pp != t {
				// Swap pairs
				up, vp := TIMP[t], TIMP[t+1]
				q := LINK[t]
				qp := LINK[q]
				LINK[qp], LINK[p] = pp, t
				TIMP[pp], TIMP[pp+1] = up, vp
				LINK[pp] = q
				TIMP[t], TIMP[t+1] = v, l^1
				LINK[t] = ppp
				pp = t
			}

			s = TSIZE[v^1] - 1
			TSIZE[v^1] = s
			t = TIMP[v^1] + 2*s

			if ppp != t {
				// Swap pairs
				up, vp := TIMP[t], TIMP[t+1]
				q := LINK[t]
				qp := LINK[q]
				LINK[qp], LINK[pp] = ppp, t
				TIMP[ppp], TIMP[ppp+1] = up, vp
				LINK[ppp] = q
				TIMP[t], TIMP[t+1] = l^1, u
				LINK[t] = p
			}
		}
	}

	if debug && stats.Verbosity > 0 {
		dump()
	}

	for i := 0; i < TSIZE[L]; i++ {
		p := TIMP[L] + 2*i
		u, v := TIMP[p], TIMP[p+1]

		//
		// @note L8 [Consider u or v.]
		//

		if debug {
			log.Printf("L8. Consider u or v")
		}

		// We have deduced that u or v must be true; five cases arise.
		// TODO: don't calculate these values until necessary

		uFixed := VAL[u>>1] >= T
		uFixedTrue := uFixed && VAL[u>>1]&1 == u&1
		uFixedFalse := uFixed && VAL[(u^1)>>1]&1 == (u^1)&1

		vFixed := VAL[v>>1] >= T
		vFixedTrue := vFixed && VAL[v>>1]&1 == v&1
		vFixedFalse := vFixed && VAL[(v^1)>>1]&1 == (v^1)&1

		if uFixedTrue || vFixedTrue {

			// Case 1. u or v is fixed true, do nothing

		} else if uFixedFalse && vFixedFalse {

			// Case 2. u and v are fixed false
			switch CONFLICT {
			case 11:
				goto L11
			default:
				log.Panicf("Unknown value of CONFLICT: %d", CONFLICT)
			}

		} else if uFixedFalse && !vFixed {

			// Case 3. u is fixed false but v isn't fixed
			if binary_propagation(v) {
				switch CONFLICT {
				case 11:
					goto L11
				default:
					log.Panicf("Unknown value of CONFLICT: %d", CONFLICT)
				}
			}

		} else if vFixedFalse && !uFixed {

			// Case 4. v is fixed false but u isn't fixed
			if binary_propagation(u) {
				switch CONFLICT {
				case 11:
					goto L11
				default:
					log.Panicf("Unknown value of CONFLICT: %d", CONFLICT)
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
			for i := 0; i < BSIZE[u^1]; i++ {
				if BIMP[u^1][i] == v {
					vInBimp = true
				}
				if BIMP[u^1][i] == v^1 {
					notvInBimp = true
				}
			}

			if notvInBimp {
				if binary_propagation(u) {
					switch CONFLICT {
					case 11:
						goto L11
					default:
						log.Panicf("Unknown value of CONFLICT: %d", CONFLICT)
					}
				}
			} else if vInBimp {
				// do nothing, we already have the clause u or v
			} else {

				var notuInBimp bool
				for i := 0; i < BSIZE[v^1]; i++ {
					if BIMP[v^1][i] == u^1 {
						notuInBimp = true
					}
				}

				if notuInBimp {
					if binary_propagation(v) {
						switch CONFLICT {
						case 11:
							goto L11
						default:
							log.Panicf("Unknown value of CONFLICT: %d", CONFLICT)
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

	F = E

	if BRANCH[d] >= 0 {
		d += 1
		if debug {
			log.Printf("  branch[%d]=%d, incremented d to %d, going to L2", d-1, BRANCH[d-1], d)
		}
		goto L2
	} else if d > 0 {
		if debug {
			log.Printf("  branch[%d]=%d and d=%d > 0, going to L3", d, BRANCH[d], d)
		}
		goto L3
	} else { // d == 0
		if debug {
			log.Printf("  branch[%d]=%d and d=0, going to L2", d, BRANCH[d])
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

	for E > G {
		E -= 1
		VAL[R[E]>>1] = 0
	}

	//
	// @note L12 [Unfix real truths.]
	//
L12:
	if debug {
		log.Printf("L12. Unfix real truths")
	}

	if debug && stats.Verbosity > 0 {
		dump()
	}

	for E > F {
		// Implicitly restore X to the free list because N + E = n
		// (Exercise 137)
		E -= 1
		X = R[E] >> 1

		// Reactivate the TIMP pairs that involve X
		// (Exercise 137)
		for _, l = range []int{2*X + 1, 2 * X} {
			for i := TSIZE[l] - 1; i >= 0; i-- {
				p := TIMP[l] + 2*i
				u, v := TIMP[p], TIMP[p+1]

				TSIZE[v^1] += 1
				TSIZE[u^1] += 1
			}
		}
		VAL[X] = 0
	}

	//
	// @note L13 [Downdate BIMPs.]
	//

	if debug {
		log.Printf("L13. Downdate BIMPs")
	}

	if BRANCH[d] >= 0 {
		for I > BACKI[d] {
			I -= 1
			l, s := ISTACK[I][0], ISTACK[I][1]
			BSIZE[l] = s
		}
	}

	//
	// @note L14 [Try again?]
	//

	if debug {
		log.Printf("L14. Try again?")
	}

	// We've discovered that DEC[d] doesn't work
	if BRANCH[d] == 0 {
		l = DEC[d]
		DEC[d] = l ^ 1
		l = l ^ 1
		BRANCH[d] = 1 // l didn't work out, so try ^l

		if debug {
			log.Printf("  Trying again, d=%d, branch=%v, l=%d", d, BRANCH[0:d], l)
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
	E = F
	F = BACKF[d]
	goto L12
}
