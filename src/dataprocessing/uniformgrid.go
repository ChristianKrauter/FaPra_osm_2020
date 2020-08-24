package dataprocessing

import (
	"../grids"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

func createUniformGrid(xSize, ySize int, boundingTreeRoot *boundingTree, polygons *Polygons, ug *grids.UniformGrid, basicPointInPolygon bool) string {
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
			phi := 2.0 * math.Pi * n / mPhi
			nCount++
			coords := []float64{(phi / math.Pi) * 180.0, (theta/math.Pi)*180.0 - 90.0}
			if coords[0] > 180.0 {
				coords[0] = coords[0] - 360.0
			}
			if basicPointInPolygon {
				gridRow = append(gridRow, isLand(boundingTreeRoot, coords, polygons))
			} else {
				gridRow = append(gridRow, isLandSphere(boundingTreeRoot, coords, polygons))
			}
		}
		grid = append(grid, gridRow)
	}

	(*ug).N = nCount
	(*ug).FirstIndexOf = firstIndexOf
	(*ug).VertexData = grid
	(*ug).XSize = xSize
	(*ug).YSize = ySize
	(*ug).BigN = int(N)
	(*ug).A = a
	(*ug).D = d
	(*ug).MTheta = mTheta
	(*ug).DTheta = dTheta
	(*ug).DPhi = dPhi

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created Uniform Grid in          : %s\n", elapsed)

	return elapsed.String()
}

func createUniformGridNBT(xSize, ySize int, allBoundingBoxes *[]map[string]float64, polygons *Polygons, ug *grids.UniformGrid, basicPointInPolygon bool) string {
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
			phi := 2.0 * math.Pi * n / mPhi
			nCount++
			coords := []float64{(phi / math.Pi) * 180.0, (theta/math.Pi)*180.0 - 90.0}
			if coords[0] > 180.0 {
				coords[0] = coords[0] - 360.0
			}
			if basicPointInPolygon {
				gridRow = append(gridRow, isLandNBT(allBoundingBoxes, coords, polygons))
			} else {
				gridRow = append(gridRow, isLandSphereNBT(allBoundingBoxes, coords, polygons))
			}
		}
		grid = append(grid, gridRow)
	}

	(*ug).N = nCount
	(*ug).FirstIndexOf = firstIndexOf
	(*ug).VertexData = grid
	(*ug).XSize = xSize
	(*ug).YSize = ySize
	(*ug).BigN = int(N)
	(*ug).A = a
	(*ug).D = d
	(*ug).MTheta = mTheta
	(*ug).DTheta = dTheta
	(*ug).DPhi = dPhi

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created Uniform Grid in          : %s\n", elapsed)

	return elapsed.String()
}

func storeUniformGrid(ug *grids.UniformGrid, filename string) string {
	start := time.Now()
	var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(ug)
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
