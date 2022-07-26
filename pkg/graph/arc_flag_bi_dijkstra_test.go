package graph

import "testing"

//const arcFlagGraphFile = "testdata/arc_flag_graph.fmi"

func TestArcFlagBiDijkstra(t *testing.T) {
	t.Parallel()

	g := NewAdjacencyArrayFromFmi("out2.fmi")
	fg := NewFlaggedAdjacencyArrayFromFmi("out2.fmi")
	fgt := NewFlaggedAdjacencyArrayFromFmi("out2t.fmi")

	n := 0
	totalDijkstraPqPops := 0
	totalArcFlagPqPops := 0
	for orig := 0; orig < g.NodeCount(); orig = orig + 50 {
		for dest := 0; dest < g.NodeCount(); dest = dest + 40 {
			if orig == dest {
				continue
			}
			_, length1, dPqPops := Dijkstra(g, orig, dest)
			_, length2, afPqPops := ArcFlagBiDijkstra(fg, fgt, orig, dest)
			if length1 != length2 {
				t.Errorf("[Path(from=%d, dest=%d)]: Different lengths found: %d≠%d", orig, dest, length1, length2)
				return
			}
			totalDijkstraPqPops += dPqPops
			totalArcFlagPqPops += afPqPops
			n++
		}
	}
	t.Logf("Avg number of PQ-pops (Dijkstra): %d", totalDijkstraPqPops/n)
	t.Logf("Avg number of PQ-pops (Arc-Flags): %d", totalArcFlagPqPops/n)
}