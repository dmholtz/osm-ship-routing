package graph

import (
	"container/heap"
	"fmt"
	"math"
)

type PriorityQueueItem struct {
	itemId   int // node id of this item
	priority int // distance from origin to this node
	index    int // index of the item in the heap
}

// A PriorityQueue implements the heap.Interface and hold PriorityQueueItems
type PriorityQueue []PriorityQueueItem

func (h PriorityQueue) Len() int {
	return len(h)
}

func (h PriorityQueue) Less(i, j int) bool {
	// MinHeap implementation
	return h[i].priority < h[j].priority
}

func (h PriorityQueue) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

func (h *PriorityQueue) Push(item interface{}) {
	n := len(*h)
	pqItem := item.(PriorityQueueItem)
	pqItem.index = n
	*h = append(*h, pqItem)
}

func (h *PriorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	pqItem := old[n-1]
	pqItem.index = -1 // for safety
	*h = old[0 : n-1]
	return pqItem
}

func (h *PriorityQueue) update(pqItemId int, newPriority int) {
	(*h)[pqItemId].priority = newPriority
	index := (*h)[pqItemId].index
	heap.Fix(h, index)

}

type DijkstraItem struct {
	Distance    int
	Predecessor int
}

func (aag *AdjacencyArrayGraph) Dijkstra(origin, destination int) {
	dijkstraItems := make([]DijkstraItem, aag.NodeCount(), aag.NodeCount())
	for i := 0; i < len(dijkstraItems); i++ {
		dijkstraItems[i] = DijkstraItem{Distance: math.MaxInt, Predecessor: -1}
	}
	dijkstraItems[origin].Distance = 0

	pq := make(PriorityQueue, 0)
	pq.Push(PriorityQueueItem{itemId: origin, priority: 0})
	heap.Init(&pq)

	for len(pq) > 0 {
		currentNodeId := heap.Pop(&pq).(PriorityQueueItem).itemId

		if currentNodeId == destination {
			fmt.Printf("Destination reached. Distance = %d\n", dijkstraItems[currentNodeId].Distance)
			break
		}

		for _, edge := range aag.GetEdgesFrom(currentNodeId) {
			successor := edge.To
			if updatedDistance := dijkstraItems[currentNodeId].Distance + edge.Distance; updatedDistance < dijkstraItems[successor].Distance {
				dijkstraItems[successor].Distance = updatedDistance
				dijkstraItems[successor].Predecessor = currentNodeId
				pqItem := PriorityQueueItem{itemId: successor, priority: updatedDistance}
				heap.Push(&pq, pqItem)
			}
		}
	}

	// turn off for benchmarking
	//for i := destination; i != -1; i = dijkstraItems[i].Predecessor {
	//	fmt.Printf("%v <- ", i)
	//}
}
