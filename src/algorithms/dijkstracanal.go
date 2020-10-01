package algorithms

import (
	"../grids"
	"container/heap"
	"math"
)

// DijkstraCanal find route through land on uniform grid
func DijkstraCanal(from, to int, ug *grids.UniformGrid) (*[][][]float64, int, float64) {
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

			neighbours := SimpleNeighboursCanalUg(u, ug)
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

// SimpleNeighboursCanalUg gets ug neighbours cheaper
func SimpleNeighboursCanalUg(in int, ug *grids.UniformGrid) []int {
	var neighbours [][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	// lengths of rows
	var lm = len(ug.VertexData[m])

	neighbours = append(neighbours, []int{m, mod(n-1, lm)})
	neighbours = append(neighbours, []int{m, mod(n+1, lm)})

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData)-1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))
		neighbours = append(neighbours, []int{m + 1, mod(nUp, lmp)})
		neighbours = append(neighbours, []int{m + 1, mod(nUp+1.0, lmp)})
		neighbours = append(neighbours, []int{m + 1, mod(nUp-1.0, lmp)})
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))
		neighbours = append(neighbours, []int{m - 1, mod(nDown, lmm)})
		neighbours = append(neighbours, []int{m - 1, mod(nDown+1.0, lmm)})
		neighbours = append(neighbours, []int{m - 1, mod(nDown-1.0, lmm)})
	}

	var neighbours1d []int
	for _, neighbour := range neighbours {
		neighbours1d = append(neighbours1d, ug.GridToID(neighbour))
	}
	return neighbours1d
}
