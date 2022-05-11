package geometry

type BoundingBox struct {
	LatMin float64
	LatMax float64
	LonMin float64
	LonMax float64
}

func (bbox BoundingBox) Contains(point Point) bool {
	return bbox.LatMin <= point.Lat() && point.Lat() <= bbox.LatMax && bbox.LonMin <= point.Lon() && point.Lon() <= bbox.LonMax
}
