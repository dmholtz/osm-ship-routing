package graph

import "testing"

func TestShortestPathTree(t *testing.T) {
	t.Parallel()

	fg := NewFlaggedAdjacencyArrayFromFmi(arcFlagGraphFile)
	tree := ShortestPathTree(fg, 0)

	for len(tree.children) > 0 {
		t.Log(tree.id)
		tree = *tree.children[0]
	}
	// TODO improve test assertions
}

func TestMultipleShortestPathTree(t *testing.T) {
	t.Parallel()

	fg := NewFlaggedAdjacencyArrayFromFmi("two.fmi")
	tree := ShortestPathTree(fg, 0)

	for _, child1 := range tree.children {
		t.Logf("Level 1: %d", child1.id)
		if len(child1.children) < 1 {
			t.Error("No path to level 2")
		}
		for _, child2 := range child1.children {
			t.Logf("Level 2: %d", child2.id)
			if child2.id != 3 {
				t.Errorf("No path to node 3.")
			}
		}
	}
}
