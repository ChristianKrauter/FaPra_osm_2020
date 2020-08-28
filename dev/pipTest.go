package main

import (
	"../src/grids"
	"encoding/json"
	"fmt"
	"github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"sort"
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

func createPolygons(allCoastlines *[][][]float64, coastlineMap *map[int64][]int64, nodeMap *map[int64][]float64) string {
	start := time.Now()
	var coastline [][]float64

	for len(*coastlineMap) > 0 {
		var key = getSomeKey(coastlineMap)
		var nodeIDs = (*coastlineMap)[key]
		coastline = nil
		for _, x := range nodeIDs {
			coastline = append(coastline, []float64{(*nodeMap)[x][0], (*nodeMap)[x][1]})
		}
		delete(*coastlineMap, key)
		key = nodeIDs[len(nodeIDs)-1]
		for {
			if val, ok := (*coastlineMap)[key]; ok {
				for i, x := range val {
					if i != 0 {
						coastline = append(coastline, []float64{(*nodeMap)[x][0], (*nodeMap)[x][1]})
					}
				}
				delete(*coastlineMap, key)
				key = val[len(val)-1]
			} else {
				break
			}
		}
		*allCoastlines = append(*allCoastlines, coastline)
	}

	sort.Sort(arrayOfArrays(*allCoastlines))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Made all polygons in             : %s\n", elapsed)

	var temp = []float64{-400, -400}
	for _, i := range *allCoastlines {
		for _, j := range i {
			if temp[0] == j[0] && temp[1] == j[1] {
				fmt.Printf("bad")
			}
			temp = j
		}
	}

	return elapsed.String()
}

func readFile(pbfFileName string, coastlineMap *map[int64][]int64, nodeMap *map[int64][]float64) string {
	start := time.Now()

	// Read coastlines
	f, err := os.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)
	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(runtime.NumCPU()))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				(*nodeMap)[v.ID] = []float64{v.Lon, v.Lat}
			case *osmpbf.Way:
				for _, value := range v.Tags {
					if value == "coastline" {
						(*coastlineMap)[v.NodeIDs[0]] = v.NodeIDs
					}
				}
			case *osmpbf.Relation:
				continue
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Read file in                     : %s\n", elapsed)
	return elapsed.String()
}

func transformLon(newNorth, point []float64) float64 {
	var transformedLon float64

	var dtr = math.Pi / 180.0

	// New north is already the north pole
	if newNorth[1] == 90.0 {
		transformedLon = point[0]
	} else {
		var t = math.Sin((point[0]-newNorth[0])*dtr) * math.Cos(point[1]*dtr)
		var b = math.Sin(dtr*point[1])*math.Cos(newNorth[1]*dtr) - math.Cos(point[1]*dtr)*math.Sin(newNorth[1]*dtr)*math.Cos((point[0]-newNorth[0])*dtr)
		transformedLon = math.Atan2(t, b) / dtr
	}

	return transformedLon
}

// Direction of the shortest path from a to b
// 1 = east, -1 = west, 0 = neither
func eastOrWest(aLon, bLon float64) int {
	var out int
	var del = bLon - aLon
	if del > 180.0 {
		del = del - 360.0
	}
	if del < -180.0 {
		del = del + 360.0
	}
	if del > 0 && del != 180.0 {
		out = -1
	} else if del < 0 && del != -180.0 {
		out = 1
	} else {
		out = 0
	}
	return out
}

