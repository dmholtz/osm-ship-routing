package server

import (
	"fmt"
	"math"
	"time"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	heur "github.com/dmholtz/graffiti/examples/heuristics"
	g "github.com/dmholtz/graffiti/graph"
)

type IGeoPoint interface {
	g.GeoPoint | g.PartGeoPoint | g.TwoLevelPartGeoPoint
}

type ShipRouter interface {
	ProcessRequest(req RouteRequest, showSearchSpace bool) RouteResponse
	String() string
}

type ShipRouter1[N IGeoPoint, E g.IHalfEdge] struct {
	Graph  g.Graph[N, E]
	Router sp.Router[int]
}

func (sr ShipRouter1[N, E]) ProcessRequest(req RouteRequest, showSearchSpace bool) RouteResponse {

	source := findClosestNode(sr.Graph, req.Origin)
	target := findClosestNode(sr.Graph, req.Destination)

	startTime := time.Now()
	res := sr.Router.Route(source, target, showSearchSpace)
	elapsed := time.Since(startTime).Milliseconds()

	var path Path
	if res.Length > 0 {
		waypoints := make([]Point, 0)
		for _, nodeId := range res.Path {
			node := sr.Graph.GetNode(nodeId)
			waypoints = append(waypoints, getPoint(node))
		}
		path = Path{Length: res.Length, Waypoints: waypoints}
	}

	var searchSpace []Point
	if showSearchSpace {
		searchSpace = make([]Point, 0)
		for _, nodeId := range res.SearchSpace {
			node := sr.Graph.GetNode(nodeId)
			searchSpace = append(searchSpace, getPoint(node))
		}
	}

	return RouteResponse{Exists: res.Length > 0, Time: elapsed, Path: path, SearchSpace: searchSpace}
}

func (sr ShipRouter1[N, E]) String() string {
	return sr.Router.(fmt.Stringer).String()
}

func getPoint[N IGeoPoint](n N) Point {
	var point Point
	switch p := any(n).(type) {
	case g.GeoPoint:
		point = Point{Lat: p.Lat, Lon: p.Lon}
	case g.PartGeoPoint:
		point = Point{Lat: p.Lat, Lon: p.Lon}
	case g.TwoLevelPartGeoPoint:
		point = Point{Lat: p.Lat, Lon: p.Lon}
	default:
		panic("Node type is not a IGeoPoint.")
	}
	return point
}

func findClosestNode[N IGeoPoint, E g.IHalfEdge](graph g.Graph[N, E], p Point) g.NodeId {
	gp := g.GeoPoint{Lat: p.Lat, Lon: p.Lon}
	minDist := math.MaxInt
	var closestNode g.NodeId
	for nodeId := 0; nodeId < graph.NodeCount(); nodeId++ {
		var otherGp g.GeoPoint
		switch other := any(graph.GetNode(nodeId)).(type) {
		case g.GeoPoint:
			otherGp = other
		case g.PartGeoPoint:
			otherGp = g.GeoPoint{Lat: other.Lat, Lon: other.Lon}
		case g.TwoLevelPartGeoPoint:
			otherGp = g.GeoPoint{Lat: other.Lat, Lon: other.Lon}
		}
		if dist := heur.Haversine(otherGp, gp); dist < minDist {
			minDist = dist
			closestNode = nodeId
		}
	}
	return closestNode
}
