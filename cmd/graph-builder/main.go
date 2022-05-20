package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dmholtz/osm-ship-routing/pkg/geometry"
	"github.com/dmholtz/osm-ship-routing/pkg/graph"
	"github.com/dmholtz/osm-ship-routing/pkg/grid"
)

const density = 20  // parameter for SimpleSphereGrid
const nTarget = 1e4 // parameter for EquiSphereGrid

func main() {

	//arg := loadPolyJsonPolygons("antarctica.poly.json")
	arg := loadPolyJsonPolygons("planet-coastlines.poly.json")

	//grid := grid.NewSimpleSphereGrid(2*density, density, arg)
	grid := grid.NewEquiSphereGrid(nTarget, arg)

	gridGraph := grid.ToGraph()
	jsonObj, err := json.Marshal(gridGraph)
	if err != nil {
		panic(err)
	}

	wErr := os.WriteFile("graph.json", jsonObj, 0644)
	if wErr != nil {
		panic(err)
	}

	graph.WriteFmi(gridGraph, "graph.fmi")
}

func loadPolyJsonPolygons(file string) []geometry.Polygon {

	start := time.Now()
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("[TIME] Read file: %s\n", elapsed)

	start = time.Now()
	var polygons []geometry.Polygon
	err = json.Unmarshal(bytes, &polygons)
	if err != nil {
		fmt.Println(err)
	}
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Unmarshal: %s\n", elapsed)

	return polygons
}
