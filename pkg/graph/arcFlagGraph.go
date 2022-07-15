package graph

import "fmt"

// range of possible values must be restricted to [0,64) due to size of arc flags
type PartitionId = uint8

type ArcFlag = uint64

type FlaggedHalfEdge struct {
	To     NodeId
	Weight int
	Flag   ArcFlag
}

// Returns true iff the half edge flags the partition p
func (fhe FlaggedHalfEdge) IsFlagged(p PartitionId) bool {
	if !(p < 64) {
		panic(fmt.Sprintf("Partion %d out of range [0,64).", p))
	}
	return (fhe.Flag & (1 << p)) > 0
}

// Adds a flag for partition p to the half edge
func (fhe *FlaggedHalfEdge) AddFlag(p PartitionId) {
	if !(p < 64) {
		panic(fmt.Sprintf("Partion %d out of range [0,64).", p))
	}
	fhe.Flag = fhe.Flag | (1 << p)
}

type FlaggedGraph interface {
	GetNode(id NodeId) Node
	GetHalfEdgesFrom(id NodeId) []FlaggedHalfEdge
	NodeCount() int
	EdgeCount() int

	GetPartition(id NodeId) PartitionId
}

// Implementation for static graphs using an adjacency array
type ArcFlagGraph struct {
	Nodes      []Node
	Partitions []PartitionId
	Edges      []FlaggedHalfEdge
	Offsets    []int
}

func (afg *ArcFlagGraph) GetNode(id NodeId) Node {
	if id < 0 || id >= afg.NodeCount() {
		panic(fmt.Sprintf("NodeId %d is not contained in the graph.", id))
	}
	return afg.Nodes[id]
}

func (afg *ArcFlagGraph) GetHalfEdgesFrom(id NodeId) []FlaggedHalfEdge {
	if id < 0 || id >= afg.NodeCount() {
		panic(fmt.Sprintf("NodeId %d is not contained in the graph.", id))
	}
	return afg.Edges[afg.Offsets[id]:afg.Offsets[id+1]]
}

func (afg *ArcFlagGraph) NodeCount() int {
	return len(afg.Nodes)
}

func (afg *ArcFlagGraph) EdgeCount() int {
	return len(afg.Edges)
}
