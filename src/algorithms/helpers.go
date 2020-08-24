package algorithms

import (
	"math"
)

func haversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func distance(start, end []float64) float64 {
	var fLat, fLng, fLat2, fLng2, radius float64
	fLng = start[0] * math.Pi / 180.0
	fLat = start[1] * math.Pi / 180.0
	fLng2 = end[0] * math.Pi / 180.0
	fLat2 = end[1] * math.Pi / 180.0

	radius = 6378100
	h := haversin(fLat2-fLat) + math.Cos(fLat)*math.Cos(fLat2)*haversin(fLng2-fLng)
	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
	return (c * radius)
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
