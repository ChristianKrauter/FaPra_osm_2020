package testserver

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
)

var port int = 8082
var algoStr = []string{"Dij   ", "A*    ", "BiDij ", "BiA*  ", "A*-JPS"}

// TestData ...
type TestData struct {
	Point []float64
	Nbs   [][]float64
	Nnbs  [][]float64
}

// Start test server
func Start(xSize, ySize int) {
	var ug1D []bool
	var ug grids.UniformGrid
	var filename = fmt.Sprintf("data/output/uniformGrid_%v_%v.json", xSize, ySize)

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

		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" || fileEnding == "ico" {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else if strings.Contains(r.URL.Path, "/grid") {
			pointsJSON, err := json.Marshal(points)
			if err != nil {
				panic(err)
			}
			w.Write(pointsJSON)
		} else {
			http.ServeFile(w, r, "src/testserver/globe.html")
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("on localhost%s\n\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

// StartBg test server
func StartBg(xSize, ySize int) {
	var bg grids.BasicGrid
	var bg2D [][]bool
	var filename string

	bg.XSize = xSize
	bg.YSize = ySize
	bg.XFactor = float64(xSize) / 360.0
	bg.YFactor = float64(ySize) / 360.0

	filename = fmt.Sprintf("data/output/meshgrid_%v_%v.json", xSize, ySize)

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
			var routes = make([][][][]float64, 5)
			var lengths = make([]float64, 5)

			routes[0], _, lengths[0] = algorithms.DijkstraBg(int(x), int(y), &bg)
			routes[1], _, lengths[1] = algorithms.AStarBg(int(x), int(y), &bg)
			routes[2], _, lengths[2] = algorithms.BiDijkstraBg(int(x), int(y), &bg)
			routes[3], _, lengths[3] = algorithms.BiAStarBg(int(x), int(y), &bg)
			routes[4], _, lengths[4] = algorithms.AStarJPSBg(int(x), int(y), &bg)

			for i := 0; i < 5; i++ {
				result[i] = toGeojson(routes[i])
				fmt.Printf("(%s) length: %v\n", algoStr[i], lengths[i])
			}

			tdJSON, err := json.Marshal(result)
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

			var gridIDXxCoord = bg.GridToCoord(bg.IDToGrid(int(gridIDXx)))
			var gridIDXyCoord = bg.GridToCoord(bg.IDToGrid(int(gridIDXy)))
			var td = make([][]float64, 2)
			td[0] = gridIDXxCoord
			td[1] = gridIDXyCoord
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
			var gridNodeCoord = bg.GridToCoord(bg.IDToGrid(id))
			var td TestData
			var nbs = algorithms.NeighboursBg(id, &bg)
			for _, i := range nbs {
				td.Nbs = append(td.Nbs, bg.GridToCoord(bg.IDToGrid(i)))
				td.Nnbs = append(td.Nnbs, []float64{})
			}
			td.Point = gridNodeCoord
			tdJSON, err := json.Marshal(td)
			if err != nil {
				panic(err)
			}

			w.Write(tdJSON)

		} else if fileEnding == "js" || fileEnding == "html" || fileEnding == "css" || fileEnding == "ico" {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else if strings.Contains(r.URL.Path, "/grid") {
			pointsJSON, err := json.Marshal(points)
			if err != nil {
				panic(err)
			}
			w.Write(pointsJSON)
		} else {
			http.ServeFile(w, r, "src/testserver/globe.html")
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("on localhost%s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func toGeojson(route [][][]float64) *geojson.FeatureCollection {
	routes := geojson.NewFeatureCollection()
	for _, j := range route {
		routes = routes.AddFeature(geojson.NewFeature(geojson.NewLineStringGeometry(j)))
	}
	return routes
}
