package geometry

import (
	"fmt"
	"math"
)

var hit = 0
var miss = 0

// based on: Some Algorithms for Polygons on a Sphere (Robert.G .Chamberlain)
// with code here: https://github.com/kellydunn/golang-geo/blob/master/polygon.go

/* Some possible interface definitions
type Area interface {
	Contains(point *Point) bool
}

type Polygon interface {
	At(index int) *Point
	Add(point *Point)
	Size() int
	IsClosed() bool
	BoundingBox() BoundingBox
}
*/

type Polygon []*Point

func NewPolygon(points []*Point) *Polygon {
	p := Polygon(points)
	return &p
}

func (p *Polygon) Points() []*Point {
	return *p
}

func (p *Polygon) At(index int) *Point {
	return (*p)[index]
}

func (p *Polygon) Add(point *Point) {
	*p = append(*p, point)
}

func (p *Polygon) Size() int {
	return len(p.Points())
}

func (p *Polygon) IsClosed() bool {
	if p.Size() < 3 || p.At(0).Lat() != p.At(p.Size()-1).Lat() || p.At(0).Lon() != p.At(p.Size()-1).Lon() {
		return false
	}
	return true
}

func (p *Polygon) LatLonBoundingBox() BoundingBox {
	latMin, lonMin := math.Inf(1), math.Inf(1)
	latMax, lonMax := math.Inf(-1), math.Inf(-1)
	for _, point := range p.Points() {
		if point.Lat() < latMin {
			latMin = point.Lat()
		}
		if point.Lat() > latMax {
			latMax = point.Lat()
		}
		if point.Lon() < lonMin {
			lonMin = point.Lon()
		}
		if point.Lon() > lonMax {
			lonMax = point.Lon()
		}
	}
	return BoundingBox{LatMin: latMin, LatMax: latMax, LonMin: lonMin, LonMax: lonMax}
}

func (p *Polygon) GreatCircleBoundingBox() BoundingBox {
	latMin, lonMin := math.Inf(1), math.Inf(1)
	latMax, lonMax := math.Inf(-1), math.Inf(-1)
	tempLatMin, tempLatMax := math.Inf(1), math.Inf(-1)

	fullyNorthern, fullySouthern := true, true

	for _, p1 := range p.Points() {
		if p1.Lon() < lonMin {
			lonMin = p1.Lon()
		}
		if p1.Lon() > lonMax {
			lonMax = p1.Lon()
		}
		if p1.Lat() < tempLatMin {
			tempLatMin = p1.Lat()
		}
		if p1.Lat() > tempLatMax {
			tempLatMax = p1.Lat()
		}
		if p1.Lat() < 0 {
			fullyNorthern = false
		}
		if p1.Lat() > 0 {
			fullySouthern = false
		}
	}
	if fullyNorthern == fullySouthern {
		panic("Polygon seems to be misformed.")
	}
	if fullyNorthern {
		latMin = tempLatMin
	}
	if fullySouthern {
		latMax = tempLatMax
	}

	phiMin, phiMax := math.Inf(1), math.Inf(-1)
	for i := 0; i < p.Size()-1; i++ {
		p1, p2 := p.At(i), p.At(i+1)
		if p1.Lon() < lonMin {
			lonMin = p1.Lon()
		}
		if p1.Lon() > lonMax {
			lonMax = p1.Lon()
		}
		if p1.Lat() < tempLatMin {
			tempLatMin = p1.Lat()
		}
		if p1.Lat() > tempLatMax {
			tempLatMax = p1.Lat()
		}
		if p1.Lat() < 0 {
			fullyNorthern = false
		} else if p1.Lat() > 0 {
			fullySouthern = false
		}
		azimuth := 0.0
		if p1.Lambda() != p2.Lambda() {
			tanAzimuth := (math.Sin(p1.Lambda()-p2.Lambda()) * math.Cos(p2.Phi())) / (math.Cos(p1.Phi())*math.Cos(p2.Phi()) - math.Sin(p1.Phi())*math.Cos(p2.Phi())*math.Cos(p2.Lambda()-p1.Lambda()))
			azimuth = math.Atan(tanAzimuth) // maybe use math.Atan2()
			bearing := Deg2Rad(p1.Bearing(p2))
			if math.Abs(azimuth-bearing) < 0.1 {
				fmt.Printf("Difference in calculating azimuth: %v", azimuth-bearing)
			}
		} else if p1.Lambda() == p2.Lambda() && p1.Phi() < p2.Phi() {
			azimuth = 0
		} else if p1.Lambda() == p2.Lambda() && p1.Phi() > p2.Phi() {
			azimuth = math.Pi
		} else if p1.Lambda() == p2.Lambda() && p1.Phi() == p2.Phi() {
			panic("Identical points: azimuth is undefined.")
		}

		phi := math.Acos(math.Abs(math.Sin(azimuth) * math.Cos(p1.Phi())))
		//phiMin = phiMax
		if phi > phiMax {
			phiMax = phi
		}
		if phi < phiMin {
			phiMin = phi
		}
		if p1.Phi() > phiMax {
			phiMax = p1.Phi()
		}
		if p1.Phi() < phiMin {
			phiMin = p1.Phi()
		}
		if p2.Phi() > phiMax {
			phiMax = p2.Phi()
		}
		if p2.Phi() < phiMin {
			phiMin = p2.Phi()
		}
	}
	if fullyNorthern == fullySouthern {
		panic("Polygon seems to be misformed.")
	}
	if fullyNorthern {
		latMin = tempLatMin
		latMax = Rad2Deg(phiMax)
	} else if fullySouthern {
		latMax = tempLatMax
		latMin = Rad2Deg(phiMin)
	}

	return BoundingBox{LatMin: latMin, LatMax: latMax, LonMin: lonMin, LonMax: lonMax}
}

