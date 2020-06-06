package main

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	//"net/http/httputil"
)

var port int = 8081
var meshgrid [][]bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toGeojson(route [][]float64) []byte {
	var rawJson []byte
	g := geojson.NewLineStringGeometry(route)
	rawJson, err4 := g.MarshalJSON()
	check(err4)
	return rawJson
}

func neigbourghs(point []int64) [][]int64 {
	var neigbourghs [][]int64

	var xPlus1  = int64(math.Mod(float64(point[0] + 1),360.0))
	var xMinus1 = int64(math.Mod(float64(point[0] - 1),360.0))
	var yPlus1  = int64(math.Mod(float64(point[1] + 1),360.0))
	var yMinus1 = int64(math.Mod(float64(point[1] - 1),360.0))

	neigbourghs = append(neigbourghs, []int64{xPlus1, point[1]})
	neigbourghs = append(neigbourghs, []int64{xPlus1, yPlus1})
	neigbourghs = append(neigbourghs, []int64{xPlus1, yMinus1})
	neigbourghs = append(neigbourghs, []int64{xMinus1, point[1]})
	neigbourghs = append(neigbourghs, []int64{xMinus1, yPlus1})
	neigbourghs = append(neigbourghs, []int64{xMinus1, yMinus1})
	neigbourghs = append(neigbourghs, []int64{point[0], yPlus1})
	neigbourghs = append(neigbourghs, []int64{point[0], yMinus1})
	return neigbourghs
}

func gridToCoord(in []int64) []float64 {
	var out []float64
	out = append(out, float64(in[0] - 180))
	out = append(out, float64((in[1] / 2) - 90))
	return out
}

func coordToGrid(in []float64) []int64 {
	var out []int64
	out = append(out, int64(math.Round(in[0])) + 180)
	out = append(out, (int64(math.Round(in[1])) + 90) * 2)
	return out
}

func dijkstra(startLng float64, startLat float64, endLng float64, endLat float64) [][]float64 {

	var route [][]float64
	route = append(route, []float64{math.Round(startLng), math.Round(startLat)})
	route = append(route, []float64{math.Round(endLng), math.Round(endLat)})

	//var start_mesh []int64 = []int64{int64(math.Round(startLng)) + 180, (int64(math.Round(startLat)) + 90) * 2}
	//var end_mesh []int64 = []int64{int64(math.Round(endLng)) + 180, (int64(math.Round(endLat)) + 90) * 2}

	//for _,k := range neigbourghs(coordToGrid([]float64{startLng, startLat})) {
	//	route = append(route, gridToCoord(k))
	//}

	fmt.Printf("%v/", meshgrid[int64(math.Round(startLng))+180][(int64(math.Round(startLat))+90)*2])
	fmt.Printf("%v\n", meshgrid[int64(math.Round(endLng))+180][(int64(math.Round(endLat))+90)*2])

	return route
}

func main() {

	meshgridRaw, errJson := os.Open("tmp/meshgrid.json")
	if errJson != nil {
		panic(errJson)
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &meshgrid)
	//fmt.Printf("%v", meshgrid)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/point") {
			query := r.URL.Query()
			var startLat, err = strconv.ParseFloat(query.Get("startLat"), 10)
			if err != nil {
				panic(err)
			}
			var startLng, err1 = strconv.ParseFloat(query.Get("startLng"), 10)
			if err1 != nil {
				panic(err1)
			}
			var endLat, err2 = strconv.ParseFloat(query.Get("endLat"), 10)
			if err2 != nil {
				panic(err2)
			}
			var endLng, err3 = strconv.ParseFloat(query.Get("endLng"), 10)
			if err3 != nil {
				panic(err3)
			}

			fmt.Printf("start: %d / %d ", int64(math.Round(startLat)), int64(math.Round(startLng)))
			fmt.Printf("end: %d / %d\n", int64(math.Round(endLat)), int64(math.Round(endLng)))

			var route = dijkstra(startLng, startLat, endLng, endLat)
			var result = toGeojson(route)

			w.Write(result)
		} else {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
