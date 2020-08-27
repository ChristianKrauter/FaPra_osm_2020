package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// BiDijkstra implementation on uniform grid
func BiDijkstraBg(fromIDX, toIDX []int, bg *grids.BasicGrid) ([][][]float64, int) {
	var prev [][]int
	dist := make([][]float64, 2)
	pq := []priorityQueue{make(priorityQueue, 1), make(priorityQueue, 1)}
	proc := []map[int]bool{make(map[int]bool), make(map[int]bool)}
	var meeting int

	// Init
	for i := 0; i < len(bg.VertexData); i++ {
		dist[0] = append(dist[0], math.Inf(1))
		dist[1] = append(dist[1], math.Inf(1))
		prev = append(prev, []int{-1, -1})
	}

	dist[0][bg.FlattenIndex(fromIDX)] = 0
	pq[0][0] = &Item{
		value:    bg.FlattenIndex(fromIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[0])

	dist[1][bg.FlattenIndex(toIDX)] = 0
	pq[1][0] = &Item{
		value:    bg.FlattenIndex(toIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[1])

	var dir = 0 // Direction in which next step should be taken. 0=forward, 1=backward
	// Main loop
	for {
		if len(pq[0]) == 0 || len(pq[1]) == 0 {
			break
		} else {
			u := heap.Pop(&pq[dir]).(*Item).value
			proc[dir][u] = true

			if proc[1-dir][u] {
				var bestDist = dist[0][u] + dist[1][u]
				meeting = u

				for _, k := range pq[0] {
					if proc[1][k.value] {
						if dist[0][k.value]+dist[1][k.value] < bestDist {
							bestDist = dist[0][k.value] + dist[1][k.value]
							meeting = k.value
						}
					}
				}

				for _, k := range pq[1] {
					if proc[0][k.value] {
						if dist[0][k.value]+dist[1][k.value] < bestDist {
							bestDist = dist[0][k.value] + dist[1][k.value]
							meeting = k.value
						}
					}
				}
				break
			}

			neighbours := neighboursBg(u, bg)
			for _, j := range neighbours {
				var alt = dist[dir][u] + distance(bg.GridToCoord(bg.ExpandIndex(u)), bg.GridToCoord(bg.ExpandIndex(j)))
				if alt < dist[dir][j] {
					dist[dir][j] = alt
					prev[j][dir] = u
					item := &Item{
						value:    j,
						priority: -dist[dir][j],
					}
					heap.Push(&pq[dir], item)
				}
			}
		}
		dir = 1 - dir // Change direction
	}

	var route = ExtractRouteBi(&prev, meeting, bg)
	return route, len(proc[0]) + len(proc[1])
}

// BiDijkstraAllNodes additionally returns all visited nodes on uniform grid
func BiDijkstraAllNodesBg(fromIDX, toIDX []int, bg *grids.BasicGrid) ([][][]float64, [][]float64) {
	var prev [][]int
	dist := make([][]float64, 2)
	pq := []priorityQueue{make(priorityQueue, 1), make(priorityQueue, 1)}
	proc := []map[int]bool{make(map[int]bool), make(map[int]bool)}
	var meeting int

	// Init
	for i := 0; i < len(bg.VertexData); i++ {
		dist[0] = append(dist[0], math.Inf(1))
		dist[1] = append(dist[1], math.Inf(1))
		prev = append(prev, []int{-1, -1})
	}

	dist[0][bg.FlattenIndex(fromIDX)] = 0
	pq[0][0] = &Item{
		value:    bg.FlattenIndex(fromIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[0])

	dist[1][bg.FlattenIndex(toIDX)] = 0
	pq[1][0] = &Item{
		value:    bg.FlattenIndex(toIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[1])

	var dir = 0 // Direction in which next step should be taken. 0=forward, 1=backward
	// Main loop
	for {
		if len(pq[0]) == 0 || len(pq[1]) == 0 {
			break
		} else {
			u := heap.Pop(&pq[dir]).(*Item).value
			proc[dir][u] = true

			if proc[1-dir][u] {
				var bestDist = dist[0][u] + dist[1][u]
				meeting = u

				for _, k := range pq[0] {
					if proc[1][k.value] {
						if dist[0][k.value]+dist[1][k.value] < bestDist {
							bestDist = dist[0][k.value] + dist[1][k.value]
							meeting = k.value
						}
					}
				}

				for _, k := range pq[1] {
					if proc[0][k.value] {
						if dist[0][k.value]+dist[1][k.value] < bestDist {
							bestDist = dist[0][k.value] + dist[1][k.value]
							meeting = k.value
						}
					}
				}
				break
			}

			neighbours := neighboursBg(u, bg)
			for _, j := range neighbours {
				var alt = dist[dir][u] + distance(bg.GridToCoord(bg.ExpandIndex(u)), bg.GridToCoord(bg.ExpandIndex(j)))
				if alt < dist[dir][j] {
					dist[dir][j] = alt
					prev[j][dir] = u
					item := &Item{
						value:    j,
						priority: -dist[dir][j],
					}
					heap.Push(&pq[dir], item)
				}
			}
		}
		dir = 1 - dir // Change direction
	}

	keys := make([]int, len(proc[0])+len(proc[1]))
	i := 0
	for k := range proc[0] {
		keys[i] = k
		i++
	}
	for k := range proc[1] {
		keys[i] = k
		i++
	}
	var processedNodes = extractNodes(&keys, bg)
	var route = ExtractRouteBi(&prev, meeting, bg)
	return route, processedNodes
}
