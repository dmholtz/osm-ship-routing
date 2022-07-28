package graph

import (
	"container/heap"
	"math"
)

type DjkstraItem struct {
	Distance    int
	Predecessor NodeId
}

func DijkstraOneToAll(g Graph, origin int) []DjkstraItem {
	dijkstraItems := make([]*PriorityQueueItem, g.NodeCount(), g.NodeCount())
	originItem := PriorityQueueItem{itemId: origin, priority: 0, predecessor: -1, index: -1}
	dijkstraItems[origin] = &originItem

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[origin])

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*PriorityQueueItem)
		currentNodeId := currentPqItem.itemId
		pqPops++

		for _, edge := range g.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To

			if dijkstraItems[successor] == nil {
				newPriority := dijkstraItems[currentNodeId].priority + edge.Distance
				pqItem := PriorityQueueItem{itemId: successor, priority: newPriority, predecessor: currentNodeId, index: -1}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedDistance := dijkstraItems[currentNodeId].priority + edge.Distance; updatedDistance < dijkstraItems[successor].priority {
					pq.update(dijkstraItems[successor], updatedDistance)
					dijkstraItems[successor].predecessor = currentNodeId
				}
			}
		}
	}

	// explicitly mark nil-entries
	result := make([]DjkstraItem, 0, g.NodeCount())
	for _, pqItem := range dijkstraItems {
		var dijkstraItem DjkstraItem
		if pqItem != nil {
			dijkstraItem = DjkstraItem{Distance: pqItem.priority, Predecessor: pqItem.predecessor}
		} else {
			dijkstraItem = DjkstraItem{Distance: math.MaxInt, Predecessor: -1}
		}
		result = append(result, dijkstraItem)
	}
	return result
}
