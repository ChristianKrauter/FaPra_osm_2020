package algorithms

import (
    "../grids"
    "fmt"
    "math"
)

func neighboursBg(idx, xSize int, bg *grids.BasicGrid) []int {
    var neighbours []int
    var result []int

    neighbours = append(neighbours, idx-xSize-1) // top left
    neighbours = append(neighbours, idx-xSize)   // top
    neighbours = append(neighbours, idx-xSize+1) // top right
    neighbours = append(neighbours, idx-1)       // left
    neighbours = append(neighbours, idx+1)       // right
    neighbours = append(neighbours, idx+xSize-1) // bottom left
    neighbours = append(neighbours, idx+xSize)   // bottom
    neighbours = append(neighbours, idx+xSize+1) // bottom right

    for _, j := range neighbours {
        if j >= 0 && j < len((*bg).VertexData) {
            if !(*bg).VertexData[j] {
                result = append(result, j)
            }
        }
    }
    return result
}

func extractRoute(prev *[]int, end int, bg *grids.BasicGrid) [][][]float64 {
    var route [][][]float64
    var tempRoute [][]float64
    temp := bg.ExpandIndex(end)
    for {
        x := bg.ExpandIndex(end)
        if math.Abs(float64(temp[0]-x[0])) > 1 {
            fmt.Printf("helpersBG\n")
            route = append(route, tempRoute)
            tempRoute = make([][]float64, 0)
        }
        tempRoute = append(tempRoute, bg.GridToCoord([]int{x[0], x[1]}))

        if (*prev)[end] == -1 {
            break
        }
        end = (*prev)[end]
        temp = x
    }
    route = append(route, tempRoute)
    return route
}

func extractNodes(nodesProcessed *[]int, bg *grids.BasicGrid) [][]float64 {
    var nodesExtended [][]float64
    for _, node := range *nodesProcessed {
        x := bg.ExpandIndex(node)
        coord := bg.GridToCoord([]int{x[0], x[1]})
        nodesExtended = append(nodesExtended, coord)
    }
    return nodesExtended
}
