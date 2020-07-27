package main

import (
	//"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	//"github.com/qedus/osmpbf"
	//"io"
	//"log"
	"math"
	"os"
	//"runtime"
	"sort"
	//"sync"
	//"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createPoint(theta float64, phi float64) []float64 {
	//x := 57.296 * math.Sin(theta)*math.Cos(phi)
	//y := 57.296 * math.Sin(theta)*math.Sin(phi)
	//z := 57.296 * math.Cos(theta)
	return []float64{phi/math.Pi*180,theta/math.Pi*180-90}
}

func main(){
	var points [][]float64
	var grid [][][]float64

	N := 10.0*500.0
	nCount := 0
	a := 4.0*math.Pi/N 
	d := math.Sqrt(a)
	M_theta  := math.Round(math.Pi/d)
	d_theta := math.Pi/M_theta
	d_phi := a/d_theta
	
	for m :=0.0; m < M_theta ; m += 1.0{
		theta := math.Pi*(m+0.5)/M_theta
		M_phi := math.Round(2.0*math.Pi*math.Sin(theta)/d_phi)
		var gridRow [][]float64	
		for n := 0.0; n < M_phi;n+=1.0{
			phi := 2*math.Pi *n / M_phi
			nCount +=1
			points = append(points,createPoint(theta,phi))
			gridRow = append(gridRow,createPoint(theta,phi))
		}
		fmt.Printf("%v\n", gridRow)	
		grid = append(grid,gridRow)
	}
	
	fmt.Printf("Points created: %v\n",nCount)
	var rawJson []byte
	g := geojson.NewMultiPointGeometry(points...)
	rawJson, err4 := g.MarshalJSON()
	check(err4)
	var testgeojsonFilename = fmt.Sprintf("../data/output/gridTest.geojson")
	f, err5 := os.Create(testgeojsonFilename)
	check(err5)
	_, err6 := f.Write(rawJson)
	check(err6)
	f.Sync()


	dict := make(map[float64][][]float64)



	//var newGrid [][][]float64 
	/*
	for _,point := range points {
		if val,ok := dict[math.Round(point[0])]; ok {
			dict[math.Round(point[0])] = append(val,point)
		} else{
			dict[math.Round(point[0])] = [][]float64{point}
		}
	}
	*/
	for _,point := range points {
		if val,ok := dict[point[0]]; ok {
			dict[point[1]] = append(val,point)
		} else{
			dict[point[1]] = [][]float64{point}
		}
	}
	/*for _,v := range dict {
		fmt.Printf("%v\n",len(v))
	}*/
	keys := make([]float64, 0, len(dict))
    for k := range dict {
        keys = append(keys, k)
    }
    sort.Float64s(keys)
	//fmt.Printf("%v\n",keys)	 
	/*fmt.Printf("%v\n",len(newGrid))
	for _,j := range newGrid{
		fmt.Printf("%v\n",len(j))
	}*/

	/*var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(newGrid)
	check(err1)
	var filename = fmt.Sprintf("data/output/equiGridTest.json")
	f, err2 := os.Create(filename)
	check(err2)
	_, err3 := f.Write(meshgridBytes)
	check(err3)
	f.Sync()*/

}