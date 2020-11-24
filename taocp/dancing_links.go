package taocp

import (
	"fmt"
	"strconv"
	"strings"
)

// Explore Dancing Links from The Art of Computer Programming, Volume 4,
// Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
// Dancing Links, 2020
//
// §7.2.2.1 Dancing Links

// Stats is a struct for tracking ExactCover statistics and reporting
// runtime progress
type Stats struct {
	Progress  bool  // Display runtime progress
	MaxLevel  int   // Maximum level reached
	Delta     int   // Display progress every Delta number of Nodes
	Theta     int   // Display progress at next Theta number of Nodes
	Levels    []int // Count of times each level is entered
	Nodes     int   // Count of nodes processed
	Solutions int   // Count of solutions returned
	Debug     bool  // Enable debug logging
	Verbosity int   // Debug verbosity level (0 or 1)
}

func (s Stats) String() string {
	// Find first non-zero level count
	i := len(s.Levels)
	for s.Levels[i-1] == 0 && i > 0 {
		i--
	}

	return fmt.Sprintf("nodes=%d, solutions=%d, levels=%v", s.Nodes,
		s.Solutions, s.Levels[:i])
}

// ExactCover implements Algorithm X, exact cover via dancing links.
//
// Arguments:
// items     -- sorted list of primary items
// options   -- list of list of options; every option must contain at least one
// 			    primary item
// secondary -- sorted list of secondary items
// stats     -- structure to capture runtime statistics and provide feedback on
//              progress
// visit     -- function called with each discovered solution, returns true
//              if the search should resume
//
func ExactCover(items []string, options [][]string, secondary []string,
	stats *Stats, visit func(solution [][]string) bool) {

	var (
		n1    int      // number of primary items
		n2    int      // number of secondary items
		n     int      // total number of items
		name  []string // name of the item
		llink []int    // right link of the item
		rlink []int    // left link of the item
		top   []int
		llen  []int
		ulink []int
		dlink []int
		level int
		state []int // search state
		debug bool  // is debug enabled?
	)

	dump := func() {
		i := 0
		for rlink[i] != 0 {
			i = rlink[i]
			fmt.Print("  ", name[i])
			x := i
			for dlink[x] != i {
				x = dlink[x]
				fmt.Print(" ", name[top[x]])
			}
			fmt.Println()
		}
		fmt.Println("---")
	}

	initialize := func() {
		n1 = len(items)
		n2 = len(secondary)
		n = n1 + n2

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
		}

		// Fill out the item tables
		name = make([]string, n+2)
		llink = make([]int, n+2)
		rlink = make([]int, n+2)

		for j, item := range append(items, secondary...) {
			i := j + 1
			name[i] = item
			llink[i] = i - 1
			rlink[i-1] = i
		}

		// two doubly linked lists, primary and secondary
		// head of the primary list is at i=0
		// head of the secondary list is at i=n+1
		llink[n+1] = n
		rlink[n] = n + 1
		llink[n1+1] = n + 1
		rlink[n+1] = n1 + 1
		llink[0] = n1
		rlink[n1] = 0

		if debug {
			fmt.Println("name", name)
			fmt.Println("llink", llink)
			fmt.Println("rlink", rlink)
		}

		// Fill out the option tables
		nOptions := len(options)
		nOptionItems := 0
		for _, option := range options {
			nOptionItems += len(option)
		}
		size := n + 1 + nOptions + 1 + nOptionItems

		top = make([]int, size)
		llen = top[0 : n+1] // first n+1 elements of top
		ulink = make([]int, size)
		dlink = make([]int, size)

		// Set empty list for each item
		for i := 1; i <= n; i++ {
			llen[i] = 0
			ulink[i] = i
			dlink[i] = i
		}

		// Insert each of the options and their items
		x := n + 1
		spacer := 0
		top[x] = spacer
		spacerX := x

		// Iterate over each option
		for _, option := range options {
			// Iterate over each item in this option
			for _, item := range option {
				x++
				i := 0
				for _, value := range name {
					if value == item {
						break
					}
					i++
				}
				top[x] = i

				// Insert into the option list for this item
				llen[i]++ // increase the size by one
				head := i
				tail := i
				for dlink[tail] != head {
					tail = dlink[tail]
				}

				dlink[tail] = x
				ulink[x] = tail

				ulink[head] = x
				dlink[x] = head
			}

			// Insert spacer at end of each option
			dlink[spacerX] = x
			x++
			ulink[x] = spacerX + 1

			spacer--
			top[x] = spacer
			spacerX = x
		}

		if debug {
			fmt.Println("top", top)
			fmt.Println("llen", llen)
			fmt.Println("ulink", ulink)
			fmt.Println("dlink", dlink)
		}

		level = 0
		state = make([]int, nOptions)

		if debug {
			dump()
		}
	}

	showProgress := func() {

		est := 0.0 // estimate of percentage done
		tcum := 1

		fmt.Printf("Current level %d of max %d\n", level, stats.MaxLevel)

		// Iterate over the options
		for _, p := range state[0:level] {
			// Cyclically gather the items in the option, beginning at p
			fmt.Print("  ")
			q := p
			for {
				fmt.Print(name[top[q]] + " ")
				q++
				if top[q] <= 0 {
					q = ulink[q]
				}
				if q == p {
					break
				}
			}

			// Get position stats for this option
			i := top[p]
			q = dlink[i]
			k := 1
			for q != p && q != i {
				q = dlink[q]
				k++
			}

			if q != i {
				fmt.Printf(" %d of %d\n", k, llen[i])
				tcum *= llen[i]
				est += float64(k-1) / float64(tcum)
			} else {
				fmt.Println(" not in this list")
			}
		}

		est += 1.0 / float64(2*tcum)

		fmt.Printf("  solutions=%d, nodes=%d, est=%4.4f\n",
			stats.Solutions, stats.Nodes, est)
		fmt.Println("---")
	}

	lvisit := func() bool {
		// Iterate over the options
		options := make([][]string, 0)
		for i, p := range state[0:level] {
			options = append(options, make([]string, 0))
			// Move back to first element in the option
			for top[p-1] > 0 {
				p--
			}
			// Iterate over elements in the option
			q := p
			for top[q] > 0 {
				options[i] = append(options[i], name[top[q]])
				q++
			}
		}

		return visit(options)
	}

	// mrv selects the next item to try using the Minimum Remaining
	// Values heuristic.
	mrv := func() int {

		i := 0
		theta := -1
		p := rlink[0]
		for p != 0 {
			lambda := llen[p]
			if lambda < theta || theta == -1 {
				theta = lambda
				i = p
				if theta == 0 {
					return i
				}
			}
			p = rlink[p]
		}

		return i
	}

	hide := func(p int) {
		q := p + 1
		for q != p {
			x := top[q]
			u, d := ulink[q], dlink[q]
			if x <= 0 {
				q = u // q was a spacer
			} else {
				dlink[u], ulink[d] = d, u
				llen[x]--
				q++
			}
		}
	}

	cover := func(i int) {
		p := dlink[i]
		for p != i {
			hide(p)
			p = dlink[p]
		}
		l, r := llink[i], rlink[i]
		rlink[l], llink[r] = r, l
	}

	unhide := func(p int) {
		q := p - 1
		for q != p {
			x := top[q]
			u, d := ulink[q], dlink[q]
			if x <= 0 {
				q = d // q was a spacer
			} else {
				dlink[u], ulink[d] = q, q
				llen[x]++
				q--
			}
		}
	}

	uncover := func(i int) {
		l, r := llink[i], rlink[i]
		rlink[l], llink[r] = i, i
		p := ulink[i]
		for p != i {
			unhide(p)
			p = ulink[p]
		}
	}

	// X1 [Initialize.]
	initialize()

	var (
		i int
		j int
		p int
	)

