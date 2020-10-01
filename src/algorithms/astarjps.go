package algorithms

import (
	"../grids"
	"container/heap"
	"log"
	"math"
)

// AStarJPS Jump Point Search implementation on uniform grid
func AStarJPS(from, to int, ug *grids.UniformGrid) (*[][][]float64, int, float64) {
	var popped int
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	var pq = make(pqJPS, 1)
	var toCoord = ug.GridToCoord(ug.IDToGrid(to))

	for i := 0; i < ug.N; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[from] = 0
	pq[0] = &NodeJPS{
		IDX:      from,
		grid:     ug.IDToGrid(from),
		priority: 0,
		index:    0,
		dir:      -1,
	}
	heap.Init(&pq)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*NodeJPS)
			popped++

			if u.IDX == to {
				return extractRouteUg(&prev, to, ug), popped, dist[to]
			}

			neighbours := prune(u, SimpleNeighboursUgJPS(*u, ug))
			uCoord := ug.GridToCoord(u.grid)
			for _, n := range *neighbours {
				j := jump(u.IDX, u, n.dir, from, to, ug)
				if j != nil {
					j.grid = ug.IDToGrid(j.IDX)
					alt := dist[u.IDX] + distance(uCoord, ug.GridToCoord(j.grid))
					if alt < dist[j.IDX] {
						dist[j.IDX] = alt
						prev[j.IDX] = u.IDX
						item := &NodeJPS{
							grid:     j.grid,
							IDX:      j.IDX,
							dir:      j.dir,
							priority: -(dist[j.IDX] + distance(ug.GridToCoord(j.grid), toCoord)),
						}
						heap.Push(&pq, item)
					}
				}
			}
		}
	}
	return extractRouteUg(&prev, to, ug), popped, dist[to]
}

// AStarJPSAllNodes also returns visited nodes on uniform grid
func AStarJPSAllNodes(from, to int, ug *grids.UniformGrid) (*[][][]float64, *[][]float64, float64) {
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	var pq = make(pqJPS, 1)
	var nodesProcessed []int
	var toCoord = ug.GridToCoord(ug.IDToGrid(to))

	for i := 0; i < ug.N; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[from] = 0
	pq[0] = &NodeJPS{
		IDX:      from,
		grid:     ug.IDToGrid(from),
		priority: 0,
		index:    0,
		dir:      -1,
	}
	heap.Init(&pq)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*NodeJPS)
			nodesProcessed = append(nodesProcessed, u.IDX)

			if u.IDX == to {
				return extractRouteUg(&prev, to, ug), extractNodesUg(&nodesProcessed, ug), dist[to]
			}

			neighbours := prune(u, SimpleNeighboursUgJPS(*u, ug))
			uCoord := ug.GridToCoord(u.grid)
			for _, n := range *neighbours {
				j := jump(u.IDX, u, n.dir, from, to, ug)
				if j != nil {
					j.grid = ug.IDToGrid(j.IDX)
					alt := dist[u.IDX] + distance(uCoord, ug.GridToCoord(j.grid))
					if alt < dist[j.IDX] {
						dist[j.IDX] = alt
						prev[j.IDX] = u.IDX
						item := &NodeJPS{
							grid:     j.grid,
							IDX:      j.IDX,
							dir:      j.dir,
							priority: -(dist[j.IDX] + distance(ug.GridToCoord(j.grid), toCoord)),
						}
						heap.Push(&pq, item)
					}
				}
			}
		}
	}
	return extractRouteUg(&prev, to, ug), extractNodesUg(&nodesProcessed, ug), dist[to]
}

func jump(u int, nn *NodeJPS, dir, from, to int, ug *grids.UniformGrid) *NodeJPS {
	n := step(nn, dir, ug)
	n.grid = ug.IDToGrid(n.IDX)

	if u == n.IDX || ug.VertexData[n.grid[0]][n.grid[1]] {
		return nil
	}

	if n.IDX == to {
		return n
	}

	for _, i := range *prune(n, SimpleNeighboursUgJPS(*n, ug)) {
		if i.forced {
			return n
		}
	}

	switch dir {
	case 0:
		if jump(n.IDX, n, 1, from, to, ug) != nil || jump(n.IDX, n, 3, from, to, ug) != nil {
			return n
		}
	case 2:
		if jump(n.IDX, n, 1, from, to, ug) != nil || jump(n.IDX, n, 4, from, to, ug) != nil {
			return n
		}
	case 5:
		if jump(n.IDX, n, 3, from, to, ug) != nil || jump(n.IDX, n, 6, from, to, ug) != nil {
			return n
		}
	case 7:
		if jump(n.IDX, n, 4, from, to, ug) != nil || jump(n.IDX, n, 6, from, to, ug) != nil {
			return n
		}
	}
	return jump(u, n, n.dir, from, to, ug)
}

