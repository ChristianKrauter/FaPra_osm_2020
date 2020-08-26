package algorithms

import(
	"../grids"
	"container/heap"
	"math"
)


// Astern implementation on uniform grid
func BiAstern(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, int) {

	var popped1 int
	var dist1 []float64
	var prev1 []int
	pq1 := make(priorityQueue, 1)
	nodesProcessed1 := make(map[int]bool)

	var popped2 int
	var dist2 []float64
	var prev2 []int
	pq2 := make(priorityQueue, 1)
	nodesProcessed2 := make(map[int]bool)

	for i := 0; i < (*ug).N; i++ {
		dist1 = append(dist1, math.Inf(1))
		prev1 = append(prev1, -1)
		dist2 = append(dist2, math.Inf(1))
		prev2 = append(prev2, -1)
	}

	dist1[(*ug).GridToID(fromIDX)] = 0
	dist2[(*ug).GridToID(toIDX)] = 0

	pq1[0] = &Item{
		value:    (*ug).GridToID(fromIDX),
		priority: 0,
		index:    0,
	}

	pq2[0] = &Item{
		value:    (*ug).GridToID(toIDX),
		priority: 0,
		index:    0,
	}

	heap.Init(&pq1)
	heap.Init(&pq2)
	

	for {
		if len(pq1) == 0 || len(pq2) == 0 {
			break
		} else {
			u1 := heap.Pop(&pq1).(*Item).value
			nodesProcessed1[u1] = true;
			popped1++
			if _, ok := nodesProcessed2[u1]; ok {
				
				return [][][]float64{ExtractRouteBiUg(&prev1, &prev2, u1, ug)}, popped1+popped2
			}
			
			
			neighbours1 := NeighboursUg(u1, ug)
			for _, j := range neighbours1 {
				var alt1 = dist1[u1] + distance((*ug).GridToCoord((*ug).IDToGrid(u1)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt1 < dist1[j] {
					dist1[j] = alt1
					prev1[j] = u1
					item := &Item{
						value:    j,
						priority: -dist1[j] - distance((*ug).GridToCoord((*ug).IDToGrid(j)), (*ug).GridToCoord(toIDX)),
					}
					heap.Push(&pq1, item)
				}
			}

			u2 := heap.Pop(&pq2).(*Item).value
			nodesProcessed2[u2] = true;
			popped2++
			if _, ok := nodesProcessed1[u2]; ok {
				
				return [][][]float64{ExtractRouteBiUg(&prev1, &prev2, u2, ug)}, popped1+popped2
			}

			neighbours2 := NeighboursUg(u2, ug)
			for _, j := range neighbours2 {
				var alt2 = dist2[u2] + distance((*ug).GridToCoord((*ug).IDToGrid(u2)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt2 < dist2[j] {
					dist2[j] = alt2
					prev2[j] = u2
					item := &Item{
						value:    j,
						priority: -dist2[j] - distance((*ug).GridToCoord((*ug).IDToGrid(j)), (*ug).GridToCoord(fromIDX)),
					}
					heap.Push(&pq2, item)
				}
			}
		}
	}
	return ExtractRouteUg(&prev1, (*ug).GridToID(toIDX), ug), popped1
}

// AsternAllNodes additionally returns all visited nodes on uniform grid
func BiAsternAllNodes(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {

	
	var popped1 int
	var dist1 []float64
	var prev1 []int
	pq1 := make(priorityQueue, 1)
	nodesProcessed1 := make(map[int]bool)

	var popped2 int
	var dist2 []float64
	var prev2 []int
	pq2 := make(priorityQueue, 1)
	nodesProcessed2 := make(map[int]bool)

	for i := 0; i < (*ug).N; i++ {
		dist1 = append(dist1, math.Inf(1))
		prev1 = append(prev1, -1)
		dist2 = append(dist2, math.Inf(1))
		prev2 = append(prev2, -1)
	}

	dist1[(*ug).GridToID(fromIDX)] = 0
	dist2[(*ug).GridToID(toIDX)] = 0

	pq1[0] = &Item{
		value:    (*ug).GridToID(fromIDX),
		priority: 0,
		index:    0,
	}

	pq2[0] = &Item{
		value:    (*ug).GridToID(toIDX),
		priority: 0,
		index:    0,
	}

	heap.Init(&pq1)
	heap.Init(&pq2)
	

	for {
		if len(pq1) == 0 || len(pq2) == 0 {
			break
		} else {
			u1 := heap.Pop(&pq1).(*Item).value
			nodesProcessed1[u1] = true;
			u2 := heap.Pop(&pq2).(*Item).value
			nodesProcessed2[u2] = true;
			popped1++
			popped2++

			if _, ok := nodesProcessed2[u1]; ok {
				var route = [][][]float64{ExtractRouteBiUg(&prev1, &prev2, u1, ug)}
				for k,v := range nodesProcessed1 {
					nodesProcessed2[k] = v
				}
				var allNodes []int 
				for k,_ := range nodesProcessed2 {
					allNodes = append(allNodes,k)
				}
				var processedNodes = ExtractNodesUg(&allNodes, ug)
				return route,processedNodes
			}
			if _, ok := nodesProcessed1[u2]; ok {
				
				var route = [][][]float64{ExtractRouteBiUg(&prev1, &prev2, u2, ug)}
				for k,v := range nodesProcessed1 {
					nodesProcessed2[k] = v
				}
				var allNodes []int 
				for k,_ := range nodesProcessed2 {
					allNodes = append(allNodes,k)
				}
				var processedNodes = ExtractNodesUg(&allNodes, ug)
				return route,processedNodes
			}
			
			/*if u1 == (*ug).GridToID(toIDX) {
				fmt.Printf("u1 done")
				return ExtractRouteUg(&prev1, (*ug).GridToID(toIDX), ug), popped1
			}
			if u2 == (*ug).GridToID(fromIDX) {
				fmt.Printf("u2 done")
				return ExtractRouteUg(&prev2, (*ug).GridToID(toIDX), ug), popped2
			}*/
			
			neighbours1 := NeighboursUg(u1, ug)
			for _, j := range neighbours1 {
				var alt1 = dist1[u1] + distance((*ug).GridToCoord((*ug).IDToGrid(u1)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt1 < dist1[j] {
					dist1[j] = alt1
					prev1[j] = u1
					item := &Item{
						value:    j,
						priority: -dist1[j] - distance((*ug).GridToCoord((*ug).IDToGrid(j)), (*ug).GridToCoord(toIDX)),
					}
					heap.Push(&pq1, item)
				}
			}

			neighbours2 := NeighboursUg(u2, ug)
			for _, j := range neighbours2 {
				var alt2 = dist2[u2] + distance((*ug).GridToCoord((*ug).IDToGrid(u2)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt2 < dist2[j] {
					dist2[j] = alt2
					prev2[j] = u2
					item := &Item{
						value:    j,
						priority: -dist2[j] - distance((*ug).GridToCoord((*ug).IDToGrid(j)), (*ug).GridToCoord(fromIDX)),
					}
					heap.Push(&pq2, item)
				}
			}
		}
	}
	var route = ExtractRouteUg(&prev1, (*ug).GridToID(toIDX), ug)
	var allNodes []int 
	for k,_ := range nodesProcessed1 {
		allNodes = append(allNodes,k)
	}
	var processedNodes = ExtractNodesUg(&allNodes, ug)
	return route, processedNodes
}



