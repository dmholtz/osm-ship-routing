package graph

type NodeId = int

type Node struct {
	Lon float64
	Lat float64
}

type Edge struct {
	From     NodeId
	To       NodeId
	Distance int
}

func (e Edge) Invert() Edge {
	return Edge{From: e.To, To: e.From, Distance: e.Distance}
}

func (e Edge) toOutgoingEdge() outgoingEdge {
	return outgoingEdge{To: e.To, Distance: e.Distance}
}

type outgoingEdge struct {
	To       NodeId
	Distance int
}

func (oe outgoingEdge) toEdge(from NodeId) Edge {
	return Edge{From: from, To: oe.To, Distance: oe.Distance}
}

type outgoingEdges = []outgoingEdge

type Graph interface {
	GetNode(id NodeId) Node
	GetEdgesFrom(id NodeId) []Edge
	NodeCount() int
	EdgeCount() int
}

type DynamicGraph interface {
	Graph
	AddNode(n Node)
	AddEdge(e Edge)
}
