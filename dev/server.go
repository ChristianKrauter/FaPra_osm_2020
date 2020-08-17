package main

import (
	//"../src/dataprocessing"
	"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var port int = 8081
var meshgrid []bool
var meshgrid2d [][]bool
var meshWidth int
var uniformGrid UniformGrid

// UniformGrid ...
type UniformGrid struct {
	N            int
	VertexData   [][]bool
	FirstIndexOf []int
}

func (sg UniformGrid) gridToID(m, n int) int {
	return sg.FirstIndexOf[m] + n
}

func (sg UniformGrid) idToGrid(id int) (int, int) {
	m := sort.Search(len(sg.FirstIndexOf)-1, func(i int) bool { return sg.FirstIndexOf[i] > id })
	n := id - sg.FirstIndexOf[m-1]
	return m - 1, n
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func mod(a, b int) int {
	a = a % b
	if a >= 0 {
		return a
	}
	if b < 0 {
		return a - b
	}
	return a + b
}

// UniformGridToCoord returns lat, lng for grid coordinates
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

// M = LATITUTE, N = LONGITUDE

// UniformCoordToGrid returns grid coordinates given lat, lng
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

	//fmt.Printf("in0 %f, in1 %f\n", in[0], in[1])
	//fmt.Printf("THETA %f, PHI %f\n", theta, phi)
	//fmt.Printf("m %d, mod mtheta %d = %d\n", int(m), int(mTheta), mod(int(m), int(mTheta)))
	//fmt.Printf("n %d, mod mPhi %d = %d\n", int(n), int(mPhi), mod(int(n), int(mPhi)))
	return []int{mod(int(m), int(mTheta)), mod(int(n), int(mPhi))}
	// return []int{mod2(int(m), int(mTheta)), mod2(int(n), int(mPhi))}
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

// Test if it still works with less than 3 points in one grid row
func UniformNeighboursRow(in []float64, xSize, ySize int) [][]int {
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

func GetNeighboursUniformGrid(in []int) []int {
	var neighbours [][]int
	m := in[0]
	n := in[1]
	neighbours = append(neighbours, []int{m, mod(n-1, len(uniformGrid.VertexData[m]))})
	neighbours = append(neighbours, []int{m, mod(n+1, len(uniformGrid.VertexData[m]))})

	coord := UniformGridToCoord(in, 100, 500)
	
	if m > 0 {
		
		coordDown := UniformGridToCoord([]int{m - 1, n}, 100, 500)
		neighbours = append(neighbours, getNeighbours([]float64{coordDown[0], coord[1]}, 100, 500)...)
	}

	if m < len(uniformGrid.VertexData)-1 {
		coordUp := UniformGridToCoord([]int{m + 1, n}, 100, 500)
		neighbours = append(neighbours, getNeighbours([]float64{coordUp[0], coord[1]}, 100, 500)...)
	}
	var neighbours1d []int
	for _, neighbour := range neighbours {
		if !uniformGrid.VertexData[neighbour[0]][neighbour[1]] {
			neighbours1d = append(neighbours1d, uniformGrid.gridToID(neighbour[0], neighbour[1]))
		}
	}
	return neighbours1d
}

func testNeighboursUniformGrid(in []int) [][]int {
	var neighbours [][]int
	m := in[0]
	n := in[1]
	//fmt.Printf("height: %v\n", len(uniformGrid.VertexData))
	//fmt.Printf("len: %v\n", len(uniformGrid.VertexData[m]))
	//fmt.Printf("len: %v\n", len(uniformGrid.VertexData[m+1]))
	//fmt.Printf("len: %v\n", len(uniformGrid.VertexData[m-1]))
	neighbours = append(neighbours, []int{m, mod(n-1, len(uniformGrid.VertexData[m]))})
	neighbours = append(neighbours, []int{m, mod(n+1, len(uniformGrid.VertexData[m]))})

	coord := UniformGridToCoord(in, 100, 500)
	//fmt.Printf("coord: %v\n", coord)
	if m > 0 {
		//fmt.Printf("m: %v\n", m)
		coordDown := UniformGridToCoord([]int{m - 1, n}, 100, 500)
		neighbours = append(neighbours, getNeighbours([]float64{coordDown[0], coord[1]}, 100, 500)...)
	}

	if m < len(uniformGrid.VertexData)-1 {
		coordUp := UniformGridToCoord([]int{m + 1, n}, 100, 500)
		neighbours = append(neighbours, getNeighbours([]float64{coordUp[0], coord[1]}, 100, 500)...)
	}
	return neighbours
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
	fLng = start[1] * math.Pi / 180.0
	fLat = start[0] * math.Pi / 180.0
	fLng2 = end[1] * math.Pi / 180.0
	fLat2 = end[0] * math.Pi / 180.0

	radius = 6378100
	h := haversin(fLat2-fLat) + math.Cos(fLat)*math.Cos(fLat2)*haversin(fLng2-fLng)
	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
	return (c * radius)
}

func extractRoute(prev *[]int, end int) [][][]float64 {
	//print("started extracting route\n")
	var route [][][]float64
	var tempRoute [][]float64
	m, n := uniformGrid.idToGrid(int(end))
	temp := []int{m, n}
	for {
		m, n = uniformGrid.idToGrid(int(end))
		x := []int{m, n}
		if math.Abs(float64(temp[0]-x[0])) > 1 {
			route = append(route, tempRoute)
			tempRoute = make([][]float64, 0)
		}
		tempRoute = append(tempRoute, UniformGridToCoord([]int{x[0], x[1]}, 100, 500))

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

	for i := 0; i < uniformGrid.N; i++ {
		dist = append(dist, math.Inf(1))
		prev = append(prev, -1)
	}
	startID := uniformGrid.gridToID(startLngInt, startLatInt)

	dist[startID] = 0
	vertices[startID] = true

	for {
		if len(vertices) == 0 {
			break
		} else {
			//var expansion = make([]int,0)
			var u = min(&dist, &vertices)
			m, n := uniformGrid.idToGrid(u)
			gridPos := []int{m, n}
			neighbours := neighboursUniformGrid(gridPos)
			delete(vertices, u)
			//fmt.Printf("after:\n%v\n", neighbours2d)

			for _, j := range neighbours {
				//fmt.Printf("j: %v, land:%v\n", j, meshgrid[j])
				if j == uniformGrid.gridToID(endLngInt, endLatInt) {
					prev[j] = u
					return extractRoute(&prev, uniformGrid.gridToID(endLngInt, endLatInt))
					//return testExtractRoute(&expansions)
				}
				//fmt.Printf("Dist[u]: %v\n", dist[u])
				//fmt.Printf("u/j: %v/%v\n", u,j)
				//fmt.Printf("Distance u-j: %v\n", distance(expandIndx(u), expandIndx(j)))
				//fmt.Printf("Summe: %v\n", (dist[u] + distance(expandIndx(u), expandIndx(j))))
				//fmt.Scanln()
				mu, nu := uniformGrid.idToGrid(u)
				p1 := []int{mu, nu}
				mj, nj := uniformGrid.idToGrid(j)
				p2 := []int{mj, nj}
				var alt = dist[u] + distance(UniformGridToCoord(p1, 100, 500), UniformGridToCoord(p2, 100, 500))
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

	return extractRoute(&prev, uniformGrid.gridToID(endLngInt, endLatInt))
	//return testExtractRoute(&expansions)
}

func main() {
	//meshgridRaw, errJSON := os.Open("data/output/meshgrid__planet_big.json")
	/*meshgridRaw, errJSON := os.Open("../data/output/meshgrid.json")
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
	}*/
	uniformGridRaw, errJSON2 := os.Open("../data/output/uniformGrid_10_500.json")
	if errJSON2 != nil {
		panic(errJSON2)
	}
	defer uniformGridRaw.Close()
	uniformByteValue, _ := ioutil.ReadAll(uniformGridRaw)
	json.Unmarshal(uniformByteValue, &uniformGrid)
	//fmt.Printf("%v\n",uniformGrid)

	var points [][][]float64
	for i := 0; i < len(uniformGrid.VertexData); i++ {
		var pointRow [][]float64
		for j := 0; j < len(uniformGrid.VertexData[i]); j++ {
			if uniformGrid.VertexData[i][j] {
				pointRow = append(pointRow, UniformGridToCoord([]int{i, j}, 100, 500))
			}
		}
		points = append(points, pointRow)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/testneighbours") {

			query := r.URL.Query()
			var lat, err = strconv.ParseFloat(query.Get("lat"), 10)
			if err != nil {
				panic(err)
			}
			var lng, err1 = strconv.ParseFloat(query.Get("lng"), 10)
			if err1 != nil {
				panic(err1)
			}
			fmt.Printf("Input:%v,%v\n", lat, lng)
			gridPoint := UniformCoordToGrid([]float64{lat, lng}, 100, 500)
			fmt.Printf("gridPoint %v\n", gridPoint)
			neighbours := testNeighboursUniformGrid(gridPoint)
			var coords [][]float64
			for _, x := range neighbours {
				fmt.Printf("neighbour: %v\n", x)
				coords = append(coords, UniformGridToCoord(x, 100, 500))
			}
			byteCoords, errByteCoords := json.Marshal(coords)
			if errByteCoords != nil {
				panic(errByteCoords)
			}
			w.Write(byteCoords)

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
			fmt.Printf("\n \n nat, lng %v,%v\n", lat, lng)
			gridPoint := UniformCoordToGrid([]float64{lat, lng}, 100, 500)
			var x = uniformGrid.gridToID(gridPoint[0], gridPoint[1])
			fmt.Printf("x %v\n", x)
			var m, n = uniformGrid.idToGrid(x)
			fmt.Printf("gridpoint %v\n", gridPoint)
			fmt.Printf("m,n %v %v\n", m, n)

			coords := UniformGridToCoord(gridPoint, 100, 500)
			fmt.Printf("coords %v,%v\n", coords[0], coords[1])
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

			var start = UniformCoordToGrid([]float64{startLat, startLng}, 100, 500)
			var startM = start[0]
			var startN = start[1]

			var end = UniformCoordToGrid([]float64{endLat, endLng}, 100, 500)
			var endM = end[0]
			var endN = end[1]

			fmt.Printf("\nstartLngInt, startLatInt, endLngInt, endLatInt\n")
			fmt.Printf("%v/%v, %v/%v\n", startM, startN, endM, endN)
			//fmt.Printf("%v/%v\n", meshgrid2d[startLngInt][startLatInt], meshgrid2d[endLngInt][endLatInt])

			if !uniformGrid.VertexData[startM][startN] && !uniformGrid.VertexData[endM][endN] {
				//fmt.Printf("start: %d / %d ", int(math.Round(startLat)), int(math.Round(startLng)))
				//fmt.Printf("end: %d / %d\n", int(math.Round(endLat)), int(math.Round(endLng)))
				var start = time.Now()
				var route = dijkstra(startM, startN, endM, endN)

				t := time.Now()
				elapsed := t.Sub(start)
				fmt.Printf("time: %s\n", elapsed)
				//mt.Printf("%v\n", route)
				var result = toGeojson(route)
				w.Write(result)
			} else {
				w.Write([]byte("false"))
			}

		} else if strings.Contains(r.URL.Path, "/grid") {
			pointsJSON, err := json.Marshal(points)
			if err != nil {
				panic(err)
			}
			w.Write(pointsJSON)
		} else {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
