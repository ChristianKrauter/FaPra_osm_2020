package algorithms

import (
	"../grids"
	"fmt"
	"math"
)

// ExtractRouteUg ...
func ExtractRouteUg(prev *[]int, end int, ug *grids.UniformGrid) [][][]float64 {
	var route [][][]float64
	var tempRoute [][]float64
	temp := ug.IDToGrid(end)
	for {
		x := ug.IDToGrid(end)
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			fmt.Printf("helpersUG\n")
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

// NeighboursUg gets up to 8 neighbours
func NeighboursUg(in int, ug *grids.UniformGrid) []int {
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
func ExtractRouteUgBi(prev *[][]int, end int, ug *grids.UniformGrid) [][][]float64 {
	var routes = make([][][]float64, 2)

	var secondEnd = end
	fmt.Printf("prev[end][0] %v\n", (*prev)[end][0])
	fmt.Printf("prev[end][1] %v\n", (*prev)[end][1])
	for {
		routes[0] = append(routes[0], ug.GridToCoord(ug.IDToGrid(end)))
		if (*prev)[end][0] == -1 {
			break
		}
		end = (*prev)[end][0]
	}

	end = secondEnd
	for {
		routes[1] = append(routes[1], ug.GridToCoord(ug.IDToGrid(end)))

		if (*prev)[end][1] == -1 {

			break
		}
		end = (*prev)[end][1]
	}
	fmt.Printf("route %v\n", routes)
	return routes
}
