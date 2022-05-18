package geometry

import (
	"testing"
)

var p1 Point = Point{0, 0}
var p2 Point = Point{0, 179}
var p3 Point = Point{0, -179}
var p4 Point = Point{0, 2}

func TestHaversine(t *testing.T) {
	d1 := p1.IntHaversine(&p2)
	d2 := p2.IntHaversine(&p1)
	if d1 != d2 {
		t.Errorf("p1.Haversine(p2) ≠ p2.Haversine(p1): %d≠%d", d1, d2)
	}

	d1 = p2.IntHaversine(&p3)
	d2 = p3.IntHaversine(&p2)
	if d1 != d2 {
		t.Errorf("p2.Haversine(p3) ≠ p3.Haversine(p2): %d≠%d", d1, d2)
	}

	d1 = p1.IntHaversine(&p4) // delta-Lon = 2 degrees
	d2 = p3.IntHaversine(&p2) // delta-Lon = 2 degrees
	if d1 != d2 {
		t.Errorf("p1.Haversine(p4) ≠ p3.Haversine(p2): %d≠%d", d1, d2)
	}
}
