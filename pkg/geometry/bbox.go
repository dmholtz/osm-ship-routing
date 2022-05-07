package geometry

type BoundingBox struct {
	LatMin float64
	LatMax float64
	LonMin float64
	LonMax float64
}

func (bbox BoundingBox) Contains(point Point) bool {
	return bbox.LatMin <= point.lat && point.lat <= bbox.LatMax && bbox.LonMin <= point.lon && point.lon <= bbox.LonMax
}
