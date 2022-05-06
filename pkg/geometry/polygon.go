package geometry

import "math"

// mostly based on this: https://github.com/kellydunn/golang-geo/blob/master/polygon.go
// and partly on this (not tested, very old): Locating a Point on a Spherical Surface Relative to a Spherical Polygon of Arbitrary Shape
// the Fortran translation can maybe get removed

type Polygon struct {
	points []*Point
}

func NewPolygon(points []*Point) *Polygon {
	return &Polygon{points: points}
}

func (p *Polygon) Points() []*Point {
	return p.points
}

func (p *Polygon) Add(point *Point) {
	p.points = append(p.points, point)
}

func (p *Polygon) IsClosed() bool {
	if len(p.points) < 3 || p.points[0].lat != p.points[len(p.points)-1].lat || p.points[0].lon != p.points[len(p.points)-1].lon {
		return false
	}
	return true
}

func (p *Polygon) BoundingBox() (float64, float64, float64, float64) {
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

func (p *Polygon) Contains(point *Point) bool {
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

func (p *Polygon) IntersectsWithRaycast(point *Point, start *Point, end *Point) bool {
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

func LocatePointRelBoundary(p *Point, xc *Point, boundary int64, nv_c int64, tlonv []float64) int {
	var dellon float64
	var crossCounter int
	var polygon *Polygon
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
			brngAB := EastOrWest(&Point{lon: tlonA}, &Point{lon: tlonB})
			brngAP := EastOrWest(&Point{lon: tlonA}, &Point{lon: tlonP})
			brngPB := EastOrWest(&Point{lon: tlonP}, &Point{lon: tlonB})
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
			brng_BX := EastOrWest(&Point{lon: tlon_B}, &Point{lon: tlon_X})
			brng_BP := EastOrWest(&Point{lon: tlon_B}, &Point{lon: tlon_X})
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
func EastOrWest(c *Point, d *Point) int {
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
