package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// Dijkstra implementation on uniform grid
func Dijkstra(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, int) {

	var popped int
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	pq := make(priorityQueue, 1)

	for i := 0; i < ug.N; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[ug.GridToID(fromIDX)] = 0
	pq[0] = &Item{
		value:    ug.GridToID(fromIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*Item).value
			popped++
			if u == ug.GridToID(toIDX) {
				return ExtractRouteUg(&prev, ug.GridToID(toIDX), ug), popped
			}

			neighbours := neighboursUg(u, ug)
			for _, j := range neighbours {
				var alt = dist[u] + distance(ug.GridToCoord(ug.IDToGrid(u)), ug.GridToCoord(ug.IDToGrid(j)))
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
	return ExtractRouteUg(&prev, ug.GridToID(toIDX), ug), popped
}

// DijkstraAllNodes additionally returns all visited nodes on uniform grid
func DijkstraAllNodes(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {

	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	var nodesProcessed []int
	pq := make(priorityQueue, 1)

	for i := 0; i < ug.N; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[ug.GridToID(fromIDX)] = 0
	pq[0] = &Item{
		value:    ug.GridToID(fromIDX),
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

			if u == ug.GridToID(toIDX) {
				var route = ExtractRouteUg(&prev, ug.GridToID(toIDX), ug)
				var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
				return route, processedNodes
			}

			neighbours := neighboursUg(u, ug)

			for _, j := range neighbours {
				var alt = dist[u] + distance(ug.GridToCoord(ug.IDToGrid(u)), ug.GridToCoord(ug.IDToGrid(j)))
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
	var route = ExtractRouteUg(&prev, ug.GridToID(toIDX), ug)
	var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
	return route, processedNodes
}
