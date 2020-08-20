package algorithms

import (
	"container/heap"
	//"fmt"
	"math"
)

var meshgrid []bool

// Dijkstra implementation
func Dijkstra(startLngInt, startLatInt, endLngInt, endLatInt, xSize, ySize int, meshgridPointer *[]bool) [][][]float64 {

	meshgrid = *meshgridPointer
	var dist []float64
	var prev []int
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
func DijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, xSize, ySize int, meshgridPointer *[]bool) ([][][]float64, [][]float64) {

	meshgrid = *meshgridPointer
	var dist []float64
	var prev []int
	var nodesProcessed []int
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
