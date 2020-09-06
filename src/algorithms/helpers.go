package algorithms

import (
	"math"
)

func distance(start, end []float64) float64 {
	var lat1 = start[1] * math.Pi / 180
	var lat2 = end[1] * math.Pi / 180
	var dlat = (end[1] - start[1]) * math.Pi / 180
	var dlon = (end[0] - start[0]) * math.Pi / 180

	var a = math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	return 6371000.0 * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func mod(a, b int) int {
	a = a % b
	if a >= 0 {
		return a
	}
	if b < 0 {
		return a - b
	}
	return a + b
}

// Item of priority queue
type Item struct {
	value    int     // The value of the item; arbitrary.
	priority float64 // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*Item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push item into priority queue
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop item from priority queue
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
