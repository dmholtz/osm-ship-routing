package grid

import (
	"sync"

	"github.com/dmholtz/osm-ship-routing/pkg/geometry"
	gr "github.com/dmholtz/osm-ship-routing/pkg/graph"
)

type SphereGridGraph struct {
	nLon      int // discretization of longitude, i.e number of points in [lonMin, lonMax]
	nLat      int // discretization of latitude, i.e number of points in [latMin, latMax]
	GridGraph gr.DynamicGraph
	isWater   []bool
}

func NewSphereGridGraph(nLon int, nLat int) *SphereGridGraph {

	if nLon < 1 {
		panic(nLon)
	}
	if nLat < 2 {
		panic(nLat)
	}
	alg := gr.AdjacencyListGraph{}
	sgg := SphereGridGraph{nLon: nLon, nLat: nLat, GridGraph: &alg, isWater: make([]bool, 0)}
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

func (sgg *SphereGridGraph) LandWaterTest(polygons []geometry.Polygon) {
	sgg.isWater = make([]bool, sgg.GridGraph.NodeCount(), sgg.GridGraph.NodeCount())

	// pre-compute bounding boxes for every polygon
	bboxes := make([]geometry.BoundingBox, len(polygons), len(polygons))

	var wg sync.WaitGroup
	for i, polygon := range polygons {
		wg.Add(1)
		go func(i int, polygon geometry.Polygon) {
			bbox := polygon.BoundingBox()
			bboxes[i] = bbox
			wg.Done()
		}(i, polygon)
	}
	wg.Wait()

	wg.Add(sgg.GridGraph.NodeCount())
	for nodeId := 0; nodeId < sgg.GridGraph.NodeCount(); nodeId++ {
		go func(nodeId int) {
			sgg.isWater[nodeId] = true
			testPoint := geometry.NewPoint(sgg.GridGraph.GetNode(nodeId).Lat, sgg.GridGraph.GetNode(nodeId).Lon)
			for i, pol := range polygons {
				// roughly check, whether the point is contained in the bounding box of the polygon
				if bboxes[i].Contains(*testPoint) {
					// precisely check, whether the polygon contains the point
					if pol.Contains(testPoint) {
						sgg.isWater[nodeId] = false
						break
					}
				}
			}
			wg.Done()
		}(nodeId)
	}
	wg.Wait()

	// invert map (just for demo)
	for i := range sgg.isWater {
		sgg.isWater[i] = !sgg.isWater[i]
	}
}

func (sgg *SphereGridGraph) CreateEdges() {
	for nodeId := 0; nodeId < sgg.GridGraph.NodeCount(); nodeId++ {
		neighbors := sgg.neighborsOf(nodeId)
		for _, neighbor := range neighbors {
			if sgg.isWater[nodeId] == true && sgg.isWater[neighbor] == true {
				edge := gr.Edge{From: nodeId, To: neighbor, Distance: 1} // todo: compute distance
				sgg.GridGraph.AddEdge(edge)
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