X2:
	// X2. [Enter level l.]
	if debug {
		fmt.Printf("X2. level=%d, x=%v\n", level, state[0:level])
	}

	if stats != nil {
		stats.Levels[level]++
		stats.Nodes++

		if stats.Progress {
			if level > stats.MaxLevel {
				stats.MaxLevel = level
			}
			if stats.Nodes >= stats.Theta {
				showProgress()
				stats.Theta += stats.Delta
			}
		}
	}

	if rlink[0] == 0 {
		// visit the solution
		if debug {
			fmt.Println("X2. Visit the solution")
		}
		if stats != nil {
			stats.Solutions++
		}
		resume := lvisit()
		if !resume {
			if debug {
				fmt.Println("X2. Halting the search")
			}
			return
		}
		goto X8
	}

	// X3. [Choose i.]
	i = mrv()

	if debug {
		fmt.Printf("X3. Choose i=%d (%s)\n", i, name[i])
	}

	// X4. [Cover i.]
	if debug {
		fmt.Printf("X4. Cover i=%d (%s)\n", i, name[i])
	}
	cover(i)
	state[level] = dlink[i]

X5:
	// X5. [Try x_l.]
	if debug {
		fmt.Printf("X5. Try l=%d, x[l]=%d\n", level, state[level])
	}
	if state[level] == i {
		goto X7
	}
	p = state[level] + 1
	for p != state[level] {
		j := top[p]
		if j <= 0 {
			p = ulink[p]
		} else {
			cover(j)
			p++
		}
	}
	level++
	goto X2

