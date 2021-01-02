package taocp

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestPolyominoShapes(t *testing.T) {

	// Build YAML struct
	shapes := PolyominoShapes{
		PieceSets: map[string]map[string]string{
			"1": {"A": "00"},
			"3": {"C": "0[012]", "D": "00 01 11"},
		},
		Boards: map[string]string{
			"3x20": "[0-2][0-j]",
			"1x1":  "00",
		},
	}

	// Serialize to YAML
	data, err := yaml.Marshal(shapes)
	if err != nil {
		t.Errorf("Error serializing PolyominoShapes: %v", err)
		return
	}

	// Deserialize from YAML
	var shapes2 PolyominoShapes
	err = yaml.Unmarshal([]byte(data), &shapes2)
	if err != nil {
		t.Errorf("Error deserializing PolyominoShapes: %v", err)
	}

	// Test the round trip
	if !reflect.DeepEqual(shapes, shapes2) {
		t.Errorf("Got back %v; want %v", shapes2, shapes)
	}
}

func TestNewPolyominoShapes(t *testing.T) {

	// Build YAML struct
	shapes := &PolyominoShapes{
		PieceSets: map[string]map[string]string{},
		Boards:    map[string]string{},
	}

	shapes2 := NewPolyominoShapes()

	// Compare
	if !reflect.DeepEqual(shapes, shapes2) {
		t.Errorf("Got %v; want %v", shapes2, shapes)
	}
}

func TestGeneratePolyominoShapes(t *testing.T) {

	cases := []struct {
		n      int          // size
		count  int          // number of shapes generated
		shapes []Polyomino2 // generated shapes
	}{
		{
			1,
			1,
			[]Polyomino2{{{0, 0}}},
		},
		{
			2,
			1,
			[]Polyomino2{{{0, 0}, {0, 1}}},
		},
		{
			3,
			2,
			[]Polyomino2{
				{{0, 0}, {0, 1}, {0, 2}},
				{{0, 0}, {0, 1}, {1, 0}},
			},
		},
		{
			4,
			5,
			nil,
		},
		{
			5,
			12,
			nil,
		},
		{
			6,
			35,
			nil,
		},
		{
			7,
			108,
			nil,
		},
		{
			8,
			369,
			nil,
		},
		{
			9,
			1285,
			nil,
		},
		{
			10,
			4655,
			nil,
		},
		// { // too slow
		// 	11,
		// 	17073,
		// 	nil,
		// },
	}

	for _, c := range cases {
		shapes := GeneratePolyominoShapes(c.n)

		if count := len(shapes); count != c.count {
			t.Errorf("for n=%d, got number of shapes %d; want %d", c.n, count, c.count)
		}

		if c.shapes != nil && !reflect.DeepEqual(shapes, c.shapes) {
			t.Errorf("for n=%d, got shapes %v; want %v", c.n, shapes, c.shapes)
		}
	}
}

func TestBounds(t *testing.T) {

	cases := []struct {
		po   Polyomino2
		xMin int
		yMin int
		xMax int
		yMax int
	}{
		{
			Polyomino2{{0, 0}},
			0,
			0,
			0,
			0,
		},
		{
			Polyomino2{},
			-1,
			-1,
			-1,
			-1,
		},
		{
			Polyomino2{{0, 1}, {1, 1}, {1, 2}},
			0,
			1,
			1,
			2,
		},
	}

	for _, c := range cases {
		xMin, yMin, xMax, yMax := c.po.Bounds()

		if xMin != c.xMin {
			t.Errorf("for po=%v, got xMin=%d; want %d", c.po, xMin, c.xMin)
		}
		if yMin != c.yMin {
			t.Errorf("for po=%v, got yMin=%d; want %d", c.po, yMin, c.yMin)
		}
		if xMax != c.xMax {
			t.Errorf("for po=%v, got xMax=%d; want %d", c.po, xMax, c.xMax)
		}
		if yMax != c.yMax {
			t.Errorf("for po=%v, got yMax=%d; want %d", c.po, yMax, c.yMax)
		}
	}
}

func TestIsConvex(t *testing.T) {

	cases := []struct {
		po       Polyomino2
		isConvex bool
	}{
		{
			Polyomino2{},
			true,
		},
		{
			Polyomino2{{0, 1}, {1, 1}, {1, 2}, {2, 2}},
			true,
		},
		{
			Polyomino2{{0, 1}, {1, 1}, {1, 2}, {2, 2}, {3, 2}, {3, 1}},
			false,
		},
		{
			Polyomino2{{0, 0}, {1, 0}, {1, 1}, {1, 2}, {0, 2}},
			false,
		},
	}

	for _, c := range cases {
		isConvex := c.po.IsConvex()

		if isConvex != c.isConvex {
			t.Errorf("for po=%v, got isConvex=%t; want %t", c.po, isConvex, c.isConvex)
		}
	}
}
