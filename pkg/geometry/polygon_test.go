package geometry

import (
	"fmt"
	"testing"
)

func TestPolygonTest(t *testing.T) {
	t.Parallel()

	p1 := Point2D{80, 80}
	fmt.Println(p1)
	p2 := Point2D{80, 40}
	p3 := Point2D{100, 60}
	p4 := Point2D{100, 10}
	p5 := Point2D{10, 10}
	p6 := Point2D{10, 80}
	p7 := Point2D{80, 80}
	testPointsInside := []Point2D{{70, 40}}
	testPointsOutside := []Point2D{{0, 50}}

	points := []*Point2D{&p1, &p2, &p3, &p4, &p5, &p6, &p7}
	polygon := Polygon2D{points}

	for _, testPoint := range testPointsInside {
		isInPolygon := polygon.IsInPolygon(&testPoint)
		if !isInPolygon {
			t.Errorf("want 'isInPolygon==true', got 'isInPolygon==%t' for point %v", isInPolygon, testPoint)
		}
	}
	for _, testPoint := range testPointsOutside {
		isInPolygon := polygon.IsInPolygon(&testPoint)
		if isInPolygon {
			t.Errorf("want 'isInPolygon==false', got 'isInPolygon==%t' for point %v", isInPolygon, testPoint)
		}
	}
}
