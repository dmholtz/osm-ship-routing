package geometry

import "math"

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
