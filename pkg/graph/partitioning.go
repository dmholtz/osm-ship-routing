package graph

import (
	"fmt"
	"math"
	"sort"
)

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

type piNode struct {
	node      Node
	id        NodeId
	partition PartitionId
}

func KdPartitioning(g Graph, depth int) *FlaggedAdjacencyArrayGraph {
	piNodes := make([]piNode, 0, g.NodeCount())
	for i := 0; i < g.NodeCount(); i++ {
		piNodes = append(piNodes, piNode{node: g.GetNode(i), id: i, partition: 0})
	}

	queue := make([][]piNode, 0)
	queue = append(queue, piNodes)

	for d := 0; d < depth; d++ {
		end := len(queue)
		fmt.Printf("Depth = %d, length = %d\n", d, end)

		for i := 0; i < end; i++ {
			piNodes = queue[0]
			queue = queue[1:]

			if d%2 == 0 {
				// north-south split (lat)
				sort.Slice(piNodes, func(i, j int) bool {
					return piNodes[i].node.Lat < piNodes[j].node.Lat
				})
			} else {
				// east-west split (lon)
				sort.Slice(piNodes, func(i, j int) bool {
					return piNodes[i].node.Lon < piNodes[j].node.Lon
				})
			}

			first := piNodes[:len(piNodes)/2]
			for j := 0; j < len(first); j++ {
				first[j].partition = first[j].partition << 1
			}

			second := piNodes[len(piNodes)/2:]
			for j := 0; j < len(second); j++ {
				second[j].partition = (second[j].partition << 1) + 1
			}

			queue = append(queue, first)
			queue = append(queue, second)
		}
	}

	piNodes = make([]piNode, 0)
	for _, s := range queue {
		piNodes = append(piNodes, s...)
	}
	sort.Slice(piNodes, func(i, j int) bool {
		return piNodes[i].id < piNodes[j].id
	})

	nodes := make([]Node, g.NodeCount(), g.NodeCount())
	partitions := make([]PartitionId, g.NodeCount(), g.NodeCount())
	halfEdges := make([]FlaggedHalfEdge, 0, g.EdgeCount())
	offsets := make([]int, g.NodeCount()+1, g.NodeCount()+1)
	for i := 0; i < g.NodeCount(); i++ {
		nodes[i] = piNodes[i].node
		partitions[i] = piNodes[i].partition
		offsets[i] = len(halfEdges)
		for _, he := range g.GetHalfEdgesFrom(i) {
			halfEdges = append(halfEdges, FlaggedHalfEdge{To: he.To, Weight: he.Distance})
		}
	}

	offsets[g.NodeCount()] = len(halfEdges)

	return &FlaggedAdjacencyArrayGraph{Nodes: nodes, Edges: halfEdges, Partitions: partitions, Offsets: offsets}
}
