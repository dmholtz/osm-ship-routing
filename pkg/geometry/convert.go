package geometry

import "math"

const ratio = math.Pi / 180

// Convert degree to radian
func Deg2Rad(degree float64) float64 {
	return degree * ratio
}

// Convert radian to degree
func Rad2Deg(radian float64) float64 {
	return radian / ratio
}
