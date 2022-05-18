package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dmholtz/osm-ship-routing/pkg/graph"
)

func main() {
	//aag := graph.NewAdjacencyArrayFromFmi("simpleGraph.fmi")
	start := time.Now()
	aag := graph.NewAdjacencyArrayFromFmi("ocean.fmi")
	//aag := graph.NewAdjacencyArrayFromFmi("graph.fmi")
	elapsed := time.Since(start)
	fmt.Printf("[TIME-Import] = %s\n", elapsed)

	benchmark(aag, 100)
}

func benchmark(aag *graph.AdjacencyArrayGraph, n int) {

	runtime := 0
	for i := 0; i < n; i++ {
		origin := rand.Intn(aag.NodeCount())
		destination := rand.Intn(aag.NodeCount())

		start := time.Now()
		graph.Dijkstra(aag, origin, destination)
		elapsed := time.Since(start)
		fmt.Printf("[TIME-Navigate] = %s\n", elapsed)

		runtime += int(elapsed)
	}
	fmt.Printf("Average runtime: %.3fms", float64(runtime/n)/1000000)
}
