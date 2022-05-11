package geometry

import "math"

const earthRadius = 6371e3

type Point [2]float64

// Create a new Point with latitude and longitude (both in degree)
func NewPoint(lat, lon float64) *Point {
	return &Point{lat, lon}
}

func NewPointFromBearing(initialPoint *Point, bearing float64, distance float64) *Point {
	phi := math.Asin(math.Sin(initialPoint.Phi())*math.Cos(distance/earthRadius) + math.Cos(initialPoint.Phi())*math.Sin(distance/earthRadius)*math.Cos(bearing))
	lambda := initialPoint.Lambda() + math.Atan2(math.Sin(bearing)*math.Sin(distance/earthRadius)*math.Cos(initialPoint.Phi()),
		math.Cos(distance/earthRadius)-math.Sin(initialPoint.Phi())*math.Sin(phi))
	return NewPoint(Rad2Deg(phi), Rad2Deg(lambda))
}

// TODO Testing:  For the Cartesian Coordinates (1, 2, 3), the Spherical-Equivalent Coordinates are (√(14), 36.7°, 63.4°).
// TODO: Avoid spherical trigonometrical computaitons by first checking the bounding box (cf. some algorithms for polygon on a sphere)

// latitude in degree
func (p *Point) Lat() float64 {
	return p[0]
}

// longitude in degree
func (p *Point) Lon() float64 {
	return p[1]
}

// Latitude in radian
func (p *Point) Phi() float64 {
	return Deg2Rad(p.Lat())
}

// Longitude in radian
func (p *Point) Lambda() float64 {
	return Deg2Rad(p.Lon())
}

func (p *Point) X() float64 {
	return earthRadius * math.Sin(p.Phi()) * math.Cos(p.Lambda())
}

func (p *Point) Y() float64 {
	return earthRadius * math.Sin(p.Phi()) * math.Sin(p.Lambda())
}

func (p *Point) Z() float64 {
	return earthRadius * math.Cos(p.Phi())
}

// The great circle distance
func (first *Point) Haversine(second *Point) float64 {
	// one can reduce one function call/calculation by directly substraction the latitudes/longitudes and then convert to radian:
	// (point.Lat() - p.Lat()) * math.Pi / 180
	// (point.Lon() - p.Lon()) * math.Pi / 180
	// But this is probably not worth to improve
	deltaPhi := second.Phi() - first.Phi()
	deltaLambda := second.Lambda() - first.Lambda()

	a := math.Pow(math.Sin(deltaPhi/2), 2) + math.Cos(first.Phi())*math.Cos(second.Phi())*math.Pow(math.Sin(deltaLambda/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// Calculate the distance with the Spherical Law of Cosines.
// This is a roughly simpler formula (which may improve the performance).
// In other tests however, it took a  bit longer
func (first *Point) SphericalCosineDistance(second *Point) float64 {
	// may be better to only calculate Phi only once for each point and store in local variable. But probably no big improvements
	return math.Acos(math.Sin(first.Phi())*math.Sin(second.Phi())+math.Cos(first.Phi())*math.Cos(second.Phi())*math.Cos(second.Lambda()-first.Lambda())) * earthRadius
}

func (first *Point) Bearing(second *Point) float64 {
	y := math.Sin(second.Lambda()-first.Lambda()) * math.Cos(second.Phi())
	x := math.Cos(first.Phi())*math.Sin(second.Phi()) -
		math.Sin(first.Phi())*math.Cos(second.Phi())*math.Cos(second.Lambda()-first.Lambda())
	theta := math.Atan2(y, x)
	bearing := math.Mod(Rad2Deg(theta)+360, 360) // bearing in degrees
	return bearing
}

// Half-way point along a great circle path between two points
func (first *Point) Midpoint(second *Point) *Point {
	Bx := math.Cos(second.Phi()) * math.Cos(second.Lambda()-first.Lambda())
	By := math.Cos(second.Phi()) * math.Sin(second.Lambda()-first.Lambda())
	phi := math.Atan2(math.Sin(first.Phi())+math.Sin(second.Phi()),
		math.Sqrt(math.Pow(math.Cos(first.Phi())+Bx, 2)+math.Pow(By, 2)))
	lambda := first.Lambda() + math.Atan2(By, math.Cos(first.Phi())+Bx)
	return NewPoint(Rad2Deg(phi), Rad2Deg(lambda))
	// The longitude can be normalised to −180…+180 using (lon+540)%360-180
}
