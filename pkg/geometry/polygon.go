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

func (p *Polygon) BoundingBox() BoundingBox {
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

func (p *Polygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	start := p.Size() - 1
	end := 0

	contains := p.intersectsWithRaycast(point, p.At(start), p.At(end))
	for i := 1; i < p.Size(); i++ {
		if p.intersectsWithRaycast(point, p.At(i-1), p.At(i)) {
			contains = !contains
		}
	}
	return contains
}

func (p *Polygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
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
	var polygon *Polygon
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
	for i := 0; i < polygon.Size()-1; i++ {
		vALat := polygon.At(i).Lat()
		vALon := polygon.At(i).Lon()
		tlonA := transformedLon[i]
		vBLat := polygon.At(i + 1).Lat()
		vBLon := polygon.At(i + 1).Lon()
		tlonB := transformedLon[i+1]

		strike := 0
		if tlonP == tlonA {
			strike = 1
		} else {
			brngAB := eastOrWest(NewPoint(0, tlonA), NewPoint(0, tlonB))
			brngAP := eastOrWest(NewPoint(0, tlonA), NewPoint(0, tlonP))
			brngPB := eastOrWest(NewPoint(0, tlonP), NewPoint(0, tlonB))
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
			brng_BX := eastOrWest(NewPoint(0, tlon_B), NewPoint(0, tlon_X))
			brng_BP := eastOrWest(NewPoint(0, tlon_B), NewPoint(0, tlon_P))
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
