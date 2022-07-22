package graph

func TransposeGraph(fg FlaggedGraph) FlaggedGraph {
	transpose := FlaggedAdjacencyListGraph{}

	for i := 0; i < fg.NodeCount(); i++ {
		transpose.AddNode(fg.GetNode(i))
	}

	for i := 0; i < fg.NodeCount(); i++ {
		for _, halfEdge := range fg.GetHalfEdgesFrom(i) {
			reversedEdge := FlaggedHalfEdge{To: i, Weight: halfEdge.Weight, Flag: 0}
			transpose.AddHalfEdge(halfEdge.To, reversedEdge)
		}
	}

	return &transpose
}
