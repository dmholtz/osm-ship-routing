package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/dmholtz/osm-ship-routing/pkg/graph"
)

const flagggedGraphFile = "graphs/ocean_equi_4_grid_arcflags.fmi"
const n = 5000

func main() {
	rand.Seed(2022)
	fmt.Printf("Loading %s ... ", flagggedGraphFile)
	faag := graph.NewFlaggedAdjacencyArrayFromFmi(flagggedGraphFile)
	fmt.Println("Done")

	fmt.Print("Run random searches ... ")
	traces := make([][]int, 0)
	totalTraceLength := 0
	for i := 0; i < n; i++ {
		origin := rand.Intn(faag.NodeCount())
		destination := rand.Intn(faag.NodeCount())
		trace, _, _ := graph.ArcFlagBiDijkstra(faag, faag, origin, destination)
		traces = append(traces, trace)
		totalTraceLength += len(trace)
	}
	fmt.Println("Done")

	fmt.Printf("Average trace length: %d\n", totalTraceLength/n)

	jsonObj, err := json.Marshal(traces)
	if err != nil {
		panic(err)
	}

	wErr := os.WriteFile("traces.json", jsonObj, 0644)
	if wErr != nil {
		panic(err)
	}
}