func pointInPolygonSphere(polygon *[][]float64, point []float64, strikes *[][][]float64) bool {
	//fmt.Printf("%v\n","New PP")
	var inside = false
	var strike = false
	// Point is the south-pole
	// Pontentially antipodal check
	fmt.Printf("point: %v\n", point)
	if point[0] <= -80 {
		fmt.Printf("Tried to check point antipodal to the north pole.")
		return true
	}

	//fmt.Printf("%v\n",point)

	// Point is the north-pole
	if point[1] == 90 {
		fmt.Printf("south pole.")
		return false
	}
	for i := 0; i < len(*polygon); i++ {
		var a = (*polygon)[i]
		var b = (*polygon)[(i+1)%len(*polygon)]

		var nortPole = []float64{0.0, 90.0}
		if a[0] == b[0] {
			//fmt.Printf("a&b great circle\n")
			nortPole = []float64{0.01, 90.0}
			//point[0] += 0.001
		}

		strike = false

		if point[0] == a[0] && point[0] == b[0] {
			strike = true
		} else {

			var aToB = eastOrWest(a[0], b[0])
			var aToP = eastOrWest(a[0], point[0])
			var pToB = eastOrWest(point[0], b[0])

			if aToP == aToB && pToB == aToB {
				strike = true
			}
		}

		/*if strike && point[1] > a[1] && point[1] > b[1] {
			strike = false
		}*/
		if strike {
			*strikes = append(*strikes, [][]float64{a, b})
			if point[1] == a[1] && point[0] == a[0] {
				fmt.Printf("p=a\n")
				return true
			}

			// Possible to calculate once at polygon creation
			var northPoleLonTransformed = transformLon(a, nortPole)
			var bLonTransformed = transformLon(a, b)
			// Not possible
			var pLonTransformed = transformLon(a, point)

			if bLonTransformed == pLonTransformed {
				fmt.Printf("blon = plon\n")
				return true
			}

			var bToX = eastOrWest(bLonTransformed, northPoleLonTransformed)
			var bToP = eastOrWest(bLonTransformed, pLonTransformed)
			if bToX == -bToP {

				inside = !inside
			}
		}
	}

	return inside
}

func isLandSphere(tree *boundingTree, point []float64, allCoastlines *[][][]float64, strikes *[][][]float64) bool {
	land := false
	if boundingContains(&tree.boundingBox, point) {
		if (*tree).id >= 0 {
			land = pointInPolygonSphere(&(*allCoastlines)[(*tree).id], point, strikes)
			if land {
				return land
			}
		}
		for _, child := range (*tree).children {
			land = isLandSphere(&child, point, allCoastlines, strikes)
			if land {
				return land
			}
		}
	}
	return land
}

// TestData ...
type TestData struct {
	IsLand  bool
	Strikes [][][]float64
}

func main() {

	pbfFileName := "antarctica-latest.osm.pbf"
	fmt.Printf("\nStarting processing of %s\n\n", pbfFileName)
	pbfFileName = fmt.Sprintf("../data/%s", pbfFileName)

	// Read the pbf file
	var coastlineMap = make(map[int64][]int64)
	var nodeMap = make(map[int64][]float64)
	readFile(pbfFileName, &coastlineMap, &nodeMap)

	// Create coastline polygons
	var allCoastlines [][][]float64
	createPolygons(&allCoastlines, &coastlineMap, &nodeMap)

	// Create bounding boxes
	var boundingTreeRoot boundingTree
	createBoundingTree(&boundingTreeRoot, &allCoastlines)

	xSize := 100
	ySize := 500
	basicPointInPolygon := false
	//log.Fatal("Server with unidistant grid not implemented.")

	var uniformgrid []bool
	var uniformgrid2d grids.UniformGrid

	var filename string
	if basicPointInPolygon {
		filename = fmt.Sprintf("../data/output/uniformgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("../data/output/uniformgrid_%v_%v.json", xSize, ySize)
	}

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &uniformgrid2d)

	for i := 0; i < len(uniformgrid2d.VertexData); i++ {
		for j := 0; j < len(uniformgrid2d.VertexData[i]); j++ {
			uniformgrid = append(uniformgrid, uniformgrid2d.VertexData[i][j])
		}
	}

	var points [][]float64
	for i := 0; i < len(uniformgrid2d.VertexData); i++ {
		for j := 0; j < len(uniformgrid2d.VertexData[i]); j++ {
			if !uniformgrid2d.VertexData[i][j] {
				points = append(points, uniformgrid2d.GridToCoord([]int{int(i), int(j)}))
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
			var strikes [][][]float64
			var td TestData
			if startLat <= -80 {
				fmt.Printf("Tried to check point antipodal to the north pole.")
				td.IsLand = true
			} else {
				td.IsLand = isLandSphere(&boundingTreeRoot, []float64{startLng, startLat}, &allCoastlines, &strikes)
			}
			td.Strikes = strikes
			//fmt.Printf("%v\n", td)
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
