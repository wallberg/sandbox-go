package taocp

import (
	"fmt"
	"log"
)

// Explore Dancing Links from The Art of Computer Programming, Volume 4,
// Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
// Dancing Links, 2020
//
// §7.2.2.1 Dancing Links - Cubes (Exercise 145-148)

type Cube string

// Cubes generates cubes with six colors, one per each face
//
// Exercise 7.2.2.1-146
func Cubes() (cubes []Cube) {
	colors := []string{"a", "b", "c", "d", "e", "f"}

	x := []int{0, 1, 2, 3, 4, 5}
	Permutations(x, func() bool {
		if x[0] < x[1] && // u < u'
			x[0] < x[2] && // u < v
			x[2] < x[3] && // v < v'
			x[2] < x[4] && // v < w
			x[2] < x[5] { // v < w'

			cube := (Cube)(colors[x[0]] + colors[x[1]] + colors[x[2]] + colors[x[3]] + colors[x[4]] + colors[x[5]])
			cubes = append(cubes, cube)
		}
		return true
	})

	return cubes
}

// Rotations returns all rotations of the provided cube, including itself
func (c Cube) Rotations() (rotations []Cube) {

	cube := c

	rotations = append(rotations, cube)

	for _, action := range []string{
		"normalRight", "normalRight", "normalRight", "longRight",
		"normalRight", "normalRight", "normalRight", "longRight",
		"normalRight", "normalRight", "normalRight", "longLeft",
		"normalRight", "normalRight", "normalRight", "longLeft",
		"normalRight", "normalRight", "normalRight", "longRight",
		"normalRight", "normalRight", "normalRight",
	} {
		switch action {
		case "normalRight":
			cube = (Cube)(cube[0:1] + cube[1:2] + cube[5:6] + cube[4:5] + cube[2:3] + cube[3:4])
		case "longRight":
			cube = (Cube)(cube[4:5] + cube[5:6] + cube[2:3] + cube[3:4] + cube[1:2] + cube[0:1])
		case "longLeft":
			cube = (Cube)(cube[5:6] + cube[4:5] + cube[2:3] + cube[3:4] + cube[0:1] + cube[1:2])
		}
		rotations = append(rotations, cube)
	}

	return rotations
}

// Brick generates the items, options, and secondary items to assemble
// cubes into an l x m x n size brick, with each brick face having
// a single color, using XCC.
func Brick(l, m, n int) ([]string, [][]string, []string) {

	var (
		items   []string
		options [][]string
		sitems  []string
	)

	// Validation
	if l < 1 {
		log.Fatalf("invalid value l=%d", l)
	}
	if m < 1 {
		log.Fatalf("invalid value m=%d", m)
	}
	if n < 1 {
		log.Fatalf("invalid value n=%d", n)
	}

	// Primary Items - cube positions
	for i := 0; i < l; i++ {
		for j := 0; j < m; j++ {
			for k := 0; k < n; k++ {
				items = append(items, fmt.Sprintf("%d-%d-%d", 2*i+1, 2*j+1, 2*k+1))
			}
		}
	}

	// Primary Items - brick faces
	for _, brickFace := range []string{"top", "bottom", "left", "right", "front", "back"} {
		items = append(items, brickFace)

		// Options, one for each color on the cube faces matching the brick face
		for _, color := range []string{"a", "b", "c", "d", "e", "f"} {
			option := []string{brickFace}

			// Iterate over cube faces which compose the brick faces
			for x := 0; x <= 2*l; x++ {
				for y := 0; y <= 2*m; y++ {
					for z := 0; z <= 2*n; z++ {
						if x%2+y%2+z%2 == 2 {
							if (brickFace == "top" && y == 0) ||
								(brickFace == "bottom" && y == 2*m) ||
								(brickFace == "left" && z == 0) ||
								(brickFace == "right" && z == 2*n) ||
								(brickFace == "front" && x == 0) ||
								(brickFace == "back" && x == 2*l) {

								cubeFace := fmt.Sprintf("%d-%d-%d", x, y, z)
								option = append(option, fmt.Sprintf("%s:%s", cubeFace, color))

							}
						}
					}
				}
			}

			options = append(options, option)
		}
	}

	// Secondary Items - cubes
	for cubeI, cube := range Cubes() {
		sitems = append(sitems, (string)(cube))

		// Cube placement options, for each cube position
		for i := 0; i < l; i++ {
			for j := 0; j < m; j++ {
				for k := 0; k < n; k++ {
					cubePosition := fmt.Sprintf("%d-%d-%d", 2*i+1, 2*j+1, 2*k+1)

					// For each rotation of the cube
					for rotationI, rotation := range cube.Rotations() {

						// in position 1-1-1: only place the first rotation of the first cube, to reduce symmetries
						if i == 0 && j == 0 && k == 0 && (cubeI > 0 || (cubeI == 0 && rotationI > 0)) {
							continue
						}

						option := []string{cubePosition, (string)(cube)}

						// For each face of the cube
						option = append(option, fmt.Sprintf("%d-%d-%d:%s", 2*i+0, 2*j+1, 2*k+1, rotation[0:1])) // top
						option = append(option, fmt.Sprintf("%d-%d-%d:%s", 2*i+1, 2*j+0, 2*k+1, rotation[2:3])) // front
						option = append(option, fmt.Sprintf("%d-%d-%d:%s", 2*i+1, 2*j+1, 2*k+0, rotation[5:6])) // right
						option = append(option, fmt.Sprintf("%d-%d-%d:%s", 2*i+1, 2*j+1, 2*k+2, rotation[4:5])) // left
						option = append(option, fmt.Sprintf("%d-%d-%d:%s", 2*i+1, 2*j+2, 2*k+1, rotation[3:4])) // back
						option = append(option, fmt.Sprintf("%d-%d-%d:%s", 2*i+2, 2*j+1, 2*k+1, rotation[1:2])) // bottom

						options = append(options, option)
					}
				}
			}
		}
	}

	// Secondary Items - cube faces
	for x := 0; x <= 2*l; x++ {
		for y := 0; y <= 2*m; y++ {
			for z := 0; z <= 2*n; z++ {
				if x%2+y%2+z%2 == 2 {
					cubeFace := fmt.Sprintf("%d-%d-%d", x, y, z)
					sitems = append(sitems, cubeFace)
				}
			}
		}
	}

	return items, options, sitems

}
