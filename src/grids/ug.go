package grids

import (
    "math"
    "sort"
)

// UniformGrid structure
type UniformGrid struct {
    XSize        int
    YSize        int
    N            int
    BigN         int
    A            float64
    D            float64
    MTheta       float64
    DTheta       float64
    DPhi         float64
    VertexData   [][]bool
    FirstIndexOf []int
}

// GridToCoord takes grid indices and outputs lng lat
func (ug UniformGrid) GridToCoord(in []int) []float64 {
    theta := math.Pi * (float64(in[0]) + 0.5) / float64(ug.MTheta)
    mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
    phi := 2 * math.Pi * float64(in[1]) / mPhi
    return []float64{(phi / math.Pi) * 180, (theta/math.Pi)*180 - 90}
}

// CoordToGrid takes lng lat and outputs grid indices
func (ug UniformGrid) CoordToGrid(lng, lat float64) []int {
    theta := (lat + 90) * math.Pi / 180
    m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
    phi := lng * math.Pi / 180
    mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
    n := math.Round(phi * mPhi / (2 * math.Pi))
    return []int{mod(int(m), int(ug.MTheta)), mod(int(n), int(mPhi))}
}

// GridToID ...
func (ug UniformGrid) GridToID(IDX []int) int {
    return ug.FirstIndexOf[IDX[0]] + IDX[1]
}

// IDToGrid ...
func (ug UniformGrid) IDToGrid(id int) []int {
    m := sort.Search(len(ug.FirstIndexOf)-1, func(i int) bool { return ug.FirstIndexOf[i] > id })
    n := id - ug.FirstIndexOf[m-1]
    return []int{m - 1, n}
}

func mod(a, b int) int {
    a = a % b
    if a >= 0 {
        return a
    }
    if b < 0 {
        return a - b
    }
    return a + b
}
