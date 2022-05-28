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

const density = 710 // parameter for SimpleSphereGrid
const nTarget = 1e6 // parameter for EquiSphereGrid

func main() {

	//arg := loadPolyJsonPolygons("antarctica.poly.json")
	arg := loadPolyJsonPolygons("planet-coastlines.poly.json")

	/*
		var max_distances []float64
		for _, p := range arg {
			max_distance := 0.0
			for i := 1; i < p.Size(); i++ {
				p1 := p.At(i - 1)
				p2 := p.At(i)
				distance := p1.Haversine(p2)
				if distance > max_distance {
					max_distance = distance
				}
			}
			max_distances = append(max_distances, max_distance)
		}
		fmt.Printf("Distance: %v\n", max_distances)
		max := 0.0
		for i, v := range max_distances {
			if v > max && v < 56000 {
				max = v
			}
			if v > 56000 {
				fmt.Printf("%v\n", arg[i].At(0))
			}
		}
		fmt.Printf("Max: %v\n", max)
		fmt.Printf("Count: %v\n", len(arg))
		os.Exit(0)
	*/

	//grid := grid.NewSimpleSphereGrid(2*density, density, arg)
	grid := grid.NewEquiSphereGrid(nTarget, grid.SIX_NEIGHBORS, arg)

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
