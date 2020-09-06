package grids

import (
    "math"
)

// BasicGrid structure
type BasicGrid struct {
    XSize      int
    YSize      int
    XFactor    float64
    YFactor    float64
    VertexData []bool
}

// GridToCoord transforms a grid index to lat lon coordinates
func (bg BasicGrid) GridToCoord(in []int) []float64 {
    return []float64{float64(in[0])/bg.XFactor - 180.0, float64(in[1])/bg.YFactor/2.0 - 90.0}
}

// CoordToGrid transforms lat lon coordinates to a grid index
func (bg BasicGrid) CoordToGrid(in []float64) []int {
    return []int{int(((math.Round(in[0]*bg.XFactor) / bg.XFactor) + 180.0) * bg.XFactor), int(((math.Round(in[1]*bg.YFactor) / bg.YFactor) + 90.0) * 2 * bg.YFactor)}
}

// GridToID 2D to 1D Index
func (bg BasicGrid) GridToID(IDX []int) int {
    return bg.XSize*IDX[1] + IDX[0]
}

// IDToGrid 1D to 2D Index
func (bg BasicGrid) IDToGrid(ID int) []int {
    return []int{ID % bg.XSize, ID / bg.XSize}
}
