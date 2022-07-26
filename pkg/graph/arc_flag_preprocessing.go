package graph

import (
	"fmt"
	"sync"
)

type addFlagJob struct {
	from      NodeId
	to        NodeId
	partition PartitionId
}

const MAX_GOROUTINES = 8

// Input: partitioned graph (max 64 partitions) with zero arc-flags
// Output: partitioned graph with non-trivial arc-flags

// TODO re-think parameter type
func ComputeArcFlags(fg *FlaggedAdjacencyArrayGraph) FlaggedGraph {

	// preprocessing: remove existing arc flag
	for _, halfEdge := range fg.Edges {
		halfEdge.ResetFlag()
	}

	// determine boundary nodes for each region
	boundaryNodeSets := make([](map[NodeId]bool), 0)
	for i := 0; i < 64; i++ {
		boundaryNodeSets = append(boundaryNodeSets, make(map[NodeId]bool, 0))
	}
	for tailNodeId := range fg.Nodes {
		for _, halfEdge := range fg.GetHalfEdgesFrom(tailNodeId) {
			tailPartition := fg.GetPartition((tailNodeId))
			headPartition := fg.GetPartition(halfEdge.To)
			if tailPartition != headPartition {
				boundaryNodeSets[headPartition][halfEdge.To] = true
			}
		}
	}

	//for partition, set := range boundaryNodeSets {
	//	fmt.Printf("%d: ", partition)
	//	for k := range set {
	//		fmt.Printf("%d, ", k)
	//	}
	//	fmt.Printf("\n")
	//}

	// compute transposed graph
	transposedGraph := TransposeGraph(fg)

	jobs := make(chan addFlagJob, 1<<16)
	done := make(chan bool)
	guard := make(chan struct{}, MAX_GOROUTINES)
	wg := sync.WaitGroup{}

	// start consumer
	go addFlag(jobs, fg, done)

	for partition, set := range boundaryNodeSets {
		setSize := len(set)
		wg.Add(setSize)
		fmt.Printf("Partition: %d, size=%d\n", partition, setSize)
		for boundaryNodeId := range set {
			guard <- struct{}{}
			go backwardSearch(jobs, fg, transposedGraph, uint8(partition), boundaryNodeId, &wg, guard)
		}
	}

	wg.Wait()
	close(jobs)
	<-done

	// revise edges within the same partition
	for i := 0; i < fg.NodeCount(); i++ {
		for _, halfEdge := range fg.GetHalfEdgesFrom(i) {
			if fg.GetPartition(i) == fg.GetPartition(halfEdge.To) {
				AddFlag(fg, i, halfEdge.To, fg.GetPartition(i))
			}
		}
	}

	return fg
}

// producer function
func backwardSearch(jobs chan<- addFlagJob, graph *FlaggedAdjacencyArrayGraph, transposedGraph FlaggedGraph, partition PartitionId, boundaryNodeId NodeId, wg *sync.WaitGroup, guard <-chan struct{}) {
	// calculate in reverse graph
	tree := ShortestPathTree(transposedGraph, boundaryNodeId)

	stack := make([]*TreeNode, 0)
	for _, child := range tree.children {
		if graph.GetPartition(child.id) != partition {
			// add edge
			tailRev := tree.id
			headRev := child.id
			jobs <- addFlagJob{from: headRev, to: tailRev, partition: partition}
			stack = append(stack, child)
			child.visited = true
		}
	}

	for len(stack) > 0 {
		// pop
		node := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]

		for _, child := range node.children {
			if child.visited {
				continue
			}
			// add edge
			tailRev := node.id
			headRev := child.id
			jobs <- addFlagJob{from: headRev, to: tailRev, partition: partition}
			stack = append(stack, child)
			child.visited = true
		}
	}
	<-guard
	wg.Done()
}

// (single) consumer
func addFlag(jobs <-chan addFlagJob, g *FlaggedAdjacencyArrayGraph, done chan<- bool) {
	for job := range jobs {
		for i := g.Offsets[job.from]; i < g.Offsets[job.from+1]; i++ {
			edge := &g.Edges[i]
			if edge.To == job.to {
				edge.AddFlag(job.partition)
				break
			}
		}
	}
	done <- true
}

func AddFlag(fg *FlaggedAdjacencyArrayGraph, from NodeId, to NodeId, partition PartitionId) {
	for i := fg.Offsets[from]; i < fg.Offsets[from+1]; i++ {
		edge := &fg.Edges[i]
		if edge.To == to {
			edge.AddFlag(partition)
			break
		}
	}
}
