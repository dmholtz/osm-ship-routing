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
	alg.Edges[e.From] = append(alg.Edges[e.From], e.toHalfEdge())
	alg.edgeCount++
}
