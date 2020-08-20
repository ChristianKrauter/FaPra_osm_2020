package algorithms

import (
	"container/heap"
	//"fmt"
	"math"
)

// UniformDijkstra implementation on uniform grid
func UniformDijkstra(startLngInt, startLatInt, endLngInt, endLatInt int, ug *UniformGrid) [][][]float64 {

	var dist []float64
	var prev []int
	pq := make(PriorityQueue, 1)

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
				return UniformExtractRoute(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
			}

			neighbours := GetNeighboursUniformGrid(u, ug)
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
	return UniformExtractRoute(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
}

// UniformDijkstraAllNodes additionally returns all visited nodes on uniform grid
func UniformDijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt int, ug *UniformGrid) ([][][]float64, [][]float64) {

	var dist []float64
	var prev []int
	var nodesProcessed []int
	pq := make(PriorityQueue, 1)

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
				var route = UniformExtractRoute(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
				var processedNodes = UniformExtractNodes(&nodesProcessed, ug)
				return route, processedNodes
			}

			neighbours := GetNeighboursUniformGrid(u, ug)

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
	var route = UniformExtractRoute(&prev, (*ug).GridToID(endLngInt, endLatInt), ug)
	var processedNodes = UniformExtractNodes(&nodesProcessed, ug)
	return route, processedNodes
}
