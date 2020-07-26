package dataprocessing

import (
	"encoding/json"
	//"flag"
	"fmt"
	"github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getSomeKey(m *map[int64][]int64) int64 {
	for k := range *m {
		return k
	}
	return 0
}

// After https://github.com/paulmach/orb
func rayCast(point, s, e []float64) (bool, bool) {
	if s[0] > e[0] {
		s, e = e, s
	}

	if point[0] == s[0] {
		if point[1] == s[1] {
			// point == start
			return false, true
		} else if s[0] == e[0] {
			// vertical segment (s -> e)
			// return true if within the line, check to see if start or end is greater.
			if s[1] > e[1] && s[1] >= point[1] && point[1] >= e[1] {
				return false, true
			}

			if e[1] > s[1] && e[1] >= point[1] && point[1] >= s[1] {
				return false, true
			}
		}

		// Move the y coordinate to deal with degenerate case
		point[0] = math.Nextafter(point[0], math.Inf(1))
	} else if point[0] == e[0] {
		if point[1] == e[1] {
			// matching the end point
			return false, true
		}

		point[0] = math.Nextafter(point[0], math.Inf(1))
	}

	if point[0] < s[0] || point[0] > e[0] {
		return false, false
	}

	if s[1] > e[1] {
		if point[1] > s[1] {
			return false, false
		} else if point[1] < e[1] {
			return true, false
		}
	} else {
		if point[1] > e[1] {
			return false, false
		} else if point[1] < s[1] {
			return true, false
		}
	}

	rs := (point[1] - s[1]) / (point[0] - s[0])
	ds := (e[1] - s[1]) / (e[0] - s[0])

	if rs == ds {
		return false, true
	}

	return rs <= ds, false
}

// After https://github.com/paulmach/orb
func polygonContains(polygon *[][]float64, point []float64) bool {
	b, on := rayCast(point, (*polygon)[0], (*polygon)[len(*polygon)-1])
	if on {
		return true
	}

	for i := 0; i < len(*polygon)-1; i++ {
		inter, on := rayCast(point, (*polygon)[i], (*polygon)[i+1])
		if on {
			return true
		}
		if inter {
			b = !b
		}
	}
	return b
}

func rayCastSphere(point, s, e []float64) (bool,bool) {
	if s[0] > e[0] {
		s, e = e, s
	}
	return true,true
}

