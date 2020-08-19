package server

import (
	"../algorithms"
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
	"time"
)

var port int = 8081

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

func testExtractRoute(points *[][]int, xSize, ySize int) [][][]float64 {
	var route [][][]float64

	for _, j := range *points {
		coordPoints := make([][]float64, 0)
		for _, l := range j {
			point := algorithms.ExpandIndex(int(l), xSize)
			coordPoints = append(coordPoints, algorithms.GridToCoord([]int{point[0], point[1]}, xSize, ySize))
		}
		route = append(route, coordPoints)
	}
	return route
}

// Run the server with the basic grid
func Run(xSize, ySize int, basicPointInPolygon bool) {
	var meshgrid []bool
	var meshgrid2d [][]bool
	var filename string
	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v.json", xSize, ySize)
	}

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

	var points [][]float64
	for i := 0; i < xSize; i++ {
		for j := 0; j < ySize; j++ {
			if !meshgrid2d[i][j] {
				points = append(points, algorithms.GridToCoord([]int{int(i), int(j)}, int(xSize), int(ySize)))
			}
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

			var start = algorithms.CoordToGrid([]float64{startLng, startLat}, int(xSize), int(ySize))
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = algorithms.CoordToGrid([]float64{endLng, endLat}, int(xSize), int(ySize))
			var endLngInt = end[0]
			var endLatInt = end[1]

			//fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)
			//fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !meshgrid2d[startLngInt][startLatInt] && !meshgrid2d[endLngInt][endLatInt] {
				//fmt.Printf("start: %d / %d ", int(math.Round(startLat)), int(math.Round(startLng)))
				//fmt.Printf("end: %d / %d\n", int(math.Round(endLat)), int(math.Round(endLng)))

				if strings.Contains(r.URL.Path, "/dijkstraAllNodes") {
					var start = time.Now()
					var route, nodesProcessed = algorithms.DijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, int(xSize), int(ySize), &meshgrid)
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
					var route = algorithms.Dijkstra(startLngInt, startLatInt, endLngInt, endLatInt, int(xSize), int(ySize), &meshgrid)
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
		} else if strings.Contains(r.URL.Path, "/basicGrid") {
			pointsJSON, err := json.Marshal(points)
			if err != nil {
				panic(err)
			}
			w.Write(pointsJSON)
		} else {
			http.ServeFile(w, r, "src/server/globe.html")
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf(" on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

// RunUnidistant server
func RunUnidistant(xSize, ySize int, basicPointInPolygon bool) {
	//log.Fatal("Server with unidistant grid not implemented.")

	var uniformgrid []bool
	var uniformgrid2d algorithms.UniformGrid

	var filename string
	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/uniformgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/uniformgrid_%v_%v.json", xSize, ySize)
	}

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &uniformgrid2d)

	for i := 0; i < len(uniformgrid2d.VertexData); i++ {
		for j := 0; j < len(uniformgrid2d.VertexData[i]); j++ {
			uniformgrid = append(uniformgrid, uniformgrid2d.VertexData[i][j])
		}
	}

	var points [][]float64
	for i := 0; i < len(uniformgrid2d.VertexData); i++ {
		for j := 0; j < len(uniformgrid2d.VertexData[i]); j++ {
			if !uniformgrid2d.VertexData[i][j] {
				points = append(points, uniformgrid2d.GridToCoord([]int{int(i), int(j)}))
			}
		}
	}

	uniformgrid2d.XSize = xSize
	uniformgrid2d.YSize = ySize
	uniformgrid2d.BigN = xSize * ySize
	uniformgrid2d.A = 4.0 * math.Pi / float64(uniformgrid2d.BigN)
	uniformgrid2d.D = math.Sqrt(uniformgrid2d.A)
	uniformgrid2d.MTheta = math.Round(math.Pi / uniformgrid2d.D)
	uniformgrid2d.DTheta = math.Pi / uniformgrid2d.MTheta
	uniformgrid2d.DPhi = uniformgrid2d.A / uniformgrid2d.DTheta

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

			var start = uniformgrid2d.CoordToGrid(startLng, startLat)
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = uniformgrid2d.CoordToGrid(endLng, endLat)
			var endLngInt = end[0]
			var endLatInt = end[1]

			//fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)
			//fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !uniformgrid2d.VertexData[startLngInt][startLatInt] && !uniformgrid2d.VertexData[endLngInt][endLatInt] {
				//fmt.Printf("start: %d / %d ", int(math.Round(startLat)), int(math.Round(startLng)))
				//fmt.Printf("end: %d / %d\n", int(math.Round(endLat)), int(math.Round(endLng)))

				if strings.Contains(r.URL.Path, "/dijkstraAllNodes") {
					var start = time.Now()
					var route, nodesProcessed = algorithms.UniformDijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, int(xSize), int(ySize), &uniformgrid2d)
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
					var route = algorithms.UniformDijkstra(startLngInt, startLatInt, endLngInt, endLatInt, int(xSize), int(ySize), &uniformgrid2d)
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
		} else if strings.Contains(r.URL.Path, "/basicGrid") {
			pointsJSON, err := json.Marshal(points)
			if err != nil {
				panic(err)
			}
			w.Write(pointsJSON)
		} else {
			http.ServeFile(w, r, "src/server/globe.html")
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf(" on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
