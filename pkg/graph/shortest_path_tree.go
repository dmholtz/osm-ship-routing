package graph

import (
	"container/heap"
)

type TreeNode struct {
	id       NodeId
	children []*TreeNode
}

// TODO make graph generic
func ShortestPathTree(g FlaggedGraph, origin int) TreeNode {
	dijkstraItems := make([]*PriorityQueueItem, g.NodeCount(), g.NodeCount())
	originItem := PriorityQueueItem{itemId: origin, priority: 0, predecessor: -1, index: -1}
	dijkstraItems[origin] = &originItem

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[origin])

	successors := make([]*TreeNode, g.NodeCount(), g.NodeCount())

	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*PriorityQueueItem)
		currentNodeId := currentPqItem.itemId

		if currentNodeId != origin {
			predecessor := currentPqItem.predecessor
			if successors[predecessor] == nil {
				successors[predecessor] = &TreeNode{id: predecessor, children: make([]*TreeNode, 0)}
			}
			if successors[currentNodeId] == nil {
				successors[currentNodeId] = &TreeNode{id: currentNodeId, children: make([]*TreeNode, 0)}
			}
			successors[predecessor].children = append(successors[predecessor].children, successors[currentNodeId])
		}

		for _, edge := range g.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To

			if dijkstraItems[successor] == nil {
				newPriority := dijkstraItems[currentNodeId].priority + edge.Weight
				pqItem := PriorityQueueItem{itemId: successor, priority: newPriority, predecessor: currentNodeId, index: -1}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedDistance := dijkstraItems[currentNodeId].priority + edge.Weight; updatedDistance < dijkstraItems[successor].priority {
					pq.update(dijkstraItems[successor], updatedDistance)
					dijkstraItems[successor].predecessor = currentNodeId
				}
			}
		}
	}

	return *successors[origin]
}