X6:
	// X6. [Try again.]
	if debug {
		fmt.Println("X6. Try again")
	}

	if stats != nil {
		stats.Nodes++
	}

	p = state[level] - 1
	for p != state[level] {
		j = top[p]
		if j <= 0 {
			p = dlink[p]
		} else {
			uncover(j)
			p--
		}
	}
	i = top[state[level]]
	state[level] = dlink[state[level]]
	goto X5

X7:
	// X7. [Backtrack.]
	if debug {
		fmt.Println("X7. Backtrack")
	}
	uncover(i)

X8:
	// X8. [Leave level l.]
	if debug {
		fmt.Printf("X8. Leaving level %d\n", level)
	}
	if level == 0 {
		return
	}
	level--
	goto X6
}

// ExactCoverColors implements Algorithm C, exact covering with colors via
// dancing links.
//
// Arguments:
// items     -- sorted list of primary items
// options   -- list of list of options; every option must contain at least one
// 			    primary item
// secondary -- sorted list of secondary items; can contain an optional
//              "color" appended after a colon, eg "sitem:color"
// stats     -- structure to capture runtime statistics and provide feedback on
//              progress
// visit     -- function called with each discovered solution, returns true
//              if the search should continue
//
func ExactCoverColors(items []string, options [][]string, secondary []string,
	stats *Stats, visit func(solution [][]string) bool) {

	var (
		n1       int      // number of primary items
		n2       int      // number of secondary items
		n        int      // total number of items
		name     []string // name of the item
		llink    []int    // right link of the item
		rlink    []int    // left link of the item
		top      []int
		llen     []int
		ulink    []int
		dlink    []int
		color    []int    // color of a particular item in option
		colors   []string // map of color names, key is the index starting at 1
		level    int
		state    []int // search state
		debug    bool  // is debug enabled?
		progress bool  // is progress enabled?
	)

	dump := func() {
		fmt.Println("----------------------")

		// Tables
		fmt.Println("  name", name)
		fmt.Println("  llink", llink)
		fmt.Println("  rlink", rlink)
		fmt.Println("  top", top)
		fmt.Println("  llen", llen)
		fmt.Println("  ulink", ulink)
		fmt.Println("  dlink", dlink)
		fmt.Println("  color", color)
		fmt.Print("  colors")
		for i, colorName := range colors {
			if i > 0 {
				fmt.Printf(" %d=%s", i, colorName)
			}
		}
		fmt.Println()

		// Remaining items
		fmt.Print("  items:")
		i := 0
		for rlink[i] != 0 {
			i = rlink[i]
			fmt.Print(" ", name[i])
		}
		fmt.Println()

		// Selected options
		for i, p := range state[0:level] {
			fmt.Printf("  option: i=%d, p=%d (", i, p)
			// Move back to first element in the option
			for top[p-1] > 0 {
				p--
			}
			// Iterate over elements in the option
			q := p
			for top[q] > 0 {
				fmt.Print(" ", name[top[q]])
				q++
			}
			fmt.Println(" )")
		}
		fmt.Println("----------------------")

	}

	initialize := func() {
		n1 = len(items)
		n2 = len(secondary)
		n = n1 + n2

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

		// Fill out the item tables
		name = make([]string, n+2)
		llink = make([]int, n+2)
		rlink = make([]int, n+2)

		for j, item := range append(items, secondary...) {
			i := j + 1
			name[i] = item
			llink[i] = i - 1
			rlink[i-1] = i
		}

		// two doubly linked lists, primary and secondary
		// head of the primary list is at i=0
		// head of the secondary list is at i=n+1
		llink[n+1] = n
		rlink[n] = n + 1
		llink[n1+1] = n + 1
		rlink[n+1] = n1 + 1
		llink[0] = n1
		rlink[n1] = 0

		// Fill out the option tables
		nOptions := len(options)
		nOptionItems := 0
		for _, option := range options {
			nOptionItems += len(option)
		}
		size := n + 1 + nOptions + 1 + nOptionItems

		top = make([]int, size)
		llen = top[0 : n+1] // first n+1 elements of top
		ulink = make([]int, size)
		dlink = make([]int, size)
		color = make([]int, size)
		colors = make([]string, 1)

		// Set empty list for each item
		for i := 1; i <= n; i++ {
			llen[i] = 0
			ulink[i] = i
			dlink[i] = i
		}

		// Insert each of the options and their items
		x := n + 1
		spacer := 0
		top[x] = spacer
		spacerX := x

		// Iterate over each option
		for _, option := range options {
			// Iterate over each item in this option
			for _, item := range option {
				x++

				// Extract the color
				itemColor := 0 // 0 if there is no color
				iColon := strings.Index(item, ":")
				if iColon > -1 {
					itemColorName := item[iColon+1:]
					item = item[:iColon]

					// Insert the color name into color[]
					for itemColor = 1; itemColor < len(colors); itemColor++ {
						if itemColorName == colors[itemColor] {
							break
						}
					}
					if itemColor == len(colors) {
						// Not found, add new color name entry
						colors = append(colors, itemColorName)
					}
				}

				// Insert the item into name[]
				i := 0
				for _, value := range name {
					if value == item {
						break
					}
					i++
				}
				top[x] = i
				color[x] = itemColor

				// Insert into the option list for this item
				llen[i]++ // increase the size by one
				head := i
				tail := i
				for dlink[tail] != head {
					tail = dlink[tail]
				}

				dlink[tail] = x
				ulink[x] = tail

				ulink[head] = x
				dlink[x] = head
			}

			// Insert spacer at end of each option
			dlink[spacerX] = x
			x++
			ulink[x] = spacerX + 1

			spacer--
			top[x] = spacer
			spacerX = x
		}

		level = 0
		state = make([]int, nOptions)

		if debug {
			dump()
		}
	}

	showProgress := func() {

		est := 0.0 // estimate of percentage done
		tcum := 1

		fmt.Printf("Current level %d of max %d\n", level, stats.MaxLevel)

		// Iterate over the options
		for _, p := range state[0:level] {
			// Cyclically gather the items in the option, beginning at p
			fmt.Print("  ")
			q := p
			for {
				fmt.Print(name[top[q]] + " ")
				q++
				if top[q] <= 0 {
					q = ulink[q]
				}
				if q == p {
					break
				}
			}

			// Get position stats for this option
			i := top[p]
			q = dlink[i]
			k := 1
			for q != p && q != i {
				q = dlink[q]
				k++
			}

			if q != i {
				fmt.Printf(" %d of %d\n", k, llen[i])
				tcum *= llen[i]
				est += float64(k-1) / float64(tcum)
			} else {
				fmt.Println(" not in this list")
			}
		}

		est += 1.0 / float64(2*tcum)

		fmt.Printf("  solutions=%d, nodes=%d, est=%4.4f\n",
			stats.Solutions, stats.Nodes, est)
		fmt.Println("---")
	}

	lvisit := func() bool {

		if debug && stats.Verbosity > 0 {
			dump()
		}

		// Only one of the secondary items will have it's color value, the
		// others will have -1. Save the color and add it to all the matching
		// secondary items at the end.
		sitemColor := make(map[string]string)

		// Iterate over the options
		options := make([][]string, 0)
		for i, p := range state[0:level] {
			options = append(options, make([]string, 0))
			// Move back to first element in the option
			for top[p-1] > 0 {
				p--
			}
			// Iterate over elements in the option
			q := p
			for top[q] > 0 {
				name := name[top[q]]
				if color[q] > 0 {
					sitemColor[name] = colors[color[q]]
				}
				options[i] = append(options[i], name)
				q++
			}
		}

		// Add the secondary item colors
		for i, option := range options {
			for j, item := range option {
				if color, ok := sitemColor[item]; ok {
					options[i][j] += ":" + color
				}
			}
		}

		return visit(options)
	}

	// mrv selects the next item to try using the Minimum Remaining
	// Values heuristic.
	mrv := func() int {

		i := 0
		theta := -1
		p := rlink[0]
		for p != 0 {
			lambda := llen[p]
			if lambda < theta || theta == -1 {
				theta = lambda
				i = p
				if theta == 0 {
					return i
				}
			}
			p = rlink[p]
		}

		return i
	}

	// hide removes an option from further consideration
	hide := func(p int) {
		// iterate over the items in this option
		q := p + 1
		for q != p {
			x := top[q]
			u, d := ulink[q], dlink[q]
			if x <= 0 {
				q = u // q was a spacer
			} else {
				if color[q] >= 0 {
					dlink[u], ulink[d] = d, u
					llen[x]--
				}
				q++
			}
		}
	}

	unhide := func(p int) {
		q := p - 1
		for q != p {
			x := top[q]
			u, d := ulink[q], dlink[q]
			if x <= 0 {
				q = d // q was a spacer
			} else {
				if color[q] >= 0 {
					dlink[u], ulink[d] = q, q
					llen[x]++
				}
				q--
			}
		}
	}

	// cover removes i from the list of items needing to be covered removes and
	// hides all of the item's options
	cover := func(i int) {
		// hide all of the item's options
		p := dlink[i]
		for p != i {
			hide(p)
			p = dlink[p]
		}
		// remove the item from the list
		l, r := llink[i], rlink[i]
		rlink[l], llink[r] = r, l
	}

	uncover := func(i int) {
		l, r := llink[i], rlink[i]
		rlink[l], llink[r] = i, i
		p := ulink[i]
		for p != i {
			unhide(p)
			p = ulink[p]
		}
	}

	purify := func(p int) {
		c := color[p]
		i := top[p]
		color[i] = c
		q := dlink[i]
		for q != i {
			if color[q] == c {
				color[q] = -1
			} else {
				hide(q)
			}
			q = dlink[q]
		}
	}

	unpurify := func(p int) {
		c := color[p]
		i := top[p]
		q := ulink[i]
		for q != i {
			if color[q] < 0 {
				color[q] = c
			} else {
				unhide(q)
			}
			q = ulink[q]
		}
	}

	commit := func(p int, j int) {
		if color[p] == 0 {
			cover(j)
		}
		if color[p] > 0 {
			purify(p)
		}
	}

	uncommit := func(p int, j int) {
		if color[p] == 0 {
			uncover(j)
		}
		if color[p] > 0 {
			unpurify(p)
		}
	}

	// C1 [Initialize.]
	initialize()

	var (
		i int
		j int
		p int
	)

	if progress {
		showProgress()
	}

C2:
	// C2. [Enter level l.]
	if debug {
		fmt.Printf("C2. level=%d, x=%v\n", level, state[0:level])
	}

	if stats != nil {
		stats.Levels[level]++
		stats.Nodes++

		if progress {
			if level > stats.MaxLevel {
				stats.MaxLevel = level
			}
			if stats.Nodes >= stats.Theta {
				showProgress()
				stats.Theta += stats.Delta
			}
		}
	}

	if rlink[0] == 0 {
		// visit the solution
		if debug {
			fmt.Println("C2. Visit the solution")
		}
		if stats != nil {
			stats.Solutions++
		}
		resume := lvisit()
		if !resume {
			if debug {
				fmt.Println("C2. Halting the search")
			}
			if progress {
				showProgress()
			}
			return
		}
		goto C8
	}

	// C3. [Choose i.]
	i = mrv()

	if debug {
		fmt.Printf("C3. Choose i=%d (%s)\n", i, name[i])
	}

	// C4. [Cover i.]
	if debug {
		fmt.Printf("C4. Cover i=%d (%s)\n", i, name[i])
	}
	cover(i)
	state[level] = dlink[i]

C5:
	// C5. [Try x_l.]
	if debug {
		fmt.Printf("C5. Try l=%d, x[l]=%d\n", level, state[level])
	}
	if state[level] == i {
		goto C7
	}
	// Commit each of the items in this option
	p = state[level] + 1
	for p != state[level] {
		j := top[p]
		if j <= 0 {
			// spacer, go back to the first option
			p = ulink[p]
		} else {
			commit(p, j)
			p++
		}
	}
	level++
	goto C2

C6:
	// C6. [Try again.]
	if debug {
		fmt.Println("C6. Try again")
	}

	if stats != nil {
		stats.Nodes++
	}

	// Uncommit each of the items in this option
	p = state[level] - 1
	for p != state[level] {
		j = top[p]
		if j <= 0 {
			p = dlink[p]
		} else {
			uncommit(p, j)
			p--
		}
	}
	i = top[state[level]]
	state[level] = dlink[state[level]]
	goto C5

C7:
	// C7. [Backtrack.]
	if debug {
		fmt.Println("C7. Backtrack")
	}
	uncover(i)

C8:
	// C8. [Leave level l.]
	if debug {
		fmt.Printf("C8. Leaving level %d\n", level)
	}
	if level == 0 {
		if progress {
			showProgress()
		}
		return
	}
	level--
	goto C6
}

