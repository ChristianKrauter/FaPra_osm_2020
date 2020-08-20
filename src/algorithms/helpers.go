package algorithms

import (
	"math"
	"sort"
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

// ArrayOfArrays sorting by length
type ArrayOfArrays [][][]float64

func (p ArrayOfArrays) Len() int {
	return len(p)
}

func (p ArrayOfArrays) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ArrayOfArrays) Less(i, j int) bool {
	return len(p[i]) > len(p[j])
}

// UniformGrid ...
type UniformGrid struct {
	XSize        int
	YSize        int
	N            int
	BigN         int
	A            float64
	D            float64
	MTheta       float64
	DTheta       float64
	DPhi         float64
	VertexData   [][]bool
	FirstIndexOf []int
}

// GridToCoord ...
func (ug UniformGrid) GridToCoord(in []int) []float64 {
	theta := math.Pi * (float64(in[0]) + 0.5) / float64(ug.MTheta)
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
	phi := 2 * math.Pi * float64(in[1]) / mPhi
	return []float64{(phi / math.Pi) * 180, (theta/math.Pi)*180 - 90}
}

// CoordToGrid returns grid coordinates given lng,lat
func (ug UniformGrid) CoordToGrid(lng, lat float64) []int {
	theta := (lat + 90) * math.Pi / 180
	m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
	phi := lng * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
	n := math.Round(phi * mPhi / (2 * math.Pi))
	return []int{mod(int(m), int(ug.MTheta)), mod(int(n), int(mPhi))}
}

// GridToID ...
func (ug UniformGrid) GridToID(m, n int) int {
	return ug.FirstIndexOf[m] + n
}

// IDToGrid ...
func (ug UniformGrid) IDToGrid(id int) []int {
	m := sort.Search(len(ug.FirstIndexOf)-1, func(i int) bool { return ug.FirstIndexOf[i] > id })
	n := id - ug.FirstIndexOf[m-1]
	return []int{m - 1, n}
}


// UniformExtractRoute ...
func UniformExtractRoute(prev *[]int, end, xSize, ySize int, uniformGrid *UniformGrid) [][][]float64 {
	var route [][][]float64
	var tempRoute [][]float64
	temp := uniformGrid.IDToGrid(end)
	for {
		x := uniformGrid.IDToGrid(end)
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			route = append(route, tempRoute)
			tempRoute = make([][]float64, 0)
		}
		tempRoute = append(tempRoute, uniformGrid.GridToCoord([]int{x[0], x[1]}))

		if (*prev)[end] == -1 {
			break
		}
		end = (*prev)[end]
		temp = x
	}
	route = append(route, tempRoute)
	return route
}

// UniformExtractNodes ...
func UniformExtractNodes(nodesProcessed *[]int, xSize, ySize int, uniformGrid *UniformGrid) [][]float64 {
	var nodesExtended [][]float64
	for _, node := range *nodesProcessed {
		x := uniformGrid.IDToGrid(node)
		coord := uniformGrid.GridToCoord([]int{x[0], x[1]})
		nodesExtended = append(nodesExtended, coord)
	}
	return nodesExtended
}

// Test if it still works with less than 3 points in one grid row
func uniformNeighboursRow(in []float64, xSize, ySize int) [][]int {
	N := float64(xSize * ySize)
	a := 4.0 * math.Pi / N
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta

	theta := (in[1] + 90) * math.Pi / 180
	m := math.Round((theta * mTheta / math.Pi) - 0.5)

	theta = math.Pi * (m + 0.5) / mTheta

	phi := in[0] * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
	n1 := math.Round(phi * mPhi / (2 * math.Pi))
	n2 := n1 + 1
	n3 := n1 - 1
	p1 := []int{mod(int(m), int(mTheta)), mod(int(n1), int(mPhi))}
	p2 := []int{mod(int(m), int(mTheta)), mod(int(n2), int(mPhi))}
	p3 := []int{mod(int(m), int(mTheta)), mod(int(n3), int(mPhi))}
	return [][]int{p1, p2, p3}
}

// GetNeighboursUniformGrid ...
func GetNeighboursUniformGrid(in, xSize, ySize int, uniformGrid *UniformGrid) []int {
	var neighbours [][]int
	var inGrid = uniformGrid.IDToGrid(in)
	m := inGrid[0]
	n := inGrid[1]
	neighbours = append(neighbours, []int{m, mod(n-1, len(uniformGrid.VertexData[m]))})
	neighbours = append(neighbours, []int{m, mod(n+1, len(uniformGrid.VertexData[m]))})

	coord := uniformGrid.GridToCoord(inGrid)

	if m > 0 {
		coordDown := uniformGrid.GridToCoord([]int{m - 1, n})
		neighbours = append(neighbours, uniformNeighboursRow([]float64{coord[0], coordDown[1]}, xSize, ySize)...)
	}

	if m < len(uniformGrid.VertexData)-1 {
		coordUp := uniformGrid.GridToCoord([]int{m + 1, n})
		neighbours = append(neighbours, uniformNeighboursRow([]float64{coord[0], coordUp[1]}, xSize, ySize)...)
	}
	var neighbours1d []int
	for _, neighbour := range neighbours {
		if !uniformGrid.VertexData[neighbour[0]][neighbour[1]] {
			neighbours1d = append(neighbours1d, uniformGrid.GridToID(neighbour[0], neighbour[1]))
		}
	}
	return neighbours1d
}
