package graph

import (
	"math/rand"
	"testing"
)

const graphFile = "../../graphs/ocean_equi_4.fmi"
const numberOfRuns = 200

// Compare bidirectional dijkstra with textbook Dijkstra's Algorithm
func TestBidirectionalDijkstra(t *testing.T) {
	t.Logf("Loading graph from file %s.", graphFile)
	aag := NewAdjacencyArrayFromFmi(graphFile)
	t.Log("Loading done.")

	t.Logf("Compare both algorithms' results on %d random routes.", numberOfRuns)
	for i := 0; i < numberOfRuns; i++ {
		origin := rand.Intn(aag.NodeCount())
		destination := rand.Intn(aag.NodeCount())

		path1, length1, _ := Dijkstra(aag, origin, destination)
		path2, length2, _ := BidirectionalDijkstra(aag, origin, destination)

		if length1 != length2 {
			t.Errorf("Incorrect result: bidirectional dijkstra computes length %d, dijkstra computes length %d.", length1, length2)
			return
		}
		if len(path1) != len(path2) {
			t.Errorf("Incorrect result: bidirectional dijkstra computes path with %d nodes, dijkstra computes path with %d nodes.", len(path1), len(path2))
			return
		}

		// comparing the paths is not reasonable, because there might be multiple optimal routes.
	}
}
