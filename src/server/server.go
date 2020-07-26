package server

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"math"
	"strconv"
	"strings"
	"time"
	"../algorithms/dijkstra"
	"../algorithms/astern"
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

func gridToCoord(in []int64) []float64 {
	var out []float64
	// big grid
	//out = append(out, float64((float64(in[0])/10)-180))
	//out = append(out, float64(((float64(in[1])/10)/2)-90))
	// small grid
	out = append(out, float64(in[0]-180))
	out = append(out, float64((float64(in[1])/2)-90))
	return out
}

func coordToGrid(in []float64) []int64 {
	var out []int64
	// big grid
	//out = append(out, int64(((math.Round(in[0]*10)/10)+180)*10))
	//out = append(out, int64(((math.Round(in[1]*10)/10)+90)*2*10))
	// small grid
	out = append(out, int64(math.Round(in[0]))+180)
	out = append(out, (int64(math.Round(in[1]))+90)*2)
	return out
}

func expandIndx(indx int64) []int64 {
	var x = indx % meshWidth
	var y = indx / meshWidth
	return []int64{x, y}
}

func toGeojson(route [][][]float64) []byte {
	var rawJSON []byte
	routes := geojson.NewFeatureCollection()
	for _, j := range route {
		//fmt.Printf("%v\n", geojson.NewFeature(geojson.NewLineStringGeometry(j)))
		routes = routes.AddFeature(geojson.NewFeature(geojson.NewLineStringGeometry(j)))
	}
	rawJSON, err := routes.MarshalJSON()
	check(err)
	return rawJSON
}


func testExtractRoute(points *[][]int64) [][][]float64 {
	//print("started extracting route\n")
	var route [][][]float64

	for _, j := range *points {
		coordPoints := make([][]float64, 0)
		for _, l := range j {
			point := expandIndx(int64(l))
			coordPoints = append(coordPoints, gridToCoord([]int64{point[0], point[1]}))
		}
		route = append(route, coordPoints)
	}
	return route
}



// Run ...
func Run(xSize,ySize int) {
	//meshgridRaw, errJSON := os.Open("tmp/meshgrid__planet_big.json")
	filename := fmt.Sprintf("data/output/meshgrid_%v_%v.json",xSize,ySize)
	meshgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &meshgrid2d)

	meshWidth = int64(len(meshgrid2d[0]))
	for i := 0; i < len(meshgrid2d[0]); i++ {
		for j := 0; j < len(meshgrid2d); j++ {
			meshgrid = append(meshgrid, meshgrid2d[j][i])
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

			var start = coordToGrid([]float64{startLng, startLat})
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = coordToGrid([]float64{endLng, endLat})
			var endLngInt = end[0]
			var endLatInt = end[1]

			//fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)
			//fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !meshgrid2d[startLngInt][startLatInt] && !meshgrid2d[endLngInt][endLatInt] {
				//fmt.Printf("start: %d / %d ", int64(math.Round(startLat)), int64(math.Round(startLng)))
				//fmt.Printf("end: %d / %d\n", int64(math.Round(endLat)), int64(math.Round(endLng)))
				
				//
				if strings.Contains(r.URL.Path, "/dijkstraAllNodes") {
					var start = time.Now()
					var route,nodesProcessed = dijkstra.DijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, meshWidth, &meshgrid)
					t := time.Now()
					elapsed := t.Sub(start)
					fmt.Printf("time: %s\n", elapsed)

					var result,errUnmarsch = geojson.UnmarshalFeatureCollection(toGeojson(route))
					if errUnmarsch != nil {
						panic(errUnmarsch)
					}

					data := dijkstraData{
						Route: result,
						AllNodes: nodesProcessed,
					}
					
					var jsonData,errJd = json.Marshal(data)
					if errJd != nil {
						panic(errJd)
					}
					
					w.Write(jsonData)
				} else {
					var start = time.Now()
					var route = dijkstra.Dijkstra(startLngInt, startLatInt, endLngInt, endLatInt, meshWidth, &meshgrid)
					t := time.Now()
					elapsed := t.Sub(start)
					fmt.Printf("time: %s\n", elapsed)
					var result = toGeojson(route)
					w.Write(result)
				}
				

			} else {
				w.Write([]byte("false"))
			}

		} else if strings.Contains(r.URL.Path, "/astern") {
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

			var start = coordToGrid([]float64{startLng, startLat})
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = coordToGrid([]float64{endLng, endLat})
			var endLngInt = end[0]
			var endLatInt = end[1]

			//fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)
			//fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !meshgrid2d[startLngInt][startLatInt] && !meshgrid2d[endLngInt][endLatInt] {
				//fmt.Printf("start: %d / %d ", int64(math.Round(startLat)), int64(math.Round(startLng)))
				//fmt.Printf("end: %d / %d\n", int64(math.Round(endLat)), int64(math.Round(endLng)))
				
				//
				if strings.Contains(r.URL.Path, "/asternAllNodes") {
					var start = time.Now()
					var route,nodesProcessed = astern.AsternAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, meshWidth, &meshgrid)
					t := time.Now()
					elapsed := t.Sub(start)
					fmt.Printf("time: %s\n", elapsed)

					var result,errUnmarsch = geojson.UnmarshalFeatureCollection(toGeojson(route))
					if errUnmarsch != nil {
						panic(errUnmarsch)
					}

					data := dijkstraData{
						Route: result,
						AllNodes: nodesProcessed,
					}
					
					var jsonData,errJd = json.Marshal(data)
					if errJd != nil {
						panic(errJd)
					}
					
					w.Write(jsonData)
				} else {
					var start = time.Now()
					var route = astern.Astern(startLngInt, startLatInt, endLngInt, endLatInt, meshWidth, &meshgrid)
					t := time.Now()
					elapsed := t.Sub(start)
					fmt.Printf("time: %s\n", elapsed)
					var result = toGeojson(route)
					w.Write(result)
				}
				

			} else {
				w.Write([]byte("false"))
			}

		} else if(strings.Contains(r.URL.Path[1:],".") && (strings.Split(r.URL.Path[1:],".")[len(strings.Split(r.URL.Path[1:],"."))-1] == "js" || strings.Split(r.URL.Path[1:],".")[len(strings.Split(r.URL.Path[1:],"."))-1] == "html" || strings.Split(r.URL.Path[1:],".")[len(strings.Split(r.URL.Path[1:],"."))-1] == "css")) {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else {
			//fmt.Printf("%v\n", "default")
			http.ServeFile(w, r, "src/server/globe.html")
		}

	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
	//log.Fatal(http.ListenAndServe(portStr, http.FileServer(http.Dir("src/server"))))
	//http.FileServer(http.Dir("/Users/sergiotapia/go"))
}
