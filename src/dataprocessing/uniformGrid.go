package dataprocessing

import (
	"fmt"
	"github.com/paulmach/go.geojson"
	"math"
	"os"
)

func createPoint(theta float64, phi float64) []float64 {
	//x := 57.296 * math.Sin(theta)*math.Cos(phi)
	//y := 57.296 * math.Sin(theta)*math.Sin(phi)
	//z := 57.296 * math.Cos(theta)
	return []float64{phi / math.Pi * 180, theta/math.Pi*180 - 90}
}

func createGrid(saveToFile bool) {
	var points [][]float64
	var grid [][][]float64

	N := 10.0 * 500.0
	nCount := 0
	a := 4.0 * math.Pi / N
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta

	for m := 0.0; m < mTheta; m += 1.0 {
		theta := math.Pi * (m + 0.5) / mTheta
		mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
		var gridRow [][]float64
		for n := 0.0; n < mPhi; n += 1.0 {
			phi := 2 * math.Pi * n / mPhi
			nCount++
			points = append(points, createPoint(theta, phi))
			gridRow = append(gridRow, createPoint(theta, phi))
		}
		// fmt.Printf("%v\n", gridRow)
		grid = append(grid, gridRow)
	}

	if saveToFile {
		fmt.Printf("Points created: %v\n", nCount)
		var rawJSON []byte
		g := geojson.NewMultiPointGeometry(points...)
		rawJSON, err4 := g.MarshalJSON()
		check(err4)
		var testgeojsonFilename = fmt.Sprintf("tmp/gridTest.geojson")
		f, err5 := os.Create(testgeojsonFilename)
		check(err5)
		_, err6 := f.Write(rawJSON)
		check(err6)
		f.Sync()
	}
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
