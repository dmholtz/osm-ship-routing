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

func (fhe *FlaggedHalfEdge) ResetFlag() {
	fhe.Flag = 0
}

type FlaggedGraph interface {
	GetNode(id NodeId) Node
	GetHalfEdgesFrom(id NodeId) []FlaggedHalfEdge
	NodeCount() int
	EdgeCount() int

	GetPartition(id NodeId) PartitionId
}

type DynamicFlaggedGraph interface {
	FlaggedGraph

	AddNode(node Node)
	AddHalfEdge(fhe FlaggedHalfEdge)
}
