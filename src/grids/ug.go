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

// NewUG inits UniformGrid
func NewUG(xSize, ySize int) *UniformGrid {
    ug := UniformGrid{XSize: xSize, YSize: ySize}

    ug.BigN = xSize * ySize
    ug.A = 4.0 * math.Pi / float64(ug.BigN)
    ug.D = math.Sqrt(ug.A)
    ug.MTheta = math.Round(math.Pi / ug.D)
    ug.DTheta = math.Pi / ug.MTheta
    ug.DPhi = ug.A / ug.DTheta
    return &ug
}

// GridToCoord takes grid indices and outputs lng lat
func (ug UniformGrid) GridToCoord(in []int) []float64 {
    theta := math.Pi * (float64(in[0]) + 0.5) / float64(ug.MTheta)
    mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
    phi := 2 * math.Pi * float64(in[1]) / mPhi
    return []float64{(phi / math.Pi) * 180.0, (theta/math.Pi)*180.0 - 90.0}
}

// CoordToGrid takes lng lat and outputs grid indices
func (ug UniformGrid) CoordToGrid(lng, lat float64) []int {
    theta := (lat + 90.0) * math.Pi / 180.0
    m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
    theta = math.Pi * (float64(m) + 0.5) / float64(ug.MTheta)
    var phi float64
    if lng < 0 {
        phi = float64(lng+360.0) * math.Pi / 180.0
    } else {
        phi = lng * math.Pi / 180.0
    }
    mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
    n := math.Round(phi * mPhi / (2.0 * math.Pi))
    return []int{mod(int(m), int(ug.MTheta)), mod(int(n), int(mPhi))}
}

// GridToID 2D to 1D Index
func (ug UniformGrid) GridToID(IDX []int) int {
    return ug.FirstIndexOf[IDX[0]] + IDX[1]
}

// IDToGrid 1D to 2D Index
func (ug UniformGrid) IDToGrid(ID int) []int {
    m := sort.Search(len(ug.FirstIndexOf)-1, func(i int) bool { return ug.FirstIndexOf[i] > ID })
    n := ID - ug.FirstIndexOf[m-1]
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
