package graph

import (
	"container/heap"
	"math"
)

func ArcFlagBiDijkstra(g FlaggedGraph, origin, destination int) ([]int, int, int) {
	// reference: https://www.homepages.ucl.ac.uk/~ucahmto/math/2020/05/30/bidirectional-dijkstra.html

	dijkstraItemsForward := make([]*PriorityQueueItem, g.NodeCount(), g.NodeCount())
	originItem := PriorityQueueItem{itemId: origin, priority: 0, predecessor: -1, index: -1}
	dijkstraItemsForward[origin] = &originItem

	dijkstraItemsBackward := make([]*PriorityQueueItem, g.NodeCount(), g.NodeCount())
	targetItem := PriorityQueueItem{itemId: destination, priority: 0, predecessor: -1, index: -1}
	dijkstraItemsBackward[destination] = &targetItem

	pqForward := make(PriorityQueue, 0)
	heap.Init(&pqForward)
	heap.Push(&pqForward, dijkstraItemsForward[origin])

	pqBackward := make(PriorityQueue, 0)
	heap.Init(&pqBackward)
	heap.Push(&pqBackward, dijkstraItemsBackward[destination])

	mu := math.MaxInt // will contain the shortest distance once the loop terminates
	middleNodeId := 0

	origPart := g.GetPartition(origin)
	destPart := g.GetPartition(destination)

	pqPops := 0
	// works only on undirected graphs
	for len(pqForward) > 0 && len(pqBackward) > 0 {
		forwardPqItem := heap.Pop(&pqForward).(*PriorityQueueItem)
		forwardNodeId := forwardPqItem.itemId
		backwardPqItem := heap.Pop(&pqBackward).(*PriorityQueueItem)
		backwardNodeId := backwardPqItem.itemId
		pqPops += 2

		// stopping criterion
		if dijkstraItemsForward[forwardNodeId].priority+dijkstraItemsBackward[backwardNodeId].priority >= mu {
			break
		}

		// forward search
		for _, edge := range g.GetHalfEdgesFrom(forwardNodeId) {
			if !edge.IsFlagged(destPart) {
				continue
			}
			successor := edge.To

			if dijkstraItemsForward[successor] == nil {
				newPriority := dijkstraItemsForward[forwardNodeId].priority + edge.Weight
				pqItem := PriorityQueueItem{itemId: successor, priority: newPriority, predecessor: forwardNodeId, index: -1}
				dijkstraItemsForward[successor] = &pqItem
				heap.Push(&pqForward, &pqItem)
			} else {
				if updatedDistance := dijkstraItemsForward[forwardNodeId].priority + edge.Weight; updatedDistance < dijkstraItemsForward[successor].priority {
					pqForward.update(dijkstraItemsForward[successor], updatedDistance)
					dijkstraItemsForward[successor].predecessor = forwardNodeId
				}
			}

			if x := dijkstraItemsBackward[successor]; x != nil && dijkstraItemsForward[forwardNodeId].priority+edge.Weight+x.priority < mu {
				mu = dijkstraItemsForward[forwardNodeId].priority + edge.Weight + x.priority
				dijkstraItemsForward[successor].predecessor = forwardNodeId
				middleNodeId = successor
			}
		}

		// backward search
		for _, edge := range g.GetHalfEdgesFrom(backwardNodeId) {
			if !edge.IsFlagged(origPart) {
				continue
			}
			successor := edge.To

			if dijkstraItemsBackward[successor] == nil {
				newPriority := dijkstraItemsBackward[backwardNodeId].priority + edge.Weight
				pqItem := PriorityQueueItem{itemId: successor, priority: newPriority, predecessor: backwardNodeId, index: -1}
				dijkstraItemsBackward[successor] = &pqItem
				heap.Push(&pqBackward, &pqItem)
			} else {
				if updatedDistance := dijkstraItemsBackward[backwardNodeId].priority + edge.Weight; updatedDistance < dijkstraItemsBackward[successor].priority {
					pqBackward.update(dijkstraItemsBackward[successor], updatedDistance)
					dijkstraItemsBackward[successor].predecessor = backwardNodeId
				}
			}

			if x := dijkstraItemsForward[successor]; x != nil && dijkstraItemsBackward[backwardNodeId].priority+edge.Weight+x.priority < mu {
				mu = dijkstraItemsBackward[backwardNodeId].priority + edge.Weight + x.priority
				dijkstraItemsBackward[successor].predecessor = backwardNodeId
				middleNodeId = successor
			}
		}
	}

	length := -1           // by default a non-existing path has length -1
	path := make([]int, 0) // by default, a non-existing path is an empty slice

	// check if path exists
	if mu < math.MaxInt {
		length = mu
		// sanity check: length == dijkstraItemsForward[middleNodeId].priority + dijkstraItemsBackward[middleNodeId].priority
		if dijkstraItemsForward[middleNodeId] != nil && dijkstraItemsBackward[middleNodeId] != nil {
			for nodeId := middleNodeId; nodeId != -1; nodeId = dijkstraItemsForward[nodeId].predecessor {
				path = append([]int{nodeId}, path...)
			}
			if path[len(path)-1] == middleNodeId {
				path = path[0 : len(path)-1]
			}
			for nodeId := middleNodeId; nodeId != -1; nodeId = dijkstraItemsBackward[nodeId].predecessor {
				path = append(path, nodeId)
			}
		}
	}

	return path, length, pqPops
}
