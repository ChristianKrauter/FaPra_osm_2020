package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// AStar implementation on uniform grid
func AStar(from, to int, ug *grids.UniformGrid) ([][][]float64, int) {
	var popped int
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	pq := make(priorityQueue, 1)

	for i := 0; i < ug.N; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[from] = 0
	pq[0] = &Item{
		value:    from,
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
			if u == to {
				return ExtractRouteUg(&prev, to, ug), popped
			}

			neighbours := NeighboursUg(u, ug)
			for _, j := range neighbours {
				var alt = dist[u] + distance(ug.GridToCoord(ug.IDToGrid(u)), ug.GridToCoord(ug.IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -(dist[j] + distance(ug.GridToCoord(ug.IDToGrid(j)), ug.GridToCoord(ug.IDToGrid(to)))),
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	return ExtractRouteUg(&prev, to, ug), popped
}

// AStarAllNodes additionally returns all visited nodes on uniform grid
func AStarAllNodes(from, to int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	var nodesProcessed []int
	pq := make(priorityQueue, 1)

	for i := 0; i < ug.N; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[from] = 0

	pq[0] = &Item{
		value:    from,
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

			if u == to {
				var route = ExtractRouteUg(&prev, to, ug)
				var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
				return route, processedNodes
			}

			neighbours := NeighboursUg(u, ug)

			for _, j := range neighbours {
				var alt = dist[u] + distance(ug.GridToCoord(ug.IDToGrid(u)), ug.GridToCoord(ug.IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -(dist[j] + distance(ug.GridToCoord(ug.IDToGrid(j)), ug.GridToCoord(ug.IDToGrid(to)))),
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	var route = ExtractRouteUg(&prev, to, ug)
	var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
	return route, processedNodes
}
