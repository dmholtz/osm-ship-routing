package graph

import "fmt"

// Implementation for static graphs
type AdjacencyArrayGraph struct {
	Nodes   []Node
	Edges   []HalfEdge
	Offsets []int
}

func NewAdjacencyArrayFromGraph(g Graph) *AdjacencyArrayGraph {
	nodes := make([]Node, 0)
	edges := make([]HalfEdge, 0)
	offsets := make([]int, g.NodeCount()+1, g.NodeCount()+1)

	for i := 0; i < g.NodeCount(); i++ {
		// add node
		nodes = append(nodes, g.GetNode(i))

		// add all edges of node
		for _, edge := range g.GetEdgesFrom(i) {
			edges = append(edges, edge.toOutgoingEdge())
		}

		// set stop-offset
		offsets[i+1] = len(edges)
	}

	aag := AdjacencyArrayGraph{Nodes: nodes, Edges: edges, Offsets: offsets}
	return &aag
}

func (aag *AdjacencyArrayGraph) GetNode(id NodeId) Node {
	if id < 0 || id >= aag.NodeCount() {
		panic(fmt.Sprintf("NodeId %d is not contained in the graph.", id))
	}
	return aag.Nodes[id]
}

func (aag *AdjacencyArrayGraph) GetEdgesFrom(id NodeId) []Edge {
	if id < 0 || id >= aag.NodeCount() {
		panic(fmt.Sprintf("NodeId %d is not contained in the graph.", id))
	}
	edges := make([]Edge, 0)
	for i := aag.Offsets[id]; i < aag.Offsets[id+1]; i++ {
		edges = append(edges, aag.Edges[i].toEdge(id))
	}
	return edges
}

func (aag *AdjacencyArrayGraph) GetHalfEdgesFrom(id NodeId) []HalfEdge {
	if id < 0 || id >= aag.NodeCount() {
		panic(fmt.Sprintf("NodeId %d is not contained in the graph.", id))
	}
	return aag.Edges[aag.Offsets[id]:aag.Offsets[id+1]]
}

func (aag *AdjacencyArrayGraph) NodeCount() int {
	return len(aag.Nodes)
}

func (aag *AdjacencyArrayGraph) EdgeCount() int {
	return len(aag.Edges)
}
