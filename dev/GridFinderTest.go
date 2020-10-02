package main

import (
	//"../src/algorithms"
	"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	//"github.com/qedus/osmpbf"
	//"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	//"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var port int = 8081

// UniformGrid structure
type UniformGrid struct {
	XSize        int
	YSize        int
	N            int
	BigN         int
	A            float64
	D            float64
	MTheta       float64
	DTheta       float64
	DPhi         float64
	VertexData   [][]bool
	FirstIndexOf []int
}

// GridToCoord takes grid indices and outputs lng lat
func (ug UniformGrid) GridToCoord(in []int) []float64 {
	theta := math.Pi * (float64(in[0]) + 0.5) / float64(ug.MTheta)
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
	phi := 2 * math.Pi * float64(in[1]) / mPhi
	return []float64{(phi / math.Pi) * 180.0, (theta/math.Pi)*180.0 - 90.0}
}

// CoordToGrid takes lng lat and outputs grid indices
func (ug UniformGrid) CoordToGrid(lng, lat float64) []int {
	theta := (lat + 90.0) * math.Pi / 180.0
	m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
	theta = math.Pi * (float64(m) + 0.5) / float64(ug.MTheta)
	var phi float64
	if lng < 0 {
		phi = float64(lng+360.0) * math.Pi / 180.0
	} else {
		phi = lng * math.Pi / 180.0
	}
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
	n := math.Round(phi * mPhi / (2.0 * math.Pi))
	return []int{mod(int(m), int(ug.MTheta)), mod(int(n), int(mPhi))}
}

// GridToID ...
func (ug UniformGrid) GridToID(IDX []int) int {
	return ug.FirstIndexOf[IDX[0]] + IDX[1]
}

// IDToGrid ...
func (ug UniformGrid) IDToGrid(id int) []int {
	m := sort.Search(len(ug.FirstIndexOf), func(i int) bool { return ug.FirstIndexOf[i] > id }) - 1
	n := id - ug.FirstIndexOf[m]
	return []int{m, n}
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

func boundingContains(bounding *map[string]float64, point []float64) bool {
	if (*bounding)["minX"] <= point[0] && point[0] <= (*bounding)["maxX"] {
		if (*bounding)["minY"] <= point[1] && point[1] <= (*bounding)["maxY"] {
			return true
		}
	}
	return false
}

func createBoundingBox(polygon *[][]float64) map[string]float64 {
	minX := math.Inf(1)
	maxX := math.Inf(-1)
	minY := math.Inf(1)
	maxY := math.Inf(-1)
	for _, coord := range *polygon {
		if coord[0] < minX {
			minX = coord[0]
		} else if coord[0] > maxX {
			maxX = coord[0]
		}
		if coord[1] < minY {
			minY = coord[1]
		} else if coord[1] > maxY {
			maxY = coord[1]
		}
	}
	return map[string]float64{"minX": minX, "maxX": maxX, "minY": minY, "maxY": maxY}
}

//Check if a bounding box is inside another bounding box
func checkBoundingBoxes(bb1 map[string]float64, bb2 map[string]float64) bool {
	return bb1["minX"] >= bb2["minX"] && bb1["maxX"] <= bb2["maxX"] && bb1["minY"] >= bb2["minY"] && bb1["maxY"] <= bb2["maxY"]
}

func addBoundingTree(tree *boundingTree, boundingBox *map[string]float64, id int) boundingTree {
	for i, child := range (*tree).children {
		if checkBoundingBoxes(*boundingBox, child.boundingBox) {
			child = addBoundingTree(&child, boundingBox, id)
			(*tree).children[i] = child
			return *tree
		}
	}
	(*tree).children = append((*tree).children, boundingTree{*boundingBox, id, make([]boundingTree, 0)})
	return *tree
}
func createBoundingTree(boundingTreeRoot *boundingTree, allCoastlines *[][][]float64) string {
	start := time.Now()

	*boundingTreeRoot = boundingTree{map[string]float64{"minX": math.Inf(-1), "maxX": math.Inf(1), "minY": math.Inf(-1), "maxY": math.Inf(1)}, -1, make([]boundingTree, 0)}

	for j, i := range *allCoastlines {
		bb := createBoundingBox(&i)
		*boundingTreeRoot = addBoundingTree(boundingTreeRoot, &bb, j)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created bounding tree            : %s\n", elapsed)
	return elapsed.String()
}

type boundingTree struct {
	boundingBox map[string]float64
	id          int
	children    []boundingTree
}

// Sorting arrays by length
type arrayOfArrays [][][]float64

func (p arrayOfArrays) Len() int {
	return len(p)
}

func (p arrayOfArrays) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p arrayOfArrays) Less(i, j int) bool {
	return len(p[i]) > len(p[j])
}

func getSomeKey(m *map[int64][]int64) int64 {
	for k := range *m {
		return k
	}
	return 0
}

// TestData ...
type TestData struct {
	Point []float64
	Nbs   [][]float64
	Nnbs  [][]float64
}

func jump(x, dir int, fromIDX, toIDX []int, ug *UniformGrid) []int{
	n := step(x,dir,ug)
	//xCoord := ug.IDToGrid(x)
	//fmt.Printf("xdiff: %v,\nydiff: %v\n\n", xCoord, n)
	if(n == nil || ug.VertexData[n[0]][n[1]]){
		/*if(n != nil){
			*nodesProcessed = append(*nodesProcessed, (*ug).GridToID(n))
		}*/
		return nil
	}
	if(n[0] == toIDX[0] && n[1] == toIDX[1]){
		return n
	}
	if(isForced(ug.GridToID(n),dir, ug)){
		return n 
	}
	if(dir == 0 || dir == 2 || dir == 4 || dir == 6){
		temp := jump(ug.GridToID(n),mod(dir-1,8),fromIDX, toIDX, ug)
		if(temp != nil){
			return n
		}
		temp = jump(ug.GridToID(n),mod(dir+1,8),fromIDX, toIDX, ug)
		if(temp != nil){
			return n
		}
	}
	return jump(ug.GridToID(n),dir,fromIDX, toIDX, ug)
}

func step(in,dir int, ug *UniformGrid) []int{

	var allNeighbours [8][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	// lengths of rows
	var lm = len(ug.VertexData[m])

	
		allNeighbours[7] = []int{m, mod(n-1, lm)}	
		allNeighbours[3] = []int{m, mod(n+1, lm)}	
	
	

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData) -1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))
		allNeighbours[5] = []int{m + 1, mod(nUp, lmp)}
		allNeighbours[4] = []int{m + 1, mod(nUp+1.0, lmp)}
		allNeighbours[6] = []int{m + 1, mod(nUp-1.0, lmp)}
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))
		allNeighbours[1] = []int{m - 1, mod(nDown, lmm)}
		allNeighbours[2] = []int{m - 1, mod(nDown+1.0, lmm)}
		allNeighbours[0] = []int{m - 1, mod(nDown-1.0, lmm)}
	}
	 if(len(allNeighbours[dir]) < 1 || m >= len(ug.VertexData) -1 || n < 1 || n >= len(ug.VertexData[m]) -1){
	 	return nil
	 } else {
	 	return allNeighbours[dir]
	 }
	
}


