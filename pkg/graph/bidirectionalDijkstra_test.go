package graph

import (
	"math/rand"
	"testing"
)

const graphFile = "../../graphs/ocean_equi_4.fmi"
const numberOfRuns = 20

// Compare bidirectional dijkstra with textbook Dijkstra's Algorithm
func TestBidirectionalDijkstra(t *testing.T) {
	t.Logf("Loading graph from file %s.", graphFile)
	aag := NewAdjacencyArrayFromFmi(graphFile)
	t.Log("Loading done.")

	t.Logf("Compare both algorithms' results on %d random routes.", numberOfRuns)
	for i := 0; i < numberOfRuns; i++ {
		origin := rand.Intn(aag.NodeCount())
		destination := rand.Intn(aag.NodeCount())

		_, length1 := Dijkstra(aag, origin, destination)
		_, length2 := BidirectionalDijkstra(aag, origin, destination)

		if length1 != length2 {
			t.Errorf("Incorrect result: bidirectional dijkstra computes length %d, dijkstra computes length %d.", length1, length2)
		}

		// comparing the paths is not reasonable, because there might be multiple optimal routes.
	}
}