func polygonContainsSphere(polygon *[][]float64, point []float64) bool{
	b, on := rayCastSphere(point, (*polygon)[0], (*polygon)[len(*polygon)-1])
	if on {
		return true
	}

	for i := 0; i < len(*polygon)-1; i++ {
		inter, on := rayCastSphere(point, (*polygon)[i], (*polygon)[i+1])
		if on {
			return true
		}
		if inter {
			b = !b
		}
	}
	return b
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

type boundingTree struct {
	boundingBox map[string]float64
	id          int
	children    []boundingTree
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

func countBoundingTree(bt boundingTree) int {
	count := 1
	for _, j := range bt.children {
		count = count + countBoundingTree(j)
	}
	return count
}

func isLand(tree *boundingTree, point []float64, allCoastlines *[][][]float64) bool {
	land := false
	if boundingContains(&tree.boundingBox, point) {
		for _, child := range (*tree).children {
			land = isLand(&child, point, allCoastlines)
			if land {
				return land
			}
		}
		if (*tree).id >= 0 {
			land = polygonContains(&(*allCoastlines)[(*tree).id], point)
		}
	}
	return land
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

func readFileLessMemory(pbfFileName string, coastlineMap *map[int64][]int64, nodeMap *map[int64][]float64) string {
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
			case *osmpbf.Way:
				for _, value := range v.Tags {
					if value == "coastline" {
						(*coastlineMap)[v.NodeIDs[0]] = v.NodeIDs
						for _, id := range v.NodeIDs {
							(*nodeMap)[id] = []float64{}
						}
					}
				}
			default:
				continue
			}
		}
	}

	// Read nodes
	f, err = os.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d = osmpbf.NewDecoder(f)
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
				if _, ok := (*nodeMap)[v.ID]; ok {
					(*nodeMap)[v.ID] = []float64{v.Lon, v.Lat}
				}
			default:
				continue
			}
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Read file in                     : %s\n", elapsed)
	return elapsed.String()
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
				for _, x := range val {
					coastline = append(coastline, []float64{(*nodeMap)[x][0], (*nodeMap)[x][1]})
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
	return elapsed.String()
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

func createMeshgrid(xSize int, ySize int, boundingTreeRoot *boundingTree, allCoastlines *[][][]float64, testGeoJSON *[][]float64, meshgrid *[][]bool, createTestGeoJSON bool) string {
	start := time.Now()
	var xStepSize = float64(360 / xSize)
	var yStepSize = float64(360 / ySize)

	var wg sync.WaitGroup
	for x := 0.0; x < 360; x += xStepSize {
		for y := 0.0; y < 360; y += yStepSize {
			wg.Add(1)
			go func(x float64, y float64) {
				defer wg.Done()
				var xs = x - 180
				var ys = (y / 2) - 90
				if isLand(boundingTreeRoot, []float64{xs, ys}, allCoastlines) {
					(*meshgrid)[int(x/xStepSize)][int(y/yStepSize)] = true
					if createTestGeoJSON {
						*testGeoJSON = append(*testGeoJSON, []float64{xs, ys})
					}
				}
			}(x, y)
		}
	}

	wg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created Meshrid in               : %s\n", elapsed)
	return elapsed.String()
}

func storeMeshgrid(meshgrid *[][]bool, filename string) string {
	start := time.Now()
	var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(meshgrid)
	check(err1)
	f, err2 := os.Create(filename)
	check(err2)
	_, err3 := f.Write(meshgridBytes)
	check(err3)
	f.Sync()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Stored Meshrid to disc in        : %s\n", elapsed)
	return elapsed.String()
}

func storeTestGeoJSON(testGeoJSON *[][]float64, filename string) string {
	start := time.Now()
	fmt.Printf("Points in test geojson: %d\n", len(*testGeoJSON))
	var rawJSON []byte

	g := geojson.NewMultiPointGeometry(*testGeoJSON...)
	rawJSON, err4 := g.MarshalJSON()
	check(err4)

	f, err5 := os.Create(filename)
	check(err5)

	_, err6 := f.Write(rawJSON)
	check(err6)
	f.Sync()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Stored test geojson in           : %s\n", elapsed)
	return elapsed.String()
}

func createAndStoreCoastlineGeoJSON(allCoastlines *[][][]float64, filename string) string {
	start := time.Now()
	for _, i := range *allCoastlines {
		var polygon [][][]float64
		polygon = append(polygon, i)
		g := geojson.NewPolygonGeometry(polygon)
		rawJSON, err := g.MarshalJSON()
		f, err := os.Create(filename)
		check(err)

		_, err1 := f.Write(rawJSON)
		check(err1)
		f.Sync()
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created and stored coastlines in : %s\n", elapsed)
	return elapsed.String()
}

// Start processing a pbf file to create a meshgrid
func Start(pbfFileName string, xSize, ySize int, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree bool) map[string]string {
	fmt.Printf("\nStarting processing of %s\n\n", pbfFileName)
	logging := make(map[string]string)
	pbfFileName = fmt.Sprintf("data/%s", pbfFileName)
	start := time.Now()

	// Read the pbf file
	var coastlineMap = make(map[int64][]int64)
	var nodeMap = make(map[int64][]float64)
	if lessMemory {
		var readTime = readFileLessMemory(pbfFileName, &coastlineMap, &nodeMap)
		logging["time_read"] = string(readTime)
		logging["filename"] += "_lm"
	} else {
		var readTime = readFile(pbfFileName, &coastlineMap, &nodeMap)
		logging["time_read"] = string(readTime)
	}

	// Create coastline polygons
	var allCoastlines [][][]float64
	var polygonTime = createPolygons(&allCoastlines, &coastlineMap, &nodeMap)
	logging["time_poly"] = string(polygonTime)

	// Create bounding boxes
	var boundingTreeRoot boundingTree
	if noBoundingTree {
		log.Fatal("Not implemented")
		logging["filename"] += "_nbt"
	} else {
		var boundingTreeTime = createBoundingTree(&boundingTreeRoot, &allCoastlines)
		logging["time_boundingTree"] = string(boundingTreeTime)
	}

	// Create and store meshgrid
	var testGeoJSON [][]float64
	meshgrid := make([][]bool, xSize) // inits with false!
	for i := range meshgrid {
		meshgrid[i] = make([]bool, ySize)
	}

	var meshgridTime = createMeshgrid(xSize, ySize, &boundingTreeRoot, &allCoastlines, &testGeoJSON, &meshgrid, createTestGeoJSON)
	logging["time_meshgrid"] = string(meshgridTime)

	var meshgridStoreTime = storeMeshgrid(&meshgrid, fmt.Sprintf("data/output/meshgrid_%d_%d.json", xSize, ySize))
	logging["time_meshgrid_store"] = string(meshgridStoreTime)

	// Create and safe additional files
	if createTestGeoJSON {
		var testGeoJSONTime = storeTestGeoJSON(&testGeoJSON, fmt.Sprintf("data/output/test_for_%d_%d.geojson", xSize, ySize))
		logging["time_testGeoJSON"] = string(testGeoJSONTime)
	}

	if createCoastlineGeoJSON {
		var coastlineGeoJSONTime = createAndStoreCoastlineGeoJSON(&allCoastlines, fmt.Sprintf("data/output/coastlines_for_%d_%d.geojson", xSize, ySize))
		logging["time_coastlineGeoJSON"] = string(coastlineGeoJSONTime)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("\nProgram finished after           : %s\n", elapsed)
	logging["time_total"] = elapsed.String()
	return logging
}
