package dataprocessing

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

func createMeshgrid(xSize int, ySize int, boundingTreeRoot *boundingTree, allCoastlines *[][][]float64, testGeoJSON *[][]float64, meshgrid *[][]bool, createTestGeoJSON bool) string {
	start := time.Now()
	var xStepSize = float64(360 / xSize)
	var yStepSize = float64(360 / ySize)

	var wg sync.WaitGroup
	for x := 0.0; x < 360; x += xStepSize {
		for y := 0.0; y < 360; y += yStepSize {
			wg.Add(1)
			go func(x float64, y float64) {
				defer wg.Done()
				var xs = x - 180
				var ys = (y / 2) - 90
				if isLand(boundingTreeRoot, []float64{xs, ys}, allCoastlines) {
					(*meshgrid)[int(x/xStepSize)][int(y/yStepSize)] = true
					if createTestGeoJSON {
						*testGeoJSON = append(*testGeoJSON, []float64{xs, ys})
					}
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

func storeMeshgrid(meshgrid *[][]bool, filename string) string {
	start := time.Now()
	var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(meshgrid)
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
