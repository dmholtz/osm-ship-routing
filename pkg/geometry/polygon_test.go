package geometry

import (
	"testing"
)

// lat, lon coordinates for brunei
var brunei [][]float64 = [][]float64{
	{4.525874, 114.204017},
	{4.900011, 114.599961},
	{5.44773, 115.45071},
	{4.955228, 115.4057},
	{4.316636, 115.347461},
	{4.348314, 114.869557},
	{4.007637, 114.659596},
	{4.525874, 114.204017},
}

// lat, lon coordinates for capital of brunei(lies inside brunei)
var coordinatesOfBruneiCapital []float64 = []float64{
	4.9402900, 114.9480600,
}

var coordinatesOutsideBrunei []float64 = []float64{
	47.45, 122.30,
}

// lat, lon coordinates for a "low rendered antarctis"
var lowRenderedAntarctis [][]float64 = [][]float64{
	{-78.34941069014627, -30.234375},
	{-77.76758238272801, -57.65624999999999},
	{-75.67219739055291, -126.91406249999999},
	{-81.03861703916249, -163.4765625},
	{-80.05804956215623, 160.3125},
	{-69.162557908105, 149.0625},
	{-71.41317683396565, 11.6015625},
	{-78.34941069014627, -30.234375},
}

// lat, lon coordinates for a point inside "low rendered antarctis"
var coordinatesInAntarctis []float64 = []float64{
	-77.54209596075546, 38.67187499999999,
}

var coordinatesOutsideAntarctis [][]float64 = [][]float64{
	{38.8225909761771, -26.015625},
	{-66.93006025862445, -34.80468749999999},
}

// an area (at antarctis), where the sphere really matters. A direct line is not necesarrily the shortest path between two points
var curvedArea [][]float64 = [][]float64{
	{-78.20656311074711, -35.15625},
	{-78.9039293885709, -49.5703125},
	{-77.23507365492469, 48.515625},
	{-69.9001176266854, 37.265625},
	{-78.20656311074711, -35.15625},
}

var pointInsideCurvedArea []float64 = []float64{
	-78.9039293885709, 0,
}

var pointOutsideCurvedArea []float64 = []float64{
	-76.88077457250164, -17.9296875,
}

func TestPointInPolygon(t *testing.T) {
	bruneiPoints := make([]*Point, len(brunei))
	for i := range brunei {
		bruneiPoints[i] = NewPoint(brunei[i][0], brunei[i][1])
	}
	bruneiArea := NewPolygon(bruneiPoints)
	if !bruneiArea.Contains(NewPoint(coordinatesOfBruneiCapital[0], coordinatesOfBruneiCapital[1])) {
		t.Errorf("Point should lie in polygon, but isn't")
	}

	antarctisPoints := make([]*Point, len(lowRenderedAntarctis))
	for i := range lowRenderedAntarctis {
		antarctisPoints[i] = NewPoint(lowRenderedAntarctis[i][0], lowRenderedAntarctis[i][1])
	}
	antarctisArea := NewPolygon(antarctisPoints)
	if !antarctisArea.Contains(NewPoint(coordinatesInAntarctis[0], coordinatesInAntarctis[1])) {
		t.Errorf("Point should lie in polygon, but isn't")
	}
}

func TestPointNotInPolygon(t *testing.T) {
	bruneiPoints := make([]*Point, len(brunei))
	for i := range brunei {
		bruneiPoints[i] = NewPoint(brunei[i][0], brunei[i][1])
	}
	bruneiArea := NewPolygon(bruneiPoints)
	if bruneiArea.Contains(NewPoint(coordinatesOutsideBrunei[0], coordinatesOutsideBrunei[1])) {
		t.Errorf("Point should not lie in polygon, but does so")
	}

	antarctisPoints := make([]*Point, len(lowRenderedAntarctis))
	for i := range lowRenderedAntarctis {
		antarctisPoints[i] = NewPoint(lowRenderedAntarctis[i][0], lowRenderedAntarctis[i][1])
	}
	antarctisArea := NewPolygon(antarctisPoints)
	if antarctisArea.Contains(NewPoint(coordinatesOutsideAntarctis[0][0], coordinatesOutsideAntarctis[0][1])) {
		t.Errorf("Point should not lie in polygon, but does so")
	}
	if antarctisArea.Contains(NewPoint(coordinatesOutsideAntarctis[1][0], coordinatesOutsideAntarctis[1][1])) {
		t.Errorf("Point should not lie in polygon, but does so")
	}
}

/*
func TestPointInCurvedPolygon(t *testing.T) {
	points := make([]*Point, len(curvedArea))
	for i := range curvedArea {
		points[i] = NewPoint(curvedArea[i][0], curvedArea[i][1])
	}
	curvedPolygon := NewPolygon(points)

	if curvedPolygon.Contains(NewPoint(pointOutsideCurvedArea[0], pointOutsideCurvedArea[1])) {
		t.Errorf("Point should not lie in polygon, but does so")
	}

	if !curvedPolygon.Contains(NewPoint(pointInsideCurvedArea[0], pointInsideCurvedArea[1])) {
		t.Errorf("Point should lie in polygon, but isn't")
	}

}
*/
