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
var meshgrid []bool
var meshgrid2d [][]bool
var meshWidth int64
var dist []float64
var prev []int64
var vertices []int64

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

	var xPlus1 = int64(math.Mod(float64(point[0]+1), 360.0))
	var xMinus1 = int64(math.Mod(float64(point[0]-1), 360.0))
	var yPlus1 = int64(math.Mod(float64(point[1]+1), 360.0))
	var yMinus1 = int64(math.Mod(float64(point[1]-1), 360.0))

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

func neigbourghs1d(point int64) []int64 {
	var neigbourghs []int64

	neigbourghs = append(neigbourghs, point-meshWidth-1) // top left
	neigbourghs = append(neigbourghs, point-meshWidth)   // top
	neigbourghs = append(neigbourghs, point-meshWidth+1) // top right
	neigbourghs = append(neigbourghs, point-1)           // left
	neigbourghs = append(neigbourghs, point+1)           // right
	neigbourghs = append(neigbourghs, point+meshWidth-1) // bottom left
	neigbourghs = append(neigbourghs, point+meshWidth)   // bottom
	neigbourghs = append(neigbourghs, point+meshWidth+1) // bottom right
	return neigbourghs
}

func gridToCoord(in []int64) []float64 {
	var out []float64
	out = append(out, float64(in[0]-180))
	out = append(out, float64((in[1]/2)-90))
	return out
}

func coordToGrid(in []float64) []int64 {
	var out []int64
	out = append(out, int64(math.Round(in[0]))+180)
	out = append(out, (int64(math.Round(in[1]))+90)*2)
	return out
}

func flattenIndx(lat int64, lng int64) int64 {
	return ((int64(meshWidth) * lat) + lng)
}

func expandIndx(indx int64) []float64 {
	var lat = float64(indx / meshWidth)
	var lng = float64(indx % meshWidth)
	return []float64{lat, lng}
}

func remove(s []int64, i int64) []int64 {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func min(dist *[]float64, vertices *[]int64) []int64 {
	var min = math.Inf(1)
	var argmin int64
	var argminIndx int64

	for j, i := range *vertices {
		if (*dist)[i] < min {
			min = (*dist)[i]
			argmin = i
			argminIndx = int64(j)
		}
	}
	return []int64{argmin, argminIndx}
}

func haversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func distance(start, end []float64) float64 {
	var fLat, fLng, fLat2, fLng2, radius float64
	fLat  = start[0] * math.Pi / 100
	fLng  = start[1] * math.Pi / 100
	fLat2 = end[0] * math.Pi / 100
	fLng2 = end[1] * math.Pi / 100

	radius = 6378100
	h := haversin(fLat2-fLat) + math.Cos(fLat)*math.Cos(fLat2)*haversin(fLng2-fLng)

	return (2 * radius * math.Asin(math.Sqrt(h)))
}

func extractRoute(prev []int64, end int64) [][]float64{
	print("started expand route\n")
	var route [][]float64

	for{
		route = append(route, expandIndx(end))
		if prev[end] == -1 {
			break
		}
		end = prev[end]
	}
	fmt.Printf("%v\n", route)
	return route
}

func dijkstra(startLng, startLat, endLng, endLat float64, startLngInt, startLatInt, endLngInt, endLatInt int64) [][]float64 {

	print("started Dijkstra\n")
	var route [][]float64
	route = append(route, []float64{math.Round(startLng), math.Round(startLat)})
	route = append(route, []float64{math.Round(endLng), math.Round(endLat)})

	//var start_mesh []int64 = []int64{int64(math.Round(startLng)) + 180, (int64(math.Round(startLat)) + 90) * 2}
	//var end_mesh []int64 = []int64{int64(math.Round(endLng)) + 180, (int64(math.Round(endLat)) + 90) * 2}

	//for _,k := range neigbourghs(coordToGrid([]float64{startLng, startLat})) {
	//	route = append(route, gridToCoord(k))
	//}

	fmt.Printf("%v/", meshgrid2d[int64(math.Round(startLng))+180][(int64(math.Round(startLat))+90)*2])
	fmt.Printf("%v\n", meshgrid2d[int64(math.Round(endLng))+180][(int64(math.Round(endLat))+90)*2])

	for i := range meshgrid {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
		vertices = append(vertices, int64(i))
	}

	//// Flatten / expand test
	//var tempFlattened int64
	//tempFlattened = flattenIndx(startLatInt, startLngInt)
	//fmt.Printf("lat/lng: %v / %v\n", startLatInt, startLngInt)
	//fmt.Printf("flatten: %v\n", tempFlattened)
	//fmt.Printf("expand : %v\n", expandIndx(tempFlattened))
	//os.Exit(-1)

	dist[flattenIndx(startLatInt, startLngInt)] = 0
	print("starting loop\n")
	for {
		if len(vertices) == 0 {
			break
		} else {
			var whatever = min(&dist, &vertices)
			var u = whatever[0]
			//fmt.Printf("u: %v\n", u)
			remove(vertices, whatever[1])

			for _, j := range neigbourghs1d(u) {
				//fmt.Printf("j: %v, land:%v\n", j, meshgrid[j])
				if !meshgrid[j] {
					var alt = dist[u] + distance(expandIndx(u), expandIndx(j))
					if alt < dist[j] {
						dist[j] = alt
						prev[j] = u
					}
				}
			}
		}
	}

	return extractRoute(prev, flattenIndx(endLatInt, endLngInt))
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
	for i := 0; i < len(meshgrid2d); i++ {
		meshgrid = append(meshgrid, meshgrid2d[i]...)
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

			fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)

			if !meshgrid2d[startLatInt][startLngInt] && !meshgrid2d[endLatInt][endLngInt] {
				fmt.Printf("start: %d / %d ", int64(math.Round(startLat)), int64(math.Round(startLng)))
				fmt.Printf("end: %d / %d\n", int64(math.Round(endLat)), int64(math.Round(endLng)))

				var route = dijkstra(startLng, startLat, endLng, endLat, startLngInt, startLatInt, endLngInt, endLatInt)
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
	fmt.Printf("Starting server on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
