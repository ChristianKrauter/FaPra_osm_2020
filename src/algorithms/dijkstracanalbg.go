package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// DijkstraCanalBg implementation on basic grid
func DijkstraCanalBg(from, to int, bg *grids.BasicGrid) ([][][]int, int, float64) {
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
				return extractRouteCanal(&prev, to, bg), popped, dist[to]
			}

			neighbours := NeighboursCanalBg(u, bg)
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
	return extractRouteCanal(&prev, to, bg), popped, dist[to]
}
