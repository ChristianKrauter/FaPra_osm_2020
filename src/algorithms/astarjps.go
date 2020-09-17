package algorithms

import (
	"../grids"
	"container/heap"
	"fmt"
	//"github.com/paulmach/go.geojson"
	"math"
)

// AStarJPS implementation on uniform grid
func AStarJPS(from, to int, ug *grids.UniformGrid) ([][][]float64, int) {
	var popped int
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	pq := make(pqJPS, 1)

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
			fmt.Printf("pq empty\n")
			break
		} else {
			u := heap.Pop(&pq).(*NodeJPS)
			popped++
			if u.IDX == to {
				fmt.Printf("u: %v\n", u)
				return ExtractRouteUg(&prev, to, ug), popped
			}

			neighbours := neighboursUgJPS(*u, ug)
			neighbours = prune(u, to, neighbours, ug)
			//fmt.Printf("pruned: %v\n", neighbours)
			//fmt.Printf("nb: %v\n", neighbours)
			var successors []NodeJPS
			for _, j := range *neighbours {

				n := jump(u.IDX, u, j.dir, from, to, ug)
				if n != nil {
					if n.IDX == to {
						fmt.Printf("to: %v\n", n)
					}
					successors = append(successors, *n)
				}
			}

			if len(successors) < 1 {
				for _, j := range *neighbours {
					successors = append(successors, j)
				}
			}

			for _, j := range successors {
				var alt = dist[u.IDX] + distance(ug.GridToCoord(u.grid), ug.GridToCoord(j.grid))
				//xx := false
				if j.IDX == to {
					fmt.Printf("to with dist: %v, %v\n", j, alt)
				}
				if alt < dist[j.IDX] {
					//xx = true
					// dist[j] always +inf?
					//fmt.Printf("alt < dist[j]: %v < %v\n", alt, dist[j.IDX])
					dist[j.IDX] = alt
					prev[j.IDX] = u.IDX
					item := &NodeJPS{
						grid:     j.grid,
						IDX:      j.IDX,
						dir:      j.dir,
						priority: -(dist[j.IDX] + distance(ug.GridToCoord(j.grid), ug.GridToCoord(ug.IDToGrid(to)))),
					}
					if j.IDX == to {
						fmt.Printf("to item: %v\n", item)
					}
					heap.Push(&pq, item)
				}
				/*if !xx {
					fmt.Printf("got none\n")
				}*/
			}
		}
	}
	return ExtractRouteUg(&prev, to, ug), popped
}

//var np []int

