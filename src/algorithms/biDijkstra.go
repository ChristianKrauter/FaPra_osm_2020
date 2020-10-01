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

	var met = false
	var bestDist = math.Inf(1)
	var bestDistAt int

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
		if len(pq1) == 0 && len(pq2) == 0 {
			break
		} else {
			if(len(pq1) > 0){
				u1 := heap.Pop(&pq1).(*Item).value
				nodesProcessed1[u1] = true;
			
				popped1++
		
				if _, ok := nodesProcessed2[u1]; ok {
					met = true
				}
				expandNodeBiDijkstra(&u1,!met,ug,&dist1,&dist2,&prev1,&pq1, &bestDist ,&bestDistAt )
			}

			if(len(pq2) > 0){
				u2 := heap.Pop(&pq2).(*Item).value
				nodesProcessed2[u2] = true;
				popped2++

				if _, ok := nodesProcessed1[u2]; ok {
					met = true
				}
				expandNodeBiDijkstra(&u2,!met,ug,&dist2, &dist1,&prev2,&pq2, &bestDist ,&bestDistAt)
			}			
		}
	}
	var route = [][][]float64{ExtractRouteBiUg(&prev1, &prev2, bestDistAt, ug)}
	fmt.Printf("bi dist: %v\n", bestDist)
	return route , (popped1 + popped2)

}

// DijkstraAllNodes additionally returns all visited nodes on uniform grid
func BiDijkstraAllNodes(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {
	var popped1 int
	var dist1 []float64
	var prev1 []int

	var met = false
	var bestDist = math.Inf(1)
	var bestDistAt int

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
		if len(pq1) == 0 && len(pq2) == 0 {
			break
		} else {
			if(len(pq1) > 0){
				u1 := heap.Pop(&pq1).(*Item).value
				nodesProcessed1[u1] = true;
			
				popped1++
		
				if _, ok := nodesProcessed2[u1]; ok {
					met = true
				}
				expandNodeBiDijkstra(&u1,!met,ug,&dist1,&dist2,&prev1,&pq1, &bestDist ,&bestDistAt )
			}

			if(len(pq2) > 0){
				u2 := heap.Pop(&pq2).(*Item).value
				nodesProcessed2[u2] = true;
				popped2++

				if _, ok := nodesProcessed1[u2]; ok {
					met = true
				}
				expandNodeBiDijkstra(&u2,!met,ug,&dist2, &dist1,&prev2,&pq2, &bestDist ,&bestDistAt)
			}			
		}
	}
	var route = [][][]float64{ExtractRouteBiUg(&prev1, &prev2, bestDistAt, ug)}
	for k,v := range nodesProcessed1 {
		nodesProcessed2[k] = v
	}
	var allNodes []int 
	for k,_ := range nodesProcessed2 {
		allNodes = append(allNodes,k)
	}
	var processedNodes = ExtractNodesUg(&allNodes, ug)
	fmt.Printf("bi dist: %v\n", bestDist)
	return route,processedNodes
}
