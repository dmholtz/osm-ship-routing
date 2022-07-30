package graph

import "container/heap"

func max(x, y int) int {
	if x >= y {
		return x
	} else {
		return y
	}
}

func alt_heuristic(landmarkDistancesCollection []LandmarkDistances, from NodeId, to NodeId) int {
	upper_bound := 0
	for _, landmark := range landmarkDistancesCollection {
		upper_bound = max(upper_bound, landmark.DistancesFrom[to]-landmark.DistancesFrom[from])
		upper_bound = max(upper_bound, landmark.DistancesTo[from]-landmark.DistancesFrom[to])
	}
	return upper_bound
}

func Alt(g Graph, landmarkDistancesCollection []LandmarkDistances, origin, destination int) ([]int, int, int) {
	dijkstraItems := make([]*AStarPriorityQueueItem, g.NodeCount(), g.NodeCount())
	originItem := AStarPriorityQueueItem{itemId: origin, priority: 0, predecessor: -1, index: -1}
	dijkstraItems[origin] = &originItem

	pq := make(AStarPriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[origin])

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*AStarPriorityQueueItem)
		currentNodeId := currentPqItem.itemId
		pqPops++

		for _, edge := range g.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To

			if dijkstraItems[successor] == nil {
				newDistance := currentPqItem.distance + edge.Distance
				newPriority := newDistance + alt_heuristic(landmarkDistancesCollection, successor, destination)
				pqItem := AStarPriorityQueueItem{itemId: successor, priority: newPriority, distance: newDistance, predecessor: currentNodeId, index: -1}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedPriority := currentPqItem.distance + edge.Distance + alt_heuristic(landmarkDistancesCollection, successor, destination); updatedPriority < dijkstraItems[successor].priority {
					pq.update(dijkstraItems[successor], updatedPriority, currentPqItem.distance+edge.Distance)
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
