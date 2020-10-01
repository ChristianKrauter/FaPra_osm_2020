package dataprocessing

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createMeshgrid(xSize, ySize int, boundingTreeRoot *boundingTree, polygons *Polygons, bg *[][]bool, basicPointInPolygon bool) string {
	start := time.Now()
	var xStepSize float64 = 360.0 / float64(xSize)
	var yStepSize float64 = 360.0 / float64(ySize)
	var wg sync.WaitGroup

	for x := 0.0; x < 360; x += xStepSize {
		for y := 0.0; y < 360; y += yStepSize {
			wg.Add(1)
			go func(x, y float64) {
				defer wg.Done()
				var xs = x - 180
				var ys = (y / 2) - 90
				if basicPointInPolygon {
					(*bg)[int(x/xStepSize)][int(y/yStepSize)] = isLand(boundingTreeRoot, []float64{xs, ys}, polygons)
				} else {
					(*bg)[int(x/xStepSize)][int(y/yStepSize)] = isLandSphere(boundingTreeRoot, []float64{xs, ys}, polygons)
				}
			}(x, y)
		}
	}

	wg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created Meshrid in               : %s\n", elapsed)
	return elapsed.String()
}

func createMeshgridNBT(xSize, ySize int, allBoundingBoxes *[]map[string]float64, polygons *Polygons, bg *[][]bool, basicPointInPolygon bool) string {
	start := time.Now()
	var xStepSize float64 = 360.0 / float64(xSize)
	var yStepSize float64 = 360.0 / float64(ySize)
	var wg sync.WaitGroup

	for x := 0.0; x < 360; x += xStepSize {
		for y := 0.0; y < 360; y += yStepSize {
			wg.Add(1)
			go func(x float64, y float64) {
				defer wg.Done()
				var xs = x - 180
				var ys = (y / 2) - 90
				if basicPointInPolygon {
					(*bg)[int(x/xStepSize)][int(y/yStepSize)] = isLandNBT(allBoundingBoxes, []float64{xs, ys}, polygons)
				} else {
					(*bg)[int(x/xStepSize)][int(y/yStepSize)] = isLandSphereNBT(allBoundingBoxes, []float64{xs, ys}, polygons)
				}
			}(x, y)
		}
	}

	wg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created Meshrid in               : %s\n", elapsed)
	return elapsed.String()
}

func storeMeshgrid(bg *[][]bool, filename string) string {
	start := time.Now()
	var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(bg)
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
