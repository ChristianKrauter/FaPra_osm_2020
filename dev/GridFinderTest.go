package main

import (
	//"../src/algorithms"
	"../src/algorithms"
	"../src/grids"
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
	//"sort"
	"strconv"
	"strings"
	"time"
)

var port int = 8081

var algoStr = []string{
	"Dij   ",
	"A*    ",
	"BiDij ",
	"BiA*  ",
	"A*-JPS",
}

func toGeojson(route [][][]float64) *geojson.FeatureCollection {
	routes := geojson.NewFeatureCollection()
	for _, j := range route {
		routes = routes.AddFeature(geojson.NewFeature(geojson.NewLineStringGeometry(j)))
	}
	return routes
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

// TestData ...
type TestData struct {
	Point []float64
	Nbs   [][]float64
	Nnbs  [][]float64
}

func main() {
	xSize := 360
	ySize := 360
	basicPointInPolygon := false

	var ug1D []bool
	var ug grids.UniformGrid

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
			var from = ug.CoordToGrid(startLng, startLat)
			var gridNodeCoord = ug.GridToCoord(from)
			var td TestData
			var nbs = algorithms.NeighboursUg(ug.GridToID(from), &ug)
			for _, i := range nbs {
				td.Nbs = append(td.Nbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			var nnbs = algorithms.SimpleNeighboursUg(ug.GridToID(from), &ug)
			for _, i := range nnbs {
				td.Nnbs = append(td.Nnbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			td.Point = gridNodeCoord
			tdJSON, err := json.Marshal(td)
			if err != nil {
				panic(err)
			}
			w.Write(tdJSON)

		} else if strings.Contains(r.URL.Path, "/route") {
			query := r.URL.Query()
			var x, err = strconv.ParseInt(query.Get("gridIDXx"), 10, 64)
			if err != nil {
				panic(err)
			}
			var y, err1 = strconv.ParseInt(query.Get("gridIDXy"), 10, 64)
			if err1 != nil {
				panic(err1)
			}

			fmt.Printf("GridPoint\n")
			fmt.Printf("%v,%v\n", x, y)
			var result = make([]*geojson.FeatureCollection, 5)
			var routes = make([]*[][][]float64, 5)
			var lengths = make([]float64, 5)

			routes[0], _, lengths[0] = algorithms.Dijkstra(int(x), int(y), &ug)
			routes[1], _, lengths[1] = algorithms.AStar(int(x), int(y), &ug)
			routes[2], _, lengths[2] = algorithms.BiDijkstra(int(x), int(y), &ug)
			routes[3], _, lengths[3] = algorithms.BiAStar(int(x), int(y), &ug)
			routes[4], _, lengths[4] = algorithms.AStarJPS(int(x), int(y), &ug)

			for i := 0; i < 5; i++ {
				result[i] = toGeojson(*routes[i])
				fmt.Printf("(%s) length: %v\n", algoStr[i], lengths[i])
			}

			tdJSON, err := json.Marshal(result)
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
			var nbs = algorithms.NeighboursUg(id, &ug)
			for _, i := range nbs {
				td.Nbs = append(td.Nbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			var nnbs = algorithms.SimpleNeighboursUg(id, &ug)
			for _, i := range nnbs {
				td.Nnbs = append(td.Nnbs, ug.GridToCoord(ug.IDToGrid(i)))
			}
			td.Point = gridNodeCoord
			tdJSON, err := json.Marshal(td)
			if err != nil {
				panic(err)
			}
			w.Write(tdJSON)

		} else if strings.Contains(r.URL.Path, "/startend") {
			query := r.URL.Query()
			var gridIDXx, err = strconv.ParseInt(query.Get("gridIDXx"), 10, 64)
			if err != nil {
				panic(err)
			}
			var gridIDXy, err2 = strconv.ParseInt(query.Get("gridIDXy"), 10, 64)
			if err2 != nil {
				panic(err2)
			}

			var gridIDXxCoord = ug.GridToCoord(ug.IDToGrid(int(gridIDXx)))
			var gridIDXyCoord = ug.GridToCoord(ug.IDToGrid(int(gridIDXy)))
			var td = make([][]float64, 2)
			td[0] = gridIDXxCoord
			td[1] = gridIDXyCoord

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
