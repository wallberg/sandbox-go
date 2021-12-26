package golang

import (
	"fmt"
	"testing"
)

// Benchmark performance for storing a 2D matrix as [][]float32 vs. []float32

func BenchmarkMatrix2DNew(b *testing.B) {
	for _, size := range []int{2, 3, 4} {
		msg := fmt.Sprintf("%dx%d", size, size)

		b.Run(msg, func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				m := make([][]float32, size)
				for i := 0; i < size; i++ {
					m[i] = make([]float32, size)
				}
			}
		})
	}
}

func BenchmarkMatrix1DNew(b *testing.B) {
	for _, size := range []int{2, 3, 4} {
		msg := fmt.Sprintf("%dx%d", size, size)

		b.Run(msg, func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				m := make([]float32, size*size)
				_ = m
			}
		})
	}
}

func BenchmarkMatrix2DAccess(b *testing.B) {
	for _, size := range []int{2, 3, 4} {
		msg := fmt.Sprintf("%dx%d", size, size)
		m := make([][]float32, size)
		for i := 0; i < size; i++ {
			m[i] = make([]float32, size)
		}
		b.Run(msg, func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				f := (float32)(b.N)
				for i := 0; i < size; i++ {
					for j := 0; j < size; j++ {
						_ = m[i][j]
						m[i][j] = f
					}
				}
			}
		})
	}
}

func BenchmarkMatrix1DAccess(b *testing.B) {
	for _, size := range []int{2, 3, 4} {
		msg := fmt.Sprintf("%dx%d", size, size)
		m := make([]float32, size*size)

		b.Run(msg, func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				f := (float32)(b.N)
				for i := 0; i < size; i++ {
					for j := 0; j < size; j++ {
						k := i*size + j
						_ = m[k]
						m[k] = f
					}
				}
			}
		})
	}
}

func get(m *[]float32, i, j, size int) float32 {
	k := i*size + j
	return (*m)[k]
}

func set(m *[]float32, i, j, size int, value float32) {
	k := i*size + j
	(*m)[k] = value
}

func BenchmarkMatrix1DFunc(b *testing.B) {
	for _, size := range []int{2, 3, 4} {
		msg := fmt.Sprintf("%dx%d", size, size)
		m := make([]float32, size*size)

		b.Run(msg, func(b *testing.B) {
			for repeat := 0; repeat < b.N; repeat++ {
				f := (float32)(b.N)
				for i := 0; i < size; i++ {
					for j := 0; j < size; j++ {
						_ = get(&m, i, j, size)
						set(&m, i, j, size, f)
					}
				}
			}
		})
	}
}
