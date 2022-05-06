package geometry

import "math"

type Point struct {
	lat float64 // latitude in degree
	lon float64 // longitude in degree
}

func NewPoint(lat, lon float64) *Point {
	return &Point{lat: lat, lon: lon}
}

// TODO Testing:  For the Cartesian Coordinates (1, 2, 3), the Spherical-Equivalent Coordinates are (√(14), 36.7°, 63.4°).
// TODO: Avoid spherical trigonometrical computaitons by first checking the bounding box (cf. some algorithms for polygon on a sphere)

// Latitude in radian
func (p *Point) Phi() float64 {
	return p.lat * math.Pi / 180 // latitude in radian
}

// Longitude in radian
func (p *Point) Lambda() float64 {
	return p.lon * math.Pi / 180 // longitude in radian
}

func (p *Point) X() float64 {
	R := 6371e3 // earth radius
	return R * math.Sin(p.Phi()) * math.Cos(p.Lambda())
}

func (p *Point) Y() float64 {
	R := 6371e3 // earth radius
	return R * math.Sin(p.Phi()) * math.Sin(p.Lambda())
}

func (p *Point) Z() float64 {
	R := 6371e3 // earth radius
	return R * math.Cos(p.Phi())
}

// The great circle distance
func (p *Point) Haversine(point *Point) float64 {
	R := 6371e3 // earth radius
	phi1 := p.lat * math.Pi / 180
	phi2 := point.lat * math.Pi / 180
	deltaPhi := (point.lat - p.lat) * math.Pi / 180
	deltaLambda := (point.lon - p.lon) * math.Pi / 180

	a := math.Pow(math.Sin(deltaPhi/2), 2) + math.Cos(phi1)*math.Cos(phi2)*math.Pow(math.Sin(deltaLambda/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
