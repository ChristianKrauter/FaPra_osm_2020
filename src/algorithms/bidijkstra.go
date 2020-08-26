package algorithms

import (
	"../grids"
	"container/heap"
	"fmt"
	"math"
)

// BiDijkstra implementation on uniform grid
func BiDijkstra(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, int) {

	var prev [][]int
	var dist []float64
	var distR []float64
	pq := make(priorityQueue, 1)
	pqR := make(priorityQueue, 1)
	proc := make(map[int]bool)
	procR := make(map[int]bool)

	var meeting int
	var bestDist = math.Inf(1)
	var minF = math.Inf(1)
	var minR = math.Inf(1)

	for i := 0; i < (*ug).N; i++ {
		dist = append(dist, math.Inf(1))
		distR = append(dist, math.Inf(1))
		prev = append(prev, []int{-1, -1})
	}

	dist[(*ug).GridToID(fromIDX)] = 0
	pq[0] = &Item{
		value:    (*ug).GridToID(fromIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq)

	// Reverse
	distR[(*ug).GridToID(toIDX)] = 0
	pqR[0] = &Item{
		value:    (*ug).GridToID(toIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pqR)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*Item).value

			if procR[u] {
				fmt.Printf("bestDist         %v\n", bestDist)
				fmt.Printf("dist[u]+distR[u] %v\n", dist[u]+distR[u])
				bestDist = math.Min(bestDist, dist[u]+distR[u])
				//var minF = math.MaxFloat64
				//var minR = math.MaxFloat64
				for _, k := range pqR {
					if proc[k.value] {
						//fmt.Printf("[k]: %v\n", k.value)
						if distR[k.value] < minR {
							//fmt.Printf("distR < minR: %v < %v\n", distR[k.value], minR)
							minR = distR[k.value]
						} /*else  {
							fmt.Printf("distR >= minR: %v >= %v\n", distR[k.value], minR)
						}*/
					}
				}

				for _, k := range pq {
					if procR[k.value] {
						//fmt.Printf("[k]: %v\n", k.value)
						if dist[k.value] < minF {
							//fmt.Printf("dist < min: %v < %v\n", distR[k.value], minR)
							minF = dist[k.value]
						} /*else {
							fmt.Printf("dist >= minF: %v >= %v\n", dist[k.value], minF)
						} */
					}
				}
				fmt.Printf("minR+minF - bestDist: %v - %v\n", minR+minF, bestDist)
				if minR+minF >= bestDist {
					//meeting = u
					break
				}
			}

			neighbours := NeighboursUg(u, ug)
			for _, j := range neighbours {
				var alt = dist[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					/*if dist[j] < minF {
						minF = dist[j]
					}*/
					prev[j][0] = u
					item := &Item{
						value:    j,
						priority: -dist[j],
					}
					heap.Push(&pq, item)
				}
			}
			proc[u] = true

		}

		// Reverse
		if len(pqR) == 0 {
			break
		} else {
			u := heap.Pop(&pqR).(*Item).value

			if proc[u] {
				fmt.Printf("bestDist    %v\n", bestDist)
				bestDist = math.Min(bestDist, dist[u]+distR[u])
				fmt.Printf("bestDistNEW %v\n", bestDist)
				//var minF = math.MaxFloat64
				//var minR = math.MaxFloat64
				for _, k := range pqR {
					if proc[k.value] {
						//fmt.Printf("[k]: %v\n", k.value)
						if distR[k.value] < minR {
							//fmt.Printf("distR < minR: %v < %v\n", distR[k.value], minR)
							minR = distR[k.value]
						} /*else  {
							fmt.Printf("distR >= minR: %v >= %v\n", distR[k.value], minR)
						}*/
					}
				}

				for _, k := range pq {
					if procR[k.value] {
						//fmt.Printf("[k]: %v\n", k.value)
						if dist[k.value] < minF {
							//fmt.Printf("dist < min: %v < %v\n", distR[k.value], minR)
							minF = dist[k.value]
						} /*else {
							fmt.Printf("dist >= minF: %v >= %v\n", dist[k.value], minF)
						} */
					}
				}

				fmt.Printf("minR+minF - bestDist: %v - %v\n", minR+minF, bestDist)
				if minR+minF >= bestDist {
					meeting = u
					break
				}
			}

			neighbours := NeighboursUg(u, ug)
			for _, j := range neighbours {
				var alt = distR[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt < distR[j] {
					distR[j] = alt
					prev[j][1] = u
					item := &Item{
						value:    j,
						priority: -distR[j],
					}
					heap.Push(&pqR, item)
				}
			}
			procR[u] = true
		}

	}
	var route = ExtractRouteUgBi(&prev, meeting, ug)
	return route, len(procR) + len(proc)
}

