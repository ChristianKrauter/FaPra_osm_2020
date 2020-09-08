package algorithms

import (
	"../grids"
	"container/heap"
	"fmt"
	"github.com/paulmach/go.geojson"
	"math"
	"os"
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
			break
		} else {
			u := heap.Pop(&pq).(*NodeJPS)
			popped++
			if u.IDX == to {
				return ExtractRouteUg(&prev, to, ug), popped
			}
			// ToDo: Neighbour pruning
			//if u.dir == 3 || u.dir == 4 {
			//	fmt.Printf("hor: %v\n", hor(u, to, &dist, ug))
			//}

			neighbours := neighboursUgJPS(*u, ug)
			neighbours = prune(u, to, neighbours, ug)
			fmt.Printf("pruned: %v\n", neighbours)
			//fmt.Printf("nb: %v\n", neighbours)

			for _, j := range *neighbours {
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
	return ExtractRouteUg(&prev, to, ug), popped
}

func prune(i *NodeJPS, to int, nbs *map[int]NodeJPS, ug *grids.UniformGrid) *map[int]NodeJPS {
	var res = make(map[int]NodeJPS)
	fmt.Printf("dir: %v\n", i.dir)
	if i.dir == -1 {
		return nbs
	}
	// forced neighbours to check
	var n [][]int
	// natural neighbours to check
	var m []int
	m = append(m, i.dir)
	//res[i.dir] = (*nbs)[i.dir]
	// Right
	switch i.dir {
	case 4:
		n = append(n, []int{1, 2})
		n = append(n, []int{6, 7})
		/*if val, ok := (*nbs)[1]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				if val2, ok := (*nbs)[2]; ok {
					if ug.VertexData[val2.grid[0]][val2.grid[1]] {
						res[2] = (*nbs)[2]
					}
				}
			}
		}
		if val, ok := (*nbs)[6]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[7] = (*nbs)[7]
			}
		}*/

	case 3:
		// Left
		n = append(n, []int{1, 0})
		n = append(n, []int{6, 5})
		/*if val, ok := (*nbs)[1]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[0] = (*nbs)[0]
			}
		}
		if val, ok := (*nbs)[6]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[5] = (*nbs)[5]
			}
		}*/

	case 1:
		// Up
		n = append(n, []int{3, 0})
		n = append(n, []int{4, 2})
		/*if val, ok := (*nbs)[3]; ok {

			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[0] = (*nbs)[0]
			}
		}
		if val, ok := (*nbs)[4]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[2] = (*nbs)[2]
			}
		}*/

	case 6:
		// Down
		n = append(n, []int{4, 7})
		n = append(n, []int{3, 5})
		/*if val, ok := (*nbs)[4]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[7] = (*nbs)[7]
			}
		}
		if val, ok := (*nbs)[3]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[5] = (*nbs)[5]
			}
		}*/

	case 2:
		// Top Right
		// res[1] = (*nbs)[1]
		// res[4] = (*nbs)[4]
		m = append(m, 1)
		m = append(m, 4)

		n = append(n, []int{3, 0})
		n = append(n, []int{6, 7})
		/*if val, ok := (*nbs)[3]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[0] = (*nbs)[0]
			}
		}
		if val, ok := (*nbs)[6]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[7] = (*nbs)[7]
			}
		}*/

	case 0:
		// Top Left
		// res[1] = (*nbs)[1]
		// res[3] = (*nbs)[3]
		m = append(m, 1)
		m = append(m, 3)

		n = append(n, []int{6, 0})
		n = append(n, []int{4, 2})
		/*if val, ok := (*nbs)[6]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[0] = (*nbs)[0]
			}
		}
		if val, ok := (*nbs)[4]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[2] = (*nbs)[2]
			}
		}*/

	case 7:
		// Bottom Right
		// res[4] = (*nbs)[4]
		// res[6] = (*nbs)[6]
		m = append(m, 4)
		m = append(m, 6)

		n = append(n, []int{1, 2})
		n = append(n, []int{3, 5})
		/*if val, ok := (*nbs)[1]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[2] = (*nbs)[2]
			}
		}
		if val, ok := (*nbs)[3]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[5] = (*nbs)[5]
			}
		}*/

	case 5:
		// Bottom Left
		// res[3] = (*nbs)[3]
		// res[6] = (*nbs)[6]
		m = append(m, 3)
		m = append(m, 6)

		n = append(n, []int{1, 0})
		n = append(n, []int{4, 7})
		/*if val, ok := (*nbs)[1]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[0] = (*nbs)[0]
			}
		}
		if val, ok := (*nbs)[4]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				res[7] = (*nbs)[7]
			}
		}*/
	}
	for _, i := range m {
		if val2, ok := (*nbs)[i]; ok {
			if !ug.VertexData[val2.grid[0]][val2.grid[1]] {
				res[i] = (*nbs)[i]
			}
		}
	}
	for _, i := range n {
		if val, ok := (*nbs)[i[0]]; ok {
			if ug.VertexData[val.grid[0]][val.grid[1]] {
				if val2, ok := (*nbs)[i[1]]; ok {
					if !ug.VertexData[val2.grid[0]][val2.grid[1]] {
						res[i[1]] = (*nbs)[i[1]]
					}
				}
			}
		}
	}
	return &res
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func hor(i *NodeJPS, to int, dist *[]float64, ug *grids.UniformGrid) []NodeJPS {
	//var x0 = i.grid[0]
	//var y0 = i.grid[1]
	var nodes []NodeJPS
	var nbs *map[int]NodeJPS
	//var d = (*dist)[i.IDX]
	for {
		nbs = neighboursUgJPS(*i, ug)
		fmt.Printf("i       :: %v\n", i)
		fmt.Printf("nbs     :: %v\n", nbs)
		//fmt.Printf("nbs[dir]:: %v\n", nbs[i.dir])
		var nxt = (*nbs)[i.dir]
		fmt.Printf("next: %v\n", nxt)

		// Back at the beginning?
		if nxt.IDX == i.IDX {
			return make([]NodeJPS, 0)
		}
		if ug.VertexData[nxt.grid[0]][nxt.grid[1]] {
			return make([]NodeJPS, 0)
		}
		if nxt.IDX == to {
			return []NodeJPS{nxt}
		}
		nodes = append(nodes, hor(&nxt, to, dist, ug)...)
	}

	var fc = geojson.NewFeatureCollection()
	for _, i := range nodes {
		//fmt.Printf("i         : %v\n", i)
		//fmt.Printf("dir, coord: %v, %v\n", i.dir, ug.GridToCoord(i.grid))
		fc.AddFeature(geojson.NewPointFeature(ug.GridToCoord(i.grid)))
	}

	rawJSON, err := fc.MarshalJSON()
	f, err := os.Create("./nbs.geojson")
	check(err)

	_, err1 := f.Write(rawJSON)
	check(err1)
	f.Sync()
	return make([]NodeJPS, 0)
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
