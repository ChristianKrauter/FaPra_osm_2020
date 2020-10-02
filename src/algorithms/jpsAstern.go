package algorithms

import(
	"../grids"
	"container/heap"
	"math"
	"fmt"
)



// JPSAstern implementation on uniform grid
func JPSAstern(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, int) {
	var popped int
	var dist []float64
	var prev []int
	pq := make(jpsPriorityQueue, 1)

	for i := 0; i < (*ug).N; i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[(*ug).GridToID(fromIDX)] = 0
	pq[0] = &JPSItem{
		value:    (*ug).GridToID(fromIDX),
		priority: 0,
		index:    0,
		direction: -1,
	}
	heap.Init(&pq)
	for {
		if len(pq) == 0 {
			break
		} else {
			item := heap.Pop(&pq).(*JPSItem)
			u := item.value
			direction := item.direction
			//u := heap.Pop(&pq).(*Item).value
			popped++

			if u == (*ug).GridToID(toIDX) {
				fmt.Printf("JPS astern dist: %v\n", dist[u]+ distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord(toIDX)))
				return ExtractRouteUg(&prev, (*ug).GridToID(toIDX), ug), popped
			}
			succesors, directions := IdentifySuccessors(u, direction, fromIDX, toIDX, ug)
			//neighbours := NeighboursUg(&u, ug)

			for i, j := range succesors {
				var alt = dist[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &JPSItem{
						value:    j,
						priority: -dist[j] - distance((*ug).GridToCoord((*ug).IDToGrid(j)), (*ug).GridToCoord(toIDX)),
						direction: directions[i],
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	return ExtractRouteUg(&prev, (*ug).GridToID(toIDX), ug), popped
}

// AsternAllNodes additionally returns all visited nodes on uniform grid
func JPSAsternAllNodes(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {
	var dist []float64
	var prev []int
	var nodesProcessed []int
	pq := make(jpsPriorityQueue, 1)

	for i := 0; i < (*ug).N; i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}

	dist[(*ug).GridToID(fromIDX)] = 0
	pq[0] = &JPSItem{
		value:    (*ug).GridToID(fromIDX),
		priority: 0,
		index:    0,
		direction: -1,
	}
	heap.Init(&pq)
	
	for {
		if len(pq) == 0 {
			break
		} else {
			item := heap.Pop(&pq).(*JPSItem)
			u := item.value
			direction := item.direction
			nodesProcessed = append(nodesProcessed, u)

			if u == (*ug).GridToID(toIDX) {
				fmt.Printf("JPS astern dist: %v\n", dist[u]+ distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord(toIDX)))
				var route = ExtractRouteUg(&prev, (*ug).GridToID(toIDX), ug)
				var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
				return route, processedNodes
			}
			
			succesors, directions := IdentifySuccessors(u, direction, fromIDX, toIDX, ug)
			
			for i, j := range succesors {
				var alt = dist[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					item := &JPSItem{
						value:    j,
						priority: -dist[j] - distance((*ug).GridToCoord((*ug).IDToGrid(j)), (*ug).GridToCoord(toIDX)),
						direction: directions[i],
					}
					heap.Push(&pq, item)
				}
			}
		}
	}
	var route = ExtractRouteUg(&prev, (*ug).GridToID(toIDX), ug)
	var processedNodes = ExtractNodesUg(&nodesProcessed, ug)
	return route, processedNodes
}


func IdentifySuccessors(x, dir int, fromIDX, toIDX []int, ug *grids.UniformGrid) ([]int,[]int){
	dirs := JPSNeighboursUg(x, dir, ug);
	var succesors []int
	var directions []int
	var sGrid [][]int
	for _, j := range dirs{
		succ := jump(x, x, &j, fromIDX, toIDX, ug)
		if(succ != nil){
			sGrid = append(sGrid,succ)
			succesors = append(succesors, ug.GridToID(succ))
			directions = append(directions, j)
		}
	}
	//fmt.Printf("in:%v\nn: %v\ndir: %v\n\n",ug.IDToGrid(x),sGrid,directions)
	return succesors,directions
}


func jump(origin,x int, dir *int, fromIDX, toIDX []int, ug *grids.UniformGrid) []int{
	n := step(origin, x, dir, ug)
	//fmt.Printf("- ")
	//*jumpLookAts= append(*jumpLookAts, ug.GridToCoord(ug.IDToGrid(x)))
	//xCoord := ug.IDToGrid(x)
	//fmt.Printf("xdiff: %v,\nydiff: %v\n\n", xCoord, n)
	if(n == nil || ug.VertexData[n[0]][n[1]]){		
		return nil
	}
	
	if(n[0] == toIDX[0] && n[1] == toIDX[1]){
		return n
	}
	if(isForced(ug.GridToID(n),*dir, ug)){
		return n 
	}
	if(*dir == 0 || *dir == 2 || *dir == 4 || *dir == 6){
		newDir1 := mod((*dir)-1,8)
		temp := jump(ug.GridToID(n),ug.GridToID(n),&newDir1,fromIDX, toIDX, ug)
		if(temp != nil){
			return n
		}
		newDir2 := mod((*dir)+1,8)
		temp = jump(ug.GridToID(n),ug.GridToID(n),&newDir2,fromIDX, toIDX, ug)
		if(temp != nil){
			return n
		}
	}
	return jump(origin, ug.GridToID(n), dir, fromIDX, toIDX, ug)
}

func step(origin, in int,dir *int, ug *grids.UniformGrid) []int{
	var allNeighbours [8][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	var lm = len(ug.VertexData[m])
	allNeighbours[7] = []int{m, mod(n-1, lm)}	
	allNeighbours[3] = []int{m, mod(n+1, lm)}	
	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData) -1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))
		allNeighbours[5] = []int{m + 1, mod(nUp, lmp)}
		allNeighbours[4] = []int{m + 1, mod(nUp+1.0, lmp)}
		allNeighbours[6] = []int{m + 1, mod(nUp-1.0, lmp)}
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))
		allNeighbours[1] = []int{m - 1, mod(nDown, lmm)}
		allNeighbours[2] = []int{m - 1, mod(nDown+1.0, lmm)}
		allNeighbours[0] = []int{m - 1, mod(nDown-1.0, lmm)}
	}
	if allNeighbours[*dir][0] >= len(ug.VertexData) -1 {
		*dir = mod((*dir)+4,8)
	} 
	if len(allNeighbours[*dir]) < 1 || ug.GridToID(allNeighbours[*dir]) == origin {
	 	return nil
	} else {
	 	return allNeighbours[*dir]
	}
	
}

