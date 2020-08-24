package server

import (
	"../algorithms"
	"../grids"
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
	var bg grids.BasicGrid
	var bg2D [][]bool
	var filename string

	bg.XSize = xSize
	bg.YSize = ySize

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v.json", xSize, ySize)
	}

	meshgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		log.Fatal(fmt.Sprintf("\nThe meshgrid '%s'\ncould not be found. Please create it first.\n", filename))
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &bg2D)

	for i := 0; i < len(bg2D[0]); i++ {
		for j := 0; j < len(bg2D); j++ {
			bg.VertexData = append(bg.VertexData, bg2D[j][i])
		}
	}

	var points [][]float64
	for i := 0; i < xSize; i++ {
		for j := 0; j < ySize; j++ {
			if !bg2D[i][j] {
				points = append(points, bg.GridToCoord([]int{int(i), int(j)}))
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

			var from = bg.CoordToGrid([]float64{startLng, startLat})
			var to = bg.CoordToGrid([]float64{endLng, endLat})

			if !bg2D[from[0]][from[1]] && !bg2D[to[0]][to[1]] {
				if strings.Contains(r.URL.Path, "/dijkstraAllNodes") {
					var start = time.Now()
					var route, nodesProcessed = algorithms.DijkstraAllNodesBg(from, to, &bg)
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
					var route, _ = algorithms.DijkstraBg(from, to, &bg)
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
			var grid = bg.CoordToGrid([]float64{lng, lat})
			if bg2D[grid[0]][grid[1]] {
				w.Write([]byte("false"))
			} else {
				var coord = bg.GridToCoord(grid)
				rawJSON, err := geojson.NewPointGeometry(coord).MarshalJSON()
				check(err)
				w.Write(rawJSON)
			}
		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" || fileEnding == "ico" {
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
	fmt.Printf("on localhost%s\n\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

// RunUnidistant server
func RunUnidistant(xSize, ySize, algorithm int, basicPointInPolygon bool) {
	var ug1D []bool
	var ug grids.UniformGrid

	var filename string
	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/uniformgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/uniformgrid_%v_%v.json", xSize, ySize)
	}

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		log.Fatal(fmt.Sprintf("\nThe meshgrid '%s'\ncould not be found. Please create it first.\n", filename))
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

	var algoStr string
	switch algorithm {
	case 0:
		algoStr = "_dij"
	default:
		algoStr = "_dij"
	}
	fmt.Printf("Algo: %v\n", algoStr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var fileEnding = strings.Split(r.URL.Path[1:], ".")[len(strings.Split(r.URL.Path[1:], "."))-1]

		if strings.Contains(r.URL.Path, "/astern") {
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

			var from = ug.CoordToGrid(startLng, startLat)
			var to = ug.CoordToGrid(endLng, endLat)

			if !ug.VertexData[from[0]][from[1]] && !ug.VertexData[to[0]][to[1]] {
				if strings.Contains(r.URL.Path, "/asternAllNodes") {
					var start = time.Now()
					var route, nodesProcessed = algorithms.AsternAllNodes(from, to, &ug)
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
					var route, _ = algorithms.Astern(from, to, &ug)
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

		} else  if strings.Contains(r.URL.Path, "/dijkstra") {
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

			var from = ug.CoordToGrid(startLng, startLat)
			var to = ug.CoordToGrid(endLng, endLat)

			if !ug.VertexData[from[0]][from[1]] && !ug.VertexData[to[0]][to[1]] {
				if strings.Contains(r.URL.Path, "/dijkstraAllNodes") {
					var start = time.Now()
					var route, nodesProcessed = algorithms.DijkstraAllNodes(from, to, &ug)
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
					var route, _ = algorithms.Dijkstra(from, to, &ug)
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
		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" || fileEnding == "ico" {
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
	fmt.Printf("on localhost%s\n\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
