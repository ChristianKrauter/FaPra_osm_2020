package algorithms

import (
	"math"
	"sort"
)

// UniformGrid structure
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

// GridToCoord takes grid indices and outputs lng lat
func (ug UniformGrid) GridToCoord(in []int) []float64 {
	theta := math.Pi * (float64(in[0]) + 0.5) / float64(ug.MTheta)
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
	phi := 2 * math.Pi * float64(in[1]) / mPhi
	return []float64{(phi / math.Pi) * 180, (theta/math.Pi)*180 - 90}
}

// CoordToGrid takes lng lat and outputs grid indices
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
func UniformExtractRoute(prev *[]int, end int, ug *UniformGrid) [][][]float64 {
	var route [][][]float64
	var tempRoute [][]float64
	temp := ug.IDToGrid(end)
	for {
		x := ug.IDToGrid(end)
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			route = append(route, tempRoute)
			tempRoute = make([][]float64, 0)
		}
		tempRoute = append(tempRoute, ug.GridToCoord([]int{x[0], x[1]}))

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
func UniformExtractNodes(nodesProcessed *[]int, ug *UniformGrid) [][]float64 {
	var nodesExtended [][]float64
	for _, node := range *nodesProcessed {
		x := ug.IDToGrid(node)
		coord := ug.GridToCoord([]int{x[0], x[1]})
		nodesExtended = append(nodesExtended, coord)
	}
	return nodesExtended
}

// Gets neighours left and right in the same row
func uniformNeighboursRow(in []float64, ug *UniformGrid) [][]int {
	// Test if it still works with less than 3 points in one grid row
	theta := (in[1] + 90) * math.Pi / 180
	m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
	theta = math.Pi * (m + 0.5) / ug.MTheta
	phi := in[0] * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)

	n1 := math.Round(phi * mPhi / (2 * math.Pi))
	p1 := []int{mod(int(m), int(ug.MTheta)), mod(int(n1), int(mPhi))}
	p2 := []int{mod(int(m), int(ug.MTheta)), mod(int(n1+1), int(mPhi))}
	p3 := []int{mod(int(m), int(ug.MTheta)), mod(int(n1-1), int(mPhi))}
	return [][]int{p1, p2, p3}
}

// GetNeighboursUniformGrid gets up to 8 neighbours
func GetNeighboursUniformGrid(in int, ug *UniformGrid) []int {
	var neighbours [][]int
	var inGrid = ug.IDToGrid(in)
	m := inGrid[0]
	n := inGrid[1]
	neighbours = append(neighbours, []int{m, mod(n-1, len(ug.VertexData[m]))})
	neighbours = append(neighbours, []int{m, mod(n+1, len(ug.VertexData[m]))})

	coord := ug.GridToCoord(inGrid)

	if m > 0 {
		coordDown := ug.GridToCoord([]int{m - 1, n})
		neighbours = append(neighbours, uniformNeighboursRow([]float64{coord[0], coordDown[1]}, ug)...)
	}

	if m < len(ug.VertexData)-1 {
		coordUp := ug.GridToCoord([]int{m + 1, n})
		neighbours = append(neighbours, uniformNeighboursRow([]float64{coord[0], coordUp[1]}, ug)...)
	}
	var neighbours1d []int
	for _, neighbour := range neighbours {
		if !ug.VertexData[neighbour[0]][neighbour[1]] {
			neighbours1d = append(neighbours1d, ug.GridToID(neighbour[0], neighbour[1]))
		}
	}
	return neighbours1d
}
