package taocp

import (
	"fmt"
	"sort"
)

// A Polyomino is a plane geometric figure formed by joining one or more equal
// squares edge to edge. Free polyominoes are distinct when none is a
// translation, rotation, reflection or glide reflection of another polyomino.
// https://rosettacode.org/wiki/Free_polyominoes_enumeration#Go

// Point represents a single square in a polyomino
type Point struct{ X, Y int }

// Polyomino represents a single polyomino of multiple points
type Polyomino []Point

type pointset map[Point]bool

// PolyominoShape holds a single shape
type PolyominoShape struct {
	Shape      string      `yaml:",omitempty"`
	Points     Polyomino   `yaml:",omitempty"`
	Placements []Polyomino `yaml:",omitempty"`
}

// PolyominoShapes holds PolyominoShape piece sets and boards
type PolyominoShapes struct {
	PieceSets map[string]map[string]*PolyominoShape `yaml:""` // Piece Sets
	Boards    map[string]*PolyominoShape            `yaml:""` // Boards
}

func (p Point) rotate90() Point  { return Point{p.Y, -p.X} }
func (p Point) rotate180() Point { return Point{-p.X, -p.Y} }
func (p Point) rotate270() Point { return Point{-p.Y, p.X} }
func (p Point) reflect() Point   { return Point{-p.X, p.Y} }

func (p Point) String() string {
	return fmt.Sprintf("%c%c", valueMap[p.X], valueMap[p.Y])
}

// Bounds returns the bounding box of (x, y) coordinates of a Polyomino2 piece
func (po Polyomino) Bounds() (int, int, int, int) {
	xMin, yMin, xMax, yMax := -1, -1, -1, -1
	for _, point := range po {
		if xMin == -1 || point.X < xMin {
			xMin = point.X
		}
		if yMin == -1 || point.Y < yMin {
			yMin = point.Y
		}
		if xMax == -1 || point.X > xMax {
			xMax = point.X
		}
		if yMax == -1 || point.Y > yMax {
			yMax = point.Y
		}
	}
	return xMin, yMin, xMax, yMax
}

// IsConvex tests whether a shape is convex, ie if it contains all of the
// squares between any two of its squares that lie in the same row of the same
// column.
func (po Polyomino) IsConvex() bool {
	pset := po.toPointset()
	xMin, yMin, xMax, yMax := po.Bounds()

	// Check each row
	for x := xMin; x <= xMax; x++ {
		changes := 0     // Record changes from in the shape to out
		inShape := false // Track whether we are currently in the shape or not
		for y := yMin; y <= yMax; y++ {
			p := Point{x, y}
			if pset[p] {
				if !inShape {
					inShape = true
					changes++
				}
			} else {
				if inShape {
					inShape = false
					changes++
				}
			}
			if changes == 3 {
				// We changed from out to in, in to out, and back to in
				// So not convex
				return false
			}
		}
	}

	// Check each column
	for y := yMin; y <= yMax; y++ {
		changes := 0     // Record changes from in the shape to out
		inShape := false // Track whether we are currently in the shape or not
		for x := xMin; x <= xMax; x++ {
			p := Point{x, y}
			if pset[p] {
				if !inShape {
					inShape = true
					changes++
				}
			} else {
				if inShape {
					inShape = false
					changes++
				}
			}
			if changes == 3 {
				// We changed from out to in, in to out, and back to in
				// So not convex
				return false
			}
		}
	}

	// Shape is convex
	return true
}

// All four points in Von Neumann neighborhood
func (p Point) contiguous() Polyomino {
	return Polyomino{Point{p.X - 1, p.Y}, Point{p.X + 1, p.Y},
		Point{p.X, p.Y - 1}, Point{p.X, p.Y + 1}}
}

// Finds the min x and y coordinates of a Polyomino.
func (po Polyomino) minima() (int, int) {
	minx := po[0].X
	miny := po[0].Y
	for i := 1; i < len(po); i++ {
		if po[i].X < minx {
			minx = po[i].X
		}
		if po[i].Y < miny {
			miny = po[i].Y
		}
	}
	return minx, miny
}

func (po Polyomino) translateToOrigin() Polyomino {
	minx, miny := po.minima()
	res := make(Polyomino, len(po))
	for i, p := range po {
		res[i] = Point{p.X - minx, p.Y - miny}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].X < res[j].X || (res[i].X == res[j].X && res[i].Y < res[j].Y)
	})
	return res
}

