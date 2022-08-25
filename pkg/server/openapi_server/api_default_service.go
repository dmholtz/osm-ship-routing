// SPDX-License-Identifier: MIT

package openapi_server

import (
	"context"
	"math"
	"net/http"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	heur "github.com/dmholtz/graffiti/examples/heuristics"
	g "github.com/dmholtz/graffiti/graph"
)

// DefaultApiService is a service that implements the logic for the DefaultApiServicer
// This service should implement the business logic for every endpoint for the DefaultApi API.
// Include any external packages or services that will be required by this service.
type DefaultApiService struct {
	Graph  g.Graph[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]]
	Router sp.Router[int]
}

// ComputeRoute - Compute a new route
func (s *DefaultApiService) ComputeRoute(ctx context.Context, routeRequest RouteRequest) (ImplResponse, error) {
	origin := g.PartGeoPoint{GeoPoint: g.GeoPoint{Lat: float64(routeRequest.Origin.Lat), Lon: float64(routeRequest.Origin.Lon)}}
	destination := g.PartGeoPoint{GeoPoint: g.GeoPoint{Lat: float64(routeRequest.Destination.Lat), Lon: float64(routeRequest.Destination.Lon)}}

	source := findClosestNode(s.Graph, origin)      // find closest source node
	target := findClosestNode(s.Graph, destination) // find closest target node

	shortestPathResult := s.Router.Route(source, target, false)

	routeResult := RouteResult{Origin: routeRequest.Origin, Destination: routeRequest.Destination}
	if shortestPathResult.Length < math.MaxInt {
		routeResult.Reachable = true
		waypoints := make([]Point, 0)
		for _, nodeId := range shortestPathResult.Path {
			geoPoint := s.Graph.GetNode(nodeId)
			p := Point{Lat: float32(geoPoint.Lat), Lon: float32(geoPoint.Lon)}
			waypoints = append(waypoints, p)
		}
		routeResult.Path = Path{Length: int32(shortestPathResult.Length), Waypoints: waypoints}
	} else {
		routeResult.Reachable = false
	}

	return Response(http.StatusOK, routeResult), nil
}

func findClosestNode(graph g.Graph[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]], node g.PartGeoPoint) g.NodeId {
	minDist := math.MaxInt
	var closestNode g.NodeId
	for nodeId := 0; nodeId < graph.NodeCount(); nodeId++ {
		other := graph.GetNode(nodeId)
		if dist := heur.Haversine(node.GeoPoint, other.GeoPoint); dist < minDist {
			minDist = dist
			closestNode = nodeId
		}
	}
	return closestNode
}
