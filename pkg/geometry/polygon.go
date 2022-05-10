package geometry

import "math"

// mostly based on this: https://github.com/kellydunn/golang-geo/blob/master/polygon.go
// and partly on this (not tested, very old): Locating a Point on a Spherical Surface Relative to a Spherical Polygon of Arbitrary Shape
// the Fortran translation can maybe get removed

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

type StandardPolygon struct {
	points []*Point
}

func NewStandardPolygon(points []*Point) *StandardPolygon {
	return &StandardPolygon{points: points}
}

func (p *StandardPolygon) Points() []*Point {
	return p.points
}

func (p *StandardPolygon) At(index int) *Point {
	return p.points[index]
}

func (p *StandardPolygon) Add(point *Point) {
	p.points = append(p.points, point)
}

func (p *StandardPolygon) Size() int {
	return len(p.points)
}

func (p *StandardPolygon) IsClosed() bool {
	if len(p.points) < 3 || p.At(0).Lat() != p.At(len(p.points)-1).Lat() || p.At(0).Lon() != p.At(len(p.points)-1).Lon() {
		return false
	}
	return true
}

func (p *StandardPolygon) BoundingBox() BoundingBox {
	latMin, lonMin := math.Inf(1), math.Inf(1)
	latMax, lonMax := math.Inf(-1), math.Inf(-1)
	for _, point := range p.points {
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

func (p *StandardPolygon) BoundingBox_() (float64, float64, float64, float64) {
	// TODO: don't convert to Phi / Lambda
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
		// check lambda east / west definition
		if lambdaWest < p.points[i].Lambda() {
			lambdaWest = p.points[i].Lambda()
		}
		if lambdaEast > p.points[i].Lambda() {
			lambdaEast = p.points[i].Lambda()
		}
	}
	return phiNorth, phiSouth, lambdaWest, lambdaEast
}

func (p *StandardPolygon) BBoxContains(point *Point) bool {
	phiNorth, phiSouth, lambdaWest, lambdaEast := p.BoundingBox_()
	// TODO: check how lambda is defined
	return point.Phi() <= phiNorth && point.Phi() >= phiSouth && point.Lambda() >= lambdaEast && point.Lambda() <= lambdaWest
}

func (p *StandardPolygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	start := len(p.points) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, p.points[start], p.points[end])
	for i := 1; i < len(p.points); i++ {
		if p.intersectsWithRaycast(point, p.points[i-1], p.points[i]) {
			contains = !contains
		}
	}
	return contains
}

func (p *StandardPolygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	if start.Lon() > end.Lon() {
		start, end = end, start
	}
	for point.Lon() == start.Lon() || point.Lon() == end.Lon() {
		newLon := math.Nextafter(point.Lon(), math.Inf(1))
		point = NewPoint(point.Lat(), newLon)
	}
	if point.Lon() < start.Lon() || point.Lon() > end.Lon() {
		return false
	}
	if start.Lat() > end.Lat() {
		if point.Lat() > start.Lat() {
			return false
		}
		if point.Lat() < end.Lat() {
			return true
		}
	} else {
		if point.Lat() > end.Lat() {
			return false
		}
		if point.Lat() < start.Lat() {
			return true
		}
	}
	raySlope := (point.Lon() - start.Lon()) / (point.Lat() - start.Lat())
	diagSlope := (end.Lon() - start.Lon()) / (end.Lat() - start.Lat())

	return raySlope >= diagSlope
}

func locatePointRelBoundary(p *Point, xc *Point, boundary int64, nv_c int64, tlonv []float64) int {
	var dellon float64
	var crossCounter int
	var polygon *StandardPolygon
	var transformedLon []float64
	if boundary == 0 {
		panic("Boundary not defined")
	}
	if p.Lat() == -xc.Lat() {
		dellon = p.Lon() - xc.Lon()
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

	if p.Lat() == xc.Lat() && p.Lon() == xc.Lon() {
		return 1
	}

	tlonP := transformLon(xc.Lat(), xc.Lon(), p.Lat(), p.Lon())
	for i := 0; i < len(polygon.points)-1; i++ {
		vALat := polygon.points[i].Lat()
		vALon := polygon.points[i].Lon()
		tlonA := transformedLon[i]
		vBLat := polygon.points[i+1].Lat()
		vBLon := polygon.points[i+1].Lon()
		tlonB := transformedLon[i+1]

		strike := 0
		if tlonP == tlonA {
			strike = 1
		} else {
			brngAB := eastOrWest(&Point{lon: tlonA}, &Point{lon: tlonB})
			brngAP := eastOrWest(&Point{lon: tlonA}, &Point{lon: tlonP})
			brngPB := eastOrWest(&Point{lon: tlonP}, &Point{lon: tlonB})
			if brngAP == brngAB && brngPB == brngAB {
				strike = 1
			}
		}
		if strike == 1 {
			if p.Lat() == vALat && p.Lon() == vALon {
				return 2 // P lies on a vertex of S
			}
			tlon_X := transformLon(vALat, vALon, xc.Lat(), xc.Lon())
			tlon_B := transformLon(vALat, vALon, vBLat, vBLon)
			tlon_P := transformLon(vALat, vALon, p.Lat(), p.Lon())
			if tlon_P == tlon_B {
				return 2 // P lies on side of S
			}
			brng_BX := eastOrWest(&Point{lon: tlon_B}, &Point{lon: tlon_X})
			brng_BP := eastOrWest(&Point{lon: tlon_B}, &Point{lon: tlon_X})
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
func transformLon(plat, plon, qlat, qlon float64) float64 {
	dtr := math.Pi / 180.0
	if plat == 90 {
		return qlon
	}
	t := math.Sin((qlon-plon)*dtr) * math.Cos(qlat*dtr)
	b := math.Sin(dtr*qlat)*math.Cos(plat*dtr) - math.Cos(qlat*dtr)*math.Sin(plat*dtr)*math.Cos((qlon-plon)*dtr)
	return math.Atan2(t, b) / dtr

}

// Determine if the shorted path form c to d is east or west
func eastOrWest(c *Point, d *Point) int {
	delta := d.Lon() - c.Lon()
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
