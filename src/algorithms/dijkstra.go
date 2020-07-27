package algorithms

import (
	"container/heap"
	"math"
)

var meshgrid []bool

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

func gridToCoord(in []int64) []float64 {
	var out []float64
	// big grid
	//out = append(out, float64((float64(in[0])/10)-180))
	//out = append(out, float64(((float64(in[1])/10)/2)-90))
	// small grid
	out = append(out, float64(in[0]-180))
	out = append(out, float64((float64(in[1])/2)-90))
	return out
}

func coordToGrid(in []float64) []int64 {
	var out []int64
	// big grid
	//out = append(out, int64(((math.Round(in[0]*10)/10)+180)*10))
	//out = append(out, int64(((math.Round(in[1]*10)/10)+90)*2*10))
	// small grid
	out = append(out, int64(math.Round(in[0]))+180)
	out = append(out, (int64(math.Round(in[1]))+90)*2)
	return out
}

func flattenIndx(x, y, xSize int64) int64 {
	return ((xSize * y) + x)
}

func expandIndx(indx, xSize int64) []int64 {
	var x = indx % xSize
	var y = indx / xSize
	return []int64{x, y}
}

func haversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func neighbours1d(indx, xSize int64) []int64 {
	var neighbours []int64
	var temp []int64

	neighbours = append(neighbours, indx-xSize-1) // top left
	neighbours = append(neighbours, indx-xSize)   // top
	neighbours = append(neighbours, indx-xSize+1) // top right
	neighbours = append(neighbours, indx-1)       // left
	neighbours = append(neighbours, indx+1)       // right
	neighbours = append(neighbours, indx+xSize-1) // bottom left
	neighbours = append(neighbours, indx+xSize)   // bottom
	neighbours = append(neighbours, indx+xSize+1) // bottom right

	for _, j := range neighbours {
		if j >= 0 && j < int64(len(meshgrid)) {
			if !meshgrid[j] {
				temp = append(temp, j)
			}
		}
	}
	return temp
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

func extractRoute(prev *[]int64, end, xSize, ySize int64) [][][]float64 {
	var route [][][]float64
	var tempRoute [][]float64
	temp := expandIndx(end, xSize)
	for {
		x := expandIndx(end, xSize)
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			route = append(route, tempRoute)
			tempRoute = make([][]float64, 0)
		}
		tempRoute = append(tempRoute, gridToCoord([]int64{x[0], x[1]}))

		if (*prev)[end] == -1 {
			break
		}
		end = (*prev)[end]
		temp = x
	}
	route = append(route, tempRoute)
	return route
}

func extractNodes(nodesProcessed *[]int64, xSize, ySize int64) [][]float64 {
	var nodesExtended [][]float64
	for _, node := range *nodesProcessed {
		coord := gridToCoord([]int64{x[0], x[1]})
		x := expandIndx(node, xSize)
		nodesExtended = append(nodesExtended, coord)
	}
	return nodesExtended
}

// Dijkstra implementation
func Dijkstra(startLngInt, startLatInt, endLngInt, endLatInt, xSize, ySize int64, meshgridPointer *[]bool) [][][]float64 {

	meshgrid = *meshgridPointer
	var dist []float64
	var prev []int64
	pq := make(PriorityQueue, 1)

	for i := 0; i < len(meshgrid); i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[flattenIndx(startLngInt, startLatInt, xSize)] = 0
	pq[0] = &Item{
		value:    flattenIndx(startLngInt, startLatInt, xSize),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*Item).value

			if u == flattenIndx(endLngInt, endLatInt, xSize) {
				return extractRoute(&prev, flattenIndx(endLngInt, endLatInt, xSize), xSize, ySize)
			}

			neighbours := neighbours1d(u, xSize)

			for _, j := range neighbours {
				var alt = dist[u] + distance(gridToCoord(expandIndx(u)), gridToCoord(expandIndx(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -dist[j],
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	return extractRoute(&prev, flattenIndx(endLngInt, endLatInt, xSize), xSize, ySize)
}

// DijkstraAllNodes additionally returns all visited nodes
func DijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, xSize, ySize int64, meshgridPointer *[]bool) ([][][]float64, [][]float64) {

	meshgrid = *meshgridPointer
	var dist []float64
	var prev []int64
	var nodesProcessed []int64
	pq := make(PriorityQueue, 1)

	for i := 0; i < len(meshgrid); i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[flattenIndx(startLngInt, startLatInt, xSize)] = 0
	pq[0] = &Item{
		value:    flattenIndx(startLngInt, startLatInt, xSize),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq)

	for {
		if len(pq) == 0 {
			break
		} else {

			u := heap.Pop(&pq).(*Item).value
			nodesProcessed = append(nodesProcessed, u)

			if u == flattenIndx(endLngInt, endLatInt, xSize) {
				var route = extractRoute(&prev, flattenIndx(endLngInt, endLatInt, xSize), xSize, ySize)
				var processedNodes = extractNodes(&nodesProcessed, xSize, ySize)
				return route, processedNodes
			}

			neighbours := neighbours1d(u, xSize)

			for _, j := range neighbours {
				var alt = dist[u] + distance(gridToCoord(expandIndx(u)), gridToCoord(expandIndx(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -dist[j],
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	var route = extractRoute(&prev, flattenIndx(endLngInt, endLatInt, xSize), xSize, ySize)
	var processedNodes = extractNodes(&nodesProcessed, xSize, ySize)
	return route, processedNodes
}
