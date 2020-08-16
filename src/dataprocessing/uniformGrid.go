package dataprocessing

import (
	"../algorithms"
	"fmt"
	//"github.com/paulmach/go.geojson"
	"encoding/json"
	"math"
	"os"
	"sort"
	"time"
)

// UniformGrid ...
type UniformGrid struct {
	N            int
	VertexData   [][]bool
	FirstIndexOf []int
}

func (sg UniformGrid) gridToID(m, n int) int {
	return sg.FirstIndexOf[m] + n
}

func (sg UniformGrid) idToGrid(id int) (int, int) {
	m := sort.Search(len(sg.FirstIndexOf)-1, func(i int) bool { return sg.FirstIndexOf[i] > id })
	n := id - sg.FirstIndexOf[m-1]
	return m - 1, n
}

func createPoint(theta float64, phi float64) []float64 {
	//x := 57.296 * math.Sin(theta)*math.Cos(phi)
	//y := 57.296 * math.Sin(theta)*math.Sin(phi)
	//z := 57.296 * math.Cos(theta)
	return []float64{theta/math.Pi*180 - 90, phi / math.Pi * 180}
}

func createUniformGrid(xSize, ySize int, boundingTreeRoot *boundingTree, allCoastlines *[][][]float64, testGeoJSON *[][]float64, uniformGrid *UniformGrid, createTestGeoJSON, basicPointInPolygon bool) string {
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
			coords := algorithms.UniformGridToCoord([]int{int(m), int(n)}, xSize, ySize)
			if coords[0] > 90 {
				fmt.Printf("coords: %v\n", coords)
			}
			if basicPointInPolygon {
				if isLand(boundingTreeRoot, []float64{coords[1], coords[0]}, allCoastlines) {
					gridRow = append(gridRow, true)
					if createTestGeoJSON {
						*testGeoJSON = append(*testGeoJSON, []float64{coords[0], coords[1]})
					}
				} else {
					gridRow = append(gridRow, false)
				}
			} else {
				if isLandSphere(boundingTreeRoot, []float64{coords[1], coords[0]}, allCoastlines) {
					gridRow = append(gridRow, true)
					if createTestGeoJSON {
						*testGeoJSON = append(*testGeoJSON, []float64{coords[0], coords[1]})
					}
				} else {
					gridRow = append(gridRow, false)
				}
			}

		}
		// fmt.Printf("%v\n", gridRow)
		grid = append(grid, gridRow)
	}

	(*uniformGrid).N = nCount
	(*uniformGrid).FirstIndexOf = firstIndexOf
	(*uniformGrid).VertexData = grid

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created Uniform Grid in          : %s\n", elapsed)

	return elapsed.String()
}

func storeUniformGrid(uniformGrid *UniformGrid, filename string) string {
	start := time.Now()
	var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(uniformGrid)
	check(err1)
	f, err2 := os.Create(filename)
	check(err2)
	_, err3 := f.Write(meshgridBytes)
	check(err3)
	f.Sync()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Stored Meshrid to disc in        : %s\n", elapsed)
	return elapsed.String()
}
