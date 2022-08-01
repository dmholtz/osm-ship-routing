package graph

import (
	"container/heap"

	geo "github.com/dmholtz/osm-ship-routing/pkg/geometry"
)

type AStarPriorityQueueItem struct {
	itemId      int  // node id of this item
	priority    int  // estimated distance from this node to destination (f-value)
	distance    int  // distance from origin to this node (g-value)
	predecessor int  // node id of the predecessor
	index       int  // index of the item in the heap
	settled     bool // true iff the node has been settled by the algorithm
}

// A AStarPriorityQueue implements the heap.Interface and hold AStarPriorityQueueItems
type AStarPriorityQueue []*AStarPriorityQueueItem

func (h AStarPriorityQueue) Len() int {
	return len(h)
}

func (h AStarPriorityQueue) Less(i, j int) bool {
	// MinHeap implementation
	return h[i].priority < h[j].priority
}

func (h AStarPriorityQueue) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

func (h *AStarPriorityQueue) Push(item interface{}) {
	n := len(*h)
	pqItem := item.(*AStarPriorityQueueItem)
	pqItem.index = n
	*h = append(*h, pqItem)
}

func (h *AStarPriorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	pqItem := old[n-1]
	old[n-1] = nil
	pqItem.index = -1 // for safety
	*h = old[0 : n-1]
	return pqItem
}

func (h *AStarPriorityQueue) update(pqItem *AStarPriorityQueueItem, newPriority int, newDistance int) {
	pqItem.priority = newPriority
	pqItem.distance = newDistance
	heap.Fix(h, pqItem.index)
}

func heuristic(g Graph, nodeId int, targetId int) int {
	node, targetNode := g.GetNode(nodeId), g.GetNode(targetId)
	p1 := geo.NewPoint(node.Lat, node.Lon)
	p2 := geo.NewPoint(targetNode.Lat, targetNode.Lon)
	return int(p1.Haversine(p2) * 0.9999)
	//return p1.IntHaversine(p2)
}

func AStar(g Graph, origin, destination int) ([]int, int, int) {
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
		pqPops += 1

		for _, edge := range g.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To

			if dijkstraItems[successor] == nil {
				newDistance := currentPqItem.distance + edge.Distance
				newPriority := newDistance + heuristic(g, successor, destination)
				pqItem := AStarPriorityQueueItem{itemId: successor, priority: newPriority, distance: newDistance, predecessor: currentNodeId, index: -1}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedPriority := currentPqItem.distance + edge.Distance + heuristic(g, successor, destination); updatedPriority < dijkstraItems[successor].priority {
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
