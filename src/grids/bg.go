package grids

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
    var xFactor = bg.XSize / 360.0
    var yFactor = bg.YSize / 360.0
    out = append(out, float64(in[0]/xFactor-180.0))
    out = append(out, float64((in[1]/yFactor)/2.0-90.0))
    return out
}

// CoordToGrid transforms lat lon coordinates to a grid index
func (bg BasicGrid) CoordToGrid(in []float64) []int {
    var out []int
    var xFactor = float64(bg.XSize / 360.0)
    var yFactor = float64(bg.YSize / 360.0)
    out = append(out, int(((math.Round(in[0]*xFactor)/xFactor)+180.0)*xFactor))
    out = append(out, int(((math.Round(in[1]*yFactor)/yFactor)+90.0)*2*yFactor))
    return out
}

// FlattenIndex from 2d to 1d
func (bg BasicGrid) FlattenIndex(IDX []int) int {
    return ((bg.XSize * IDX[1]) + IDX[0])
}

// ExpandIndex from 1d to 2d
func (bg BasicGrid) ExpandIndex(idx int) []int {
    return []int{idx % bg.XSize, idx / bg.XSize}
}
