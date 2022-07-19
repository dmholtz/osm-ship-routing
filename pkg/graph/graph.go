package graph

type NodeId = int

type Node struct {
	Lon float64
	Lat float64
}

func NewNode(lon float64, lat float64) *Node {
	return &Node{Lon: lon, Lat: lat}
}

type Edge struct {
	From     NodeId
	To       NodeId
	Distance int
}

func (e Edge) Invert() Edge {
	return Edge{From: e.To, To: e.From, Distance: e.Distance}
}

func (e Edge) toHalfEdge() HalfEdge {
	return HalfEdge{To: e.To, Distance: e.Distance}
}

type HalfEdge struct {
	To       NodeId
	Distance int
}

func (oe HalfEdge) toEdge(from NodeId) Edge {
	return Edge{From: from, To: oe.To, Distance: oe.Distance}
}

type HalfEdges = []HalfEdge

type Graph interface {
	GetNode(id NodeId) Node
	GetHalfEdgesFrom(id NodeId) []HalfEdge
	NodeCount() int
	EdgeCount() int
}

type DynamicGraph interface {
	Graph
	AddNode(n Node)
	AddEdge(e Edge)
}
