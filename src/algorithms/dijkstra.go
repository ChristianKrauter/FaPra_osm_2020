package algorithms

import (
	"container/heap"
	"math"
)

var meshgrid []bool

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

func extractRoute(prev *[]int64, end, xSize, ySize int64) [][][]float64 {
	var route [][][]float64
	var tempRoute [][]float64
	temp := ExpandIndex(end, xSize)
	for {
		x := ExpandIndex(end, xSize)
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			route = append(route, tempRoute)
			tempRoute = make([][]float64, 0)
		}
		tempRoute = append(tempRoute, GridToCoord([]int64{x[0], x[1]}, xSize, ySize))

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
		x := ExpandIndex(node, xSize)
		coord := GridToCoord([]int64{x[0], x[1]}, xSize, ySize)
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

	dist[flattenIndex(startLngInt, startLatInt, xSize)] = 0
	pq[0] = &Item{
		value:    flattenIndex(startLngInt, startLatInt, xSize),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*Item).value

			if u == flattenIndex(endLngInt, endLatInt, xSize) {
				return extractRoute(&prev, flattenIndex(endLngInt, endLatInt, xSize), xSize, ySize)
			}

			neighbours := neighbours1d(u, xSize)

			for _, j := range neighbours {
				var alt = dist[u] + distance(GridToCoord(ExpandIndex(u, xSize), xSize, ySize), GridToCoord(ExpandIndex(j, xSize), xSize, ySize))
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
	return extractRoute(&prev, flattenIndex(endLngInt, endLatInt, xSize), xSize, ySize)
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

	dist[flattenIndex(startLngInt, startLatInt, xSize)] = 0
	pq[0] = &Item{
		value:    flattenIndex(startLngInt, startLatInt, xSize),
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

			if u == flattenIndex(endLngInt, endLatInt, xSize) {
				var route = extractRoute(&prev, flattenIndex(endLngInt, endLatInt, xSize), xSize, ySize)
				var processedNodes = extractNodes(&nodesProcessed, xSize, ySize)
				return route, processedNodes
			}

			neighbours := neighbours1d(u, xSize)

			for _, j := range neighbours {
				var alt = dist[u] + distance(GridToCoord(ExpandIndex(u, xSize), xSize, ySize), GridToCoord(ExpandIndex(j, xSize), xSize, ySize))
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
	var route = extractRoute(&prev, flattenIndex(endLngInt, endLatInt, xSize), xSize, ySize)
	var processedNodes = extractNodes(&nodesProcessed, xSize, ySize)
	return route, processedNodes
}
