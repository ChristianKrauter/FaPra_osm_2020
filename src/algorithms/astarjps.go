package algorithms

import (
	"../grids"
	"container/heap"
	"fmt"
	"math"
)

var didit = 0

// AStarJPS implementation on uniform grid
func AStarJPS(from, to int, ug *grids.UniformGrid) (*[][][]float64, int, float64) {
	didit = 0
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
			break
		} else {
			u := heap.Pop(&pq).(*NodeJPS)
			popped++
			if u.IDX == to {
				fmt.Printf("u, to: %v, %v\n", u, to)
				return extractRouteUg(&prev, to, ug), popped, dist[to]
			}

			neighbours := prune(u, SimpleNeighboursUgJPS(*u, ug))
			for _, n := range *neighbours {
				j := jump(u.IDX, u, n.dir, from, to, &prev, &dist, ug)

				if j != nil {
					fmt.Printf("j: %v\n", j)
					fmt.Printf("to: %v\n", to)
					var alt = dist[u.IDX] + distance(ug.GridToCoord(u.grid), ug.GridToCoord(j.grid))
					if alt < dist[j.IDX] {
						dist[j.IDX] = alt
						prev[j.IDX] = u.IDX
						item := &NodeJPS{
							grid:     j.grid,
							IDX:      j.IDX,
							dir:      j.dir,
							priority: -(dist[j.IDX] + distance(ug.GridToCoord(j.grid), ug.GridToCoord(ug.IDToGrid(to)))),
						}
						heap.Push(&pq, item)
					}
				}
			}
		}
	}
	return extractRouteUg(&prev, to, ug), popped, dist[to]
}

//var nodesProcessed []int

// AStarJPSAllNodes implementation on uniform grid
func AStarJPSAllNodes(from, to int, ug *grids.UniformGrid) (*[][][]float64, *[][]float64, float64) {
	var popped int
	var dist = make([]float64, ug.N)
	var prev = make([]int, ug.N)
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
	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*NodeJPS)
			nodesProcessed = append(nodesProcessed, u.IDX)

			popped++
			if u.IDX == to {
				return extractRouteUg(&prev, to, ug), extractNodesUg(&nodesProcessed, ug), dist[to]
			}

			neighbours := prune(u, SimpleNeighboursUgJPS(*u, ug))
			for _, n := range *neighbours {
				j := jump(u.IDX, u, n.dir, from, to, &prev, &dist, ug)
				if j != nil {
					var alt = dist[u.IDX] + distance(ug.GridToCoord(u.grid), ug.GridToCoord(j.grid))
					if alt < dist[j.IDX] {
						dist[j.IDX] = alt
						prev[j.IDX] = u.IDX
						item := &NodeJPS{
							grid:     j.grid,
							IDX:      j.IDX,
							dir:      j.dir,
							priority: -(dist[j.IDX] + distance(ug.GridToCoord(j.grid), ug.GridToCoord(ug.IDToGrid(to)))),
						}
						heap.Push(&pq, item)
					}
				}
			}
		}
	}
	return extractRouteUg(&prev, to, ug), extractNodesUg(&nodesProcessed, ug), dist[to]
}

