package main

import (
	"../algorithms"
	"../grids"
)

// WayFinding for profiling
func WayFinding(from, to, algorithm int, ug *grids.UniformGrid) {
	switch algorithm {
	case 0:
		_, _, _ = algorithms.Dijkstra(from, to, ug)
	case 1:
		_, _, _ = algorithms.AStar(from, to, ug)
	case 2:
		_, _, _ = algorithms.BiDijkstra(from, to, ug)
	case 3:
		_, _, _ = algorithms.BiAStar(from, to, ug)
	case 4:
		_, _, _ = algorithms.AStarJPS(from, to, ug)
	default:
		_, _, _ = algorithms.Dijkstra(from, to, ug)
	}
}
