package algorithms

import (
	"../grids"
	"math"
	"container/heap"
	//"fmt"
)

func expandNodeDijkstra(node *int,addToQueue bool, ug *grids.UniformGrid, dist *[]float64, prev *[]int, pq *priorityQueue){
	neighbours := NeighboursUg(node, ug)
		for _, j := range neighbours {
			var alt = (*dist)[*node] + distance((*ug).GridToCoord((*ug).IDToGrid(*node)), (*ug).GridToCoord((*ug).IDToGrid(j)))
			if alt < (*dist)[j] {
				(*dist)[j] = alt
				(*prev)[j] = *node
				item := &Item{
					value:    j,
					priority: -(*dist)[j],
				}
				if(addToQueue){
					heap.Push(pq, item)	
				}
			}
		}
}

// ExtractRouteUg ...
func ExtractRouteUg(prev *[]int, end int, ug *grids.UniformGrid) [][][]float64 {
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

func Reverse(numbers [][]float64) [][]float64 {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

//ExtractRouteBiUg..
func ExtractRouteBiUg(prev1,prev2 *[]int, meet int, ug *grids.UniformGrid) [][]float64 {
	var route [][]float64
	var tempRoute [][]float64
	
	end := meet
	for {
		x := ug.IDToGrid(end)
		route = append(route, ug.GridToCoord([]int{x[0], x[1]}))
		if (*prev1)[end] == -1 {
			break
		}
		end = (*prev1)[end]	
	}

	end = (*prev2)[meet]
	for {
		x := ug.IDToGrid(end)
		tempRoute = append(tempRoute, ug.GridToCoord([]int{x[0], x[1]}))
		if (*prev2)[end] == -1 {
			break
		}
		end = (*prev2)[end]	
	}

	route = Reverse(route)
	route = append(route,tempRoute...)

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

// SimpleNeighboursUg gets ug neighbours cheaper
func SimpleNeighboursUg(in int, ug *grids.UniformGrid) []int {
	var neighbours [][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	// lengths of rows
	var lm = len(ug.VertexData[m])

	neighbours = append(neighbours, []int{m, mod(n-1, lm)})
	neighbours = append(neighbours, []int{m, mod(n+1, lm)})

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData)-1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))
		neighbours = append(neighbours, []int{m + 1, mod(nUp, lmp)})
		neighbours = append(neighbours, []int{m + 1, mod(nUp+1.0, lmp)})
		neighbours = append(neighbours, []int{m + 1, mod(nUp-1.0, lmp)})
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))
		neighbours = append(neighbours, []int{m - 1, mod(nDown, lmm)})
		neighbours = append(neighbours, []int{m - 1, mod(nDown+1.0, lmm)})
		neighbours = append(neighbours, []int{m - 1, mod(nDown-1.0, lmm)})
	}

	var neighbours1d []int
	for _, neighbour := range neighbours {
		if !ug.VertexData[neighbour[0]][neighbour[1]] {
			neighbours1d = append(neighbours1d, ug.GridToID(neighbour))
		}
	}
	return neighbours1d
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
func NeighboursUg(in *int, ug *grids.UniformGrid) []int {
	var neighbours [][]int
	var inGrid = ug.IDToGrid(*in)
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

func JPSNeighboursUg(in, dir int, ug *grids.UniformGrid) ([]int, []int) {
	var allNeighbours [8][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	var directions []int

	// lengths of rows
	var lm = len(ug.VertexData[m])

	allNeighbours[7] = []int{m, mod(n-1, lm)}
	allNeighbours[3] = []int{m, mod(n+1, lm)}

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData) -1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))
		allNeighbours[5] = []int{m + 1, mod(nUp, lmp)}
		allNeighbours[4] = []int{m + 1, mod(nUp+1.0, lmp)}
		allNeighbours[6] = []int{m + 1, mod(nUp-1.0, lmp)}
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))
		allNeighbours[1] = []int{m - 1, mod(nDown, lmm)}
		allNeighbours[2] = []int{m - 1, mod(nDown+1.0, lmm)}
		allNeighbours[0] = []int{m - 1, mod(nDown-1.0, lmm)}
	}
	if(dir == -1){
		var n1d []int 
		for _,j := range allNeighbours{
			n1d = append(n1d, ug.GridToID(j))
		}
		return n1d, []int{0,1,2,3,4,5,6,7}
	}
	var neighbours [][]int
	neighbours = append(neighbours, allNeighbours[dir])
	directions = append(directions, dir)

	if(dir == 0 || dir == 2 || dir == 4 || dir == 6){
		neighbours = append(neighbours, allNeighbours[mod(dir-1, 8)])
		neighbours = append(neighbours, allNeighbours[mod(dir+1, 8)])
		directions = append(directions, mod(dir-1, 8))
		directions = append(directions, mod(dir+1, 8))

		check1 := mod(dir-3, 8)
		check2 := mod(dir+3, 8)
		if(ug.VertexData[allNeighbours[check1][0]][allNeighbours[check1][1]]){
			neighbours = append(neighbours, allNeighbours[mod(dir-2, 8)])
			directions = append(directions, mod(dir-2, 8))
		}

		if(ug.VertexData[allNeighbours[check2][0]][allNeighbours[check2][1]]){
			neighbours = append(neighbours, allNeighbours[mod(dir+2, 8)])			
			directions = append(directions, mod(dir+2, 8))
		}
	} else {

		check1 := mod(dir-2, 8)
		check2 := mod(dir+2, 8)
		if(ug.VertexData[allNeighbours[check1][0]][allNeighbours[check1][1]]){
			neighbours = append(neighbours, allNeighbours[mod(dir-1, 8)])
			directions = append(directions, mod(dir-1, 8))		
		}

		if(ug.VertexData[allNeighbours[check2][0]][allNeighbours[check2][1]]){
			neighbours = append(neighbours, allNeighbours[mod(dir+1, 8)])			
			directions = append(directions, mod(dir+1, 8))
		}
	}
	/*
	switch dir {

	case 0:
		neighbours = append(neighbours, allNeighbours[1])
		neighbours = append(neighbours, allNeighbours[7])

		if(ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]]){
			neighbours = append(neighbours, allNeighbours[2])			
		}

		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]]){
			neighbours = append(neighbours, allNeighbours[6])			
		}

	case 1:
		if(ug.VertexData[allNeighbours[7][0]][allNeighbours[7][1]]){
			neighbours = append(neighbours, allNeighbours[0])			
		}

		if(ug.VertexData[allNeighbours[4][0]][allNeighbours[4][1]]){
			neighbours = append(neighbours, allNeighbours[2])			
		}

	case 2:
		neighbours = append(neighbours, allNeighbours[1])
		neighbours = append(neighbours, allNeighbours[4])

		if(ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]]){
			neighbours = append(neighbours, allNeighbours[0])			
		}

		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]]){
			neighbours = append(neighbours, allNeighbours[7])			
		}

	case 3:
		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]]){
			neighbours = append(neighbours, allNeighbours[0])			
		}

		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]]){
			neighbours = append(neighbours, allNeighbours[5])			
		}

	case 4:
		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]]){
			neighbours = append(neighbours, allNeighbours[2])			
		}

		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]]){
			neighbours = append(neighbours, allNeighbours[7])			
		}

	case 5:	
		neighbours = append(neighbours, allNeighbours[3])
		neighbours = append(neighbours, allNeighbours[6])

		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]]){
			neighbours = append(neighbours, allNeighbours[0])			
		}

		if(ug.VertexData[allNeighbours[4][0]][allNeighbours[4][1]]){
			neighbours = append(neighbours, allNeighbours[7])			
		}
	case 6:
		if(ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]]){
			neighbours = append(neighbours, allNeighbours[5])			
		}

		if(ug.VertexData[allNeighbours[4][0]][allNeighbours[4][1]]){
			neighbours = append(neighbours, allNeighbours[7])			
		}
	case 7:
		neighbours = append(neighbours, allNeighbours[4])
		neighbours = append(neighbours, allNeighbours[6])

		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]]){
			neighbours = append(neighbours, allNeighbours[2])			
		}

		if(ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]]){
			neighbours = append(neighbours, allNeighbours[5])			
		}
	}*/

	var neighbours1d []int
	for _, neighbour := range neighbours {
		if !ug.VertexData[neighbour[0]][neighbour[1]] {
			neighbours1d = append(neighbours1d, ug.GridToID(neighbour))
		}
	}
	return neighbours1d,directions
}

