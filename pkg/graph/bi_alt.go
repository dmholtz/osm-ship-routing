package graph

import (
	"container/heap"
	"math"
)

func BidirectionalAlt(g Graph, landmarkDistancesCollection []LandmarkDistances, origin, destination int) ([]int, int, int) {
	dijkstraItemsForward := make([]*AStarPriorityQueueItem, g.NodeCount(), g.NodeCount())
	originItem := AStarPriorityQueueItem{itemId: origin, priority: 0, distance: 0, predecessor: -1, index: -1}
	dijkstraItemsForward[origin] = &originItem

	dijkstraItemsBackward := make([]*AStarPriorityQueueItem, g.NodeCount(), g.NodeCount())
	targetItem := AStarPriorityQueueItem{itemId: destination, priority: 0, distance: 0, predecessor: -1, index: -1}
	dijkstraItemsBackward[destination] = &targetItem

	pqForward := make(AStarPriorityQueue, 0)
	heap.Init(&pqForward)
	heap.Push(&pqForward, dijkstraItemsForward[origin])

	pqBackward := make(AStarPriorityQueue, 0)
	heap.Init(&pqBackward)
	heap.Push(&pqBackward, dijkstraItemsBackward[destination])

	mu := math.MaxInt // will contain the shortest distance once the loop terminates
	middleNodeId := 0

	// works only on undirected graphs
	pqPops := 0
	for len(pqForward) > 0 && len(pqBackward) > 0 {
		forwardPqItem := heap.Pop(&pqForward).(*AStarPriorityQueueItem)
		forwardPqItem.settled = true
		forwardNodeId := forwardPqItem.itemId
		backwardPqItem := heap.Pop(&pqBackward).(*AStarPriorityQueueItem)
		backwardNodeId := backwardPqItem.itemId
		backwardPqItem.settled = true
		pqPops += 2

		// forward search
		for _, edge := range g.GetHalfEdgesFrom(forwardNodeId) {
			successor := edge.To

			if dijkstraItemsBackward[successor] != nil && dijkstraItemsBackward[successor].settled {
				// improvement by Kwa: An admissible bidirectional staged heuristic search algorithm
				//continue
			}
			if dijkstraItemsForward[successor] == nil {
				newDistance := forwardPqItem.distance + edge.Distance
				newPriority := newDistance + alt_heuristic(landmarkDistancesCollection, successor, destination)
				pqItem := AStarPriorityQueueItem{itemId: successor, priority: newPriority, distance: newDistance, predecessor: forwardNodeId, index: -1}
				dijkstraItemsForward[successor] = &pqItem
				heap.Push(&pqForward, &pqItem)
			} else {
				if updatedPriority := forwardPqItem.distance + edge.Distance + alt_heuristic(landmarkDistancesCollection, successor, destination); updatedPriority < dijkstraItemsForward[successor].priority {
					pqForward.update(dijkstraItemsForward[successor], updatedPriority, forwardPqItem.distance+edge.Distance)
					dijkstraItemsForward[successor].predecessor = forwardNodeId
				}
			}
			if x := dijkstraItemsBackward[successor]; x != nil && dijkstraItemsForward[forwardNodeId].distance+edge.Distance+x.distance < mu {
				mu = dijkstraItemsForward[forwardNodeId].distance + edge.Distance + x.distance
				dijkstraItemsForward[successor].predecessor = forwardNodeId
				middleNodeId = successor
			}
		}

		// backward search
		for _, edge := range g.GetHalfEdgesFrom(backwardNodeId) {
			successor := edge.To

			if mu >= math.MaxInt && dijkstraItemsForward[successor] != nil && dijkstraItemsForward[successor].settled {
				// improvement by Kwa: An admissible bidirectional staged heuristic search algorithm
				//continue
			}
			if dijkstraItemsBackward[successor] == nil {
				newDistance := backwardPqItem.distance + edge.Distance
				newPriority := newDistance + alt_heuristic(landmarkDistancesCollection, successor, origin)
				pqItem := AStarPriorityQueueItem{itemId: successor, priority: newPriority, distance: newDistance, predecessor: backwardNodeId, index: -1}
				dijkstraItemsBackward[successor] = &pqItem
				heap.Push(&pqBackward, &pqItem)
			} else {
				if updatedPriority := dijkstraItemsBackward[backwardNodeId].distance + edge.Distance + alt_heuristic(landmarkDistancesCollection, successor, origin); updatedPriority < dijkstraItemsBackward[successor].priority {
					pqBackward.update(dijkstraItemsBackward[successor], updatedPriority, backwardPqItem.distance+edge.Distance)
					dijkstraItemsBackward[successor].predecessor = backwardNodeId
				}
			}

			if x := dijkstraItemsForward[successor]; x != nil && dijkstraItemsBackward[backwardNodeId].distance+edge.Distance+x.distance < mu {
				mu = dijkstraItemsBackward[backwardNodeId].distance + edge.Distance + x.distance
				dijkstraItemsBackward[successor].predecessor = backwardNodeId
				middleNodeId = successor
			}
		}

		// stopping criterion (Symmetric Approach, cf Pohl: Bi-directional Search, 1971)
		if dijkstraItemsForward[forwardNodeId].priority >= mu || dijkstraItemsBackward[backwardNodeId].priority >= mu {
			break
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