// LangfordPairs uses ExactCover to return solutions for Langford pairs
// of n values
func LangfordPairs(n int, stats *Stats, visit func(solution []int) bool) {

	// Build the list of items
	items := make([]string, 3*n)
	for i := 0; i < n; i++ {
		items[i] = strconv.Itoa(i + 1)
	}
	for i := 0; i < 2*n; i++ {
		items[n+i] = "s" + strconv.Itoa(i)
	}

	// Build the list of options
	options := make([][]string, 0)
	for i := 1; i <= n; i++ {
		j := 1
		k := j + i + 1
		for k <= 2*n {
			// Exercise 15: Omit the reversals
			x := 0
			if n%2 == 0 {
				x = 1
			}
			if i != n-x || j <= n/2 {
				options = append(options,
					[]string{
						strconv.Itoa(i),
						"s" + strconv.Itoa(j-1),
						"s" + strconv.Itoa(k-1),
					})
			}
			j++
			k++
		}
	}

	// Generate solutions
	ExactCover(items, options, []string{}, stats,
		func(solution [][]string) bool {
			x := make([]int, 2*n)
			for _, option := range solution {
				value, _ := strconv.Atoi(option[0])
				i, _ := strconv.Atoi(string(option[1][1]))
				j, _ := strconv.Atoi(string(option[2][1]))

				x[i], x[j] = value, value
			}
			return visit(x)
		})
}

