package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dmholtz/osm-ship-routing/pkg/graph"
)

const graphFile = "graphs/ocean_equi_4.fmi"
const flagggedGraphFile = "graphs/ocean_equi_4_arcflags.fmi"
const n = 100

type stats struct {
	runtime uint64
	pqPops  uint64
}

func (s1 *stats) add(s2 stats) {
	s1.runtime += s2.runtime
	s1.pqPops += s2.pqPops
}

func (s stats) average(n int) stats {
	m := uint64(n)
	return stats{runtime: s.runtime / m, pqPops: s.pqPops / m}
}

func (s stats) print() {
	fmt.Printf("Average runtime: %.3fms\n", float64(s.runtime)/1000000)
	fmt.Printf("Average number of PQ-pops: %d\n", s.pqPops)
}

func main() {
	fmt.Printf("Loading %s ... ", graphFile)
	aag := graph.NewAdjacencyArrayFromFmi(graphFile)
	fmt.Println("Done")
	fmt.Printf("Loading %s ... ", flagggedGraphFile)
	faag := graph.NewFlaggedAdjacencyArrayFromFmi(flagggedGraphFile)
	fmt.Println("Done")

	dijkstraStats := stats{}
	arcFlagDijkstraStats := stats{}
	arcFlagBiDijkstraStats := stats{}
	for i := 0; i < n; i++ {
		origin := rand.Intn(aag.NodeCount())
		destination := rand.Intn(aag.NodeCount())

		dijkstraStats.add(benchmarkDijkstra(aag, destination, origin))
		arcFlagDijkstraStats.add(benchmarkArcFlagDijkstra(faag, origin, destination))
		arcFlagBiDijkstraStats.add(benchmarkArcFlagBiDijkstra(faag, origin, destination))
	}

	dijkstraStats = dijkstraStats.average(n)
	arcFlagDijkstraStats = arcFlagDijkstraStats.average(n)
	arcFlagBiDijkstraStats = arcFlagBiDijkstraStats.average(n)

	fmt.Println("Benchmark standard dijkstra:")
	dijkstraStats.print()
	fmt.Println("Benchmark arc-flag dijkstra:")
	arcFlagDijkstraStats.print()
	fmt.Println("Benchmark arc-flag bi-dijkstra:")
	arcFlagBiDijkstraStats.print()
}

func benchmarkDijkstra(aag *graph.AdjacencyArrayGraph, origin int, destination int) stats {
	start := time.Now()
	_, _, pqPops := graph.Dijkstra(aag, origin, destination)
	elapsed := time.Since(start)
	return stats{runtime: uint64(elapsed), pqPops: uint64(pqPops)}
}

func benchmarkArcFlagDijkstra(faag *graph.FlaggedAdjacencyArrayGraph, origin int, destination int) stats {
	start := time.Now()
	_, _, pqPops := graph.ArcFlagDijkstra(faag, origin, destination)
	elapsed := time.Since(start)
	return stats{runtime: uint64(elapsed), pqPops: uint64(pqPops)}
}

func benchmarkArcFlagBiDijkstra(faag *graph.FlaggedAdjacencyArrayGraph, origin int, destination int) stats {
	start := time.Now()
	_, _, pqPops := graph.ArcFlagBiDijkstra(faag, origin, destination)
	elapsed := time.Since(start)
	return stats{runtime: uint64(elapsed), pqPops: uint64(pqPops)}
}