func jump(u int, nn *NodeJPS, dir, from, to int, prev *[]int, dist *[]float64, ug *grids.UniformGrid) *NodeJPS {
	n := step(nn, dir, ug)
	//nodesProcessed = append(nodesProcessed, n.IDX)
	//fmt.Printf("n: %v\n", n)
	if u == n.IDX || ug.VertexData[n.grid[0]][n.grid[1]] { //n == nil ||
		//nodesProcessed = append(nodesProcessed, n.IDX)
		return nil
	}
	if n.IDX == to {
		//if (n.dir == 3 && dir == 3) || (n.dir == 4 && dir == 4) /*&& math.Abs(float64(n.grid[1]-ug.IDToGrid(u)[1])) > float64(len(ug.VertexData[n.grid[0]]))*/ {
		/*(*prev)[ug.GridToID(grid)] = nn.IDX
		(*dist)[ug.GridToID(grid)] = (*dist)[nn.IDX] + distance(ug.GridToCoord(nn.grid), ug.GridToCoord(grid))
		(*dist)[n.IDX] = (*dist)[ug.GridToID(grid)] + distance(ug.GridToCoord(ug.IDToGrid(n.IDX)), ug.GridToCoord(grid))
		(*prev)[n.IDX] = ug.GridToID(grid)
		*/
		//fmt.Printf("%v\n", ug.GridToCoord(n.grid)[1])
		if /*(n.dir == 3 || n.dir == 4 || dir == 3 || dir == 4) &&*/ didit < 1 /*&& ug.GridToCoord(n.grid)[1] < -60 */ {
			//grid := []int{n.grid[0], mod(ug.IDToGrid(u)[1]+int(math.Abs(float64(n.grid[1]-ug.IDToGrid(u)[1])))/2, len(ug.VertexData[n.grid[0]]))}
			//grid := []int{n.grid[0], ug.IDToGrid(u)[1] + int(math.Abs(float64(n.grid[1]-ug.IDToGrid(u)[1])))/2}
			//var moduplus int
			var moduplus2 int
			fmt.Printf("dir: %v\n", dir)
			fmt.Printf("n.dir: %v\n", n.dir)
			if n.dir == 4 {

				leng := len(ug.VertexData[n.grid[0]])
				diff := n.grid[1] - ug.IDToGrid(u)[1]
				//diff2 := (n.grid[1] - ug.IDToGrid(u)[1]) / 2
				//moddiff := mod(diff2, leng)
				moddiff2 := mod(diff, leng) / 2
				//uplus := ug.IDToGrid(u)[1] + moddiff
				uplus2 := ug.IDToGrid(u)[1] + moddiff2
				//moduplus = mod(uplus, leng)
				moduplus2 = mod(uplus2, leng)
				//fmt.Printf("u1: %v,\nn1: %v,\ndiff: %v,\ndiff2: %v,\nmoddiff: %v,\nmoddiff2: %v,\nu+: %v,\nu+2: %v,\nlen: %v,\nmod(u+): %v\nmod(u+2): %v\n",
				/*fmt.Printf("u1: %v,\nn1: %v,\ndiff: %v,\ndiff2: %v,\nmoddiff2: %v,\nu+2: %v,\nlen: %v,\nmod(u+2): %v\n",
					ug.IDToGrid(u)[1],
					n.grid[1],
					diff,
					diff2,
					//moddiff,
					moddiff2,
					//uplus,
					uplus2,
					leng,
					//moduplus,
					moduplus2,
					//mod(ug.IDToGrid(u)[1]+mod((n.grid[1]-ug.IDToGrid(u)[1]), len(ug.VertexData[n.grid[0]]))/2, len(ug.VertexData[n.grid[0]])),
				)*/

			} else if n.dir == 3 {

				leng := len(ug.VertexData[n.grid[0]])
				diff := ug.IDToGrid(u)[1] - n.grid[1]
				//diff2 := (n.grid[1] - ug.IDToGrid(u)[1]) / 2
				//diff2 := (ug.IDToGrid(u)[1] - n.grid[1]) / 2
				/*if diff2 > 0 {
					diff2 = (ug.IDToGrid(u)[1] - n.grid[1]) / 2
				}*/

				//diff2 := (ug.IDToGrid(u)[1] - n.grid[1]) / 2
				//moddiff := mod(diff2, leng)
				moddiff2 := mod(diff, leng) / 2
				//uplus := ug.IDToGrid(u)[1] + diff2
				/*uplus := ug.IDToGrid(u)[1] + diff
				if diff > 0 {
					fmt.Printf("diff > 0: %v\n", diff)
					fmt.Printf("uplus: %v\n", uplus)
					uplus = ug.IDToGrid(u)[1] - diff
					fmt.Printf("new uplus: %v\n", uplus)
				}*/
				//uplus2 := ug.IDToGrid(u)[1] + moddiff2
				uplus2 := n.grid[1] + moddiff2
				//uplus2 := ug.IDToGrid(u)[1] + moddiff2
				//moduplus = mod(uplus, leng)
				moduplus2 = mod(uplus2, leng)
				//fmt.Printf("u1: %v,\nn1: %v,\ndiff: %v,\ndiff2: %v,\nmoddiff: %v,\nmoddiff2: %v,\nu+: %v,\nu+2: %v,\nlen: %v,\nmod(u+): %v\nmod(u+2): %v\n",
				/*fmt.Printf("u1: %v,\nn1: %v,\ndiff: %v,\ndiff2: %v,\nmoddiff2: %v,\nu+2: %v,\nlen: %v,\nmod(u+2): %v\n",
					ug.IDToGrid(u)[1],
					n.grid[1],
					diff,
					diff2,
					//moddiff,
					moddiff2,
					//uplus,
					uplus2,
					leng,
					//moduplus,
					moduplus2,
					//mod(ug.IDToGrid(u)[1]+mod((n.grid[1]-ug.IDToGrid(u)[1]), len(ug.VertexData[n.grid[0]]))/2, len(ug.VertexData[n.grid[0]])),
				)*/
			}

			//grid := []int{n.grid[0], mod(ug.IDToGrid(u)[1]+ /*mod(*/ int(math.Abs(float64(n.grid[1]-ug.IDToGrid(u)[1]))) /*, len(ug.VertexData[n.grid[0]]))*/ /2, len(ug.VertexData[n.grid[0]]))}
			grid := []int{n.grid[0], moduplus2}

			fmt.Printf("u, grid, n, len, mod: %v, %v, %v, %v\n\n",
				ug.IDToGrid(u)[1],
				grid[1], n.grid[1],
				len(ug.VertexData[n.grid[0]])) //, mod(ug.IDToGrid(u)[1]+int(math.Abs(float64(n.grid[1]-ug.IDToGrid(u)[1])))/2, len(ug.VertexData[n.grid[0]])))
			//fmt.Printf("-- %v, %v, %v, %v\n", n.grid, nn.grid, math.Abs(float64(n.grid[1]-ug.IDToGrid(u)[1])), float64(len(ug.VertexData[n.grid[0]])/100))
			fmt.Printf("return: %v\n", ug.GridToID(grid))
			didit++
			return &NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  n.dir,
			}
		}
		//}
		//nodesProcessed = append(nodesProcessed, n.IDX)
		return n
	}

	for _, i := range *prune(n, SimpleNeighboursUgJPS(*n, ug)) {
		if i.forced {
			return n
		}
	}

	switch dir {
	case 0:
		if jump(n.IDX, n, 1, from, to, prev, dist, ug) != nil || jump(n.IDX, n, 3, from, to, prev, dist, ug) != nil {
			return n
		}
	case 2:
		if jump(n.IDX, n, 1, from, to, prev, dist, ug) != nil || jump(n.IDX, n, 4, from, to, prev, dist, ug) != nil {
			return n
		}
	case 5:
		if jump(n.IDX, n, 3, from, to, prev, dist, ug) != nil || jump(n.IDX, n, 6, from, to, prev, dist, ug) != nil {
			return n
		}
	case 7:
		if jump(n.IDX, n, 4, from, to, prev, dist, ug) != nil || jump(n.IDX, n, 6, from, to, prev, dist, ug) != nil {
			return n
		}
	}
	y := jump(u, n, n.dir, from, to, prev, dist, ug)
	if y != nil && y.IDX == to {
		fmt.Printf("y %v\n", y)
		if didit < 3 {
			var moduplus2 int
			fmt.Printf("dir: %v\n", dir)
			fmt.Printf("y.dir: %v\n", y.dir)
			if y.dir == 4 {
				leng := len(ug.VertexData[y.grid[0]])
				diff := y.grid[1] - ug.IDToGrid(u)[1]
				moddiff2 := mod(diff, leng) / 2
				uplus2 := ug.IDToGrid(u)[1] + moddiff2
				moduplus2 = mod(uplus2, leng)
			} else if y.dir == 3 {
				leng := len(ug.VertexData[y.grid[0]])
				diff := ug.IDToGrid(u)[1] - y.grid[1]
				moddiff2 := mod(diff, leng) / 2
				uplus2 := y.grid[1] + moddiff2
				moduplus2 = mod(uplus2, leng)
			}

			grid := []int{y.grid[0], moduplus2}

			fmt.Printf("u, grid, n, len, mod: %v, %v, %v, %v\n\n",
				ug.IDToGrid(u)[1],
				grid[1], y.grid[1],
				len(ug.VertexData[y.grid[0]])) //, mod(ug.IDToGrid(u)[1]+int(math.Abs(float64(y.grid[1]-ug.IDToGrid(u)[1])))/2, len(ug.VertexData[y.grid[0]])))
			//fmt.Printf("-- %v, %v, %v, %v\n", y.grid, nn.grid, math.Abs(float64(y.grid[1]-ug.IDToGrid(u)[1])), float64(len(ug.VertexData[y.grid[0]])/100))
			fmt.Printf("return: %v\n", ug.GridToID(grid))
			didit++
			return &NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  y.dir,
			}
		}
	}
	return y
}

