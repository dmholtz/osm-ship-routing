package graph

import "testing"

const arcFlagGraphFile = "testdata/arc_flag_graph.fmi"

func TestArcFlagDijkstra(t *testing.T) {
	t.Parallel()

	g := NewAdjacencyArrayFromFmi(arcFlagGraphFile)
	fg := NewFlaggedAdjacencyArrayFromFmi(arcFlagGraphFile)

	for orig := 0; orig < g.NodeCount(); orig++ {
		for dest := 0; dest < g.NodeCount(); dest++ {
			_, length1, _ := Dijkstra(g, orig, dest)
			_, length2, _ := ArcFlagDijkstra(fg, orig, dest)
			if length1 != length2 {
				t.Errorf("[Path(from=%d, dest=%d)]: Different lengths found: %d≠%d", orig, dest, length1, length2)
				return
			}
		}
	}
}

func TestArcFlagDijkstra2(t *testing.T) {
	t.Parallel()

	g := NewAdjacencyArrayFromFmi("out2.fmi")
	fg := NewFlaggedAdjacencyArrayFromFmi("out2.fmi")

	n := 0
	totalDijkstraPqPops := 0
	totalArcFlagPqPops := 0
	for orig := 0; orig < g.NodeCount(); orig = orig + 50 {
		for dest := 0; dest < g.NodeCount(); dest = dest + 50 {
			_, length1, dPqPops := Dijkstra(g, orig, dest)
			_, length2, afPqPops := ArcFlagDijkstra(fg, orig, dest)
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

func TestArcFlagDijkstra3(t *testing.T) {
	t.Parallel()

	g := NewAdjacencyArrayFromFmi("out0k2.fmi")
	fg := NewFlaggedAdjacencyArrayFromFmi("out0k2.fmi")

	for orig := 0; orig < g.NodeCount(); orig++ {
		for dest := 0; dest < g.NodeCount(); dest++ {
			_, length1, _ := Dijkstra(g, orig, dest)
			_, length2, _ := ArcFlagDijkstra(fg, orig, dest)
			if length1 != length2 {
				t.Errorf("[Path(from=%d, dest=%d)]: Different lengths found: %d≠%d", orig, dest, length1, length2)
				return
			}
		}
	}
}

func TestArcFlagDijkstra4(t *testing.T) {
	t.Parallel()

	g := NewAdjacencyArrayFromFmi(arcFlagGraphFile)
	fg := NewFlaggedAdjacencyArrayFromFmi(arcFlagGraphFile)

	for orig := 0; orig < g.NodeCount(); orig++ {
		for dest := 0; dest < g.NodeCount(); dest++ {
			_, length1, _ := Dijkstra(g, orig, dest)
			_, length2, _ := ArcFlagDijkstra(fg, orig, dest)
			if length1 != length2 {
				t.Errorf("[Path(from=%d, dest=%d)]: Different lengths found: %d≠%d", orig, dest, length1, length2)
				return
			}
		}
	}
}
