package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Point [3]float64

func (p Point) String() string {
	return fmt.Sprintf("%f, %f, %f", p[0], p[1], p[2])
}

func toCartesian(r float64, theta float64, phi float64) Point {
	var cartesian Point
	cartesian[0] = r * math.Sin(theta) * math.Cos(phi)
	cartesian[1] = r * math.Sin(theta) * math.Sin(phi)
	cartesian[2] = r * math.Cos(theta)
	return cartesian
}

/*
Mapping to geopgraphic coordinates:
- longitude (lambda) \in [-180, +180] -> refers to phi \in [0, 360]
- latitude (theta) \in [-90, +90]
*/

func equidistribute(n int) {
	nCount := 0
	r := 1.0

	points := make([]Point, 0)

	a := 4.0 * math.Pi * r * r / float64(n)
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta
	for m := 0; m < int(mTheta-1); m++ {
		theta := math.Pi * (float64(m) + 0.5) / mTheta
		mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
		for n := 0; n < int(mPhi-1); n++ {
			phi := 2 * math.Pi * float64(n) / mPhi
			points = append(points, toCartesian(r, theta, phi))
			nCount++
		}
		nCount++
	}
	fmt.Printf("Number of points: %d\n", nCount)

	file, err := os.Create("points.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range points {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}

func simpleDistribution(n int) {

	points := make([]Point, 0)

	lat := 0.0 //-math.Pi / 2.0
	lon := -math.Pi

	dLat := math.Pi / float64(n)
	dLon := 2.0 * math.Pi / float64(n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			points = append(points, toCartesian(1.0, lat, lon))
			lon += dLon
		}
		lon = -math.Pi
		lat += dLat
	}

	file, err := os.Create("points.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range points {
		fmt.Fprintln(w, line)
	}
	w.Flush()

}

func main() {
	//equidistribute(500)
	simpleDistribution(20)
}