// NQueens uses ExactCover to return solutions for the n-queens problem
func NQueens(n int, stats *Stats, visit func(solution []string) bool) {

	items := make([]string, 2*n)
	sitems := make([]string, 4*n-2)
	options := make([][]string, n*n)

	k := 0
	for i := 0; i < n; i++ {
		row := "r" + strconv.Itoa(i+1)
		items[i] = row
		for j := 0; j < n; j++ {
			col := "c" + strconv.Itoa(j+1)
			if i == n-1 {
				items[j+n] = col
			}
			upDiag := "a" + strconv.Itoa(i+j+2)
			downDiag := "b" + strconv.Itoa(i-j)

			var x int

			for x = 0; sitems[x] != upDiag && sitems[x] != ""; x++ {
			}
			if sitems[x] == "" {
				sitems[x] = upDiag
			}

			for x = 0; sitems[x] != downDiag && sitems[x] != ""; x++ {
			}
			if sitems[x] == "" {
				sitems[x] = downDiag
			}

			options[k] = []string{row, col, upDiag, downDiag}
			k++
		}
	}

	// Generate solutions
	ExactCover(items, options, sitems, stats,
		func(solution [][]string) bool {
			x := make([]string, 2*n)
			i := 0
			for _, option := range solution {
				x[i] = option[0]   // row
				x[i+1] = option[1] // col
				i += 2
			}
			return visit(x)
		})
}

