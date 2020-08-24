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
	fmt.Printf("m: %v\n",in[0])
	fmt.Printf("n: %v\n",in[1])
    fmt.Printf("theta: %v\n",theta)
    fmt.Printf("mPhi: %v\n",mPhi)
    fmt.Printf("phi: %v\n",phi)
    fmt.Printf("lon,lat: %v,%v\n",(phi / math.Pi) * 180.0,(theta/math.Pi)*180.0 - 90.0)
    return []float64{(phi / math.Pi) * 180.0, (theta/math.Pi)*180.0 - 90.0}
}

// CoordToGrid takes lng lat and outputs grid indices
func (ug UniformGrid) CoordToGrid(lng, lat float64) []int {

    theta := ((lat + 90.0)/ 180.0) * math.Pi
    
    m := math.Round((theta * ug.MTheta / math.Pi) - 0.5)
    
    theta = math.Pi * (float64(m) + 0.5) / float64(ug.MTheta)

    var phi float64
    if(lng < 0){
    	phi = float64(lng+360.0) * math.Pi / 180.0	
    } else {
    	phi = lng * math.Pi / 180.0	
    }
    
    
    mPhi := math.Round(2.0 * math.Pi * math.Sin(theta) / ug.DPhi)
    
    fmt.Printf("n not rounded: %v\n", phi * mPhi / (2.0 * math.Pi))
    n := math.Round(phi * mPhi / (2.0 * math.Pi))

    fmt.Printf("theta: %v\n",theta)
    fmt.Printf("m: %v\n",m)
    fmt.Printf("phi: %v\n",phi)
    fmt.Printf("mPhi: %v\n",mPhi)
    fmt.Printf("n: %v\n",n)
    return []int{mod(int(m), int(ug.MTheta)), mod(int(n), int(mPhi))}
}

// GridToID ...
func (ug UniformGrid) GridToID(IDX []int) int {
    return ug.FirstIndexOf[IDX[0]] + IDX[1]
}

// IDToGrid ...
func (ug UniformGrid) IDToGrid(id int) []int {
    m := sort.Search(len(ug.FirstIndexOf)-1, func(i int) bool { return ug.FirstIndexOf[i] > id })
    n := id - ug.FirstIndexOf[m-1]
    return []int{m - 1, n}
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
}

func main() {
	xSize := 10
	ySize := 500
	basicPointInPolygon := false
	//log.Fatal("Server with unidistant grid not implemented.")

	var ug1D []bool
	var ug UniformGrid

	ug.XSize = xSize
	ug.YSize = ySize
	ug.BigN = xSize * ySize
	ug.A = 4.0 * math.Pi / float64(ug.BigN)
	ug.D = math.Sqrt(ug.A)
	ug.MTheta = math.Round(math.Pi / ug.D)
	ug.DTheta = math.Pi / ug.MTheta
	ug.DPhi = ug.A / ug.DTheta

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
			fmt.Printf("point\n")
			fmt.Printf("CoordToGrid\n")
			var from = ug.CoordToGrid(startLng, startLat)
			fmt.Printf("%v\n",from)
			fmt.Printf("GridToCoord\n")
			var gridNodeCoord = ug.GridToCoord(from)
			fmt.Printf("%v\n",gridNodeCoord)
			var td TestData
			td.Point = gridNodeCoord
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
			fmt.Printf("%v,%v\n",x,y)
			fmt.Printf("GridPoint\n")
			fmt.Printf("CoordToGrid\n")
			point := ug.GridToCoord([]int{int(x),int(y)})
			fmt.Printf("%v\n",point)
			var td TestData
			td.Point = point
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
