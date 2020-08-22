package algorithms

import (
    "math"
)

// GridToCoord transforms a grid index to lat lon coordinates
func GridToCoord(in []int, xSize, ySize int) []float64 {
    var out []float64
    var xFactor = xSize / 360
    var yFactor = ySize / 360
    out = append(out, float64(in[0]/xFactor-180))
    out = append(out, float64((in[1]/yFactor)/2-90))
    return out
}

// CoordToGrid transforms lat lon coordinates to a grid index
func CoordToGrid(in []float64, xSize, ySize int) []int {
    var out []int
    var xFactor, yFactor float64
    xFactor = float64(xSize / 360)
    yFactor = float64(ySize / 360)
    // TODO check
    out = append(out, int(((math.Round(in[0]*xFactor)/xFactor)+180)*xFactor))
    out = append(out, int(((math.Round(in[1]*yFactor)/yFactor)+90)*2*yFactor))
    return out
}

func flattenIndex(x, y, xSize int) int {
    return ((xSize * y) + x)
}

// ExpandIndex from 1d to 2d
func ExpandIndex(indx, xSize int) []int {
    var x = indx % xSize
    var y = indx / xSize
    return []int{x, y}
}

func neighboursBg(indx, xSize int, mg *[]bool) []int {
    var neighbours []int
    var result []int

    neighbours = append(neighbours, indx-xSize-1) // top left
    neighbours = append(neighbours, indx-xSize)   // top
    neighbours = append(neighbours, indx-xSize+1) // top right
    neighbours = append(neighbours, indx-1)       // left
    neighbours = append(neighbours, indx+1)       // right
    neighbours = append(neighbours, indx+xSize-1) // bottom left
    neighbours = append(neighbours, indx+xSize)   // bottom
    neighbours = append(neighbours, indx+xSize+1) // bottom right

    for _, j := range neighbours {
        if j >= 0 && j < int(len((*mg))) {
            if !(*mg)[j] {
                result = append(result, j)
            }
        }
    }
    return result
}

func extractRoute(prev *[]int, end, xSize, ySize int) [][][]float64 {
    var route [][][]float64
    var tempRoute [][]float64
    temp := ExpandIndex(end, xSize)
    for {
        x := ExpandIndex(end, xSize)
        if math.Abs(float64(temp[0]-x[0])) > 1 {
            route = append(route, tempRoute)
            tempRoute = make([][]float64, 0)
        }
        tempRoute = append(tempRoute, GridToCoord([]int{x[0], x[1]}, xSize, ySize))

        if (*prev)[end] == -1 {
            break
        }
        end = (*prev)[end]
        temp = x
    }
    route = append(route, tempRoute)
    return route
}

func extractNodes(nodesProcessed *[]int, xSize, ySize int) [][]float64 {
    var nodesExtended [][]float64
    for _, node := range *nodesProcessed {
        x := ExpandIndex(node, xSize)
        coord := GridToCoord([]int{x[0], x[1]}, xSize, ySize)
        nodesExtended = append(nodesExtended, coord)
    }
    return nodesExtended
}
