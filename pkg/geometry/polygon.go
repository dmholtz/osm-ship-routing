package geometry

import "math"

// mostly based on this: https://github.com/kellydunn/golang-geo/blob/master/polygon.go
// and partly on this (not tested, very old): Locating a Point on a Spherical Surface Relative to a Spherical Polygon of Arbitrary Shape
// the Fortran translation can maybe get removed

type Polygon2D struct {
	points []*Point2D
}

type PolygonSpherical struct {
	points []*PointSpherical
}

func NewPolygon(points []*PointSpherical) *PolygonSpherical {
	return &PolygonSpherical{points: points}
}

func (p *PolygonSpherical) Points() []*PointSpherical {
	return p.points
}

func (p *PolygonSpherical) Add(point *PointSpherical) {
	p.points = append(p.points, point)
}

func (p *PolygonSpherical) IsClosed() bool {
	if len(p.points) < 3 || p.points[0].lat != p.points[len(p.points)-1].lat || p.points[0].lon != p.points[len(p.points)-1].lon {
		return false
	}
	return true
}

func (p *PolygonSpherical) BoundingBox() (float64, float64, float64, float64) {
	var phiNorth, phiSouth, lambdaWest, lambdaEast float64
	phiNorth = p.points[0].Phi()
	phiSouth = p.points[0].Phi()
	lambdaWest = p.points[0].Lambda()
	lambdaEast = p.points[0].Lambda()
	for i := 1; i < len(p.points); i++ {
		if phiNorth < p.points[i].Phi() {
			phiNorth = p.points[i].Phi()
		}
		if phiSouth > p.points[i].Phi() {
			phiSouth = p.points[i].Phi()
		}
		if lambdaWest < p.points[i].Lambda() {
			lambdaWest = p.points[i].Lambda()
		}
		if lambdaEast > p.points[i].Lambda() {
			lambdaEast = p.points[i].Lambda()
		}
	}
	return phiNorth, phiSouth, lambdaWest, lambdaEast
}

func (p *PolygonSpherical) Contains(point *PointSpherical) bool {
	if !p.IsClosed() {
		return false
	}

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

func (p *PolygonSpherical) IntersectsWithRaycast(point *PointSpherical, start *PointSpherical, end *PointSpherical) bool {
	if start.lon > end.lon {
		start, end = end, start
	}
	for point.lon == start.lon || point.lon == end.lon {
		newLon := math.Nextafter(point.lon, math.Inf(1))
		point = NewPoint(point.lat, newLon)
	}
	if point.lon < start.lon || point.lon > end.lon {
		return false
	}
	if start.lat > end.lat {
		if point.lat > start.lat {
			return false
		}
		if point.lat < end.lat {
			return true
		}
	} else {
		if point.lat > end.lat {
			return false
		}
		if point.lat < start.lat {
			return true
		}
	}
	raySlope := (point.lon - start.lon) / (point.lat - start.lat)
	diagSlope := (end.lon - start.lon) / (end.lat - start.lat)

	return raySlope >= diagSlope
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

func LocatePointRelBoundary(p *PointSpherical, xc *PointSpherical, boundary int64, nv_c int64, tlonv []float64) int {
	var dellon float64
	var crossCounter int
	var polygon *PolygonSpherical
	var transformedLon []float64
	if boundary == 0 {
		panic("Boundary not defined")
	}
	if p.lat == -xc.lat {
		dellon = p.lon - xc.lon
		if dellon < -180 {
			dellon += 360
		} else if dellon > 180 {
			dellon -= 360
		}
		if math.Abs(dellon) == 180 {
			panic("P is antipodal to X: P relative to S is undertermined")
		}
	}

	crossCounter = 0

	if p.lat == xc.lat && p.lon == xc.lon {
		return 1
	}

	tlonP := TransformLon(xc.lat, xc.lon, p.lat, p.lon)
	for i := 0; i < len(polygon.points)-1; i++ {
		vALat := polygon.points[i].lat
		vALon := polygon.points[i].lon
		tlonA := transformedLon[i]
		vBLat := polygon.points[i+1].lat
		vBLon := polygon.points[i+1].lon
		tlonB := transformedLon[i+1]

		strike := 0
		if tlonP == tlonA {
			strike = 1
		} else {
			brngAB := EastOrWest(&PointSpherical{lon: tlonA}, &PointSpherical{lon: tlonB})
			brngAP := EastOrWest(&PointSpherical{lon: tlonA}, &PointSpherical{lon: tlonP})
			brngPB := EastOrWest(&PointSpherical{lon: tlonP}, &PointSpherical{lon: tlonB})
			if brngAP == brngAB && brngPB == brngAB {
				strike = 1
			}
		}
		if strike == 1 {
			if p.lat == vALat && p.lon == vALon {
				return 2 // P lies on a vertex of S
			}
			tlon_X := TransformLon(vALat, vALon, xc.lat, xc.lon)
			tlon_B := TransformLon(vALat, vALon, vBLat, vBLon)
			tlon_P := TransformLon(vALat, vALon, p.lat, p.lon)
			if tlon_P == tlon_B {
				return 2 // P lies on side of S
			}
			brng_BX := EastOrWest(&PointSpherical{lon: tlon_B}, &PointSpherical{lon: tlon_X})
			brng_BP := EastOrWest(&PointSpherical{lon: tlon_B}, &PointSpherical{lon: tlon_X})
			if brng_BX == -brng_BP {
				crossCounter++
			}
		}
	}
	if crossCounter%2 == 0 {
		return 1
	}
	return 0
}

// Determine the 'longitude' of a Point Q in a geographic coordinate system for which point P acts as a 'north pole'
func TransformLon(plat, plon, qlat, qlon float64) float64 {
	dtr := math.Pi / 180.0
	if plat == 90 {
		return qlon
	}
	t := math.Sin((qlon-plon)*dtr) * math.Cos(qlat*dtr)
	b := math.Sin(dtr*qlat)*math.Cos(plat*dtr) - math.Cos(qlat*dtr)*math.Sin(plat*dtr)*math.Cos((qlon-plon)*dtr)
	return math.Atan2(t, b) / dtr

}

// Determine if the shorted path form c to d is east or west
func EastOrWest(c *PointSpherical, d *PointSpherical) int {
	delta := d.lon - c.lon
	if delta > 180 {
		delta -= 360
	} else if delta < -180 {
		delta += 360
	}
	if delta > 0 && delta != 180 {
		return -1 // d is west of c
	}
	if delta < 0 && delta != -180 {
		return 1 // d is east of c
	}
	return 0 // neither is or west -> antipode or same longitude
}
