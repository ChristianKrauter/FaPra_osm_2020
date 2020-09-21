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

	var bg grids.BasicGrid
	var bg2D [][]bool
	var filename string

	bg.XSize = xSize
	bg.YSize = ySize
	bg.XFactor = float64(xSize) / 360.0
	bg.YFactor = float64(ySize) / 360.0

	if basicPointInPolygon {
		filename = fmt.Sprintf("../data/output/meshgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("../data/output/meshgrid_%v_%v.json", xSize, ySize)
	}

	meshgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		log.Fatal(fmt.Sprintf("\nThe meshgrid '%s'\ncould not be found. Please create it first.\n", filename))
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &bg2D)

	var points [][]float64
	bg.VertexData = make([]bool, xSize*ySize)
	k := 0
	for i := 0; i < len(bg2D[0]); i++ {
		for j := 0; j < len(bg2D); j++ {
			bg.VertexData[k] = bg2D[j][i]
			k++
			if !bg2D[i][j] {
				points = append(points, bg.GridToCoord([]int{int(i), int(j)}))
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
			var from = bg.CoordToGrid([]float64{startLng, startLat})
			var gridNodeCoord = bg.GridToCoord(from)
			var td TestData
			var nbs = algorithms.NeighboursBg(bg.GridToID(from), &bg)
			for _, i := range nbs {
				td.Nbs = append(td.Nbs, bg.GridToCoord(bg.IDToGrid(i)))
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

/* else if strings.Contains(r.URL.Path, "/id") {
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

} */
