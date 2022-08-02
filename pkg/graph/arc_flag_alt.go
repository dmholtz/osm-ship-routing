package graph

import "container/heap"

func ArcFlagAlt(g FlaggedGraph, landmarkDistancesCollection []LandmarkDistances, origin, destination int) ([]int, int, int) {
	dijkstraItems := make([]*AStarPriorityQueueItem, g.NodeCount(), g.NodeCount())
	originItem := AStarPriorityQueueItem{itemId: origin, priority: 0, predecessor: -1, index: -1}
	dijkstraItems[origin] = &originItem

	pq := make(AStarPriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[origin])

	destPart := g.GetPartition(destination)

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*AStarPriorityQueueItem)
		currentNodeId := currentPqItem.itemId
		pqPops++

		for _, edge := range g.GetHalfEdgesFrom(currentNodeId) {
			if !edge.IsFlagged(destPart) {
				continue
			}
			successor := edge.To

			if dijkstraItems[successor] == nil {
				newDistance := currentPqItem.distance + edge.Weight
				newPriority := newDistance + alt_heuristic(landmarkDistancesCollection, successor, destination)
				pqItem := AStarPriorityQueueItem{itemId: successor, priority: newPriority, distance: newDistance, predecessor: currentNodeId, index: -1}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedPriority := currentPqItem.distance + edge.Weight + alt_heuristic(landmarkDistancesCollection, successor, destination); updatedPriority < dijkstraItems[successor].priority {
					pq.update(dijkstraItems[successor], updatedPriority, currentPqItem.distance+edge.Weight)
					dijkstraItems[successor].predecessor = currentNodeId
				}
			}
		}

		if currentNodeId == destination {
			break
		}
	}

	length := -1           // by default a non-existing path has length -1
	path := make([]int, 0) // by default, a non-existing path is an empty slice
	if dijkstraItems[destination] != nil {
		length = dijkstraItems[destination].priority
		for nodeId := destination; nodeId != -1; nodeId = dijkstraItems[nodeId].predecessor {
			path = append([]int{nodeId}, path...)
		}
	}
	return path, length, pqPops
}
