package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// Dijkstra implementation on uniform grid
func Dijkstra(from, to int, ug *grids.UniformGrid) (*[][][]float64, int, float64) {
	var popped int
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	var pq = make(priorityQueue, 1)

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
				return extractRouteUg(&prev, to, ug), popped, dist[to]
			}

			neighbours := SimpleNeighboursUg(u, ug)
			uCoord := ug.GridToCoord(ug.IDToGrid(u))
			for _, j := range neighbours {
				alt := dist[u] + distance(uCoord, ug.GridToCoord(ug.IDToGrid(j)))
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
	return extractRouteUg(&prev, to, ug), popped, dist[to]
}

// DijkstraAllNodes also returns visited nodes on uniform grid
func DijkstraAllNodes(from, to int, ug *grids.UniformGrid) (*[][][]float64, *[][]float64, float64) {
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	var nodesProcessed []int
	var pq = make(priorityQueue, 1)

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
				return extractRouteUg(&prev, to, ug), extractNodesUg(&nodesProcessed, ug), dist[to]
			}

			neighbours := SimpleNeighboursUg(u, ug)
			uCoord := ug.GridToCoord(ug.IDToGrid(u))
			for _, j := range neighbours {
				alt := dist[u] + distance(uCoord, ug.GridToCoord(ug.IDToGrid(j)))
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
	return extractRouteUg(&prev, to, ug), extractNodesUg(&nodesProcessed, ug), dist[to]
}
