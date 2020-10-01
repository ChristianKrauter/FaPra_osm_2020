package algorithms

import (
	"../grids"
	"container/heap"
	//"fmt"
	"math"
)

// AStarJPSBg implementation on uniform grid
func AStarJPSBg(from, to int, bg *grids.BasicGrid) ([][][]float64, int, float64) {
	var popped int
	var dist = make([]float64, len(bg.VertexData))
	var prev = make([]int, len(bg.VertexData))
	pq := make(pqJPS, 1)
	var toCoord = bg.GridToCoord(bg.IDToGrid(to))

	for i := 0; i < len(bg.VertexData); i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[from] = 0
	pq[0] = &NodeJPS{
		IDX:      from,
		grid:     bg.IDToGrid(from),
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
				return extractRoute(&prev, to, bg), popped, dist[to]
			}

			neighbours := pruneBg(u, NeighboursBgJPS(u.IDX, bg))
			uCoord := bg.GridToCoord(u.grid)
			for _, n := range *neighbours {
				j := jumpBg(u.IDX, u, n.dir, from, to, bg)
				if j != nil {
					j.grid = bg.IDToGrid(j.IDX)
					alt := dist[u.IDX] + distance(uCoord, bg.GridToCoord(j.grid))
					if alt < dist[j.IDX] {
						dist[j.IDX] = alt
						prev[j.IDX] = u.IDX
						item := &NodeJPS{
							grid:     j.grid,
							IDX:      j.IDX,
							dir:      j.dir,
							priority: -(dist[j.IDX] + distance(bg.GridToCoord(j.grid), toCoord)),
						}
						heap.Push(&pq, item)
					}
				}
			}
		}
	}
	return extractRoute(&prev, to, bg), popped, dist[to]
}

// AStarJPSAllNodesBg implementation on uniform grid
func AStarJPSAllNodesBg(from, to int, bg *grids.BasicGrid) ([][][]float64, [][]float64, float64) {
	var popped int
	var dist = make([]float64, len(bg.VertexData))
	var prev = make([]int, len(bg.VertexData))
	pq := make(pqJPS, 1)
	var nodesProcessed []int
	var toCoord = bg.GridToCoord(bg.IDToGrid(to))

	for i := 0; i < len(bg.VertexData); i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}

	dist[from] = 0
	pq[0] = &NodeJPS{
		IDX:      from,
		grid:     bg.IDToGrid(from),
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

			popped++
			if u.IDX == to {
				return extractRoute(&prev, to, bg), extractNodes(&nodesProcessed, bg), dist[to]
			}

			neighbours := pruneBg(u, NeighboursBgJPS(u.IDX, bg))
			uCoord := bg.GridToCoord(u.grid)
			for _, n := range *neighbours {
				j := jumpBg(u.IDX, u, n.dir, from, to, bg)
				if j != nil {
					j.grid = bg.IDToGrid(j.IDX)
					alt := dist[u.IDX] + distance(uCoord, bg.GridToCoord(j.grid))
					if alt < dist[j.IDX] {
						dist[j.IDX] = alt
						prev[j.IDX] = u.IDX
						item := &NodeJPS{
							grid:     j.grid,
							IDX:      j.IDX,
							dir:      j.dir,
							priority: -(dist[j.IDX] + distance(bg.GridToCoord(j.grid), toCoord)),
						}
						heap.Push(&pq, item)
					}
				}
			}
		}
	}
	return extractRoute(&prev, to, bg), extractNodes(&nodesProcessed, bg), dist[to]
}

