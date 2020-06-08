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

func neighbours(point []int64) [][]int64 {
	var neighbours [][]int64

	var xPlus1 = int64(math.Mod(float64(point[0]+1), 360.0))
	var xMinus1 = int64(math.Mod(float64(point[0]-1), 360.0))
	var yPlus1 = int64(math.Mod(float64(point[1]+1), 360.0))
	var yMinus1 = int64(math.Mod(float64(point[1]-1), 360.0))

	neighbours = append(neighbours, []int64{xPlus1, point[1]})
	neighbours = append(neighbours, []int64{xPlus1, yPlus1})
	neighbours = append(neighbours, []int64{xPlus1, yMinus1})
	neighbours = append(neighbours, []int64{xMinus1, point[1]})
	neighbours = append(neighbours, []int64{xMinus1, yPlus1})
	neighbours = append(neighbours, []int64{xMinus1, yMinus1})
	neighbours = append(neighbours, []int64{point[0], yPlus1})
	neighbours = append(neighbours, []int64{point[0], yMinus1})
	return neighbours
}

func neighbours1d(point int64) []int64 {
	var neighbours []int64
	var temp []int64

	if !meshgrid[point-meshWidth-1] {
		neighbours = append(neighbours, point-meshWidth-1) // top left
	}
	if !meshgrid[point-meshWidth] {
		neighbours = append(neighbours, point-meshWidth) // top
	}
	if !meshgrid[point-meshWidth+1] {
		neighbours = append(neighbours, point-meshWidth+1) // top right
	}
	if !meshgrid[point-1] {
		neighbours = append(neighbours, point-1) // left
	}
	if !meshgrid[point+1] {
		neighbours = append(neighbours, point+1) // right
	}
	if !meshgrid[point+meshWidth-1] {
		neighbours = append(neighbours, point+meshWidth-1) // bottom left
	}
	if !meshgrid[point+meshWidth] {
		neighbours = append(neighbours, point+meshWidth) // bottom
	}
	if !meshgrid[point+meshWidth+1] {
		neighbours = append(neighbours, point+meshWidth+1) // bottom right
	}

	for _, j := range neighbours {
		if j >= 0 && j < 360.0*360.0 {
			//fmt.Printf("neighbours: %v\n", j)
			temp = append(temp, j)
		}
		//fmt.Printf("id: %v, vool: %v\n", j, meshgrid[j])
	}
	return neighbours
}

func gridToCoord(in []int64) []float64 {
	var out []float64
	out = append(out, float64(in[0]-180))
	out = append(out, float64((in[1]/2)-90))
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
		//if !meshgrid[i] {
		if (*dist)[i] < min {
			min = (*dist)[i]
			argmin = i
			argminIndx = int64(j)
		}
		//}
	}
	return []int64{argmin, argminIndx}
}

func haversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func distance(start, end []float64) float64 {
	var fLat, fLng, fLat2, fLng2, radius float64
	fLat = start[0] * math.Pi / 100
	fLng = start[1] * math.Pi / 100
	fLat2 = end[0] * math.Pi / 100
	fLng2 = end[1] * math.Pi / 100

	radius = 6378100
	h := haversin(fLat2-fLat) + math.Cos(fLat)*math.Cos(fLat2)*haversin(fLng2-fLng)

	return (2 * radius * math.Asin(math.Sqrt(h)))
}

func extractRoute(prev *[]int64, end int64) [][]float64 {
	print("started extracting route\n")
	var route [][]float64

	for {
		x := expandIndx(end)
		route = append(route, gridToCoord([]int64{int64(x[1]), int64(x[0])}))
		if (*prev)[end] == -1 {
			break
		}
		end = (*prev)[end]
	}
	fmt.Printf("%v\n", route)
	return route
}

func dijkstra(startLng, startLat, endLng, endLat float64, startLngInt, startLatInt, endLngInt, endLatInt int64) [][]float64 {

	var dist []float64
	var prev []int64
	var vertices []int64
	print("started Dijkstra\n")
	//var route [][]float64
	//route = append(route, []float64{math.Round(startLng), math.Round(startLat)})
	//route = append(route, []float64{math.Round(endLng), math.Round(endLat)})

	//var start_mesh []int64 = []int64{int64(math.Round(startLng)) + 180, (int64(math.Round(startLat)) + 90) * 2}
	//var end_mesh []int64 = []int64{int64(math.Round(endLng)) + 180, (int64(math.Round(endLat)) + 90) * 2}

	//for _,k := range neighbours(coordToGrid([]float64{startLng, startLat})) {
	//	route = append(route, gridToCoord(k))
	//}

	//fmt.Printf("%v/", meshgrid2d[int64(math.Round(startLng))+180][(int64(math.Round(startLat))+90)*2])
	//fmt.Printf("%v\n", meshgrid2d[int64(math.Round(endLng))+180][(int64(math.Round(endLat))+90)*2])

	for i, j := range meshgrid {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
		if !j {
			vertices = append(vertices, int64(i))
		}
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
			print("running")
			for _, j := range neighbours1d(u) {
				//fmt.Printf("j: %v, land:%v\n", j, meshgrid[j])
				//if !meshgrid[j] {
				if j == flattenIndx(endLatInt, endLngInt) {
					prev[j] = u
					return extractRoute(&prev, flattenIndx(endLatInt, endLngInt))
				}
				var alt = dist[u] + distance(expandIndx(u), expandIndx(j))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
				}
				//}
			}
		}
	}

	return extractRoute(&prev, flattenIndx(endLatInt, endLngInt))
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

			fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)
			fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !meshgrid2d[startLngInt][startLatInt] && !meshgrid2d[endLngInt][endLatInt] {
				fmt.Printf("start: %d / %d ", int64(math.Round(startLat)), int64(math.Round(startLng)))
				fmt.Printf("end: %d / %d\n", int64(math.Round(endLat)), int64(math.Round(endLng)))

				var route = dijkstra(startLng, startLat, endLng, endLat, startLngInt, startLatInt, endLngInt, endLatInt)
				fmt.Printf("%v", route)
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
