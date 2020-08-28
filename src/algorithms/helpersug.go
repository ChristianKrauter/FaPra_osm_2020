package algorithms

import (
	"../grids"
	"math"
)

// ExtractRouteUg ...
func ExtractRouteUg(prev *[]int, end int, ug *grids.UniformGrid) [][][]float64 {
	var route = make([][][]float64, 1)
	for {
		route[0] = append(route[0], ug.GridToCoord(ug.IDToGrid(end)))
		if (*prev)[end] == -1 {
			break
		}
		end = (*prev)[end]
	}
	return route
}

// ExtractNodesUg ...
func ExtractNodesUg(nodesProcessed *[]int, ug *grids.UniformGrid) [][]float64 {
	var nodesExtended [][]float64
	for _, node := range *nodesProcessed {
		x := ug.IDToGrid(node)
		coord := ug.GridToCoord([]int{x[0], x[1]})
		nodesExtended = append(nodesExtended, coord)
	}
	return nodesExtended
}

// Gets neighours left and right in the same row
func neighboursRowUg(in []float64, ug *grids.UniformGrid) [][]int {
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

// neighboursUg gets up to 8 neighbours
func neighboursUg(in int, ug *grids.UniformGrid) []int {
	var neighbours [][]int
	var inGrid = ug.IDToGrid(in)
	m := inGrid[0]
	n := inGrid[1]
	neighbours = append(neighbours, []int{m, mod(n-1, len(ug.VertexData[m]))})
	neighbours = append(neighbours, []int{m, mod(n+1, len(ug.VertexData[m]))})

	coord := ug.GridToCoord(inGrid)

	if m > 0 {
		coordDown := ug.GridToCoord([]int{m - 1, n})
		neighbours = append(neighbours, neighboursRowUg([]float64{coord[0], coordDown[1]}, ug)...)
	}

	if m < len(ug.VertexData)-1 {
		coordUp := ug.GridToCoord([]int{m + 1, n})
		neighbours = append(neighbours, neighboursRowUg([]float64{coord[0], coordUp[1]}, ug)...)
	}
	var neighbours1d []int
	for _, neighbour := range neighbours {
		if !ug.VertexData[neighbour[0]][neighbour[1]] {
			neighbours1d = append(neighbours1d, ug.GridToID(neighbour))
		}
	}
	return neighbours1d
}

// ExtractRouteUgBi ...
func ExtractRouteUgBi(prev *[][]int, meeting int, ug *grids.UniformGrid) [][][]float64 {
	var routes = make([][][]float64, 2)

	var secondMeeting = meeting
	for {
		routes[0] = append(routes[0], ug.GridToCoord(ug.IDToGrid(meeting)))
		if (*prev)[meeting][0] == -1 {
			break
		}
		meeting = (*prev)[meeting][0]
	}

	meeting = secondMeeting
	for {
		routes[1] = append(routes[1], ug.GridToCoord(ug.IDToGrid(meeting)))

		if (*prev)[meeting][1] == -1 {

			break
		}
		meeting = (*prev)[meeting][1]
	}
	return routes
}

func hUg(dir, node int, fromIDX, toIDX []int, ug *grids.UniformGrid) float64 {
	if dir == 0 {
		return 0.5 * (distance(ug.GridToCoord(ug.IDToGrid(node)), ug.GridToCoord(toIDX)) - distance(ug.GridToCoord(ug.IDToGrid(node)), ug.GridToCoord(fromIDX)))
	}
	return 0.5 * (distance(ug.GridToCoord(ug.IDToGrid(node)), ug.GridToCoord(fromIDX)) - distance(ug.GridToCoord(ug.IDToGrid(node)), ug.GridToCoord(toIDX)))

}
