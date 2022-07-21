package graph

import "fmt"

// Input: partitioned graph (max 64 partitions) with zero arc-flags
// Output: partitioned graph with non-trivial arc-flags

// TODO re-think parameter type
func ComputeArcFlags(fg *FlaggedAdjacencyArrayGraph) FlaggedGraph {

	// preprocessing: remove existing arc flag
	for _, halfEdge := range fg.Edges {
		halfEdge.ResetFlag()
	}

	// determine boundary nodes for each region
	boundaryNodeSets := make([](map[NodeId]bool), 0)
	for i := 0; i < 64; i++ {
		boundaryNodeSets = append(boundaryNodeSets, make(map[NodeId]bool, 0))
	}
	for tailNodeId := range fg.Nodes {
		for _, halfEdge := range fg.GetHalfEdgesFrom(tailNodeId) {
			tailPartition := fg.GetPartition((tailNodeId))
			headPartition := fg.GetPartition(halfEdge.To)
			if tailPartition != headPartition {
				boundaryNodeSets[headPartition][halfEdge.To] = true
			}
		}
	}

	//for partition, set := range boundaryNodeSets {
	//	fmt.Printf("%d: ", partition)
	//	for k := range set {
	//		fmt.Printf("%d, ", k)
	//	}
	//	fmt.Printf("\n")
	//}

	// compute transposed graph
	transposedGraph := TransposeGraph(fg)

	for partition, set := range boundaryNodeSets {
		fmt.Printf("Partition %d, size=%d\n", partition, len(set))
		for boundaryNodeId := range set {
			//fmt.Printf("Check boundary node %d in region %d\n", boundaryNodeId, partition)
			// calculate in reverse graph
			tree := ShortestPathTree(transposedGraph, boundaryNodeId)

			stack := make([]*TreeNode, 0)
			for _, child := range tree.children {
				if fg.GetPartition(child.id) != uint8(partition) {
					// add edge
					tailRev := tree.id
					headRev := child.id
					AddFlag(fg, headRev, tailRev, uint8(partition))
					stack = append(stack, child)
				}
			}

			for len(stack) > 0 {
				// pop
				node := stack[len(stack)-1]
				stack = stack[0 : len(stack)-1]

				for _, child := range node.children {
					if fg.GetPartition(child.id) != uint8(partition) || true {
						// add edge
						tailRev := node.id
						headRev := child.id
						AddFlag(fg, headRev, tailRev, uint8(partition))
						stack = append(stack, child)
					}
				}
			}
		}
	}
	// revise edges within the same partition
	for i := 0; i < fg.NodeCount(); i++ {
		for _, halfEdge := range fg.GetHalfEdgesFrom(i) {
			if fg.GetPartition(i) == fg.GetPartition(halfEdge.To) {
				AddFlag(fg, i, halfEdge.To, fg.GetPartition(i))
			}
		}
	}

	return fg
}

func AddFlag(fg *FlaggedAdjacencyArrayGraph, from NodeId, to NodeId, partition PartitionId) {
	for i := fg.Offsets[from]; i < fg.Offsets[from+1]; i++ {
		edge := &fg.Edges[i]
		if edge.To == to {
			edge.AddFlag(partition)
			break
		}
	}
}
