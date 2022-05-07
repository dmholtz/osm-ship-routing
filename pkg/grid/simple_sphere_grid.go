package grid

import (
	"github.com/dmholtz/osm-ship-routing/pkg/geometry"
	gr "github.com/dmholtz/osm-ship-routing/pkg/graph"
)

type SphereGridGraph struct {
	nLon      int // discretization of longitude, i.e number of points in [lonMin, lonMax]
	nLat      int // discretization of latitude, i.e number of points in [latMin, latMax]
	GridGraph gr.DynamicGraph
}

func NewSphereGridGraph(nLon int, nLat int) *SphereGridGraph {

	if nLon < 1 {
		panic(nLon)
	}
	if nLat < 2 {
		panic(nLat)
	}
	alg := gr.AdjacencyListGraph{}
	sgg := SphereGridGraph{nLon: nLon, nLat: nLat, GridGraph: &alg}
	return &sgg
}

func (sgg *SphereGridGraph) DistributeNodes() {
	lat := LatMin
	lon := LonMin

	dLat := (LatMax - LatMin) / (float64(sgg.nLat) - 1)
	dLon := (LonMax - LonMin) / float64(sgg.nLon)
	for iLat := 0; iLat < sgg.nLat; iLat++ {
		for iLon := 0; iLon < sgg.nLon; iLon++ {
			sgg.GridGraph.AddNode(gr.Node{Lon: lon, Lat: lat})
			lon += dLon
		}
		lon = LonMin
		lat += dLat
	}
}

func (sgg *SphereGridGraph) CreateEdges(polygons []geometry.Polygon) {
	for nodeId := 0; nodeId < sgg.GridGraph.NodeCount(); nodeId++ {
		neighbors := sgg.neighborsOf(nodeId)
		for _, neighbor := range neighbors {
			lat := sgg.GridGraph.GetNode(nodeId).Lat
			lon := sgg.GridGraph.GetNode(nodeId).Lon

			from := geometry.NewPoint(lat, lon)
			to := geometry.NewPoint(sgg.GridGraph.GetNode(neighbor).Lat, sgg.GridGraph.GetNode(neighbor).Lon)

			for _, pol := range polygons {
				if !(!pol.Contains(from) && !pol.Contains(to)) {
					edge := gr.Edge{From: nodeId, To: neighbor, Distance: 1} // todo: compute distance
					sgg.GridGraph.AddEdge(edge)
				}
			}
		}
	}
}

func (sgg *SphereGridGraph) neighborsOf(nodeId int) []int {
	neighbors := make([]int, 0)
	if nodeId < sgg.nLon*(sgg.nLat-1) {
		// northern neighbor
		neighbors = append(neighbors, nodeId+sgg.nLon)
	}
	if nodeId >= sgg.nLon {
		// southern neighbor
		neighbors = append(neighbors, nodeId-sgg.nLon)
	}

	// western neighbor
	if nodeId%sgg.nLon != 0 {
		neighbors = append(neighbors, nodeId-1)
	} else {
		neighbors = append(neighbors, nodeId+sgg.nLon-1)
	}

	// eastern neighbor
	if (nodeId+1)%sgg.nLon != 0 {
		neighbors = append(neighbors, nodeId+1)
	} else {
		neighbors = append(neighbors, nodeId-sgg.nLon+1)
	}
	return neighbors
}