func step(i *NodeJPS, dir int, ug *grids.UniformGrid) *NodeJPS {
	m := i.grid[0]
	n := i.grid[1]

	// leave grid out
	if dir == 3 {
		grid := []int{m, mod(n-1, len(ug.VertexData[m]))}
		return &NodeJPS{
			grid: grid,
			IDX:  ug.GridToID(grid),
			dir:  3,
		}
	}
	if dir == 4 {
		grid := []int{m, mod(n+1, len(ug.VertexData[m]))}
		return &NodeJPS{
			grid: grid,
			IDX:  ug.GridToID(grid),
			dir:  4,
		}
	}

	ratio := float64(n) / float64(len(ug.VertexData[m]))
	if dir > 4 {
		if m > 0 {
			lmm := len(ug.VertexData[m-1])
			nDown := int(math.Round(ratio * float64(lmm)))

			switch dir {
			case 5:
				grid := []int{m - 1, mod(nDown-1.0, lmm)}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  5,
				}
			case 6:
				grid := []int{m - 1, mod(nDown, lmm)}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  6,
				}
			case 7:
				grid := []int{m - 1, mod(nDown+1.0, lmm)}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  7,
				}
			}
		} else {
			fmt.Printf("i %v", i)
			fmt.Printf("dir %v", dir)
			fmt.Printf("m %v", m)
			fmt.Printf("Oops that should not happen...\n")
		}
	} else {
		if m < len(ug.VertexData)-1 {
			lmp := len(ug.VertexData[m+1])
			nUp := int(math.Round(ratio * float64(lmp)))

			switch dir {
			case 0:
				grid := []int{m + 1, mod(nUp-1.0, lmp)}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  0,
				}
			case 1:
				grid := []int{m + 1, mod(nUp, lmp)}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  1,
				}
			case 2:
				grid := []int{m + 1, mod(nUp+1.0, lmp)}
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  2,
				}
			}
		} else {
			lmm := len(ug.VertexData[m-2])
			nDown := int(math.Round(ratio * float64(lmm)))
			// nodesProcessed = append(nodesProcessed, i.IDX)

			switch dir {
			case 0:
				grid := []int{m - 2, mod(nDown+lmm/2, lmm)}
				// nodesProcessed = append(nodesProcessed, ug.GridToID(grid))
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  7,
				}
			case 1:
				grid := []int{m - 2, mod(nDown+lmm/2, lmm)}
				// nodesProcessed = append(nodesProcessed, ug.GridToID(grid))
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  6,
				}
			case 2:
				grid := []int{m - 2, mod(nDown+lmm/2, lmm)}
				// nodesProcessed = append(nodesProcessed, ug.GridToID(grid))
				return &NodeJPS{
					grid: grid,
					IDX:  ug.GridToID(grid),
					dir:  5,
				}
			}
		}
	}
	return nil
}

