package graph

import (
	"math"

	geo "github.com/dmholtz/osm-ship-routing/pkg/geometry"
)

type Route struct {
	Origin      geo.Point
	Destination geo.Point
	Exists      bool        // true iff a route from origin to destination exists
	Waypoints   []geo.Point // sequence of points that describe the route
	Length      int         // length of the route
}

type ShipRouter struct {
	g Graph
}

func NewShipRouter(g Graph) *ShipRouter {
	return &ShipRouter{g: g}
}

func (sr ShipRouter) closestNodes(p1, p2 *geo.Point) (n1, n2 int) {
	n1, n2 = 0, 0
	d1, d2 := math.MaxInt, math.MaxInt

	for i := 0; i < sr.g.NodeCount(); i++ {
		testPoint := nodeToPoint(sr.g.GetNode(i))
		distance := p1.IntHaversine(testPoint)
		if distance < d1 {
			n1 = i
		}

	}
	return n1, n2
}

func nodeToPoint(n Node) *geo.Point {
	return geo.NewPoint(n.Lat, n.Lon)
}

func ComputeRoute(origin, destination geo.Point) Route {
	return Route{}
}