// BiDijkstraAllNodes additionally returns all visited nodes on uniform grid
func BiDijkstraAllNodes(fromIDX, toIDX []int, ug *grids.UniformGrid) ([][][]float64, [][]float64) {

	var prev [][]int
	var dist []float64
	var distR []float64
	pq := make(priorityQueue, 1)
	pqR := make(priorityQueue, 1)
	proc := make(map[int]bool)
	procR := make(map[int]bool)

	var meeting int
	var bestDist = math.Inf(1)

	for i := 0; i < (*ug).N; i++ {
		dist = append(dist, math.Inf(1))
		distR = append(dist, math.Inf(1))
		prev = append(prev, []int{-1, -1})
	}

	dist[(*ug).GridToID(fromIDX)] = 0
	pq[0] = &Item{
		value:    (*ug).GridToID(fromIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pq)

	// Reverse
	distR[(*ug).GridToID(toIDX)] = 0
	pqR[0] = &Item{
		value:    (*ug).GridToID(toIDX),
		priority: 0,
		index:    0,
	}
	heap.Init(&pqR)

	for {
		if len(pq) == 0 {
			break
		} else {
			u := heap.Pop(&pq).(*Item).value
			proc[u] = true

			if procR[u] {
				var minF = math.Inf(1)
				var minR = math.Inf(1)
				bestDist = math.Min(bestDist, dist[u]+distR[u])

				for _, k := range pqR {
					//if proc[k.value] {
					if distR[k.value] < minR {
						minR = distR[k.value]
						meeting = k.value
					}
					//}
				}

				for _, k := range pq {
					//if procR[k.value] {
					if dist[k.value] < minF {
						minF = dist[k.value]
						meeting = k.value
					}
					//}
				}

				fmt.Printf("minR+minF >=? bestDist: %v - %v\n", minR+minF, bestDist)
				if minR+minF >= bestDist {
					meeting = u
					fmt.Printf("prev[meeting][0]: %v\n", prev[meeting][0])
					fmt.Printf("prev[meeting][1]: %v\n", prev[meeting][1])
					break
				}
			}

			neighbours := NeighboursUg(u, ug)
			for _, j := range neighbours {
				var alt = dist[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j][0] = u
					item := &Item{
						value:    j,
						priority: -dist[j],
					}
					heap.Push(&pq, item)
				}
			}
		}

		// Reverse
		if len(pqR) == 0 {
			break
		} else {
			u := heap.Pop(&pqR).(*Item).value
			procR[u] = true

			if proc[u] {
				var minF = math.Inf(1)
				var minR = math.Inf(1)
				bestDist = math.Min(bestDist, dist[u]+distR[u])

				for _, k := range pqR {
					//if proc[k.value] {
					if distR[k.value] < minR {
						minR = distR[k.value]
						meeting = k.value
					}
					//}
				}

				for _, k := range pq {
					//if procR[k.value] {
					if dist[k.value] < minF {
						minF = dist[k.value]
						meeting = k.value
					}
					//}
				}

				fmt.Printf("minR+minF >=? bestDist: %v - %v\n", minR+minF, bestDist)
				if minR+minF >= bestDist {
					meeting = u
					fmt.Printf("prev[meeting][0]: %v\n", prev[meeting][0])
					fmt.Printf("prev[meeting][1]: %v\n", prev[meeting][1])
					break
				}
			}

			neighbours := NeighboursUg(u, ug)
			for _, j := range neighbours {
				var alt = distR[u] + distance((*ug).GridToCoord((*ug).IDToGrid(u)), (*ug).GridToCoord((*ug).IDToGrid(j)))
				if alt < distR[j] {
					distR[j] = alt
					prev[j][1] = u
					item := &Item{
						value:    j,
						priority: -distR[j],
					}
					heap.Push(&pqR, item)
				}
			}
		}

	}

	keys := make([]int, len(proc)+len(procR))
	i := 0
	for k := range proc {
		keys[i] = k
		i++
	}
	for k := range procR {
		keys[i] = k
		i++
	}
	var processedNodes = ExtractNodesUg(&keys, ug)
	var route = ExtractRouteUgBi(&prev, meeting, ug)
	return route, processedNodes
}
