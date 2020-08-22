package algorithms

import (
    "math"
)

// BasicGrid structure
type BasicGrid struct {
    XSize      int
    YSize      int
    VertexData []bool
}

// GridToCoord transforms a grid index to lat lon coordinates
func (bg BasicGrid) GridToCoord(in []int) []float64 {
    var out []float64
    var xFactor = bg.XSize / 360
    var yFactor = bg.YSize / 360
    out = append(out, float64(in[0]/xFactor-180))
    out = append(out, float64((in[1]/yFactor)/2-90))
    return out
}

// CoordToGrid transforms lat lon coordinates to a grid index
func (bg BasicGrid) CoordToGrid(in []float64) []int {
    var out []int
    var xFactor = float64(bg.XSize / 360)
    var yFactor = float64(bg.YSize / 360)
    out = append(out, int(((math.Round(in[0]*xFactor)/xFactor)+180)*xFactor))
    out = append(out, int(((math.Round(in[1]*yFactor)/yFactor)+90)*2*yFactor))
    return out
}

func (bg BasicGrid) flattenIndex(x, y int) int {
    return ((bg.XSize * y) + x)
}

// ExpandIndex from 1d to 2d
func (bg BasicGrid) ExpandIndex(idx int) []int {
    return []int{idx % bg.XSize, idx / bg.XSize}
}

func neighboursBg(idx, xSize int, bg *BasicGrid) []int {
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

func extractRoute(prev *[]int, end int, bg *BasicGrid) [][][]float64 {
    var route [][][]float64
    var tempRoute [][]float64
    temp := bg.ExpandIndex(end)
    for {
        x := bg.ExpandIndex(end)
        if math.Abs(float64(temp[0]-x[0])) > 1 {
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

func extractNodes(nodesProcessed *[]int, bg *BasicGrid) [][]float64 {
    var nodesExtended [][]float64
    for _, node := range *nodesProcessed {
        x := bg.ExpandIndex(node)
        coord := bg.GridToCoord([]int{x[0], x[1]})
        nodesExtended = append(nodesExtended, coord)
    }
    return nodesExtended
}
