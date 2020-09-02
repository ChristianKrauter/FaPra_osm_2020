package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// AStarBg implementation
func AStarBg(fromIDX, toIDX []int, bg *grids.BasicGrid) ([][][]float64, int) {

	var popped int
	var dist []float64
	var fScore []float64
	var prev []int
	pq := make(priorityQueue, 1)

	for i := 0; i < len(bg.VertexData); i++ {
		dist = append(dist, math.Inf(1))
		fScore = append(fScore, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[bg.FlattenIndex(fromIDX)] = 0
	fScore[bg.FlattenIndex(fromIDX)] = distance(bg.GridToCoord(fromIDX), bg.GridToCoord(toIDX))

	pq[0] = &Item{
		value:    bg.FlattenIndex(fromIDX),
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

			if u == bg.FlattenIndex(toIDX) {
				return extractRoute(&prev, bg.FlattenIndex(toIDX), bg), popped
			}

			neighbours := neighboursBg(u, bg)

			for _, j := range neighbours {
				var alt = dist[u] + distance(bg.GridToCoord(bg.ExpandIndex(u)), bg.GridToCoord(bg.ExpandIndex(j)))
				if alt < dist[j] {
					dist[j] = alt
					fScore[j] = dist[j] + distance(bg.GridToCoord(bg.ExpandIndex(j)), bg.GridToCoord(toIDX))
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -fScore[j],
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	return extractRoute(&prev, bg.FlattenIndex(toIDX), bg), popped
}

// AStarAllNodesBg additionally returns all visited nodes
func AStarAllNodesBg(fromIDX, toIDX []int, bg *grids.BasicGrid) ([][][]float64, [][]float64) {

	var dist []float64
	var fScore []float64
	var prev []int
	var nodesProcessed []int
	pq := make(priorityQueue, 1)

	for i := 0; i < len(bg.VertexData); i++ {
		dist = append(dist, math.Inf(1))
		fScore = append(fScore, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[bg.FlattenIndex(fromIDX)] = 0
	fScore[bg.FlattenIndex(fromIDX)] = distance(bg.GridToCoord(fromIDX), bg.GridToCoord(toIDX))

	pq[0] = &Item{
		value:    bg.FlattenIndex(fromIDX),
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

			if u == bg.FlattenIndex(toIDX) {
				var route = extractRoute(&prev, bg.FlattenIndex(toIDX), bg)
				var processedNodes = extractNodes(&nodesProcessed, bg)
				return route, processedNodes
			}

			neighbours := neighboursBg(u, bg)

			for _, j := range neighbours {
				var alt = dist[u] + distance(bg.GridToCoord(bg.ExpandIndex(u)), bg.GridToCoord(bg.ExpandIndex(j)))
				if alt < dist[j] {
					dist[j] = alt
					fScore[j] = dist[j] + distance(bg.GridToCoord(bg.ExpandIndex(j)), bg.GridToCoord(toIDX))
					prev[j] = u
					item := &Item{
						value:    j,
						priority: -fScore[j],
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	var route = extractRoute(&prev, bg.FlattenIndex(toIDX), bg)
	var processedNodes = extractNodes(&nodesProcessed, bg)
	return route, processedNodes
}