func jumpBg(u int, nn *NodeJPS, dir, from, to int, bg *grids.BasicGrid) *NodeJPS {
	n := stepBg(nn.IDX, dir, bg)

	if n == nil || u == n.IDX || bg.VertexData[n.IDX] {
		return nil
	}
	if n.IDX == to {
		return n
	}

	for _, i := range *pruneBg(n, NeighboursBgJPS(n.IDX, bg)) {
		if i.forced {
			return n
		}
	}

	switch dir {
	case 0:
		if jumpBg(n.IDX, n, 1, from, to, bg) != nil || jumpBg(n.IDX, n, 3, from, to, bg) != nil {
			return n
		}
	case 2:
		if jumpBg(n.IDX, n, 1, from, to, bg) != nil || jumpBg(n.IDX, n, 4, from, to, bg) != nil {
			return n
		}
	case 5:
		if jumpBg(n.IDX, n, 3, from, to, bg) != nil || jumpBg(n.IDX, n, 6, from, to, bg) != nil {
			return n
		}
	case 7:
		if jumpBg(n.IDX, n, 4, from, to, bg) != nil || jumpBg(n.IDX, n, 6, from, to, bg) != nil {
			return n
		}
	}
	return jumpBg(u, n, n.dir, from, to, bg)
}

func stepBg(IDX int, dir int, bg *grids.BasicGrid) *NodeJPS {

	switch dir {
	case 3:
		n := IDX - 1
		if n >= 0 {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 3,
				}
			}
		}
	case 4:
		n := IDX + 1
		if n < len(bg.VertexData) {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 4,
				}
			}
		}

	case 5:
		n := IDX + bg.XSize - 1
		if n < len(bg.VertexData) {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 5,
				}
			}
		}

	case 6:
		n := IDX + bg.XSize
		if n < len(bg.VertexData) {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 6,
				}
			}
		}

	case 7:
		n := IDX + bg.XSize + 1
		if n < len(bg.VertexData) {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 7,
				}
			}
		}

	case 0:
		n := IDX - bg.XSize - 1
		if n >= 0 {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 0,
				}
			}
		}

	case 1:
		n := IDX - bg.XSize
		if n >= 0 {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 1,
				}
			}
		}

	case 2:
		n := IDX - bg.XSize + 1
		if n >= 0 {
			if !bg.VertexData[n] {
				return &NodeJPS{
					IDX: n,
					dir: 2,
				}
			}
		}
	}

	return nil
}

// NeighboursBgJPS ...
func NeighboursBgJPS(IDX int, bg *grids.BasicGrid) *map[int]NodeJPS {
	var neighbours1d = make(map[int]NodeJPS)
	var lm = len(bg.VertexData)

	n := IDX - 1
	if n >= 0 {
		if !bg.VertexData[n] {
			neighbours1d[3] = NodeJPS{
				IDX: n,
				dir: 3,
			}
		}
	}

	n = IDX + 1
	if n < lm {
		if !bg.VertexData[n] {
			neighbours1d[4] = NodeJPS{
				IDX: n,
				dir: 4,
			}
		}
	}

	n = IDX - bg.XSize - 1
	if n >= 0 {
		if !bg.VertexData[n] {
			neighbours1d[0] = NodeJPS{
				IDX: n,
				dir: 0,
			}
		}
	}

	n = IDX - bg.XSize
	if n >= 0 {
		if !bg.VertexData[n] {
			neighbours1d[1] = NodeJPS{
				IDX: n,
				dir: 1,
			}
		}
	}

	n = IDX - bg.XSize + 1
	if n >= 0 {
		if !bg.VertexData[n] {
			neighbours1d[2] = NodeJPS{
				IDX: n,
				dir: 2,
			}
		}
	}

	n = IDX + bg.XSize - 1
	if n < lm {
		if !bg.VertexData[n] {
			neighbours1d[5] = NodeJPS{
				IDX: n,
				dir: 5,
			}
		}
	}

	n = IDX + bg.XSize
	if n < lm {
		if !bg.VertexData[n] {
			neighbours1d[6] = NodeJPS{
				IDX: n,
				dir: 6,
			}
		}
	}

	n = IDX + bg.XSize + 1
	if n < lm {
		if !bg.VertexData[n] {
			neighbours1d[7] = NodeJPS{
				IDX: n,
				dir: 7,
			}
		}
	}

	return &neighbours1d
}

func pruneBg(i *NodeJPS, nbs *map[int]NodeJPS) *map[int]NodeJPS {
	var res = make(map[int]NodeJPS)
	if i.dir == -1 {
		return nbs
	}
	// forced neighbours to check
	var n = make([][]int, 2)
	// natural neighbours to check
	var m = make([]int, 3)
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
