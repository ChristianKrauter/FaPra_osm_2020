package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// DijkstraBg implementation on basic grid
func DijkstraBg(from, to int, bg *grids.BasicGrid) ([][][]float64, int, float64) {
	var popped int
	var dist = make([]float64, len(bg.VertexData))
	var prev = make([]int, len(bg.VertexData))
	var pq = make(priorityQueue, 1)

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

			neighbours := NeighboursBg(u, bg)
			uCoord := bg.GridToCoord(bg.IDToGrid(u))
			for _, j := range neighbours {
				alt := dist[u] + distance(uCoord, bg.GridToCoord(bg.IDToGrid(j)))
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
	return extractRoute(&prev, to, bg), popped, dist[to]
}

// DijkstraAllNodesBg also returns visited nodes on basic grid
func DijkstraAllNodesBg(from, to int, bg *grids.BasicGrid) ([][][]float64, [][]float64, float64) {
	var dist = make([]float64, len(bg.VertexData))
	var prev = make([]int, len(bg.VertexData))
	var nodesProcessed []int
	var pq = make(priorityQueue, 1)

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

			neighbours := NeighboursBg(u, bg)
			uCoord := bg.GridToCoord(bg.IDToGrid(u))
			for _, j := range neighbours {
				alt := dist[u] + distance(uCoord, bg.GridToCoord(bg.IDToGrid(j)))
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
	return extractRoute(&prev, to, bg), extractNodes(&nodesProcessed, bg), dist[to]
}
