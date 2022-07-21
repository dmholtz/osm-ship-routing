package graph

import "math"

func GridPartitioning(g Graph) *FlaggedAdjacencyArrayGraph {
	// copy the graph
	falg := FlaggedAdjacencyListGraph{}
	for i := 0; i < g.NodeCount(); i++ {
		falg.AddNode(g.GetNode(i))
	}
	for i := 0; i < g.NodeCount(); i++ {
		for _, halfEdge := range g.GetHalfEdgesFrom(i) {
			flaggedHalfEdge := FlaggedHalfEdge{To: halfEdge.To, Weight: halfEdge.Distance, Flag: 0}
			falg.AddHalfEdge(i, flaggedHalfEdge)
		}
	}
	falg.Partitions = make([]uint8, g.NodeCount(), g.NodeCount())

	// grid partitioning
	for i := 0; i < falg.NodeCount(); i++ {
		node := falg.GetNode(i)
		col := math.Min(((node.Lon + 180) / 360 * 8), 7)
		row := math.Min(((node.Lat + 90) / 180 * 8), 7)
		p := uint8(row)*8 + uint8(col)
		falg.Partitions[i] = p
	}

	return NewFlaggedAdjacencyArrayFromGraph(&falg)

}