// SimpleNeighboursUgJPS gets ug neighbours cheaper
func SimpleNeighboursUgJPS(in NodeJPS, ug *grids.UniformGrid) *map[int]NodeJPS {
	//var neighbours []NodeJPS
	var ratio float64
	var nUp, nDown int
	var m = in.grid[0]
	var n = in.grid[1]
	var neighbours1d = make(map[int]NodeJPS)

	// lengths of rows
	var lm = len(ug.VertexData[m])

	grid := []int{m, mod(n-1, lm)}
	if !ug.VertexData[grid[0]][grid[1]] {
		neighbours1d[3] = NodeJPS{
			grid: grid,
			IDX:  ug.GridToID(grid),
			dir:  3,
		}
	}

	grid = []int{m, mod(n+1, lm)}
	if !ug.VertexData[grid[0]][grid[1]] {
		neighbours1d[4] = NodeJPS{
			grid: grid,
			IDX:  ug.GridToID(grid),
			dir:  4,
		}
	}

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData)-1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))

		grid := []int{m + 1, mod(nUp, lmp)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[1] = NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  1,
			}
		}
		grid = []int{m + 1, mod(nUp+1.0, lmp)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[2] = NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  2,
			}
		}
		grid = []int{m + 1, mod(nUp-1.0, lmp)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[0] = NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  0,
			}
		}
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))

		grid := []int{m - 1, mod(nDown, lmm)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[6] = NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  6,
			}
		}
		grid = []int{m - 1, mod(nDown+1.0, lmm)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[7] = NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  7,
			}
		}
		grid = []int{m - 1, mod(nDown-1.0, lmm)}
		if !ug.VertexData[grid[0]][grid[1]] {
			neighbours1d[5] = NodeJPS{
				grid: grid,
				IDX:  ug.GridToID(grid),
				dir:  5,
			}
		}
	}
	return &neighbours1d
}

func prune(i *NodeJPS, nbs *map[int]NodeJPS) *map[int]NodeJPS {
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

	coord := ug.GridToCoord(in.grid)

	if m > 0 {
		coordDown := ug.GridToCoord([]int{m - 1, n})
		neighbours = append(neighbours, neighboursRowUgJPS([]float64{coord[0], coordDown[1]}, ug, 1)...)
	}

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
