package server

import (
	"../algorithms"
	"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var port int = 8081
var meshWidth int64
var meshgrid []bool
var meshgrid2d [][]bool

type dijkstraData struct {
	Route    *geojson.FeatureCollection
	AllNodes [][]float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toGeojson(route [][][]float64) *geojson.FeatureCollection {
	routes := geojson.NewFeatureCollection()
	for _, j := range route {
		routes = routes.AddFeature(geojson.NewFeature(geojson.NewLineStringGeometry(j)))
	}
	return routes
}

func testExtractRoute(points *[][]int64, xSize, ySize int64) [][][]float64 {
	var route [][][]float64

	for _, j := range *points {
		coordPoints := make([][]float64, 0)
		for _, l := range j {
			point := algorithms.ExpandIndex(int64(l), xSize)
			coordPoints = append(coordPoints, algorithms.GridToCoord([]int64{point[0], point[1]}, xSize, ySize))
		}
		route = append(route, coordPoints)
	}
	return route
}

// Run the server
func Run(xSize, ySize int) {
	filename := fmt.Sprintf("data/output/meshgrid_%v_%v.json", xSize, ySize)
	meshgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &meshgrid2d)

	for i := 0; i < len(meshgrid2d[0]); i++ {
		for j := 0; j < len(meshgrid2d); j++ {
			meshgrid = append(meshgrid, meshgrid2d[j][i])
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var fileEnding = strings.Split(r.URL.Path[1:], ".")[len(strings.Split(r.URL.Path[1:], "."))-1]

		if strings.Contains(r.URL.Path, "/dijkstra") {
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

			var start = algorithms.CoordToGrid([]float64{startLng, startLat}, int64(xSize), int64(ySize))
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = algorithms.CoordToGrid([]float64{endLng, endLat}, int64(xSize), int64(ySize))
			var endLngInt = end[0]
			var endLatInt = end[1]

			//fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)
			//fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !meshgrid2d[startLngInt][startLatInt] && !meshgrid2d[endLngInt][endLatInt] {
				//fmt.Printf("start: %d / %d ", int64(math.Round(startLat)), int64(math.Round(startLng)))
				//fmt.Printf("end: %d / %d\n", int64(math.Round(endLat)), int64(math.Round(endLng)))

				if strings.Contains(r.URL.Path, "/dijkstraAllNodes") {
					var start = time.Now()
					var route, nodesProcessed = algorithms.DijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, int64(xSize), int64(ySize), &meshgrid)
					t := time.Now()
					elapsed := t.Sub(start)
					fmt.Printf("time: %s\n", elapsed)

					var result = toGeojson(route)
					data := dijkstraData{
						Route:    result,
						AllNodes: nodesProcessed,
					}

					var jsonData, errJd = json.Marshal(data)
					if errJd != nil {
						panic(errJd)
					}

					w.Write(jsonData)
				} else {
					var start = time.Now()
					var route = algorithms.Dijkstra(startLngInt, startLatInt, endLngInt, endLatInt, int64(xSize), int64(ySize), &meshgrid)
					t := time.Now()
					elapsed := t.Sub(start)
					fmt.Printf("time: %s\n", elapsed)
					var result = toGeojson(route)
					rawJSON, err := result.MarshalJSON()
					check(err)
					w.Write(rawJSON)
				}

			} else {
				w.Write([]byte("false"))
			}

		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else {
			http.ServeFile(w, r, "src/server/globe.html")
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf(" on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