func (p *Polygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	start := p.Size() - 1
	end := 0

	// check the [start,end] edge for intersection with the test ray
	contains := p.intersectsWithRaycast(point, p.At(start), p.At(end))
	// check each other edge for intersection with the test ray
	for i := 1; i < p.Size(); i++ {
		if p.intersectsWithRaycast(point, p.At(i-1), p.At(i)) {
			contains = !contains
		}
	}
	return contains
}

func (p *Polygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// based on paper: Some Algorithms for Polygons on a Sphere (Robert.G .Chamberlain)

	// ensure that start has the lower longitude
	if start.Lon() > end.Lon() {
		start, end = end, start
	}

	// Move the point a little bit to the east to avoid miscounting
	// -> those edges whose other end is westward will be counted,
	// while those whose other end is not westward will not
	for point.Lon() == start.Lon() || point.Lon() == end.Lon() {
		newLon := math.Nextafter(point.Lon(), math.Inf(1))
		point = NewPoint(point.Lat(), newLon)
	}

	// If the longitude of the ray is not between the longitudes of the ends of the edge,
	// there is no intersection
	if point.Lon() < start.Lon() || point.Lon() > end.Lon() {
		return false
	}

	// decide which point of the edge is norhterly
	if start.Lat() > end.Lat() {
		if point.Lat() > start.Lat() {
			// the point is above the edge -> it can't intersect with the edge
			return false
		}
		if point.Lat() < end.Lat() {
			// the point's ray intersects with the edge
			return true
		}
	} else {
		if point.Lat() > end.Lat() {
			// the point is above the edge -> it can't intersect with the edge
			return false
		}
		if point.Lat() < start.Lat() {
			// the point's ray intersects with the edge
			return true
		}
	}
	// Only if the test point is north of that chord is it necessary to compute the
	// latitude of the edge at the test point's longitude and compare it to the
	// latitude of Q
	//crossLat := start.LatitudeOnLineAtLon(end, point.Lon())
	crossLat := start.GreatCircleLatOfCrossingPoint(end, point.Lon())
	intersects := crossLat >= point.Lat()

	/*
		// following could be (slightly) faster
		raySlope := (point.Lon() - start.Lon()) / (point.Lat() - start.Lat())
		diagSlope := (end.Lon() - start.Lon()) / (end.Lat() - start.Lat())
		return raySlope >= diagSlope
	*/

	return intersects
}
