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
	"time"
)

var port int = 8081
var meshgrid []bool
var meshgrid2d [][]bool
var meshWidth int64

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toGeojson(route [][][]float64) []byte {
	var rawJson []byte
	routes := geojson.NewFeatureCollection()
	for _, j := range route {
		routes = routes.AddFeature(geojson.NewFeature(geojson.NewLineStringGeometry(j)))
	}
	rawJson, err := routes.MarshalJSON()
	check(err)
	return rawJson
}

func neighbours1d(indx int64) []int64 {
	var neighbours []int64
	var temp []int64

	neighbours = append(neighbours, indx-meshWidth-1) // top left
	neighbours = append(neighbours, indx-meshWidth)   // top
	neighbours = append(neighbours, indx-meshWidth+1) // top right
	neighbours = append(neighbours, indx-1)           // left
	neighbours = append(neighbours, indx+1)           // right
	neighbours = append(neighbours, indx+meshWidth-1) // bottom left
	neighbours = append(neighbours, indx+meshWidth)   // bottom
	neighbours = append(neighbours, indx+meshWidth+1) // bottom right

	for _, j := range neighbours {
		if j >= 0 && j < int64(len(meshgrid)) {
			if !meshgrid[j] {
				temp = append(temp, j)
			}
		}
	}
	return temp
}

func gridToCoord(in []int64) []float64 {
	var out []float64
	out = append(out, float64(in[0]-180))
	out = append(out, float64((float64(in[1])/2)-90))
	return out
}

func coordToGrid(in []float64) []int64 {
	var out []int64
	out = append(out, int64(math.Round(in[0]))+180)
	out = append(out, (int64(math.Round(in[1]))+90)*2)
	return out
}

func flattenIndx(x, y int64) int64 {
	return ((int64(meshWidth) * y) + x)
}

func expandIndx(indx int64) []int64 {
	var x = indx % meshWidth
	var y = indx / meshWidth
	return []int64{x, y}
}

func remove(s []int64, i int64) []int64 {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func min(dist *[]float64, vertices *map[int64]bool) int64 {
	var min = math.Inf(1)
	var argmin int64

	for i := range *vertices {
		if (*dist)[i] < min {
			min = (*dist)[i]
			argmin = i
		}
	}
	return argmin
}

func haversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func distance(start, end []float64) float64 {
	var fLat, fLng, fLat2, fLng2, radius float64
	fLng = start[0] * math.Pi / 180.0
	fLat = start[1] * math.Pi / 180.0
	fLng2 = end[0] * math.Pi / 180.0
	fLat2 = end[1] * math.Pi / 180.0

	radius = 6378100
	h := haversin(fLat2-fLat) + math.Cos(fLat)*math.Cos(fLat2)*haversin(fLng2-fLng)
	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
	return (c * radius)
}

func extractRoute(prev *[]int64, end int64) [][][]float64 {
	var route [][][]float64
	var tempRoute [][]float64
	temp := expandIndx(end)
	for {
		x := expandIndx(end)
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			route = append(route, tempRoute)
			tempRoute = make([][]float64, 0)
		}
		tempRoute = append(tempRoute, gridToCoord([]int64{x[0], x[1]}))

		if (*prev)[end] == -1 {
			break
		}
		end = (*prev)[end]
		temp = x
	}
	route = append(route, tempRoute)
	return route
}

func dijkstra(startLngInt, startLatInt, endLngInt, endLatInt int64) [][][]float64 {

	var dist []float64
	var prev []int64
	var vertices = make(map[int64]bool)

	for i := 0; i < len(meshgrid); i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}
	dist[flattenIndx(startLngInt, startLatInt)] = 0
	vertices[flattenIndx(startLngInt, startLatInt)] = true

	for {
		if len(vertices) == 0 {
			break
		} else {
			var u = min(&dist, &vertices)
			neighbours := neighbours1d(u)
			delete(vertices, u)

			for _, j := range neighbours {

				if j == flattenIndx(endLngInt, endLatInt) {
					prev[j] = u
					return extractRoute(&prev, flattenIndx(endLngInt, endLatInt))
				}
				var alt = dist[u] + distance(gridToCoord(expandIndx(u)), gridToCoord(expandIndx(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					vertices[j] = true
				}
			}
		}
	}

	return extractRoute(&prev, flattenIndx(endLngInt, endLatInt))
}

func main() {
	meshgridRaw, errJson := os.Open("tmp/meshgrid.json")
	if errJson != nil {
		panic(errJson)
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

			var start = coordToGrid([]float64{startLng, startLat})
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = coordToGrid([]float64{endLng, endLat})
			var endLngInt = end[0]
			var endLatInt = end[1]

			if !meshgrid2d[startLngInt][startLatInt] && !meshgrid2d[endLngInt][endLatInt] {
				var start = time.Now()
				var route = dijkstra(startLngInt, startLatInt, endLngInt, endLatInt)

				t := time.Now()
				elapsed := t.Sub(start)
				fmt.Printf("time: %s\n", elapsed)

				var result = toGeojson(route)
				w.Write(result)
			} else {
				w.Write([]byte("false"))
			}

		} else {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
