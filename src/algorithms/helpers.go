package algorithms

import (
	"math"
)

// GridToCoord transforms a grid index to lat lon coordinates
func GridToCoord(in []int64, xSize, ySize int64) []float64 {
	var out []float64
	var xFactor = xSize / 360
	var yFactor = ySize / 360
	out = append(out, float64(in[0]/xFactor-180))
	out = append(out, float64((in[1]/yFactor)/2-90))
	return out
}

// CoordToGrid transforms lat lon coordinates to a grid index
func CoordToGrid(in []float64, xSize, ySize int64) []int64 {
	var out []int64
	var xFactor, yFactor float64
	xFactor = float64(xSize / 360)
	yFactor = float64(ySize / 360)
	// TODO check
	out = append(out, int64(((math.Round(in[0]*xFactor)/xFactor)+180)*xFactor))
	out = append(out, int64(((math.Round(in[1]*yFactor)/yFactor)+90)*2*yFactor))
	return out
}

func flattenIndex(x, y, xSize int64) int64 {
	return ((xSize * y) + x)
}

// ExpandIndex from 1d to 2d
func ExpandIndex(indx, xSize int64) []int64 {
	var x = indx % xSize
	var y = indx / xSize
	return []int64{x, y}
}

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

// Item of priority queue
type Item struct {
	value    int64   // The value of the item; arbitrary.
	priority float64 // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push item into priority queue
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop item from priority queue
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