// Sudoku uses ExactCover to solve 9x9 sudoku puzzles
func Sudoku(grid [9][9]int, stats *Stats,
	visit func(solution [9][9]int) bool) {

	var (
		i int // row number
		j int // column number
		k int // cell value in (row,column)
		x int // 3x3 box
	)

	// Build the [p, r, c, b] option
	buildOption := func() []string {
		return []string{
			"p" + strconv.Itoa(i) + strconv.Itoa(j), // piece
			"r" + strconv.Itoa(i) + strconv.Itoa(k), // piece in row
			"c" + strconv.Itoa(j) + strconv.Itoa(k), // piece in column
			"x" + strconv.Itoa(x) + strconv.Itoa(k), // piece in 3x3 box
		}
	}

	// Get the known items (non zero) provided in the grid
	knownItems := make(map[string]bool)
	for i = 0; i < 9; i++ {
		for j = 0; j < 9; j++ {
			k = grid[i][j]
			if k > 0 {
				x = 3*(i/3) + (j / 3)
				for _, item := range buildOption() {
					knownItems[item] = true
				}
			}
		}
	}

	// Build the items and options from the unknown values
	itemSet := make(map[string]bool)
	options := make([][]string, 0)
	for i = 0; i < 9; i++ {
		for j = 0; j < 9; j++ {
			x = 3*(i/3) + (j / 3)
			for k = 1; k < 10; k++ {
				option := buildOption()
				if !(knownItems[option[0]] || knownItems[option[1]] ||
					knownItems[option[2]] || knownItems[option[3]]) {
					for _, item := range option {
						itemSet[item] = true
					}
					options = append(options, option)
				}
			}
		}
	}

	items := make([]string, len(itemSet))
	i = 0
	for item := range itemSet {
		items[i] = item
		i++
	}

	// Generate solutions
	ExactCover(items, options, []string{}, stats,
		func(solution [][]string) bool {
			// Make a copy of the original grid
			var x [9][9]int
			for i := 0; i < 9; i++ {
				for j := 0; j < 9; j++ {
					x[i][j] = grid[i][j]
				}
			}

			// Fill in the solution values
			for _, option := range solution {
				i, _ := strconv.Atoi(string(option[0][1]))
				j, _ := strconv.Atoi(string(option[0][2]))
				k, _ := strconv.Atoi(string(option[1][2]))

				x[i][j] = k
			}

			return visit(x)
		})
}