func step(i *NodeJPS, dir int, ug *grids.UniformGrid) *NodeJPS {
	i.grid = ug.IDToGrid(i.IDX)
	m := i.grid[0]
	n := i.grid[1]

	if dir == 3 {
		return &NodeJPS{
			IDX: ug.GridToID([]int{m, mod(n-1, len(ug.VertexData[m]))}),
			dir: 3,
		}
	}

	if dir == 4 {
		return &NodeJPS{
			IDX: ug.GridToID([]int{m, mod(n+1, len(ug.VertexData[m]))}),
			dir: 4,
		}
	}

	ratio := float64(n) / float64(len(ug.VertexData[m]))
	if dir > 4 {
		if m > 0 {
			lmm := len(ug.VertexData[m-1])
			nDown := int(math.Round(ratio * float64(lmm)))

			switch dir {
			case 5:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m - 1, mod(nDown-1.0, lmm)}),
					dir: 5,
				}
			case 6:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m - 1, mod(nDown, lmm)}),
					dir: 6,
				}
			case 7:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m - 1, mod(nDown+1.0, lmm)}),
					dir: 7,
				}
			}
		} else {
			log.Fatal("Error: Water at the south pole.")
		}
	} else {
		if m < len(ug.VertexData)-1 {
			lmp := len(ug.VertexData[m+1])
			nUp := int(math.Round(ratio * float64(lmp)))

			switch dir {
			case 0:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m + 1, mod(nUp-1.0, lmp)}),
					dir: 0,
				}
			case 1:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m + 1, mod(nUp, lmp)}),
					dir: 1,
				}
			case 2:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m + 1, mod(nUp+1.0, lmp)}),
					dir: 2,
				}
			}
		} else {
			lmm := len(ug.VertexData[m-2])
			nDown := int(math.Round(ratio * float64(lmm)))

			switch dir {
			case 0:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m - 2, mod(nDown+lmm/2, lmm)}),
					dir: 7,
				}
			case 1:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m - 2, mod(nDown+lmm/2, lmm)}),
					dir: 6,
				}
			case 2:
				return &NodeJPS{
					IDX: ug.GridToID([]int{m - 2, mod(nDown+lmm/2, lmm)}),
					dir: 5,
				}
			}
		}
	}
	return nil
}

// SimpleNeighboursUgJPS gets ug neighbours cheaper
func SimpleNeighboursUgJPS(in NodeJPS, ug *grids.UniformGrid) *map[int]NodeJPS {
	var ratio float64
	var nUp, nDown int
	in.grid = ug.IDToGrid(in.IDX)
	var m = in.grid[0]
	var n = in.grid[1]
	var neighbours1d = make(map[int]NodeJPS)

	// lengths of rows
	var lm = len(ug.VertexData[m])

	grid := []int{m, mod(n-1, lm)}
	if !ug.VertexData[grid[0]][grid[1]] {
		neighbours1d[3] = NodeJPS{
			IDX: ug.GridToID(grid),
			dir: 3,
		}
	}

	grid = []int{m, mod(n+1, lm)}
	if !ug.VertexData[grid[0]][grid[1]] {
		neighbours1d[4] = NodeJPS{
			IDX: ug.GridToID(grid),
			dir: 4,
		}
	}

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData)-1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))

		grid := []int{m + 1, mod(nUp, lmp)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[1] = NodeJPS{
				IDX: ug.GridToID(grid),
				dir: 1,
			}
		}
		grid = []int{m + 1, mod(nUp+1.0, lmp)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[2] = NodeJPS{
				IDX: ug.GridToID(grid),
				dir: 2,
			}
		}
		grid = []int{m + 1, mod(nUp-1.0, lmp)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[0] = NodeJPS{
				IDX: ug.GridToID(grid),
				dir: 0,
			}
		}
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))

		grid := []int{m - 1, mod(nDown, lmm)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[6] = NodeJPS{
				IDX: ug.GridToID(grid),
				dir: 6,
			}
		}
		grid = []int{m - 1, mod(nDown+1.0, lmm)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[7] = NodeJPS{
				IDX: ug.GridToID(grid),
				dir: 7,
			}
		}
		grid = []int{m - 1, mod(nDown-1.0, lmm)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[5] = NodeJPS{
				IDX: ug.GridToID(grid),
				dir: 5,
			}
		}
	}
	return &neighbours1d
}

func prune(i *NodeJPS, nbs *map[int]NodeJPS) *map[int]NodeJPS {
	if i.dir == -1 {
		return nbs
	}

	var res = make(map[int]NodeJPS)
	var n = make([][]int, 2) // forced neighbours to check
	var m = make([]int, 3)   // natural neighbours to check

	m[0] = i.dir
	// Right
	switch i.dir {
	case 4:
		// Right
		n[0] = []int{1, 2}
		n[1] = []int{6, 7}
	case 3:
		// Left
		n[0] = []int{1, 0}
		n[1] = []int{6, 5}
	case 1:
		// Up
		n[0] = []int{3, 0}
		n[1] = []int{4, 2}
	case 6:
		// Down
		n[0] = []int{4, 7}
		n[1] = []int{3, 5}
	case 2:
		// Top Right
		m[1] = 1
		m[2] = 4

		n[0] = []int{3, 0}
		n[1] = []int{6, 7}
	case 0:
		// Top Left
		m[1] = 1
		m[2] = 3

		n[0] = []int{6, 5}
		n[1] = []int{4, 2}
	case 7:
		// Bottom Right
		m[1] = 4
		m[2] = 6

		n[0] = []int{1, 2}
		n[1] = []int{3, 5}
	case 5:
		// Bottom Left
		m[1] = 3
		m[2] = 6

		n[0] = []int{1, 0}
		n[1] = []int{4, 7}
	}
	for _, i := range m {
		if _, ok := (*nbs)[i]; ok {
			res[i] = (*nbs)[i]
		}
	}
	for _, i := range n {
		if _, ok := (*nbs)[i[0]]; !ok {
			if _, ok := (*nbs)[i[1]]; ok {
				x := (*nbs)[i[1]]
				x.forced = true
				res[i[1]] = x
			}
		}
	}
	return &res
}
