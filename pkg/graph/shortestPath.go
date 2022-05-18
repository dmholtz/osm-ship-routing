package graph

import (
	"container/heap"
	"fmt"
	"math"
)

type PriorityQueueItem struct {
	itemId      int // node id of this item
	priority    int // distance from origin to this node
	predecessor int // node id of the predecessor
	index       int // index of the item in the heap
}

// A PriorityQueue implements the heap.Interface and hold PriorityQueueItems
type PriorityQueue []*PriorityQueueItem

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
	pqItem := item.(*PriorityQueueItem)
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

func (h *PriorityQueue) update(pqItem *PriorityQueueItem, newPriority int) {
	pqItem.priority = newPriority
	heap.Fix(h, pqItem.index)
}

func Dijkstra(g Graph, origin, destination int) {
	dijkstraItems := make([]*PriorityQueueItem, g.NodeCount(), g.NodeCount())
	for i := 0; i < len(dijkstraItems); i++ {
		pqItem := PriorityQueueItem{itemId: i, priority: math.MaxInt, predecessor: -1, index: -1}
		dijkstraItems[i] = &pqItem
	}
	dijkstraItems[origin].priority = 0

	pq := make(PriorityQueue, 0)
	pq.Push(dijkstraItems[origin])
	heap.Init(&pq)

	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*PriorityQueueItem)
		currentNodeId := currentPqItem.itemId

		if currentNodeId == destination {
			fmt.Printf("Destination reached. Distance = %d\n", dijkstraItems[currentNodeId].priority)
			break
		}

		for _, edge := range g.GetEdgesFrom(currentNodeId) {
			successor := edge.To

			if updatedDistance := dijkstraItems[currentNodeId].priority + edge.Distance; updatedDistance < dijkstraItems[successor].priority {

				if dijkstraItems[successor].predecessor != -1 {
					// item is already in the pq
					pq.update(dijkstraItems[successor], updatedDistance)
				} else {
					// item is not yet in the pq
					dijkstraItems[successor].priority = updatedDistance
					heap.Push(&pq, dijkstraItems[successor])
				}

				dijkstraItems[successor].predecessor = currentNodeId

			}
		}
	}

	// turn off for benchmarking
	//for i := destination; i != -1; i = dijkstraItems[i].Predecessor {
	//	fmt.Printf("%v <- ", i)
	//}
}