// AStarJPS implementation on uniform grid
func AStarJPSAllNodes(from, to int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {
	var popped int
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
	//np = make([]int, 0)
	pq := make(pqJPS, 1)
	var nodesProcessed []int
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
	fmt.Printf("hi kids\n")
	for {
		if len(pq) == 0 {
			fmt.Printf("pq empty\n")
			break
		} else {
			u := heap.Pop(&pq).(*NodeJPS)
			nodesProcessed = append(nodesProcessed, u.IDX)

			popped++
			if u.IDX == to {
				fmt.Printf("u: %v\n", u)
				//return ExtractRouteUg(&prev, to, ug), ExtractNodesUg(&np, ug)
				return ExtractRouteUg(&prev, to, ug), ExtractNodesUg(&nodesProcessed, ug)
			}

			fmt.Printf("do you like violence?\n")
			neighbours := neighboursUgJPS(*u, ug)
			neighbours = prune(u, to, neighbours, ug)
			//fmt.Printf("pruned: %v\n", neighbours)
			//fmt.Printf("nb: %v\n", neighbours)

			fmt.Printf("do you wanna stick 9 inch nails\n")
			var successors []NodeJPS
			for _, j := range *neighbours {
				fmt.Printf("u.IDX, u: %v %v\n", u.IDX, u)
				n := jump(u.IDX, u, j.dir, from, to, ug)
				if n != nil {
					if n.IDX == to {
						fmt.Printf("to: %v\n", n)
					}
					successors = append(successors, *n)
				}
			}

			if len(successors) < 1 {
				for _, j := range *neighbours {
					successors = append(successors, j)
				}
			}

			fmt.Printf("through each one of your eyelids?\n")
			for _, j := range successors {
				var alt = dist[u.IDX] + distance(ug.GridToCoord(u.grid), ug.GridToCoord(j.grid))
				//xx := false
				if j.IDX == to {
					fmt.Printf("to with dist: %v, %v\n", j, alt)
				}
				if alt < dist[j.IDX] {
					//xx = true
					// dist[j] always +inf?
					//fmt.Printf("alt < dist[j]: %v < %v\n", alt, dist[j.IDX])
					dist[j.IDX] = alt
					prev[j.IDX] = u.IDX
					item := &NodeJPS{
						grid:     j.grid,
						IDX:      j.IDX,
						dir:      j.dir,
						priority: -(dist[j.IDX] + distance(ug.GridToCoord(j.grid), ug.GridToCoord(ug.IDToGrid(to)))),
					}
					if j.IDX == to {
						fmt.Printf("to item: %v\n", item)
					}
					heap.Push(&pq, item)
				}
				/*if !xx {
					fmt.Printf("got none\n")
				}*/
			}
		}
	}
	//return ExtractRouteUg(&prev, to, ug), ExtractNodesUg(&np, ug)
	return ExtractRouteUg(&prev, to, ug), ExtractNodesUg(&nodesProcessed, ug)
}

func jump(u int, n *NodeJPS, dir, from, to int, ug *grids.UniformGrid) *NodeJPS {
	n = step(n, dir, ug)
	//fmt.Printf("nn, u: %v, %v\n", nn.IDX, u)
	// ToDo: "outside the grid"
	if n == nil || ug.VertexData[n.grid[0]][n.grid[1]] || u == n.IDX {
		//np = append(np, n.IDX)
		return nil
	}
	if n.IDX == to {
		//np = append(np, n.IDX)
		return n
	}

	nbs := prune(n, to, neighboursUgJPS(*n, ug), ug)
	for _, i := range *nbs {
		if i.forced {
			//np = append(np, n.IDX)
			return n
		}
	}

	if dir == 0 {
		if jump(n.IDX, n, 1, from, to, ug) != nil {
			//np = append(np, n.IDX)
			return n
		}
		if jump(n.IDX, n, 3, from, to, ug) != nil {
			//np = append(np, n.IDX)
			return n
		}
	}
	if dir == 2 {
		if jump(n.IDX, n, 1, from, to, ug) != nil {
			//np = append(np, n.IDX)
			return n
		}
		if jump(n.IDX, n, 4, from, to, ug) != nil {
			//np = append(np, n.IDX)
			return n
		}
	}
	if dir == 5 {
		if jump(n.IDX, n, 3, from, to, ug) != nil {
			// = append(np, n.IDX)
			return n
		}
		if jump(n.IDX, n, 6, from, to, ug) != nil {
			//np = append(np, n.IDX)
			return n
		}
	}
	if dir == 7 {
		if jump(n.IDX, n, 4, from, to, ug) != nil {
			//np = append(np, n.IDX)
			return n
		}
		if jump(n.IDX, n, 6, from, to, ug) != nil {
			//np = append(np, n.IDX)
			return n
		}
	}
	//fmt.Printf("end\n")
	return jump(u, n, n.dir, from, to, ug)
}

func prune(i *NodeJPS, to int, nbs *map[int]NodeJPS, ug *grids.UniformGrid) *map[int]NodeJPS {
	var res = make(map[int]NodeJPS)
	//fmt.Printf("dir: %v\n", i.dir)
	if i.dir == -1 {
		return nbs
	}
	// forced neighbours to check
	var n [][]int
	// natural neighbours to check
	var m []int
	m = append(m, i.dir)
	// Right
	switch i.dir {
	case 4:
		// Right
		n = append(n, []int{1, 2})
		n = append(n, []int{6, 7})
	case 3:
		// Left
		n = append(n, []int{1, 0})
		n = append(n, []int{6, 5})
	case 1:
		// Up
		n = append(n, []int{3, 0})
		n = append(n, []int{4, 2})
	case 6:
		// Down
		n = append(n, []int{4, 7})
		n = append(n, []int{3, 5})
	case 2:
		// Top Right
		m = append(m, 1)
		m = append(m, 4)

		n = append(n, []int{3, 0})
		n = append(n, []int{6, 7})
	case 0:
		// Top Left
		m = append(m, 1)
		m = append(m, 3)

		n = append(n, []int{6, 5})
		n = append(n, []int{4, 2})
	case 7:
		// Bottom Right
		m = append(m, 4)
		m = append(m, 6)

		n = append(n, []int{1, 2})
		n = append(n, []int{3, 5})
	case 5:
		// Bottom Left
		m = append(m, 3)
		m = append(m, 6)

		n = append(n, []int{1, 0})
		n = append(n, []int{4, 7})

	}
	for _, i := range m {
		if _, ok := (*nbs)[i]; ok {
			//if !ug.VertexData[val2.grid[0]][val2.grid[1]] {
			res[i] = (*nbs)[i]
			//}
		}
	}
	for _, i := range n {
		if _, ok := (*nbs)[i[0]]; !ok {
			//if ug.VertexData[val.grid[0]][val.grid[1]] {
			if _, ok := (*nbs)[i[1]]; ok {
				//if !ug.VertexData[val2.grid[0]][val2.grid[1]] {
				x := (*nbs)[i[1]]
				x.forced = true
				res[i[1]] = (*nbs)[i[1]]
				//}
			}
			//}
		}
	}
	return &res
}

// neighboursUgJPS gets up to 8 neighbours
func neighboursUgJPS(in NodeJPS, ug *grids.UniformGrid) *map[int]NodeJPS {
	var neighbours []NodeJPS
	m := in.grid[0]
	n := in.grid[1]
	neighbours = append(neighbours, NodeJPS{
		grid: []int{m, mod(n-1, len(ug.VertexData[m]))},
		IDX:  ug.GridToID([]int{m, mod(n-1, len(ug.VertexData[m]))}),
		dir:  3,
	})
	neighbours = append(neighbours, NodeJPS{
		grid: []int{m, mod(n+1, len(ug.VertexData[m]))},
		IDX:  ug.GridToID([]int{m, mod(n+1, len(ug.VertexData[m]))}),
		dir:  4,
	})
	//neighbours = append(neighbours, NodeJPS{grid: []int{m, mod(n+1, len(ug.VertexData[m]))}, dir: 4})

	coord := ug.GridToCoord(in.grid)

	if m > 0 {
		coordDown := ug.GridToCoord([]int{m - 1, n})
		neighbours = append(neighbours, neighboursRowUgJPS([]float64{coord[0], coordDown[1]}, ug, 1)...)
	}

	// Nort Pole which neighbour is up? not nothing...
	if m < len(ug.VertexData)-1 {
		coordUp := ug.GridToCoord([]int{m + 1, n})
		neighbours = append(neighbours, neighboursRowUgJPS([]float64{coord[0], coordUp[1]}, ug, 0)...)
	}

	var neighbours1d = make(map[int]NodeJPS)
	for _, j := range neighbours {
		if !ug.VertexData[j.grid[0]][j.grid[1]] {
			neighbours1d[j.dir] = j
		}
	}
	return &neighbours1d
}

// Gets neighours left and right in the same row
func neighboursRowUgJPS(in []float64, ug *grids.UniformGrid, up int) []NodeJPS {
	theta := (in[1] + 90) * math.Pi / 180
	m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
	theta = math.Pi * (m + 0.5) / ug.MTheta
	phi := in[0] * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)

	n1 := math.Round(phi * mPhi / (2 * math.Pi))
	mIDX := mod(int(m), int(ug.MTheta))

	result := make([]NodeJPS, 3)
	result[0].grid = []int{mIDX, mod(int(n1-1), int(mPhi))}
	result[0].IDX = ug.GridToID(result[0].grid)
	result[0].dir = 5 * up

	result[1].grid = []int{mIDX, mod(int(n1), int(mPhi))}
	result[1].IDX = ug.GridToID(result[1].grid)
	result[1].dir = 5*up + 1

	result[2].grid = []int{mIDX, mod(int(n1+1), int(mPhi))}
	result[2].IDX = ug.GridToID(result[2].grid)
	result[2].dir = 5*up + 2

	return result
}

