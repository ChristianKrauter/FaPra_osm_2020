package dataprocessing

import (
	"../algorithms"
	"fmt"
	"github.com/paulmach/go.geojson"
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
func Start(pbfFileName string, xSize, ySize int, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon bool) map[string]string {
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
	var allBoundingBoxes []map[string]float64
	if noBoundingTree {
		logging["filename"] += "_nbt"

		start := time.Now()
		for _, i := range allCoastlines {
			allBoundingBoxes = append(allBoundingBoxes, createBoundingBox(&i))
		}
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("Created bounding boxes           : %s\n", elapsed)
		logging["time_boundingBoxes"] = elapsed.String()

	} else {
		var boundingTreeTime = createBoundingTree(&boundingTreeRoot, &allCoastlines)
		logging["time_boundingTree"] = string(boundingTreeTime)
	}

	// Create and store meshgrid
	var testGeoJSON [][]float64
	var filenameAdditions = ""
	if basicPointInPolygon {
		logging["filename"] += "_bpip"
		filenameAdditions += "_bpip"
	}

	if basicGrid {
		logging["filename"] += "_bg"

		meshgrid := make([][]bool, xSize) // inits with false!
		for i := range meshgrid {
			meshgrid[i] = make([]bool, ySize)
		}

		var meshgridTime string
		if noBoundingTree {
			meshgridTime = createMeshgridNBT(xSize, ySize, &allBoundingBoxes, &allCoastlines, &testGeoJSON, &meshgrid, createTestGeoJSON, basicPointInPolygon)
		} else {
			meshgridTime = createMeshgrid(xSize, ySize, &boundingTreeRoot, &allCoastlines, &testGeoJSON, &meshgrid, createTestGeoJSON, basicPointInPolygon)
		}
		logging["time_meshgrid"] = string(meshgridTime)

		var meshgridStoreTime = storeMeshgrid(&meshgrid, fmt.Sprintf("data/output/meshgrid_%d_%d%s.json", xSize, ySize, filenameAdditions))
		logging["time_meshgrid_store"] = string(meshgridStoreTime)

	} else {
		var uniformGrid algorithms.UniformGrid

		var uniformGridTime = createUniformGrid(xSize, ySize, &boundingTreeRoot, &allCoastlines, &testGeoJSON, &uniformGrid, createTestGeoJSON, basicPointInPolygon)
		logging["time_meshgrid"] = string(uniformGridTime)

		var uniformGridStoreTime = storeUniformGrid(&uniformGrid, fmt.Sprintf("data/output/uniformGrid_%v_%v%s.json", xSize, ySize, filenameAdditions))
		logging["time_meshgrid_store"] = string(uniformGridStoreTime)
	}

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
