package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// Dijkstra implementation on uniform grid
func Dijkstra(startLngInt, startLatInt, endLngInt, endLatInt int, ug *grids.UniformGrid) [][][]float64 {

	var dist []float64
	var prev []int
	pq := make(priorityQueue, 1)

	for i := 0; i < (*ug).N; i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[(*ug).GridToID(startLngInt, startLatInt)] = 0
	pq[0] = &Item{
		value:    (*ug).GridToID(startLngInt, startLatInt),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*Item).value

			if u == (*ug).GridToID(endLngInt, endLatInt) {
				return ExtractRouteUg(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
			}

			neighbours := NeighboursUg(u, ug)
			for _, j := range neighbours {
				var alt = dist[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
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
	return ExtractRouteUg(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
}

// DijkstraAllNodes additionally returns all visited nodes on uniform grid
func DijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {

	var dist []float64
	var prev []int
	var nodesProcessed []int
	pq := make(priorityQueue, 1)

	for i := 0; i < (*ug).N; i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[(*ug).GridToID(startLngInt, startLatInt)] = 0
	pq[0] = &Item{
		value:    (*ug).GridToID(startLngInt, startLatInt),
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

			if u == (*ug).GridToID(endLngInt, endLatInt) {
				var route = ExtractRouteUg(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
				var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
				return route, processedNodes
			}

			neighbours := NeighboursUg(u, ug)

			for _, j := range neighbours {
				var alt = dist[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
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
	var route = ExtractRouteUg(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
	var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
	return route, processedNodes
}
