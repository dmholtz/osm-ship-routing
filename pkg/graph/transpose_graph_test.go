package graph

import (
	"testing"
)

func TestTransposeGraph(t *testing.T) {
	t.Parallel()

	g0 := NewFlaggedAdjacencyListFromFmi(arcFlagGraphFile)
	g1 := TransposeGraph(g0)
	g2 := TransposeGraph(g1)

	// check node/edge count
	if g0.NodeCount() != g1.NodeCount() {
		t.Errorf("Different number of nodes: g0=%d, g1=%d", g0.NodeCount(), g1.NodeCount())
	}
	if g0.EdgeCount() != g1.EdgeCount() {
		t.Errorf("Different number of edges: g0=%d, g1=%d", g0.EdgeCount(), g1.EdgeCount())
	}

	// check if Transpose(Transpose(*)) == id(*)
	for i := 0; i < g0.NodeCount(); i++ {
		// compare  nodes
		if g0.GetNode(i) != g2.GetNode(i) {
			t.Errorf("Different nodes with id=%d: g0(id)=%v, g2(id)=%v", i, g0.GetNode(i), g2.GetNode(i))
		}
		// compare edges
		g0Edges := g0.GetHalfEdgesFrom(i)
		g2Edges := g2.GetHalfEdgesFrom(i)
		if len(g0Edges) != len(g2Edges) {
			t.Errorf("Different number of half edges at node %d: %d≠%d", i, len(g0Edges), len(g2Edges))
		}
		for _, g0Edge := range g0Edges {
			found := false
			for _, g2Edge := range g2Edges {
				if g0Edge.To == g2Edge.To && g0Edge.Weight == g2Edge.Weight {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("No corresponding half edge for %v in g2", g0Edge)
			}
		}
	}
}
