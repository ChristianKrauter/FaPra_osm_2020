package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// BiDijkstraBg implementation on uniform grid
func BiDijkstraBg(from, to int, bg *grids.BasicGrid) ([][][]float64, int, float64) {
	var prev = make([][]int, len(bg.VertexData))
	dist := [][]float64{make([]float64, len(bg.VertexData)), make([]float64, len(bg.VertexData))}
	pq := []priorityQueue{make(priorityQueue, 1), make(priorityQueue, 1)}
	proc := []map[int]bool{make(map[int]bool), make(map[int]bool)}
	var met = false
	var meeting int
	var bestDist = math.MaxFloat64

	// Init
	for i := 0; i < len(bg.VertexData); i++ {
		dist[0][i] = math.Inf(1)
		dist[1][i] = math.Inf(1)
		prev[i] = []int{-1, -1}
	}

	dist[0][from] = 0
	pq[0][0] = &Item{
		value:    from,
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[0])

	dist[1][to] = 0
	pq[1][0] = &Item{
		value:    to,
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[1])

	var dir = 0 // Direction in which next step should be taken. 0=forward, 1=backward
	// Main loop
	for {
		if len(pq[0]) == 0 && len(pq[1]) == 0 {
			break
		} else {
			if len(pq[dir]) > 0 {

				u := heap.Pop(&pq[dir]).(*Item).value
				proc[dir][u] = true

				if proc[1-dir][u] {
					met = true
					if bestDist > dist[0][u]+dist[1][u] {
						bestDist = dist[0][u] + dist[1][u]
						meeting = u
					}
				}

				neighbours := NeighboursBg(u, bg)
				uCoord := bg.GridToCoord(bg.IDToGrid(u))
				for _, j := range neighbours {
					alt := dist[dir][u] + distance(uCoord, bg.GridToCoord(bg.IDToGrid(j)))
					if alt < dist[dir][j] {
						dist[dir][j] = alt
						prev[j][dir] = u
						item := &Item{
							value:    j,
							priority: -dist[dir][j],
						}
						if !met {
							heap.Push(&pq[dir], item)
						} else {
							if bestDist > dist[0][j]+dist[1][j] {
								bestDist = dist[0][j] + dist[1][j]
								meeting = j
							}
						}
					}
				}
			}
		}
		dir = 1 - dir // Change direction
	}

	return extractRouteBi(&prev, meeting, bg), len(proc[0]) + len(proc[1]), dist[0][meeting] + dist[1][meeting]
}

// BiDijkstraAllNodesBg additionally returns all visited nodes on uniform grid
func BiDijkstraAllNodesBg(from, to int, bg *grids.BasicGrid) ([][][]float64, [][]float64, float64) {
	var prev = make([][]int, len(bg.VertexData))
	dist := [][]float64{make([]float64, len(bg.VertexData)), make([]float64, len(bg.VertexData))}
	pq := []priorityQueue{make(priorityQueue, 1), make(priorityQueue, 1)}
	proc := []map[int]bool{make(map[int]bool), make(map[int]bool)}
	var met = false
	var meeting int
	var bestDist = math.MaxFloat64

	// Init
	for i := 0; i < len(bg.VertexData); i++ {
		dist[0][i] = math.Inf(1)
		dist[1][i] = math.Inf(1)
		prev[i] = []int{-1, -1}
	}

	dist[0][from] = 0
	pq[0][0] = &Item{
		value:    from,
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[0])

	dist[1][to] = 0
	pq[1][0] = &Item{
		value:    to,
		priority: 0,
		index:    0,
	}
	heap.Init(&pq[1])

	var dir = 0 // Direction in which next step should be taken. 0=forward, 1=backward
	// Main loop
	for {
		if len(pq[0]) == 0 && len(pq[1]) == 0 {
			break
		} else {
			if len(pq[dir]) > 0 {

				u := heap.Pop(&pq[dir]).(*Item).value
				proc[dir][u] = true

				if proc[1-dir][u] {
					met = true
					if bestDist > dist[0][u]+dist[1][u] {
						bestDist = dist[0][u] + dist[1][u]
						meeting = u
					}
				}

				neighbours := NeighboursBg(u, bg)
				uCoord := bg.GridToCoord(bg.IDToGrid(u))
				for _, j := range neighbours {
					alt := dist[dir][u] + distance(uCoord, bg.GridToCoord(bg.IDToGrid(j)))
					if alt < dist[dir][j] {
						dist[dir][j] = alt
						prev[j][dir] = u
						item := &Item{
							value:    j,
							priority: -dist[dir][j],
						}
						if !met {
							heap.Push(&pq[dir], item)
						} else {
							if bestDist > dist[0][j]+dist[1][j] {
								bestDist = dist[0][j] + dist[1][j]
								meeting = j
							}
						}
					}
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

	return extractRouteBi(&prev, meeting, bg), extractNodes(&keys, bg), dist[0][meeting] + dist[1][meeting]
}