// Exercise 7.2.2.1-66 Construct sudoku puzzles by placing nine given cars in a
// 3x3 array

// SudokuCards constructs sudoku puzzles with one solution, given nine 3x3
// cards to order.
func SudokuCards(cards [9][3][3]int, stats *Stats,
	visit func(solution []int) bool) {

	// Compare card1 and card2 for less than (-1), equal (0),
	// or greater than (1)
	cmp := func(card1 [3][3]int, card2 [3][3]int) int {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if card1[i][j] < card2[i][j] {
					return -1 // less than
				}
				if card1[i][j] > card2[i][j] {
					return 1 // greater than
				}
			}
		}
		return 0 // equal
	}

	// Iterate over permutations of the card ordering. Each card ordering has
	// 3!3! symmetric orderings which produce identical results, so use
	// ordering constraints to produce only the first ordering
	perm := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	// ordering constraint: fix card 0 to the position 0
	Permutations(perm[1:], func() bool {

		// ordering constraint: ensure card at position 4 is less than cards at
		// positions 5, 7, and 8
		if !(cmp(cards[perm[4]], cards[perm[5]]) < 0 &&
			cmp(cards[perm[4]], cards[perm[7]]) < 0 &&
			cmp(cards[perm[4]], cards[perm[8]]) < 0) {
			return true
		}

		// Build the Sudoku grid from the provided card order
		var grid [9][9]int

		for x, card := range perm {
			i, j := (x/3)*3, (x%3)*3
			for iDelta := 0; iDelta < 3; iDelta++ {
				for jDelta := 0; jDelta < 3; jDelta++ {
					grid[i+iDelta][j+jDelta] = cards[card][iDelta][jDelta]
				}
			}
		}

		// Count the number of Sudoku solutions for this card ordering
		count := 0

		Sudoku(grid, stats,
			func(solution [9][9]int) bool {
				count++
				return true
			})

		if count == 1 {
			visit(perm)
		}

		return true
	})
}

