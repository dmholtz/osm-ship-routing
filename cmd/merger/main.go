package main

import (
	"fmt"
	"math"
)

type Point2D struct {
	x float64
	y float64
}

type Polygon2D struct {
	points []*Point2D
}

func (p *Polygon2D) IsInPolygon(point *Point2D) bool {
	// TODO: check if closed
	start := len(p.points) - 1
	end := 0

	contains := p.IntersectsWithRaycast(point, p.points[start], p.points[end])

	for i := 1; i < len(p.points); i++ {
		if p.IntersectsWithRaycast(point, p.points[i-1], p.points[i]) {
			contains = !contains
		}
	}

	return contains
}

func (p *Polygon2D) IntersectsWithRaycast(point *Point2D, start *Point2D, end *Point2D) bool {
	// ensure that first point as a lower y coordinate than second point
	// is this necessary???
	if start.y > end.y {
		// switch points
		start, end = end, start
	}

	for point.y == start.y || point.y == end.y {
		newY := math.Nextafter(point.y, math.Inf(1))
		point = &Point2D{point.x, newY}
	}

	if point.y < start.y || point.y > end.y {
		return false
	}

	if start.x > end.x {
		if point.x > start.x {
			return false
		}
		if point.x < end.x {
			return true
		}
	} else {
		if point.x > end.x {
			return false
		}
		if point.x < start.x {
			return true
		}
	}

	raySlope := (point.y - start.y) / (point.x - start.x)
	diagSlope := (end.y - start.y) / (end.x - start.x)

	return raySlope >= diagSlope
}

//const pbfFile string = "antarctica.osm.pbf"

const pbfFile string = "planet-coastlines.pbf"

func main() {

	/*
		start := time.Now()

		coastlineImporter := pbf.NewImporter(pbfFile)
		coastlineImporter.Import()

		elapsed := time.Since(start)
		fmt.Printf("[TIME] Import: %s\n", elapsed)

		start = time.Now()

		merger := coastline.NewMerger(coastlineImporter.Coastlines())
		merger.Merge()

		elapsed = time.Since(start)
		fmt.Printf("[TIME] Merge: %s\n", elapsed)
		fmt.Printf("Polygon Count: %d\n", len(merger.Polygons()))
		fmt.Printf("Merge Count: %d\n", merger.MergeCount())
	*/
	p1 := Point2D{80, 80}
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
			panic(fmt.Sprintf("Point is declared wrongly: %v", testPoint))
		}
	}
	for _, testPoint := range testPointsOutside {
		isInPolygon := polygon.IsInPolygon(&testPoint)
		if isInPolygon {
			panic(fmt.Sprintf("Point is declared wrongly: %v", testPoint))
		}
	}
}
