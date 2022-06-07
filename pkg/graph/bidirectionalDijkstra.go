package graph

import (
	"container/heap"
)

func BidirectionalDijkstra(g Graph, origin, destination int) ([]int, int) {
	dijkstraItemsForward := make([]*PriorityQueueItem, g.NodeCount(), g.NodeCount())
	originItem := PriorityQueueItem{itemId: origin, priority: 0, predecessor: -1, index: -1}
	dijkstraItemsForward[origin] = &originItem

	dijkstraItemsBackward := make([]*PriorityQueueItem, g.NodeCount(), g.NodeCount())
	targetItem := PriorityQueueItem{itemId: destination, priority: 0, predecessor: -1, index: -1}
	dijkstraItemsBackward[destination] = &targetItem

	visitedNodes := make([]bool, g.NodeCount(), g.NodeCount())

	pqForward := make(PriorityQueue, 0)
	heap.Init(&pqForward)
	heap.Push(&pqForward, dijkstraItemsForward[origin])

	pqBackward := make(PriorityQueue, 0)
	heap.Init(&pqBackward)
	heap.Push(&pqBackward, dijkstraItemsBackward[destination])

	currentNodeId := 0
	for len(pqForward) > 0 && len(pqBackward) > 0 {
		currentPqItem := heap.Pop(&pqForward).(*PriorityQueueItem)
		currentNodeId = currentPqItem.itemId

		for _, edge := range g.GetEdgesFrom(currentNodeId) {
			successor := edge.To

			if dijkstraItemsForward[successor] == nil {
				newPriority := dijkstraItemsForward[currentNodeId].priority + edge.Distance
				pqItem := PriorityQueueItem{itemId: successor, priority: newPriority, predecessor: currentNodeId, index: -1}
				dijkstraItemsForward[successor] = &pqItem
				heap.Push(&pqForward, &pqItem)
			} else {
				if updatedDistance := dijkstraItemsForward[currentNodeId].priority + edge.Distance; updatedDistance < dijkstraItemsForward[successor].priority {
					pqForward.update(dijkstraItemsForward[successor], updatedDistance)
					dijkstraItemsForward[successor].predecessor = currentNodeId
				}
			}
		}

		if visitedNodes[currentNodeId] {
			break
		} else {
			visitedNodes[currentNodeId] = true
		}

		currentPqItem = heap.Pop(&pqBackward).(*PriorityQueueItem)
		currentNodeId = currentPqItem.itemId

		for _, edge := range g.GetEdgesFrom(currentNodeId) {
			successor := edge.To

			if dijkstraItemsBackward[successor] == nil {
				newPriority := dijkstraItemsBackward[currentNodeId].priority + edge.Distance
				pqItem := PriorityQueueItem{itemId: successor, priority: newPriority, predecessor: currentNodeId, index: -1}
				dijkstraItemsBackward[successor] = &pqItem
				heap.Push(&pqBackward, &pqItem)
			} else {
				if updatedDistance := dijkstraItemsBackward[currentNodeId].priority + edge.Distance; updatedDistance < dijkstraItemsBackward[successor].priority {
					pqBackward.update(dijkstraItemsBackward[successor], updatedDistance)
					dijkstraItemsBackward[successor].predecessor = currentNodeId
				}
			}
		}

		if visitedNodes[currentNodeId] {
			break
		} else {
			visitedNodes[currentNodeId] = true
		}
	}

	length := -1           // by default a non-existing path has length -1
	path := make([]int, 0) // by default, a non-existing path is an empty slice
	if dijkstraItemsForward[currentNodeId] != nil {
		length = dijkstraItemsForward[currentNodeId].priority + dijkstraItemsBackward[currentNodeId].priority
		for nodeId := currentNodeId; nodeId != -1; nodeId = dijkstraItemsForward[nodeId].predecessor {
			path = append([]int{nodeId}, path...)
		}
		if path[len(path)-1] == currentNodeId {
			path = path[0 : len(path)-1]
		}
		for nodeId := currentNodeId; nodeId != -1; nodeId = dijkstraItemsBackward[nodeId].predecessor {
			path = append(path, nodeId)
		}
	}
	return path, length
}