func isForced(in, dir int, ug *grids.UniformGrid) bool {
	var allNeighbours [8][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	// lengths of rows
	var lm = len(ug.VertexData[m])

	allNeighbours[7] = []int{m, mod(n-1, lm)}
	allNeighbours[3] = []int{m, mod(n+1, lm)}

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData) -1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))
		allNeighbours[5] = []int{m + 1, mod(nUp, lmp)}
		allNeighbours[4] = []int{m + 1, mod(nUp+1.0, lmp)}
		allNeighbours[6] = []int{m + 1, mod(nUp-1.0, lmp)}
	} else {
		return false
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))
		allNeighbours[1] = []int{m - 1, mod(nDown, lmm)}
		allNeighbours[2] = []int{m - 1, mod(nDown+1.0, lmm)}
		allNeighbours[0] = []int{m - 1, mod(nDown-1.0, lmm)}
	}

	if(dir == 0 || dir == 2 || dir == 4 || dir == 6){
		if ug.VertexData[allNeighbours[mod(dir-3, 8)][0]][allNeighbours[mod(dir-3, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir-2, 8)][0]][allNeighbours[mod(dir-2, 8)][1]] {
			return true
		}
		if ug.VertexData[allNeighbours[mod(dir+3, 8)][0]][allNeighbours[mod(dir+3, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir-2, 8)][0]][allNeighbours[mod(dir-2, 8)][1]] {
			return true
		}
	} else {
		if ug.VertexData[allNeighbours[mod(dir-2, 8)][0]][allNeighbours[mod(dir-2, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir-1, 8)][0]][allNeighbours[mod(dir-1, 8)][1]] {
			return true
		}
		if ug.VertexData[allNeighbours[mod(dir+2, 8)][0]][allNeighbours[mod(dir+2, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir-1, 8)][0]][allNeighbours[mod(dir-1, 8)][1]] {
			return true
		}
	}

	/*
	switch dir {
	case 0:
		if(ug.VertexData[allNeighbours[4][0]][allNeighbours[4][1]] ||!ug.VertexData[allNeighbours[2][0]][allNeighbours[2][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]] ||!ug.VertexData[allNeighbours[5][0]][allNeighbours[5][1]] ){
			return true
		}
		
	case 1:
		if(ug.VertexData[allNeighbours[5][0]][allNeighbours[5][1]] ||!ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[7][0]][allNeighbours[7][1]] ||!ug.VertexData[allNeighbours[4][0]][allNeighbours[4][1]] ){
			return true
		}

	case 2:
		if(ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]] ||!ug.VertexData[allNeighbours[0][0]][allNeighbours[0][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]] ||!ug.VertexData[allNeighbours[7][0]][allNeighbours[7][1]] ){
			return true
		}

	case 3:
		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]] ||!ug.VertexData[allNeighbours[0][0]][allNeighbours[0][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]] ||!ug.VertexData[allNeighbours[5][0]][allNeighbours[5][1]] ){
			return true
		}

	case 4:
		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]] ||!ug.VertexData[allNeighbours[2][0]][allNeighbours[2][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[6][0]][allNeighbours[6][1]] ||!ug.VertexData[allNeighbours[7][0]][allNeighbours[7][1]] ){
			return true
		}

	case 5:	
		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]] ||!ug.VertexData[allNeighbours[0][0]][allNeighbours[0][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[4][0]][allNeighbours[4][1]] ||!ug.VertexData[allNeighbours[7][0]][allNeighbours[7][1]] ){
			return true
		}
	case 6:
		if(ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]] ||!ug.VertexData[allNeighbours[5][0]][allNeighbours[5][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[4][0]][allNeighbours[4][1]] ||!ug.VertexData[allNeighbours[7][0]][allNeighbours[7][1]] ){
			return true
		}
	case 7:
		if(ug.VertexData[allNeighbours[1][0]][allNeighbours[1][1]] ||!ug.VertexData[allNeighbours[2][0]][allNeighbours[2][1]] ){
			return true
		}
		if(ug.VertexData[allNeighbours[3][0]][allNeighbours[3][1]] ||!ug.VertexData[allNeighbours[5][0]][allNeighbours[5][1]] ){
			return true
		}
	}*/
	return false
}

