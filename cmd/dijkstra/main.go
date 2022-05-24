package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dmholtz/osm-ship-routing/pkg/graph"
)

const graphFile = "graphs/ocean_equi_4.fmi"

func main() {
	start := time.Now()
	aag := graph.NewAdjacencyArrayFromFmi(graphFile)
	elapsed := time.Since(start)
	fmt.Printf("[TIME-Import] = %s\n", elapsed)

	benchmark(aag, 100)
}

// Run benchmarks on the provided graphs: Compute n random routes
func benchmark(aag *graph.AdjacencyArrayGraph, n int) {

	runtime := 0
	for i := 0; i < n; i++ {
		origin := rand.Intn(aag.NodeCount())
		destination := rand.Intn(aag.NodeCount())

		start := time.Now()
		path, length := graph.Dijkstra(aag, origin, destination)
		elapsed := time.Since(start)
		fmt.Printf("[TIME-Navigate] = %s\n", elapsed)

		if length > -1 {
			if path[0] != origin || path[len(path)-1] != destination {
				panic("Invalid routing result")
			}
		}

		runtime += int(elapsed)
	}
	fmt.Printf("Average runtime: %.3fms", float64(runtime/n)/1000000)
}
