package graph

import "fmt"

// Implementation for dynamic graphs
type AdjacencyListGraph struct {
	Nodes     []Node
	Edges     [][]HalfEdge
	edgeCount int
}

func (alg *AdjacencyListGraph) GetNode(id NodeId) Node {
	if id < 0 || id >= alg.NodeCount() {
		panic(id)
	}
	return alg.Nodes[id]
}

func (alg *AdjacencyListGraph) GetEdgesFrom(id NodeId) []Edge {
	if id < 0 || id >= alg.NodeCount() {
		panic(id)
	}
	outgoingEdges := make([]Edge, 0)
	for _, outgoingEdge := range alg.Edges[id] {
		outgoingEdges = append(outgoingEdges, outgoingEdge.toEdge(id))
	}
	return outgoingEdges
}

func (alg *AdjacencyListGraph) GetHalfEdgesFrom(id NodeId) HalfEdges {
	if id < 0 || id >= alg.NodeCount() {
		panic(id)
	}
	return alg.Edges[id]
}

func (alg *AdjacencyListGraph) NodeCount() int {
	return len(alg.Nodes)
}

func (alg *AdjacencyListGraph) EdgeCount() int {
	return alg.edgeCount
}

func (alg *AdjacencyListGraph) AddNode(n Node) {
	alg.Nodes = append(alg.Nodes, n)
	alg.Edges = append(alg.Edges, make([]HalfEdge, 0))
}

func (alg *AdjacencyListGraph) AddEdge(e Edge) {
	// Check if both source and target node exit
	if e.From >= alg.NodeCount() || e.To >= alg.NodeCount() {
		panic(fmt.Sprintf("Edge out of range %v", e))
	}
	// Check for duplicates
	for _, outgoingEdge := range alg.Edges[e.From] {
		if e.To == outgoingEdge.To {
			return // ignore duplicate edges
		}
	}
	alg.Edges[e.From] = append(alg.Edges[e.From], e.toOutgoingEdge())
	alg.edgeCount++
}

// Implementation for dynamic flagged graphs
type FlaggedAdjacencyListGraph struct {
	Nodes     []Node
	Edges     [][]FlaggedHalfEdge
	edgeCount int

	Partitions []PartitionId
}

func (falg *FlaggedAdjacencyListGraph) GetNode(id NodeId) Node {
	if id < 0 || id >= falg.NodeCount() {
		panic(id)
	}
	return falg.Nodes[id]
}

func (falg *FlaggedAdjacencyListGraph) GetHalfEdgesFrom(id NodeId) []FlaggedHalfEdge {
	if id < 0 || id >= falg.NodeCount() {
		panic(id)
	}
	return falg.Edges[id]
}

func (falg *FlaggedAdjacencyListGraph) NodeCount() int {
	return len(falg.Nodes)
}

func (falg *FlaggedAdjacencyListGraph) EdgeCount() int {
	return falg.edgeCount
}

func (falg *FlaggedAdjacencyListGraph) GetPartition(id NodeId) PartitionId {
	if id < 0 || id >= falg.NodeCount() {
		panic(id)
	}
	return falg.Partitions[id]
}

func (falg *FlaggedAdjacencyListGraph) AddNode(n Node) {
	falg.Nodes = append(falg.Nodes, n)
	falg.Edges = append(falg.Edges, make([]FlaggedHalfEdge, 0))
}

func (falg *FlaggedAdjacencyListGraph) AddHalfEdge(from NodeId, fhe FlaggedHalfEdge) {
	// Check if both source and target node exit
	if from >= falg.NodeCount() || fhe.To >= falg.NodeCount() {
		panic(fmt.Sprintf("HalfEdge %v from %d out of range.", fhe, from))
	}
	// Check for duplicates
	for _, halfEdge := range falg.Edges[from] {
		if fhe.To == halfEdge.To {
			return // ignore duplicate edges
		}
	}
	falg.Edges[from] = append(falg.Edges[from], fhe)
	falg.edgeCount++
}
