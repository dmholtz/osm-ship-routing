package grid

import (
	"fmt"
	"sync"
	"time"

	geo "github.com/dmholtz/osm-ship-routing/pkg/geometry"
	gr "github.com/dmholtz/osm-ship-routing/pkg/graph"
)

type SimpleSphereGrid struct {
	nLon       int // discretization of longitude, i.e number of points in [lonMin, lonMax]
	nLat       int // discretization of latitude, i.e number of points in [latMin, latMax]
	points     []geo.Point
	isWater    []bool
	grid2nodes map[int]int
	nodes2grid []int
	nodes      []gr.Node
	edges      []gr.Edge
}

func NewSimpleSphereGrid(nLon int, nLat int, coastlines []geo.Polygon) *SimpleSphereGrid {
	if nLon < 1 {
		panic(nLon)
	}
	if nLat < 2 {
		panic(nLat)
	}
	ssg := SimpleSphereGrid{nLon: nLon, nLat: nLat}

	start := time.Now()
	ssg.distributePoints()
	elapsed := time.Since(start)
	fmt.Printf("[TIME] Distribute Points on grid: %s\n", elapsed)

	start = time.Now()
	ssg.landWaterTest(coastlines)
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Land / Water test: %s\n", elapsed)

	start = time.Now()
	ssg.createNodes()
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Create Nodes: %s\n", elapsed)

	start = time.Now()
	ssg.createEdges()
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Create Edges: %s\n", elapsed)

	return &ssg
}

func (ssg *SimpleSphereGrid) distributePoints() {
	lat := LatMin
	lon := LonMin

	dLat := (LatMax - LatMin) / (float64(ssg.nLat) - 1)
	dLon := (LonMax - LonMin) / float64(ssg.nLon)

	ssg.points = make([]geo.Point, 0)
	for iLat := 0; iLat < ssg.nLat; iLat++ {
		for iLon := 0; iLon < ssg.nLon; iLon++ {
			ssg.points = append(ssg.points, geo.Point{lat, lon})
			lon += dLon
		}
		lon = LonMin
		lat += dLat
	}
}

func (ssg *SimpleSphereGrid) landWaterTest(polygons []geo.Polygon) {
	numPoints := len(ssg.points)
	ssg.isWater = make([]bool, numPoints, numPoints)

	// pre-compute bounding boxes for every polygon
	bboxes := make([]geo.BoundingBox, len(polygons), len(polygons))
	var wg sync.WaitGroup
	wg.Add(len(polygons))
	for i, polygon := range polygons {
		go func(i int, polygon geo.Polygon) {
			bboxes[i] = polygon.BoundingBox()
			wg.Done()
		}(i, polygon)
	}
	wg.Wait()

	wg.Add(numPoints)
	for idx, point := range ssg.points {
		go func(idx int, point geo.Point) {
			if point.Lat() < -84 {
				// hard-coded: make south pole continent
				ssg.isWater[idx] = false
			} else {
				// no special treatment for non south pole points
				ssg.isWater[idx] = true
				for i, polygon := range polygons {
					// roughly check, whether the point is contained in the bounding box of the polygon
					if bboxes[i].Contains(point) {
						// precisely check, whether the polygon contains the point
						contains := polygon.Contains(&point)
						if contains {
							ssg.isWater[idx] = false
							break
						}
					}
				}
			}
			wg.Done()
		}(idx, point)
	}
	wg.Wait()
}

func (ssg *SimpleSphereGrid) createNodes() {
	ssg.grid2nodes = make(map[int]int)
	ssg.nodes2grid = make([]int, 0)
	ssg.nodes = make([]gr.Node, 0)
	for cellId, point := range ssg.points {
		if ssg.isWater[cellId] {
			ssg.grid2nodes[cellId] = len(ssg.nodes)
			ssg.nodes = append(ssg.nodes, *gr.NewNode(point.Lon(), point.Lat()))
			ssg.nodes2grid = append(ssg.nodes2grid, cellId)
		}
	}
}

func (ssg *SimpleSphereGrid) createEdges() {
	for nodeId := range ssg.nodes {
		cellId := ssg.nodes2grid[nodeId]
		neighborCellIds := ssg.neighborsOf(cellId)
		for _, neighborCellId := range neighborCellIds {
			if neighborNodeId, ok := ssg.grid2nodes[neighborCellId]; ok {
				p1 := geo.NewPoint(ssg.nodes[nodeId].Lat, ssg.nodes[nodeId].Lon)
				p2 := geo.NewPoint(ssg.nodes[neighborNodeId].Lat, ssg.nodes[neighborNodeId].Lon)
				distance := p1.IntHaversine(p2)
				edge := gr.Edge{From: nodeId, To: neighborNodeId, Distance: distance}
				ssg.edges = append(ssg.edges, edge)
			}
		}
	}
}

func (ssg *SimpleSphereGrid) neighborsOf(cellId int) []int {
	neighbors := make([]int, 0)
	if cellId < ssg.nLon*(ssg.nLat-1) {
		// northern neighbor
		neighbors = append(neighbors, cellId+ssg.nLon)
	}
	if cellId >= ssg.nLon {
		// southern neighbor
		neighbors = append(neighbors, cellId-ssg.nLon)
	}

	// western neighbor
	if cellId%ssg.nLon != 0 {
		neighbors = append(neighbors, cellId-1)
	} else {
		neighbors = append(neighbors, cellId+ssg.nLon-1)
	}

	// eastern neighbor
	if (cellId+1)%ssg.nLon != 0 {
		neighbors = append(neighbors, cellId+1)
	} else {
		neighbors = append(neighbors, cellId-ssg.nLon+1)
	}
	return neighbors
}

func (ssg *SimpleSphereGrid) ToGraph() gr.Graph {
	alg := &gr.AdjacencyListGraph{}
	for _, node := range ssg.nodes {
		alg.AddNode(node)
	}
	for _, edge := range ssg.edges {
		alg.AddEdge(edge)
	}
	return alg
}