func isForced(in, dir int, ug *UniformGrid) bool {
	var allNeighbours [8][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	// lengths of rows
	var lm = len(ug.VertexData[m])

	allNeighbours[7] = []int{m, mod(n-1, lm)}
	allNeighbours[3] = []int{m, mod(n+1, lm)}

	ratio = float64(n) / float64(lm)

	if m < len(ug.VertexData) -1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))
		allNeighbours[5] = []int{m + 1, mod(nUp, lmp)}
		allNeighbours[4] = []int{m + 1, mod(nUp+1.0, lmp)}
		allNeighbours[6] = []int{m + 1, mod(nUp-1.0, lmp)}
	} else {
		return false
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))
		allNeighbours[1] = []int{m - 1, mod(nDown, lmm)}
		allNeighbours[2] = []int{m - 1, mod(nDown+1.0, lmm)}
		allNeighbours[0] = []int{m - 1, mod(nDown-1.0, lmm)}
	}

	if(dir == 0 || dir == 2 || dir == 4 || dir == 6){
		if ug.VertexData[allNeighbours[mod(dir-3, 8)][0]][allNeighbours[mod(dir-3, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir-2, 8)][0]][allNeighbours[mod(dir-2, 8)][1]] {
			return true
		}
		if ug.VertexData[allNeighbours[mod(dir+3, 8)][0]][allNeighbours[mod(dir+3, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir+2, 8)][0]][allNeighbours[mod(dir+2, 8)][1]] {
			return true
		}
	} else {
		if ug.VertexData[allNeighbours[mod(dir-2, 8)][0]][allNeighbours[mod(dir-2, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir-1, 8)][0]][allNeighbours[mod(dir-1, 8)][1]] {
			return true
		}
		if ug.VertexData[allNeighbours[mod(dir+2, 8)][0]][allNeighbours[mod(dir+2, 8)][1]] && !ug.VertexData[allNeighbours[mod(dir+1, 8)][0]][allNeighbours[mod(dir+1, 8)][1]] {
			return true
		}
	}

	
	return false
}

