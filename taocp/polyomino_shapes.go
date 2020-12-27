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
type Point struct{ x, y int }

// Polyomino2 represents a single polyomino of multiple points
type Polyomino2 []Point

type pointset map[Point]bool

func (p Point) rotate90() Point  { return Point{p.y, -p.x} }
func (p Point) rotate180() Point { return Point{-p.x, -p.y} }
func (p Point) rotate270() Point { return Point{-p.y, p.x} }
func (p Point) reflect() Point   { return Point{-p.x, p.y} }

func (p Point) String() string {
	return fmt.Sprintf("%c%c", valueMap[p.x], valueMap[p.y])
}

// All four points in Von Neumann neighborhood
func (p Point) contiguous() Polyomino2 {
	return Polyomino2{Point{p.x - 1, p.y}, Point{p.x + 1, p.y},
		Point{p.x, p.y - 1}, Point{p.x, p.y + 1}}
}

// Finds the min x and y coordinates of a Polyomino.
func (po Polyomino2) minima() (int, int) {
	minx := po[0].x
	miny := po[0].y
	for i := 1; i < len(po); i++ {
		if po[i].x < minx {
			minx = po[i].x
		}
		if po[i].y < miny {
			miny = po[i].y
		}
	}
	return minx, miny
}

func (po Polyomino2) translateToOrigin() Polyomino2 {
	minx, miny := po.minima()
	res := make(Polyomino2, len(po))
	for i, p := range po {
		res[i] = Point{p.x - minx, p.y - miny}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].x < res[j].x || (res[i].x == res[j].x && res[i].y < res[j].y)
	})
	return res
}

// All the plane symmetries of a rectangular region.
func (po Polyomino2) rotationsAndReflections() []Polyomino2 {
	rr := make([]Polyomino2, 8)
	for i := 0; i < 8; i++ {
		rr[i] = make(Polyomino2, len(po))
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

func (po Polyomino2) canonical() Polyomino2 {
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

func (po Polyomino2) String() string {
	return fmt.Sprintf("%v", []Point(po))
}

func (po Polyomino2) toPointset() pointset {
	pset := make(pointset, len(po))
	for _, p := range po {
		pset[p] = true
	}
	return pset
}

// Finds all distinct points that can be added to a Polyomino.
func (po Polyomino2) newPoints() Polyomino2 {
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
	poly := make(Polyomino2, 0, len(m))
	for k := range m {
		poly = append(poly, k)
	}
	return poly
}

func (po Polyomino2) newPolys() []Polyomino2 {
	pts := po.newPoints()
	res := make([]Polyomino2, len(pts))
	for i, pt := range pts {
		poly := make(Polyomino2, len(po))
		copy(poly, po)
		poly = append(poly, pt)
		res[i] = poly.canonical()
	}
	return res
}

var monomino = Polyomino2{Point{0, 0}}
var monominoes = []Polyomino2{monomino}

// Generates polyominoes of rank n recursively.
func rank(n int) []Polyomino2 {
	switch {
	case n < 0:
		panic("n cannot be negative. Program terminated.")
	case n == 0:
		return []Polyomino2{}
	case n == 1:
		return monominoes
	default:
		r := rank(n - 1)
		m := make(map[string]bool)
		var polys []Polyomino2
		for _, po := range r {
			for _, po2 := range po.newPolys() {
				if s := po2.String(); !m[s] {
					polys = append(polys, po2)
					m[s] = true
				}
			}
		}
		sort.Slice(polys, func(i, j int) bool {
			return polys[i].String() < polys[j].String()
		})
		return polys
	}
}

// PolyominoShapes provides YAML (de-)serialization for Exact Cover input
type PolyominoShapes struct {
	PieceSets map[string]map[string]string `yaml:""` // Piece Sets
	Boards    map[string]string            `yaml:""` // Boards
}

// NewPolyominoShapes creates a new instance of PolyominoShapes
func NewPolyominoShapes() *PolyominoShapes {
	shapes := &PolyominoShapes{
		PieceSets: map[string]map[string]string{},
		Boards:    map[string]string{},
	}

	return shapes
}

// GeneratePolyominoShapes generates shapes of size n
func GeneratePolyominoShapes(n int) []Polyomino2 {
	return rank(n)
}
