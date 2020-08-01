package main

import (
	"../src/dataprocessing"
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
	"sort"
	"time"
)

var port int = 8081
var meshgrid []bool
var meshgrid2d [][]bool
var meshWidth int
var sphereGrid SphereGrid

// SphereGrid ...
type SphereGrid struct {
	N            int
	VertexData   [][]bool
	FirstIndexOf []int
}

func (sg SphereGrid) gridToId (m,n int) int {
	return sg.FirstIndexOf[m] + n;
}

func (sg SphereGrid) idToGrid (id int) (int,int) {
	m := sort.Search(len(sg.FirstIndexOf)-1, func(i int) bool { return sg.FirstIndexOf[i] >= id })
	n := id - sg.FirstIndexOf[m]
	return m,n
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func mod(a, b int) int {
	return (a%b + b) % b
}

// UniformGridToCoord returns lat, lon for grid coordinates
func UniformGridToCoord(in []int, xSize, ySize int) []float64 {
	m := float64(in[0])
	n := float64(in[1])
	N := float64(xSize * ySize)
	a := 4.0 * math.Pi / N
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta
	theta := math.Pi * (m + 0.5) / mTheta
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
	phi := 2 * math.Pi * n / mPhi
	return []float64{(theta/math.Pi)*180 - 90, (phi / math.Pi) * 180}
}

// UniformCoordToGrid returns grid coordinates given lat, lon
func UniformCoordToGrid(in []float64, xSize, ySize int) []int {
	N := float64(xSize * ySize)
	a := 4.0 * math.Pi / N
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta

	theta := (in[0] + 90) * math.Pi / 180
	m := math.Round((theta * mTheta / math.Pi) - 0.5)

	phi := in[1] * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
	n := math.Round(phi * mPhi / (2 * math.Pi))
	return []int{mod(int(m), int(mTheta)), mod(int(n), int(mPhi))}
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

func neighbours1d(indx int) []int {
	var neighbours []int
	var temp []int

	neighbours = append(neighbours, indx-meshWidth-1) // top left
	neighbours = append(neighbours, indx-meshWidth)   // top
	neighbours = append(neighbours, indx-meshWidth+1) // top right
	neighbours = append(neighbours, indx-1)           // left
	neighbours = append(neighbours, indx+1)           // right
	neighbours = append(neighbours, indx+meshWidth-1) // bottom left
	neighbours = append(neighbours, indx+meshWidth)   // bottom
	neighbours = append(neighbours, indx+meshWidth+1) // bottom right

	for _, j := range neighbours {
		if j >= 0 && j < int(len(meshgrid)) {
			if !meshgrid[j] {
				temp = append(temp, j)
			}
		}
	}
	return temp
}

func getNeighbours(in []float64, xSize, ySize int) [][]int {
	N := float64(xSize * ySize)
	a := 4.0 * math.Pi / N
	d := math.Sqrt(a)
	mTheta := math.Round(math.Pi / d)
	dTheta := math.Pi / mTheta
	dPhi := a / dTheta

	theta := (in[0] + 90) * math.Pi / 180
	m := math.Round((theta * mTheta / math.Pi) - 0.5)

	theta = math.Pi * (m + 0.5) / mTheta

	phi := in[1] * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / dPhi)
	n1 := math.Round(phi * mPhi / (2 * math.Pi))
	n2 := n1 + 1
	n3 := n1 - 1
	p1 := []int{mod(int(m), int(mTheta)), mod(int(n1), int(mPhi))}
	p2 := []int{mod(int(m), int(mTheta)), mod(int(n2), int(mPhi))}
	p3 := []int{mod(int(m), int(mTheta)), mod(int(n3), int(mPhi))}
	return [][]int{p1, p2, p3}
}

func neighboursUniformGrid(in []int) []int {
	var neighbours [][]int
	m := in[0]
	n := in[1]
	//fmt.Printf("height: %v\n", len(sphereGrid.VertexData))
	//fmt.Printf("len: %v\n", len(sphereGrid.VertexData[m]))
	//fmt.Printf("len: %v\n", len(sphereGrid.VertexData[m+1]))
	//fmt.Printf("len: %v\n", len(sphereGrid.VertexData[m-1]))
	neighbours = append(neighbours, []int{m, mod(n-1, len(sphereGrid.VertexData[m]))})
	neighbours = append(neighbours, []int{m, mod(n+1, len(sphereGrid.VertexData[m]))})

	coord := UniformGridToCoord(in, 10, 500)
	//fmt.Printf("coord: %v\n", coord)
	if m > 0 {
		//fmt.Printf("m: %v\n", m)
		coordDown := UniformGridToCoord([]int{m - 1, n}, 10, 500)
		neighbours = append(neighbours, getNeighbours([]float64{coordDown[0], coord[1]}, 10, 500)...)
	}

	if m < len(sphereGrid.VertexData)-1 {
		coordUp := UniformGridToCoord([]int{m + 1, n}, 10, 500)
		neighbours = append(neighbours, getNeighbours([]float64{coordUp[0], coord[1]}, 10, 500)...)
	}
	var neighbours1d []int
	//fmt.Printf("%v\n", neighbours)
	for _,neighbour := range neighbours {
		neighbours1d = append(neighbours1d, sphereGrid.gridToId(neighbour[0],neighbour[1]))
	}
	return neighbours1d
}


func gridToCoord(in []int) []float64 {
	var out []float64
	// big grid
	//out = append(out, float64((float64(in[0])/10)-180))
	//out = append(out, float64(((float64(in[1])/10)/2)-90))
	// small grid
	out = append(out, float64(in[0]-180))
	out = append(out, float64((float64(in[1])/2)-90))
	return out
}

func coordToGrid(in []float64) []int {
	var out []int
	// big grid
	//out = append(out, int(((math.Round(in[0]*10)/10)+180)*10))
	//out = append(out, int(((math.Round(in[1]*10)/10)+90)*2*10))
	// small grid
	out = append(out, int(math.Round(in[0]))+180)
	out = append(out, (int(math.Round(in[1]))+90)*2)
	return out
}

func flattenIndx(x, y int) int {
	return ((int(meshWidth) * y) + x)
}

func expandIndx(indx int) []int {
	var x = indx % meshWidth
	var y = indx / meshWidth
	return []int{x, y}
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func min(dist *[]float64, vertices *map[int]bool) int {
	var min = math.Inf(1)
	var argmin int

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

func extractRoute(prev *[]int, end int) [][][]float64 {
	//print("started extracting route\n")
	var route [][][]float64
	var tempRoute [][]float64
	m,n := sphereGrid.idToGrid(int(end))
	temp := []int{m,n}
	for {
		m,n = sphereGrid.idToGrid(int(end))
		x := []int{m,n}
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			route = append(route, tempRoute)
			tempRoute = make([][]float64, 0)
		}
		tempRoute = append(tempRoute, UniformGridToCoord([]int{x[0], x[1]},10,500))

		if (*prev)[end] == -1 {
			break
		}
		end = (*prev)[end]
		temp = x
	}
	route = append(route, tempRoute)
	return route
}
func testExtractRoute(points *[][]int) [][][]float64 {
	//print("started extracting route\n")
	var route [][][]float64

	for _, j := range *points {
		coordPoints := make([][]float64, 0)
		for _, l := range j {
			point := expandIndx(int(l))
			coordPoints = append(coordPoints, gridToCoord([]int{point[0], point[1]}))
		}
		route = append(route, coordPoints)
	}
	return route
}

func dijkstra(startLngInt, startLatInt, endLngInt, endLatInt int) [][][]float64 {

	var dist []float64
	var prev []int
	//var expansions [][]int
	var vertices = make(map[int]bool)
	//print("started Dijkstra\n")

	for i := 0; i < sphereGrid.N; i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}
	startId := sphereGrid.gridToId(int(startLngInt),int(startLatInt))
	dist[startId] = 0
	vertices[startId] = true

	for {
		if len(vertices) == 0 {
			break
		} else {
			//var expansion = make([]int,0)
			var u = min(&dist, &vertices)
			m,n := sphereGrid.idToGrid(u)
			gridPos := []int{m,n}
			neighbours := neighboursUniformGrid(gridPos)
			delete(vertices, u)

			for _, j := range neighbours {
				//fmt.Printf("j: %v, land:%v\n", j, meshgrid[j])
				if j == sphereGrid.gridToId(endLngInt, endLatInt) {
					prev[j] = u
					return extractRoute(&prev, sphereGrid.gridToId(endLngInt, endLatInt))
					//return testExtractRoute(&expansions)
				}
				//fmt.Printf("Dist[u]: %v\n", dist[u])
				//fmt.Printf("u/j: %v/%v\n", u,j)
				//fmt.Printf("Distance u-j: %v\n", distance(expandIndx(u), expandIndx(j)))
				//fmt.Printf("Summe: %v\n", (dist[u] + distance(expandIndx(u), expandIndx(j))))
				//fmt.Scanln()
				m,n = sphereGrid.idToGrid(u)
				p1 := []int{m,n}
				m,n = sphereGrid.idToGrid(j)
				p2 := []int{m,n}
				var alt = dist[u] + distance(UniformGridToCoord(p1,10,500), UniformGridToCoord(p2,10,500))
				//fmt.Printf("Distance: %v\n",distance(expandIndx(u), expandIndx(j)))
				if alt < dist[j] {
					dist[j] = alt
					prev[j] = u
					vertices[j] = true
					//expansion = append(expansion,j)
				}
			}
			//expansions = append(expansions,expansion)
		}
	}

	return extractRoute(&prev, sphereGrid.gridToId(endLngInt, endLatInt))
	//return testExtractRoute(&expansions)
}

func main() {
	//meshgridRaw, errJSON := os.Open("data/output/meshgrid__planet_big.json")
	meshgridRaw, errJSON := os.Open("../data/output/meshgrid.json")
	if errJSON != nil {
		panic(errJSON)
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &meshgrid2d)

	meshWidth = int(len(meshgrid2d[0]))
	for i := 0; i < len(meshgrid2d[0]); i++ {
		for j := 0; j < len(meshgrid2d); j++ {
			meshgrid = append(meshgrid, meshgrid2d[j][i])
		}
	}
	uniformGridRaw, errJSON2 := os.Open("../data/output/uniformGrid_10_500.json")
	if errJSON2 != nil {
		panic(errJSON2)
	}
	defer uniformGridRaw.Close()
	uniformByteValue, _ := ioutil.ReadAll(uniformGridRaw)
	json.Unmarshal(uniformByteValue, &sphereGrid)
	//fmt.Printf("%v\n",sphereGrid)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/testneighbours") {

			/*query := r.URL.Query()
			var lat, err = strconv.ParseFloat(query.Get("lat"), 10)
			if err != nil {
				panic(err)
			}
			var lng, err1 = strconv.ParseFloat(query.Get("lng"), 10)
			if err1 != nil {
				panic(err1)
			}
			fmt.Printf("Input:%v,%v\n", lat, lng)
			gridPoint := dataprocessing.UniformCoordToGrid([]float64{lat, lng}, 10, 500)
			fmt.Printf("gridPoint %v\n", gridPoint)
			neighbours := neighboursUniformGrid(gridPoint)
			var coords [][]float64
			for _, x := range neighbours {
				fmt.Printf("neighbour: %v\n",x)
				coords = append(coords, UniformGridToCoord(x, 10, 500))
			}
			byteCoords, errByteCoords := json.Marshal(coords)
			if errByteCoords != nil {
				panic(errByteCoords)
			}
			w.Write(byteCoords)*/

		} else if strings.Contains(r.URL.Path, "/testpoint") {

			query := r.URL.Query()
			var lat, err = strconv.ParseFloat(query.Get("lat"), 10)
			if err != nil {
				panic(err)
			}
			var lng, err1 = strconv.ParseFloat(query.Get("lng"), 10)
			if err1 != nil {
				panic(err1)
			}
			//fmt.Printf("%v,%v\n", lat,lng)
			gridPoint := dataprocessing.UniformCoordToGrid([]float64{lat, lng}, 10, 500)
			//fmt.Printf("%v\n", gridPoint)
			coords := dataprocessing.UniformGridToCoord(gridPoint, 10, 500)
			//fmt.Printf("%v,%v\n", coords[0], coords[1])
			byteCoords, errByteCoords := json.Marshal(coords)
			if errByteCoords != nil {
				panic(errByteCoords)
			}
			w.Write(byteCoords)

		} else if strings.Contains(r.URL.Path, "/point") {
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

			var start = UniformCoordToGrid([]float64{startLng, startLat},10,500)
			var startLngInt = start[0]
			var startLatInt = start[1]

			var end = UniformCoordToGrid([]float64{endLng, endLat},10,50)
			var endLngInt = end[0]
			var endLatInt = end[1]

			//fmt.Printf("\n%v/%v, %v/%v\n", startLngInt, startLatInt, endLngInt, endLatInt)
			//fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !sphereGrid.VertexData[startLngInt][startLatInt] && !sphereGrid.VertexData[endLngInt][endLatInt] {
				//fmt.Printf("start: %d / %d ", int(math.Round(startLat)), int(math.Round(startLng)))
				//fmt.Printf("end: %d / %d\n", int(math.Round(endLat)), int(math.Round(endLng)))
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

		} else if strings.Contains(r.URL.Path, "/grid") {
			gridRaw, errJSON := os.Open("../data/output/gridTest.geojson")
			if errJSON != nil {
				panic(errJSON)
			}
			defer gridRaw.Close()
			byteValue, _ := ioutil.ReadAll(gridRaw)
			w.Write(byteValue)
		} else {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
