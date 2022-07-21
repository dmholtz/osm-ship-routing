package main

import (
	"fmt"
	"time"

	"github.com/dmholtz/osm-ship-routing/pkg/graph"
)

func main() {
	g := graph.NewAdjacencyArrayFromFmi("graphs/ocean_equi_4.fmi")

	start := time.Now()
	fg := graph.GridPartitioning(g)
	elapsed := time.Since(start)
	fmt.Printf("[TIME-Partitioning] = %s\n", elapsed)

	start = time.Now()
	fg1 := graph.ComputeArcFlags(fg)
	elapsed = time.Since(start)
	fmt.Printf("[TIME-ArcFlagComputation] = %s\n", elapsed)

	start = time.Now()
	graph.WritePartitionedFmi(fg1, "graphs/ocean_equi_4_arcflags.fmi")
	elapsed = time.Since(start)
	fmt.Printf("[TIME-WriteFMI] = %s\n", elapsed)
}
