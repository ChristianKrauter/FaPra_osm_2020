package algorithms

import (
	"../grids"
	"container/heap"
	"math"
	"fmt"
)

// Dijkstra implementation on uniform grid
func BiDijkstra(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, int) {

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
				meet := u1
				bestDist := dist1[u1]  + dist2[u1]
				for _,node := range pq1 {
					if _,ok := nodesProcessed2[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}
				for _,node := range pq2 {
					if _,ok := nodesProcessed1[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}

				fmt.Printf("bi dist: %v\n", bestDist)
				return [][][]float64{ExtractRouteBiUg(&prev1, &prev2, meet, ug)}, popped1+popped2
			}
			expandNodeDijkstra(&u1,true,ug,&dist1,&prev1,&pq1)
			

			u2 := heap.Pop(&pq2).(*Item).value
			nodesProcessed2[u2] = true;
			popped2++

			if _, ok := nodesProcessed1[u2]; ok {
				meet := u2
				bestDist := dist1[u2]  + dist2[u2]
				for _,node := range pq2 {
					if _,ok := nodesProcessed1[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}
				for _,node := range pq1 {
					if _,ok := nodesProcessed2[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}
				fmt.Printf("bi dist: %v\n", bestDist)
				return [][][]float64{ExtractRouteBiUg(&prev1, &prev2, meet, ug)}, popped1+popped2
			}
			
			expandNodeDijkstra(&u2,true,ug,&dist2,&prev2,&pq2)
		}
	}
	return ExtractRouteUg(&prev1, (*ug).GridToID(toIDX), ug), popped1
}

// DijkstraAllNodes additionally returns all visited nodes on uniform grid
func BiDijkstraAllNodes(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {

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
				meet := u1
				bestDist := dist1[u1]  + dist2[u1]
				for _,node := range pq1 {
					if _,ok := nodesProcessed2[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}
				for _,node := range pq2 {
					if _,ok := nodesProcessed1[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}

				fmt.Printf("bi dist: %v\n", bestDist)
				var route = [][][]float64{ExtractRouteBiUg(&prev1, &prev2, meet, ug)}
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

			expandNodeDijkstra(&u1,true,ug,&dist1,&prev1,&pq1)


			u2 := heap.Pop(&pq2).(*Item).value
			nodesProcessed2[u2] = true;
			popped2++

			if _, ok := nodesProcessed1[u2]; ok {
				meet := u2
				bestDist := dist1[u2]  + dist2[u2]
				for _,node := range pq2 {
					if _,ok := nodesProcessed1[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}
				for _,node := range pq1 {
					if _,ok := nodesProcessed2[node.value]; ok {
						if dist1[node.value] + dist2[node.value] < bestDist{
							bestDist = dist1[node.value] + dist2[node.value]
							meet = node.value
						}
					}
				}
				fmt.Printf("bi dist: %v\n", bestDist)
				var route = [][][]float64{ExtractRouteBiUg(&prev1, &prev2, meet, ug)}
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
			
			expandNodeDijkstra(&u2,true,ug,&dist2,&prev2,&pq2)

			
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