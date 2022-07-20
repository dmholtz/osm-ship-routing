package graph

import "testing"

func TestShortestPathTree(t *testing.T) {
	t.Parallel()

	g := NewAdjacencyArrayFromFmi(arcFlagGraphFile)
	tree := ShortestPathTree(g, 0)

	for len(tree.children) > 0 {
		t.Log(tree.id)
		tree = *tree.children[0]
	}
	// TODO improve test assertions
}
