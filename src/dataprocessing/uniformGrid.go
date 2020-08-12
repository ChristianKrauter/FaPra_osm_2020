package dataprocessing

import (
	"fmt"
	//"github.com/paulmach/go.geojson"
	"encoding/json"
	"math"
	"os"
	"sort"
	"time"
)

// SphereGrid ...
type SphereGrid struct {
	N            int
	VertexData   [][]bool
	FirstIndexOf []int
}

func (sg SphereGrid) gridToID(m, n int) int {
	return sg.FirstIndexOf[m] + n
}

func (sg SphereGrid) idToGrid(id int) (int, int) {
	m := sort.Search(len(sg.FirstIndexOf)-1, func(i int) bool { return sg.FirstIndexOf[i] > id })
	n := id - sg.FirstIndexOf[m-1]
	return m - 1, n
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

func createPoint(theta float64, phi float64) []float64 {
	//x := 57.296 * math.Sin(theta)*math.Cos(phi)
	//y := 57.296 * math.Sin(theta)*math.Sin(phi)
	//z := 57.296 * math.Cos(theta)
	return []float64{theta/math.Pi*180 - 90, phi / math.Pi * 180}
}

// UniformGridToCoord returns lat, lng for grid coordinates
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
	return []float64{(theta/math.Pi)*180 - 90, (phi / math.Pi) * 180}
}

// UniformCoordToGrid returns grid coordinates given lat, lng
func UniformCoordToGrid(in []float64, xSize, ySize int) []int {
	N := float64(xSize * ySize)
	a := 4.0 * math.Pi / N
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta

	theta := (in[0] + 90) * math.Pi / 180
	m := math.Round((theta * mTheta / math.Pi) - 0.5)

	phi := in[1] * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
	n := math.Round(phi * mPhi / (2 * math.Pi))

	return []int{mod(int(m), int(mTheta)), mod(int(n), int(mPhi))}
}

func createUniformGrid(xSize, ySize int, sphereGrid *SphereGrid, boundingTreeRoot *boundingTree, allCoastlines *[][][]float64) string {
	start := time.Now()
	var grid [][]bool
	var firstIndexOf []int
	N := float64(xSize * ySize)
	nCount := 0
	a := 4.0 * math.Pi / N
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta

	for m := 0.0; m < mTheta; m += 1.0 {
		theta := math.Pi * (m + 0.5) / mTheta
		mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
		var gridRow []bool
		firstIndexOf = append(firstIndexOf, int(nCount))
		for n := 0.0; n < mPhi; n += 1.0 {
			// phi := 2 * math.Pi * n / mPhi
			nCount++
			coords := UniformGridToCoord([]int{int(m), int(n)}, xSize, ySize)
			if coords[0] > 90 {
				fmt.Printf("coords: %v\n", coords)
			}
			if isLandSphere(boundingTreeRoot, []float64{coords[1], coords[0]}, allCoastlines) {
				gridRow = append(gridRow, true)
			} else {
				gridRow = append(gridRow, false)
			}
		}
		// fmt.Printf("%v\n", gridRow)
		grid = append(grid, gridRow)
	}

	(*sphereGrid).N = nCount
	(*sphereGrid).FirstIndexOf = firstIndexOf
	(*sphereGrid).VertexData = grid

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created Uniform Grid in          : %s\n", elapsed)
	var rawJSON []byte
	rawJSON, err4 := json.Marshal(*sphereGrid)
	check(err4)
	var jsonFilename = fmt.Sprintf("data/output/uniformGrid_%v_%v.json", xSize, ySize)
	f, err5 := os.Create(jsonFilename)
	check(err5)
	_, err6 := f.Write(rawJSON)
	check(err6)
	f.Sync()

	return elapsed.String()

	// return grid

	/*dict := make(map[float64][][]float64)
		for _,point := range points {
			if val,ok := dict[point[0]]; ok {
				dict[point[1]] = append(val,point)
			} else{
				dict[point[1]] = [][]float64{point}
			}
		}

		keys := make([]float64, 0, len(dict))
	    for k := range dict {
	        keys = append(keys, k)
	    }
	    sort.Float64s(keys)
	*/
}