func main() {
	xSize := 100
	ySize := 500
	basicPointInPolygon := false

	var ug1D []bool
	var ug UniformGrid

	var filename string
	if basicPointInPolygon {
		filename = fmt.Sprintf("../data/output/uniformgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("../data/output/uniformgrid_%v_%v.json", xSize, ySize)
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var fileEnding = strings.Split(r.URL.Path[1:], ".")[len(strings.Split(r.URL.Path[1:], "."))-1]

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
			//fmt.Printf("point\n")
			//fmt.Printf("CoordToGrid\n")
			var from = ug.CoordToGrid(startLng, startLat)
			//fmt.Printf("%v\n", from)
			//fmt.Printf("GridToCoord\n")
			var gridNodeCoord = ug.GridToCoord(from)
			//fmt.Printf("%v\n", gridNodeCoord)
			var td TestData
			var nbs = neighboursUg(ug.GridToID(from), &ug)
			for _, i := range nbs {
				td.Nbs = append(td.Nbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			var nnbs = simpleNeighbours(ug.GridToID(from), &ug)
			for _, i := range nnbs {
				td.Nnbs = append(td.Nnbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			td.Point = gridNodeCoord
			//fmt.Printf("%v\n", td)
			tdJSON, err := json.Marshal(td)
			if err != nil {
				panic(err)
			}
			w.Write(tdJSON)

		} else if strings.Contains(r.URL.Path, "/jumpPoint") {
			query := r.URL.Query()
			var x, err = strconv.ParseInt(query.Get("x"), 10, 64)
			if err != nil {
				panic(err)
			}
			var y, err1 = strconv.ParseInt(query.Get("y"), 10, 64)
			if err1 != nil {
				panic(err1)
			}
			var dir, err2 = strconv.ParseInt(query.Get("dir"), 10, 64)
			if err2 != nil {
				panic(err2)
			}
			fmt.Printf("%v, %v, %v\n", int(x), int(y), dir)
			id := ug.GridToID([]int{int(x), int(y)})
			jump := jump(id,int(dir),[]int{0,0},ug.CoordToGrid(0.0,0.0),&ug)
			point := ug.GridToCoord(jump)
			fmt.Printf("%v\n", point)
			var td TestData
			td.Point = point
			//fmt.Printf("%v\n", td)
			tdJSON, err := json.Marshal(td)
			if err != nil {
				panic(err)
			}
			w.Write(tdJSON)

		} else if strings.Contains(r.URL.Path, "/gridPoint") {
			query := r.URL.Query()
			var x, err = strconv.ParseInt(query.Get("x"), 10, 64)
			if err != nil {
				panic(err)
			}
			var y, err1 = strconv.ParseInt(query.Get("y"), 10, 64)
			if err1 != nil {
				panic(err1)
			}
			fmt.Printf("%v,%v\n", x, y)
			fmt.Printf("GridPoint\n")
			fmt.Printf("CoordToGrid\n")
			point := ug.GridToCoord([]int{int(x), int(y)})
			fmt.Printf("%v\n", point)
			var td TestData
			td.Point = point
			//fmt.Printf("%v\n", td)
			tdJSON, err := json.Marshal(td)
			if err != nil {
				panic(err)
			}
			w.Write(tdJSON)

		} else if strings.Contains(r.URL.Path, "/id") {
			query := r.URL.Query()
			var ID, err = strconv.ParseInt(query.Get("id"), 10, 64)
			if err != nil {
				panic(err)
			}
			var id = int(ID)
			var gridNodeCoord = ug.GridToCoord(ug.IDToGrid(id))
			var td TestData
			var nbs = neighboursUg(id, &ug)
			for _, i := range nbs {
				td.Nbs = append(td.Nbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			var nnbs = simpleNeighbours(id, &ug)
			for _, i := range nnbs {
				td.Nnbs = append(td.Nnbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			td.Point = gridNodeCoord
			tdJSON, err := json.Marshal(td)
			if err != nil {
				panic(err)
			}
			w.Write(tdJSON)

		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else if strings.Contains(r.URL.Path, "/basicGrid") {
			pointsJSON, err := json.Marshal(points)
			if err != nil {
				panic(err)
			}
			w.Write(pointsJSON)
		} else {
			http.ServeFile(w, r, "globe.html")
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf(" on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func simpleNeighbours(in int, ug *UniformGrid) []int {
	var neighbours [][]int
	var inGrid = ug.IDToGrid(in)
	var ratio float64
	var nUp, nDown int
	var m = inGrid[0]
	var n = inGrid[1]

	// lengths of rows
	var lm = len(ug.VertexData[m])

	neighbours = append(neighbours, []int{m, mod(n-1, lm)})
	neighbours = append(neighbours, []int{m, mod(n+1, lm)})

	ratio = float64(n) / float64(lm)

	//fmt.Printf("m,n: %v,%v\nlen(m), len(m+1), len(m-1): %v,%v,%v\nratio: %v\nnUp, nDown  : %v,%v\nxnUp, xnDown: %v,%v\n",
	//	m, n, lm, lmp, lmm, ratio,
	//	nUp, nDown, mod(nUp, lmp), mod(nDown, lmm))
	//fmt.Printf("%v\n", mod(nUp, lmp))

	if m < len(ug.VertexData)-1 {
		var lmp = len(ug.VertexData[m+1])
		nUp = int(math.Round(ratio * float64(lmp)))

		neighbours = append(neighbours, []int{m + 1, mod(nUp, lmp)})
		neighbours = append(neighbours, []int{m + 1, mod(nUp+1.0, lmp)})
		neighbours = append(neighbours, []int{m + 1, mod(nUp-1.0, lmp)})
	}

	if m > 0 {
		var lmm = len(ug.VertexData[m-1])
		nDown = int(math.Round(ratio * float64(lmm)))

		neighbours = append(neighbours, []int{m - 1, mod(nDown, lmm)})
		neighbours = append(neighbours, []int{m - 1, mod(nDown+1.0, lmm)})
		neighbours = append(neighbours, []int{m - 1, mod(nDown-1.0, lmm)})
	}

	var neighbours1d []int
	for _, neighbour := range neighbours {
		//if !ug.VertexData[neighbour[0]][neighbour[1]] {
		neighbours1d = append(neighbours1d, ug.GridToID(neighbour))
		//}
	}
	return neighbours1d
}

// Gets neighours left and right in the same row
func neighboursRowUg(in []float64, ug *UniformGrid) [][]int {
	theta := (in[1] + 90) * math.Pi / 180
	m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
	theta = math.Pi * (m + 0.5) / ug.MTheta
	phi := in[0] * math.Pi / 180
	mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)

	n1 := math.Round(phi * mPhi / (2 * math.Pi))
	p1 := []int{mod(int(m), int(ug.MTheta)), mod(int(n1), int(mPhi))}
	p2 := []int{mod(int(m), int(ug.MTheta)), mod(int(n1+1), int(mPhi))}
	p3 := []int{mod(int(m), int(ug.MTheta)), mod(int(n1-1), int(mPhi))}
	return [][]int{p1, p2, p3}
}

// neighboursUg gets up to 8 neighbours
func neighboursUg(in int, ug *UniformGrid) []int {
	var neighbours [][]int
	var inGrid = ug.IDToGrid(in)
	m := inGrid[0]
	n := inGrid[1]
	neighbours = append(neighbours, []int{m, mod(n-1, len(ug.VertexData[m]))})
	neighbours = append(neighbours, []int{m, mod(n+1, len(ug.VertexData[m]))})
	fmt.Printf("\nIngrid: %v\n ", inGrid)
	fmt.Printf("nbs same row: %v\n ", neighbours)
	coord := ug.GridToCoord(inGrid)

	if m > 0 {
		fmt.Printf("m-1: %v\n", m-1)
		coordDown := ug.GridToCoord([]int{m - 1, n})
		neighbours = append(neighbours, neighboursRowUg([]float64{coord[0], coordDown[1]}, ug)...)
	}
	fmt.Printf("m, len(VD): %v, %v\n", m, len(ug.VertexData))
	if m < len(ug.VertexData)-1 {
		coordUp := ug.GridToCoord([]int{m + 1, n})
		neighbours = append(neighbours, neighboursRowUg([]float64{coord[0], coordUp[1]}, ug)...)
	}

	var neighbours1d []int
	for _, neighbour := range neighbours {
		//if !ug.VertexData[neighbour[0]][neighbour[1]] {
		neighbours1d = append(neighbours1d, ug.GridToID(neighbour))
		//}
	}
	return neighbours1d
}
