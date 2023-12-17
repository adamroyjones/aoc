package main

import (
	"container/heap"
)

// priroityQueue adapts an example from the container/heap documentation.
type priorityQueue []*priorityQueueItem

type priorityQueueItem struct {
	node     *node
	distance int // Used for the priority.
	index    int // Maintained by the heap.Interface methods.
}

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].distance < pq[j].distance }

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*priorityQueueItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // To avoid a memory leak.
	*pq = old[:n-1]
	return item
}

func (pq *priorityQueue) update(item *priorityQueueItem, distance int) {
	item.distance = distance
	heap.Fix(pq, item.index)
}
