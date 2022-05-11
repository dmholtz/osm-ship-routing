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
	sgg := SimpleSphereGrid{nLon: nLon, nLat: nLat}

	start := time.Now()
	sgg.distributePoints()
	elapsed := time.Since(start)
	fmt.Printf("[TIME] Distribute Points on grid: %s\n", elapsed)

	start = time.Now()
	sgg.landWaterTest(coastlines)
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Land / Water test: %s\n", elapsed)

	start = time.Now()
	sgg.createNodes()
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Create Nodes: %s\n", elapsed)

	start = time.Now()
	sgg.createEdges()
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Create Edges: %s\n", elapsed)

	return &sgg
}

func (sgg *SimpleSphereGrid) distributePoints() {
	lat := LatMin
	lon := LonMin

	dLat := (LatMax - LatMin) / (float64(sgg.nLat) - 1)
	dLon := (LonMax - LonMin) / float64(sgg.nLon)

	sgg.points = make([]geo.Point, 0)
	for iLat := 0; iLat < sgg.nLat; iLat++ {
		for iLon := 0; iLon < sgg.nLon; iLon++ {
			sgg.points = append(sgg.points, geo.Point{lat, lon})
			lon += dLon
		}
		lon = LonMin
		lat += dLat
	}
}

func (sgg *SimpleSphereGrid) landWaterTest(polygons []geo.Polygon) {
	numPoints := len(sgg.points)
	sgg.isWater = make([]bool, numPoints, numPoints)

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
	for idx, point := range sgg.points {
		go func(idx int, point geo.Point) {
			sgg.isWater[idx] = true
			for i, polygon := range polygons {
				// roughly check, whether the point is contained in the bounding box of the polygon
				if bboxes[i].Contains(point) {
					// precisely check, whether the polygon contains the point
					if polygon.Contains(&point) {
						sgg.isWater[idx] = false
						break
					}
				}
			}
			wg.Done()
		}(idx, point)
	}
	wg.Wait()
}

func (sgg *SimpleSphereGrid) createNodes() {
	sgg.grid2nodes = make(map[int]int)
	sgg.nodes2grid = make([]int, 0)
	sgg.nodes = make([]gr.Node, 0)
	for cellId, point := range sgg.points {
		if sgg.isWater[cellId] {
			sgg.grid2nodes[cellId] = len(sgg.nodes)
			sgg.nodes = append(sgg.nodes, *gr.NewNode(point.Lon(), point.Lat()))
			sgg.nodes2grid = append(sgg.nodes2grid, cellId)
		}
	}
}

func (sgg *SimpleSphereGrid) createEdges() {
	for nodeId, _ := range sgg.nodes {
		cellId := sgg.nodes2grid[nodeId]
		neighborCellIds := sgg.neighborsOf(cellId)
		for _, neighborCellId := range neighborCellIds {
			if neighborNodeId, ok := sgg.grid2nodes[neighborCellId]; ok {
				edge := gr.Edge{From: nodeId, To: neighborNodeId, Distance: 1} // todo: compute distance
				sgg.edges = append(sgg.edges, edge)
			}
		}
	}
}

func (sgg *SimpleSphereGrid) neighborsOf(cellId int) []int {
	neighbors := make([]int, 0)
	if cellId < sgg.nLon*(sgg.nLat-1) {
		// northern neighbor
		neighbors = append(neighbors, cellId+sgg.nLon)
	}
	if cellId >= sgg.nLon {
		// southern neighbor
		neighbors = append(neighbors, cellId-sgg.nLon)
	}

	// western neighbor
	if cellId%sgg.nLon != 0 {
		neighbors = append(neighbors, cellId-1)
	} else {
		neighbors = append(neighbors, cellId+sgg.nLon-1)
	}

	// eastern neighbor
	if (cellId+1)%sgg.nLon != 0 {
		neighbors = append(neighbors, cellId+1)
	} else {
		neighbors = append(neighbors, cellId-sgg.nLon+1)
	}
	return neighbors
}

func (sgg *SimpleSphereGrid) ToGraph() gr.Graph {
	alg := &gr.AdjacencyListGraph{}
	for _, node := range sgg.nodes {
		alg.AddNode(node)
	}
	for _, edge := range sgg.edges {
		alg.AddEdge(edge)
	}
	return alg
}
