package main

import (
	"fmt"
	"math"
)

func main() {

	var (
		m    uint8 = 33
		j, k uint8
	)

	// taocp.OrderingTrial(3, taocp.Plain)

	// n := int(math.Ceil(math.Log2(float64(m)))) // number of bits
	// fmt.Printf("m=%d, n=%d\n", m, n)

	// b := func(i int) string {
	// 	s := ""
	// 	for j := 0; j < n; j++ {
	// 		if i&1 == 0 {
	// 			s = "0" + s
	// 		} else {
	// 			s = "1" + s
	// 		}
	// 		i = (i >> 1)
	// 	}
	// 	return s
	// }

	xx := func(j, k uint8) uint8 {
		// fmt.Printf("  j=%d, k=%d", j, k)
		// l := math.Log2(float64(j - k))
		// fmt.Printf(", l=%1.2f", l)
		// f := math.Floor(l)
		// fmt.Printf(", f=%1.2f", f)
		// p := math.Pow(2.0, f)
		// fmt.Printf(", p=%1.2f", p)
		// ip := int(p)
		// fmt.Printf(", p=%s", b(ip))
		// fmt.Printf(", j=%s, -p=%s", b(j), b(-ip))
		// r := j & -ip
		// fmt.Printf(", r=%d", r)
		// return r
		return j & (-(uint8(math.Pow(2.0, math.Floor(math.Log2(float64(j-k)))))))
		// return uint8(math.Pow(2.0, math.Floor(math.Log2(float64(j-k)))))
	}

	grid := make([][]uint8, m)
	for j = 8; j < 16; j++ {
		grid[j] = make([]uint8, m)
		for k = 0; k < j; k++ {
			// Option 3
			grid[j][k] = xx(j, k)
			fmt.Printf("%3d ", grid[j][k])
		}
		fmt.Println()
	}

	// // Option 1
	// for k := 1; k < m; k++ {
	// 	t := k & -k

	// 	min := m
	// 	if k+t < min {
	// 		min = k + t
	// 	}
	// 	for j := k; j < min; j++ {
	// 		fmt.Printf("a%d gets y%d\n", j, k)
	// 	}
	// 	for j := k - t; j < k; j++ {
	// 		fmt.Printf("b%d gets y%d\n", j, k)
	// 	}
	// }

	// // Option 2
	// for j = 0; j < m; j++ {
	// 	fmt.Printf("a%d", j)

	// 	t := j
	// 	for t > 0 {
	// 		fmt.Printf(" y%d", t)
	// 		t = t & (t - 1)
	// 	}
	// 	fmt.Println()
	// 	// for k = 0; k < j; k++ {
	// 	// 	xx(j, k)
	// 	// 	fmt.Println()
	// 	// }
	// }
	// for k = 0; k < m; k++ {
	// 	fmt.Printf("b%d", k)

	// 	t := -1 - k
	// 	for t > -m {
	// 		fmt.Printf(" y%d", -t)
	// 		t = t & (t - 1)
	// 	}
	// 	fmt.Println()
	// 	// for j = k + 1; j < m; j++ {
	// 	// 	xx(j, k)
	// 	// 	fmt.Println()
	// 	// }
	// }

	// // Option 3
	// for k := 0; k < m; k++ {
	// 	fmt.Printf("b%d\n", k)
	// 	for j := k + 1; j < m; j++ {
	// 		fmt.Printf(" gets y%d\n", xx(j, k))
	// 	}
	// 	fmt.Println()
	// }

	// e := [9][9]int{
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{2, 2, 0, 0, 0, 0, 0, 0, 0},
	// 	{2, 2, 3, 0, 0, 0, 0, 0, 0},
	// 	{4, 4, 4, 4, 0, 0, 0, 0, 0},
	// 	{4, 4, 4, 4, 5, 0, 0, 0, 0},
	// 	{4, 4, 4, 4, 6, 6, 0, 0, 0},
	// 	{4, 4, 4, 4, 6, 6, 7, 0, 0},
	// 	{8, 8, 8, 8, 8, 8, 8, 8, 0},
	// }

	// for j := 0; j < m; j++ {
	// 	for k := 0; k < j; k++ {
	// 		r := xx(j, k)
	// 		if r != e[j][k] {
	// 			fmt.Printf("j=%d, k=%d, r=%d, e=%d", j, k, r, e[j][k])
	// 			fmt.Println()
	// 		}
	// 	}
	// 	fmt.Println()
	// }
}
