package dataprocessing

import (
	"fmt"
	"github.com/paulmach/go.geojson"
	"log"
	"os"
	"sort"
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

func getSomeKey(m *map[int64][]int64) int64 {
	for k := range *m {
		return k
	}
	return 0
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

	var uniformGrid SphereGrid
	var uniformGridTime = createUniformGrid(xSize, ySize, &uniformGrid, &boundingTreeRoot, &allCoastlines)
	fmt.Printf("%v\n", uniformGridTime)

	/*
		var meshgridTime = createMeshgrid(xSize, ySize, &boundingTreeRoot, &allCoastlines, &testGeoJSON, &meshgrid, createTestGeoJSON)
		logging["time_meshgrid"] = string(meshgridTime)

		var meshgridStoreTime = storeMeshgrid(&meshgrid, fmt.Sprintf("data/output/meshgrid_%d_%d.json", xSize, ySize))
		logging["time_meshgrid_store"] = string(meshgridStoreTime)
	*/
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
