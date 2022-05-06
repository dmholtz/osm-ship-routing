package graph

// Implementation for dynamic graphs
type AdjacencyListGraph struct {
	nodes     []Node
	edges     []outgoingEdges
	edgeCount int
}

func (alg *AdjacencyListGraph) GetNode(id NodeId) Node {
	if id < 0 || id >= alg.NodeCount() {
		panic(id)
	}
	return alg.nodes[id]
}

func (alg *AdjacencyListGraph) GetEdgesFrom(id NodeId) []Edge {
	if id < 0 || id >= alg.NodeCount() {
		panic(id)
	}
	outgoingEdges := make([]Edge, 0)
	for _, outgoingEdge := range alg.edges[id] {
		outgoingEdges = append(outgoingEdges, outgoingEdge.toEdge(id))
	}
	return outgoingEdges
}

func (alg *AdjacencyListGraph) NodeCount() int {
	return len(alg.nodes)
}

func (alg *AdjacencyListGraph) EdgeCount() int {
	return alg.edgeCount
}

func (alg *AdjacencyListGraph) AddNode(n Node) {
	alg.nodes = append(alg.nodes, n)
	alg.edges = append(alg.edges, make(outgoingEdges, 0))
}

func (alg *AdjacencyListGraph) AddEdge(e Edge) {
	// Check if both source and target node exit
	if e.From >= alg.NodeCount() || e.To >= alg.NodeCount() {
		panic(e)
	}
	// Check for duplicates
	for _, outgoingEdge := range alg.edges[e.From] {
		if e.To == outgoingEdge.To {
			panic(e)
		}
	}
	alg.edges[e.From] = append(alg.edges[e.From], e.toOutgoingEdge())
	alg.edgeCount++
}
