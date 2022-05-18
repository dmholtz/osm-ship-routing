package graph

import "github.com/dmholtz/osm-ship-routing/pkg/geometry"

type Route struct {
	Origin      geometry.Point
	Destination geometry.Point
	Exists      bool             // true iff a route from origin to destination exists
	Waypoints   []geometry.Point // sequence of points that describe the route
	Length      int              // length of the route
}

type ShipRouter struct {
	g Graph
}

func NewShipRouter(g Graph) *ShipRouter {
	return &ShipRouter{g: g}
}

func ComputeRoute(origin, destination geometry.Point) Route {
	return Route{}
}
