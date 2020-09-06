package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// AStarBg implementation
func AStarBg(from, to int, bg *grids.BasicGrid) ([][][]float64, int) {

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
				return extractRoute(&prev, to, bg), popped
			}

			neighbours := neighboursBg(u, bg)

			for _, j := range neighbours {
				var alt = dist[u] + distance(bg.GridToCoord(bg.ExpandIndex(u)), bg.GridToCoord(bg.ExpandIndex(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -(dist[j] + distance(bg.GridToCoord(bg.ExpandIndex(j)), bg.GridToCoord(bg.ExpandIndex(to)))),
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	return extractRoute(&prev, to, bg), popped
}

// AStarAllNodesBg additionally returns all visited nodes
func AStarAllNodesBg(from, to int, bg *grids.BasicGrid) ([][][]float64, [][]float64) {

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
				var route = extractRoute(&prev, to, bg)
				var processedNodes = extractNodes(&nodesProcessed, bg)
				return route, processedNodes
			}

			neighbours := neighboursBg(u, bg)

			for _, j := range neighbours {
				var alt = dist[u] + distance(bg.GridToCoord(bg.ExpandIndex(u)), bg.GridToCoord(bg.ExpandIndex(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -(dist[j] + distance(bg.GridToCoord(bg.ExpandIndex(j)), bg.GridToCoord(bg.ExpandIndex(to)))),
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	var route = extractRoute(&prev, to, bg)
	var processedNodes = extractNodes(&nodesProcessed, bg)
	return route, processedNodes
}
