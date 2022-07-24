package graph

import (
	"container/heap"
)

type TreeNode struct {
	id       NodeId
	children []*TreeNode
}

type PqItem struct {
	itemId       int   // node id of this item
	priority     int   // distance from origin to this node
	predecessors []int // node id of the predecessor
	index        int   // index of the item in the heap
}

// A PriorityQueue implements the heap.Interface and hold PriorityQueueItems
type PriorityQueue1 []*PqItem

func (h PriorityQueue1) Len() int {
	return len(h)
}

func (h PriorityQueue1) Less(i, j int) bool {
	// MinHeap implementation
	return h[i].priority < h[j].priority
}

func (h PriorityQueue1) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

func (h *PriorityQueue1) Push(item interface{}) {
	n := len(*h)
	pqItem := item.(*PqItem)
	pqItem.index = n
	*h = append(*h, pqItem)
}

func (h *PriorityQueue1) Pop() interface{} {
	old := *h
	n := len(old)
	pqItem := old[n-1]
	old[n-1] = nil
	pqItem.index = -1 // for safety
	*h = old[0 : n-1]
	return pqItem
}

func (h *PriorityQueue1) update(pqItem *PqItem, newPriority int) {
	pqItem.priority = newPriority
	heap.Fix(h, pqItem.index)
}

// TODO make graph generic
func ShortestPathTree(g FlaggedGraph, origin int) TreeNode {
	dijkstraItems := make([]*PqItem, g.NodeCount(), g.NodeCount())
	originItem := PqItem{itemId: origin, priority: 0, predecessors: make([]int, 0), index: -1}
	dijkstraItems[origin] = &originItem

	pq := make(PriorityQueue1, 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[origin])

	successors := make([]*TreeNode, 0)
	for i := 0; i < g.NodeCount(); i++ {
		successors = append(successors, &TreeNode{id: i, children: make([]*TreeNode, 0)})
	}

	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*PqItem)
		currentNodeId := currentPqItem.itemId

		if currentNodeId != origin {
			for _, pred := range currentPqItem.predecessors {
				successors[pred].children = append(successors[pred].children, successors[currentNodeId])
			}

		}

		for _, edge := range g.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To

			if dijkstraItems[successor] == nil {
				newPriority := dijkstraItems[currentNodeId].priority + edge.Weight
				predecessors := make([]int, 0)
				predecessors = append(predecessors, currentNodeId)
				pqItem := PqItem{itemId: successor, priority: newPriority, predecessors: predecessors, index: -1}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedDistance := dijkstraItems[currentNodeId].priority + edge.Weight; updatedDistance < dijkstraItems[successor].priority {
					pq.update(dijkstraItems[successor], updatedDistance)
					// reset predecessors
					predecessors := make([]int, 0)
					predecessors = append(predecessors, currentNodeId)
					dijkstraItems[successor].predecessors = predecessors
				} else if updatedDistance == dijkstraItems[successor].priority {
					// append
					dijkstraItems[successor].predecessors = append(dijkstraItems[successor].predecessors, currentNodeId)
				}
			}
		}
	}

	return *successors[origin]
}