func step(i *NodeJPS, dir int, ug *grids.UniformGrid) *NodeJPS {
	m := i.grid[0]
	n := i.grid[1]

	if dir == 3 {
		return &NodeJPS{
			grid: []int{m, mod(n-1, len(ug.VertexData[m]))},
			IDX:  ug.GridToID([]int{m, mod(n-1, len(ug.VertexData[m]))}),
			dir:  3,
		}
	}
	if dir == 4 {
		//fmt.Printf("step right\n")
		return &NodeJPS{
			grid: []int{m, mod(n+1, len(ug.VertexData[m]))},
			IDX:  ug.GridToID([]int{m, mod(n+1, len(ug.VertexData[m]))}),
			dir:  4,
		}
	}
	coord := ug.GridToCoord(i.grid)

	if dir > 4 {
		if m > 0 {
			coord2 := ug.GridToCoord([]int{m - 1, n})
			theta := (coord2[1] + 90) * math.Pi / 180
			m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
			theta = math.Pi * (m + 0.5) / ug.MTheta
			phi := coord[0] * math.Pi / 180
			mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)

			n1 := math.Round(phi * mPhi / (2 * math.Pi))
			mIDX := mod(int(m), int(ug.MTheta))

			switch dir {
			case 5:
				grid := []int{mIDX, mod(int(n1-1), int(mPhi))}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  5,
				}
			case 6:
				grid := []int{mIDX, mod(int(n1), int(mPhi))}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  6,
				}
			case 7:
				grid := []int{mIDX, mod(int(n1+1), int(mPhi))}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  7,
				}
			}
		} else {
			fmt.Printf("not probable\n")
			return &NodeJPS{
				grid: []int{m, mod(n-1, len(ug.VertexData[m]))},
				IDX:  ug.GridToID([]int{m, mod(n-1, len(ug.VertexData[m]))}),
				dir:  3,
			}
		}
	} else {
		if m < len(ug.VertexData)-1 {
			coord2 := ug.GridToCoord([]int{m + 1, n})
			theta := (coord2[1] + 90) * math.Pi / 180
			m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
			theta = math.Pi * (m + 0.5) / ug.MTheta
			phi := coord[0] * math.Pi / 180
			mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)

			n1 := math.Round(phi * mPhi / (2 * math.Pi))
			mIDX := mod(int(m), int(ug.MTheta))

			switch dir {
			case 0:
				grid := []int{mIDX, mod(int(n1-1), int(mPhi))}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  0,
				}
			case 1:
				grid := []int{mIDX, mod(int(n1), int(mPhi))}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  1,
				}
			case 2:
				grid := []int{mIDX, mod(int(n1+1), int(mPhi))}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  2,
				}
			}
		} else {
			//fmt.Printf("hello Nort Pole\n")
			ratio := float64(n) / float64(len(ug.VertexData[m]))
			lmm := len(ug.VertexData[m-1])
			nDown := int(math.Round(ratio * float64(lmm)))
			// fmt.Printf("np i: %v, %v\n", i.IDX, i.dir)
			x := &NodeJPS{}
			switch dir {
			case 0:
				x = &NodeJPS{
					grid: []int{m - 1, mod(nDown+len(ug.VertexData[m-1])/2, len(ug.VertexData[m-1]))},
					IDX:  ug.GridToID([]int{m - 1, mod(nDown+len(ug.VertexData[m-1])/2, len(ug.VertexData[m-1]))}),
					dir:  7,
				}
			case 1:
				x = &NodeJPS{
					grid: []int{m - 1, mod(nDown+len(ug.VertexData[m-1])/2, len(ug.VertexData[m-1]))},
					IDX:  ug.GridToID([]int{m - 1, mod(nDown+len(ug.VertexData[m-1])/2, len(ug.VertexData[m-1]))}),
					dir:  6,
				}
			case 2:
				x = &NodeJPS{
					grid: []int{m - 1, mod(nDown+len(ug.VertexData[m-1])/2, len(ug.VertexData[m-1]))},
					IDX:  ug.GridToID([]int{m - 1, mod(nDown+len(ug.VertexData[m-1])/2, len(ug.VertexData[m-1]))}),
					dir:  5,
				}
			}
			fmt.Printf("np i: %v, %v --  x: %v, %v\n", i.grid, i.dir, x.grid, x.dir)
			return x
		}
	}
	fmt.Printf("how?")
	return nil
}

// NodeJPS of priority queue
type NodeJPS struct {
	grid     []int // The grid IDX
	IDX      int
	priority float64 // The priority of the item in the queue.
	index    int     // The index of the item in the heap.
	dir      int
	forced   bool
}

// A pqJPS implements heap.Interface and holds Items.
type pqJPS []*NodeJPS

func (pq pqJPS) Len() int { return len(pq) }

func (pq pqJPS) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq pqJPS) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push item into priority queue
func (pq *pqJPS) Push(x interface{}) {
	n := len(*pq)
	item := x.(*NodeJPS)
	item.index = n
	*pq = append(*pq, item)
}

// Pop item from priority queue
func (pq *pqJPS) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
