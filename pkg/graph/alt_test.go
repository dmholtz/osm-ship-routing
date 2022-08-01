package graph

import "testing"

func TestAlt(t *testing.T) {
	g := NewAdjacencyArrayFromFmi("out.fmi")
	landmarks := []int{9, 1000, 2000, 3000, 4000, 5000, 6000, 7000}

	landmarkDistancesCollection := AltPreprocessing(g, g, landmarks)

	n := 0
	totalDijkstraPqPops := 0
	totalAltPqPops := 0
	for orig := 0; orig < g.NodeCount(); orig = orig + 50 {
		for dest := 0; dest < g.NodeCount(); dest = dest + 50 {
			if orig == dest {
				continue
			}
			_, length1, dPqPops := Dijkstra(g, orig, dest)
			_, length2, altPqPops := Alt(g, landmarkDistancesCollection, orig, dest)
			if length1 != length2 {
				t.Errorf("[Path(from=%d, dest=%d)]: Different lengths found: %d≠%d", orig, dest, length1, length2)
				return
			}
			totalDijkstraPqPops += dPqPops
			totalAltPqPops += altPqPops
			n++
		}
	}
	t.Logf("Avg number of PQ-pops (Dijkstra): %d", totalDijkstraPqPops/n)
	t.Logf("Avg number of PQ-pops (ALT): %d", totalAltPqPops/n)

}
