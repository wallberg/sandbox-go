package math

import "fmt"

// This package supplements the built-in math package

// MinInt returns the minimum value of the provided ints
func MinInt(is ...int) int {
	min := is[0]
	for _, i := range is[1:] {
		if i < min {
			min = i
		}
	}
	return min
}

// MaxInt returns the maximum value of the provided ints
func MaxInt(is ...int) int {
	max := is[0]
	for _, i := range is[1:] {
		if i > max {
			max = i
		}
	}
	return max
}

// Monus returns max(x-y, 0)
func MonusInt(x int, y int) int {
	return MaxInt(x-y, 0)
}

// CountDigits counts the number of digits in an integer
func CountDigits(n int64) int {

	if n < 0 {
		panic(fmt.Sprintf("Expected n >= 0; got %v", n))
	}

	count := 0
	for n != 0 {
		count++
		n /= 10
	}
	return count
}
