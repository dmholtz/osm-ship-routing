package graph

import "testing"

const arcFlagGraphFile = "testdata/arc_flag_graph.fmi"

func TestArcFlagDij(t *testing.T) {
	t.Parallel()

	g := NewAdjacencyArrayFromFmi(arcFlagGraphFile)
	fg := NewFlaggedAdjacencyArrayFromFmi(arcFlagGraphFile)

	for orig := 0; orig < g.NodeCount(); orig++ {
		for dest := 0; dest < g.NodeCount(); dest++ {
			_, length1 := Dijkstra(g, orig, dest)
			_, length2 := ArcFlagDijkstra(fg, orig, dest)
			if length1 != length2 {
				t.Errorf("[Path(from=%d, dest=%d)]: Different lengths found: %d≠%d", orig, dest, length1, length2)
				return
			}
		}
	}
}