//jump function to show all visited nodes
/*
func jumpLookAt(x, dir int, fromIDX, toIDX []int, ug *grids.UniformGrid, jumpLookAts *[][]float64) []int{
	n := step(x,dir,ug, jumpLookAts)
	//fmt.Printf("- ")
	//*jumpLookAts= append(*jumpLookAts, ug.GridToCoord(ug.IDToGrid(x)))
	//xCoord := ug.IDToGrid(x)
	//fmt.Printf("xdiff: %v,\nydiff: %v\n\n", xCoord, n)
	if(n == nil || ug.VertexData[n[0]][n[1]]){		
		return nil
	}
	
	if(n[0] == toIDX[0] && n[1] == toIDX[1]){
		return n
	}
	if(isForced(ug.GridToID(n),dir, ug)){
		return n 
	}
	if(dir == 0 || dir == 2 || dir == 4 || dir == 6){
		temp := jump(ug.GridToID(n),mod(dir-1,8),fromIDX, toIDX, ug, jumpLookAts)
		if(temp != nil){
			return n
		}
		temp = jump(ug.GridToID(n),mod(dir+1,8),fromIDX, toIDX, ug, jumpLookAts)
		if(temp != nil){
			return n
		}
	}
	return jump(ug.GridToID(n),dir,fromIDX, toIDX, ug, jumpLookAts)
}*/