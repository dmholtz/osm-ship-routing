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
func (first *Point) Haversine(second *Point) float64 {
	R := 6371e3 // earth radius
	// one can reduce one function call/calculation by directly substraction the latitudes/longitudes and then convert to radian:
	// (point.lat - p.lat) * math.Pi / 180
	// (point.lon - p.lon) * math.Pi / 180
	// But this is probably not worth to improve
	deltaPhi := second.Phi() - first.Phi()
	deltaLambda := second.Lambda() - first.Lambda()

	a := math.Pow(math.Sin(deltaPhi/2), 2) + math.Cos(first.Phi())*math.Cos(second.Phi())*math.Pow(math.Sin(deltaLambda/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// Calculate the distance with the Spherical Law of Cosines.
// This is a roughly simpler formula (which may improve the performance).
// In other tests however, it took a  bit longer
func (first *Point) SphericalCosineDistance(second *Point) float64 {
	// may be better to only calculate Phi only once for each point and store in local variable. But probably no big improvements
	R := 6371e3 // earth radius
	return math.Acos(math.Sin(first.Phi())*math.Sin(second.Phi())+math.Cos(first.Phi())*math.Cos(second.Phi())*math.Cos(second.Lambda()-first.Lambda())) * R
}

func (first *Point) Bearing(second *Point) float64 {
	y := math.Sin(second.Lambda()-first.Lambda()) * math.Cos(second.Phi())
	x := math.Cos(first.Phi())*math.Sin(second.Phi()) -
		math.Sin(first.Phi())*math.Cos(second.Phi())*math.Cos(second.Lambda()-first.Lambda())
	theta := math.Atan2(y, x)
	bearing := math.Mod(theta*180/math.Pi+360, 360) // bearing in degrees
	return bearing
}
