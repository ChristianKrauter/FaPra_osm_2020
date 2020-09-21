package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// AStarBg implementation
func AStarBg(from, to int, bg *grids.BasicGrid) ([][][]float64, int, float64) {
	var popped int
	var dist = make([]float64, len(bg.VertexData))
	var prev = make([]int, len(bg.VertexData))
	pq := make(priorityQueue, 1)

	for i := 0; i < len(bg.VertexData); i++ {
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
				return extractRoute(&prev, to, bg), popped, dist[to]
			}

			neighbours := neighboursBg(u, bg)

			for _, j := range neighbours {
				var alt = dist[u] + distance(bg.GridToCoord(bg.IDToGrid(u)), bg.GridToCoord(bg.IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -(dist[j] + distance(bg.GridToCoord(bg.IDToGrid(j)), bg.GridToCoord(bg.IDToGrid(to)))),
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	return extractRoute(&prev, to, bg), popped, dist[to]
}

// AStarAllNodesBg additionally returns all visited nodes
func AStarAllNodesBg(from, to int, bg *grids.BasicGrid) ([][][]float64, [][]float64, float64) {
	var dist = make([]float64, len(bg.VertexData))
	var prev = make([]int, len(bg.VertexData))
	var nodesProcessed []int
	pq := make(priorityQueue, 1)

	for i := 0; i < len(bg.VertexData); i++ {
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
				return extractRoute(&prev, to, bg), extractNodes(&nodesProcessed, bg), dist[to]
			}

			neighbours := neighboursBg(u, bg)

			for _, j := range neighbours {
				var alt = dist[u] + distance(bg.GridToCoord(bg.IDToGrid(u)), bg.GridToCoord(bg.IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -(dist[j] + distance(bg.GridToCoord(bg.IDToGrid(j)), bg.GridToCoord(bg.IDToGrid(to)))),
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	return extractRoute(&prev, to, bg), extractNodes(&nodesProcessed, bg), dist[to]
}