// All the plane symmetries of a rectangular region.
func (po Polyomino) rotationsAndReflections() []Polyomino {
	rr := make([]Polyomino, 8)
	for i := 0; i < 8; i++ {
		rr[i] = make(Polyomino, len(po))
	}
	copy(rr[0], po)
	for j := 0; j < len(po); j++ {
		rr[1][j] = po[j].rotate90()
		rr[2][j] = po[j].rotate180()
		rr[3][j] = po[j].rotate270()
		rr[4][j] = po[j].reflect()
		rr[5][j] = po[j].rotate90().reflect()
		rr[6][j] = po[j].rotate180().reflect()
		rr[7][j] = po[j].rotate270().reflect()
	}
	return rr
}

func (po Polyomino) canonical() Polyomino {
	rr := po.rotationsAndReflections()
	minr := rr[0].translateToOrigin()
	mins := minr.String()
	for i := 1; i < 8; i++ {
		r := rr[i].translateToOrigin()
		s := r.String()
		if s < mins {
			minr = r
			mins = s
		}
	}
	return minr
}

func (po Polyomino) String() string {
	return fmt.Sprintf("%v", []Point(po))
}

func (po Polyomino) toPointset() pointset {
	pset := make(pointset, len(po))
	for _, p := range po {
		pset[p] = true
	}
	return pset
}

// Finds all distinct points that can be added to a Polyomino.
func (po Polyomino) newPoints() Polyomino {
	pset := po.toPointset()
	m := make(pointset)
	for _, p := range po {
		pts := p.contiguous()
		for _, pt := range pts {
			if !pset[pt] {
				m[pt] = true // using an intermediate set is about 15% faster!
			}
		}
	}
	poly := make(Polyomino, 0, len(m))
	for k := range m {
		poly = append(poly, k)
	}
	return poly
}

func (po Polyomino) newPolys() []Polyomino {
	pts := po.newPoints()
	res := make([]Polyomino, len(pts))
	for i, pt := range pts {
		poly := make(Polyomino, len(po))
		copy(poly, po)
		poly = append(poly, pt)
		res[i] = poly.canonical()
	}
	return res
}

// minmax finds minimum and maximum (x, y) values
func minmax(po Polyomino) (int, int, int, int) {
	xMin, yMin, xMax, yMax := -1, -1, -1, -1
	for _, point := range po {
		x, y := point.X, point.Y
		if xMin == -1 || x < xMin {
			xMin = x
		}
		if yMin == -1 || y < yMin {
			yMin = y
		}
		if xMax == -1 || x > xMax {
			xMax = x
		}
		if yMax == -1 || y > yMax {
			yMax = y
		}
	}
	return xMin, yMin, xMax, yMax
}

var monomino = Polyomino{Point{0, 0}}
var monominoes = []Polyomino{monomino}

var valueMap = []byte{'0', '1', '2', '3', '4', '5',
	'6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
	'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

// sortPoints sorts the points of a polyomino
func sortPoints(po Polyomino) {
	sort.Slice(po, func(i, j int) bool {
		return po[i].X < po[j].X || (po[i].X == po[j].X && po[i].Y < po[j].Y)
	})
}

// sortPolyominoes sorts a slice of polyominoes; assumes each polyomino
// already has sorted points
func sortPolyominoes(polys []Polyomino) {
	sort.Slice(polys, func(i, j int) bool {
		return polys[i].String() < polys[j].String()
	})
}

// Generates polyominoes of rank n recursively.
func rank(n int) []Polyomino {
	switch {
	case n < 0:
		panic("n cannot be negative. Program terminated.")
	case n == 0:
		return []Polyomino{}
	case n == 1:
		return monominoes
	default:
		r := rank(n - 1)
		m := make(map[string]bool)
		var polys []Polyomino
		for _, po := range r {
			for _, po2 := range po.newPolys() {
				if s := po2.String(); !m[s] {
					polys = append(polys, po2)
					m[s] = true
				}
			}
		}
		sortPolyominoes(polys)
		return polys
	}
}

// NewPolyominoShapes creates a new instance of PolyominoShapes
func NewPolyominoShapes() *PolyominoShapes {
	shapes := &PolyominoShapes{
		PieceSets: make(map[string]map[string]*PolyominoShape),
		Boards:    make(map[string]*PolyominoShape),
	}

	return shapes
}

// GeneratePolyominoShapes generates shapes of size n
func GeneratePolyominoShapes(n int) []Polyomino {
	return rank(n)
}
