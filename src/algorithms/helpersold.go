package algorithms

import (
    "math"
)

// UniformGridToCoord returns lng, lat for grid coordinates
func UniformGridToCoord(in []int, xSize, ySize int) []float64 {
    m := float64(in[0])
    n := float64(in[1])
    N := float64(xSize * ySize)
    a := 4.0 * math.Pi / N
    d := math.Sqrt(a)
    mTheta := math.Round(math.Pi / d)
    dTheta := math.Pi / mTheta
    dPhi := a / dTheta
    theta := math.Pi * (m + 0.5) / mTheta
    mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
    phi := 2 * math.Pi * n / mPhi
    return []float64{(phi / math.Pi) * 180, (theta/math.Pi)*180 - 90}
}

// UniformCoordToGrid returns grid coordinates given lng,lat
func UniformCoordToGrid(in []float64, xSize, ySize int) []int {
    N := float64(xSize * ySize)
    a := 4.0 * math.Pi / N
    d := math.Sqrt(a)
    mTheta := math.Round(math.Pi / d)
    dTheta := math.Pi / mTheta
    dPhi := a / dTheta

    theta := (in[1] + 90) * math.Pi / 180
    m := math.Round((theta * mTheta / math.Pi) - 0.5)

    phi := in[0] * math.Pi / 180
    mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
    n := math.Round(phi * mPhi / (2 * math.Pi))

    return []int{mod(int(m), int(mTheta)), mod(int(n), int(mPhi))}
}