// Mathematicians lists 27 people (without special characters) who were authors
// of early papers in Acta Mathematica and subsequently cited in TAOCP
var Mathematicians = []string{
	"ABEL",
	"BERTRAND",
	"BOREL",
	"CANTOR",
	"CATALAN",
	"FROBENIUS",
	"GLAISHER",
	"GRAM",
	"HADAMARD",
	"HENSEL",
	"HERMITE",
	"HILBERT",
	"HURWITZ",
	"JENSEN",
	"KIRCHHOFF",
	"KNOPP",
	"LANDAU",
	"MARKOFF",
	"MELLIN",
	"MINKOWSKI",
	"NETTO",
	"PERRON",
	"RUNGE",
	"STERN",
	"STIELTJES",
	"SYLVESTER",
	"WEIERSTRASS",
}

// WordSearch uses ExactCoverColoring to build a m x n word search, given the
// provided words
func WordSearch(m int, n int, words []string, stats *Stats,
	visit func([][]string) bool) {

	coord := func(i int, j int) string {
		return fmt.Sprintf("%02d%02d", i, j)
	}

	// secondary items
	secondary := make([]string, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			secondary[i*n+j] = coord(i, j)
		}
	}

	// options
	options := make([][]string, 0)
	for x, word := range words {

		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				// Eight directions for each starting position (i,j)
				var wordDs [8][]string // word directions
				for d := 0; d < 8; d++ {
					wordDs[d] = []string{word}
				}

				for k := 0; k < len(word); k++ {
					// right
					if j+k < n {
						wordDs[0] = append(wordDs[0], coord(i, j+k)+":"+word[k:k+1])
					}
					// right-down
					if i+k < m && j+k < n {
						wordDs[1] = append(wordDs[1], coord(i+k, j+k)+":"+word[k:k+1])
					}

					// To avoid symmetric positions, only allow 2 (of 8) directions
					// for the first word
					if x == 0 {
						continue
					}

					// down
					if i+k < m {
						wordDs[2] = append(wordDs[2], coord(i+k, j)+":"+word[k:k+1])
					}
					// left-down
					if i+k < m && j-k >= 0 {
						wordDs[3] = append(wordDs[3], coord(i+k, j-k)+":"+word[k:k+1])
					}
					// left
					if j-k >= 0 {
						wordDs[4] = append(wordDs[4], coord(i, j-k)+":"+word[k:k+1])
					}
					// left-up
					if i-k >= 0 && j-k >= 0 {
						wordDs[5] = append(wordDs[5], coord(i-k, j-k)+":"+word[k:k+1])
					}
					// up
					if i-k >= 0 {
						wordDs[6] = append(wordDs[6], coord(i-k, j)+":"+word[k:k+1])
					}
					// right-up
					if i-k >= 0 && j+k < n {
						wordDs[7] = append(wordDs[7], coord(i-k, j+k)+":"+word[k:k+1])
					}

				}

				for _, wordD := range wordDs {
					if len(wordD) == len(word)+1 {
						options = append(options, wordD)
						// fmt.Println(wordD)
					}
				}
			}
		}
	}

	ExactCoverColors(words, options, secondary, stats,
		func(solution [][]string) bool {
			return visit(solution)
		})
}
