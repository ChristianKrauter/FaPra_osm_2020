package main

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	//"github.com/qedus/osmpbf"
	//"io"
	//"log"
	"math"
	"os"
	//"runtime"
	//"sort"
	//"sync"
	//"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createPoint(theta float64, phi float64) []float64 {
	x := 6378100 * math.Sin(theta)*math.Cos(phi)
	y := 6378100 * math.Sin(theta)*math.Sin(phi)
	z := 6378100 * math.Cos(theta)
	r := math.Sqrt(math.Pow(x,2)+math.Pow(y,2)+math.Pow(z,2))
	lat := math.Asin(z/r)
	lon := math.Atan2(y,x)
	return []float64{lat*180.0,lon*90.0}
}

func main(){
	var points [][]float64
	N := 5000.0
	nCount := 0
	a := 4.0*math.Pi/N 
	d := math.Sqrt(a)
	M_theta  := math.Round(math.Pi/d)
	d_theta := math.Pi/M_theta
	d_phi := a/d_theta
	
	for m :=0.0; m < M_theta ; m += 1.0{
		theta := math.Pi*(m+0.5)/M_theta
		M_phi := math.Round(2.0*math.Pi*math.Sin(theta)/d_phi)
		
		for n := 0.0; n < M_phi;n+=1.0{
			phi := 2*math.Pi *n / M_phi
			nCount +=1
			points = append(points,createPoint(theta,phi))
		}
	}

	fmt.Printf("Points created: %v\n",nCount)
	var rawJson []byte
	g := geojson.NewMultiPointGeometry(points...)
	rawJson, err4 := g.MarshalJSON()
	check(err4)
	var testgeojsonFilename = fmt.Sprintf("tmp/gridTest.geojson")
	f, err5 := os.Create(testgeojsonFilename)
	check(err5)
	_, err6 := f.Write(rawJson)
	check(err6)
	f.Sync()



	var newGrid [][][]float64 
	var line  [][]float64
	c := points[0][0]
	for _,point := range points {
		
		
		if(point[0] != c){
			newGrid = append(newGrid,line)
			line = [][]float64{}
			c = point[0]
		}
		line = append(line,point)
	}
		 
	fmt.Printf("%v\n",len(newGrid))
	for _,j := range newGrid{
		fmt.Printf("%v\n",len(j))
	}

	var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(newGrid)
	check(err1)
	var filename = fmt.Sprintf("tmp/equiGridTest.json")
	f, err2 := os.Create(filename)
	check(err2)
	_, err3 := f.Write(meshgridBytes)
	check(err3)
	f.Sync()

}