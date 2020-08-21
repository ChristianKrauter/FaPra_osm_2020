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

			if !meshgrid2d[startLngInt][startLatInt] && !meshgrid2d[endLngInt][endLatInt] {
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

		} else if strings.Contains(r.URL.Path, "/getGridPoint") {
			query := r.URL.Query()
			var lat, err = strconv.ParseFloat(query.Get("lat"), 10)
			if err != nil {
				panic(err)
			}
			var lng, err1 = strconv.ParseFloat(query.Get("lng"), 10)
			if err1 != nil {
				panic(err1)
			}
			var grid = algorithms.CoordToGrid([]float64{lng, lat}, int(xSize), int(ySize))
			if meshgrid2d[grid[0]][grid[1]] {
				w.Write([]byte("false"))
			} else {
				var coord = algorithms.GridToCoord(grid, int(xSize), int(ySize))
				rawJSON, err := geojson.NewPointGeometry(coord).MarshalJSON()
				check(err)
				w.Write(rawJSON)
			}
		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else if strings.Contains(r.URL.Path, "/Grid") {
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
	var ug1D []bool
	var ug algorithms.UniformGrid

	ug.XSize = xSize
	ug.YSize = ySize
	ug.BigN = xSize * ySize
	ug.A = 4.0 * math.Pi / float64(ug.BigN)
	ug.D = math.Sqrt(ug.A)
	ug.MTheta = math.Round(math.Pi / ug.D)
	ug.DTheta = math.Pi / ug.MTheta
	ug.DPhi = ug.A / ug.DTheta

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
	json.Unmarshal(byteValue, &ug)

	for i := 0; i < len(ug.VertexData); i++ {
		for j := 0; j < len(ug.VertexData[i]); j++ {
			ug1D = append(ug1D, ug.VertexData[i][j])
		}
	}

	var points [][]float64
	for i := 0; i < len(ug.VertexData); i++ {
		for j := 0; j < len(ug.VertexData[i]); j++ {
			if !ug.VertexData[i][j] {
				points = append(points, ug.GridToCoord([]int{int(i), int(j)}))
			}
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var fileEnding = strings.Split(r.URL.Path[1:], ".")[len(strings.Split(r.URL.Path[1:], "."))-1]

		if strings.Contains(r.URL.Path, "/dijkstra") {
			query := r.URL.Query()
			var startLat, err = strconv.ParseFloat(query.Get("startLat"), 15)
			if err != nil {
				panic(err)
			}
			var startLng, err1 = strconv.ParseFloat(query.Get("startLng"), 15)
			if err1 != nil {
				panic(err1)
			}
			var endLat, err2 = strconv.ParseFloat(query.Get("endLat"), 15)
			if err2 != nil {
				panic(err2)
			}
			var endLng, err3 = strconv.ParseFloat(query.Get("endLng"), 15)
			if err3 != nil {
				panic(err3)
			}

			var start = ug.CoordToGrid(startLng, startLat)
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = ug.CoordToGrid(endLng, endLat)
			var endLngInt = end[0]
			var endLatInt = end[1]

			if !ug.VertexData[startLngInt][startLatInt] && !ug.VertexData[endLngInt][endLatInt] {
				if strings.Contains(r.URL.Path, "/dijkstraAllNodes") {
					var start = time.Now()
					var route, nodesProcessed = algorithms.UniformDijkstraAllNodes(startLngInt, startLatInt, endLngInt, endLatInt, &ug)
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
					var route = algorithms.UniformDijkstra(startLngInt, startLatInt, endLngInt, endLatInt, &ug)
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

		} else if strings.Contains(r.URL.Path, "/getGridPoint") {
			query := r.URL.Query()
			var lat, err = strconv.ParseFloat(query.Get("lat"), 10)
			if err != nil {
				panic(err)
			}
			var lng, err1 = strconv.ParseFloat(query.Get("lng"), 10)
			if err1 != nil {
				panic(err1)
			}
			var grid = ug.CoordToGrid(lng, lat)
			if ug.VertexData[grid[0]][grid[1]] {
				w.Write([]byte("false"))
			} else {
				var coord = ug.GridToCoord(grid)
				rawJSON, err := geojson.NewPointGeometry(coord).MarshalJSON()
				check(err)
				w.Write(rawJSON)
			}
		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else if strings.Contains(r.URL.Path, "/Grid") {
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
