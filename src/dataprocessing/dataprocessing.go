package dataprocessing

import (
	"../grids"
	"fmt"
	"github.com/paulmach/go.geojson"
	"log"
	"os"
	"sort"
	"time"
)

// arrayOfArrays sorting by length
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

func createAndStoreCoastlineGeoJSON(allCoastlines *[][][]float64, filename string) string {
	start := time.Now()
	var polygons [][][]float64
	for _, i := range *allCoastlines {
		polygons = append(polygons, i)
	}
	var fc = geojson.NewMultiPolygonGeometry(polygons)
	rawJSON, err := fc.MarshalJSON()
	f, err := os.Create(filename)
	check(err)

	_, err1 := f.Write(rawJSON)
	check(err1)
	f.Sync()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Created and stored coastlines in : %s\n", elapsed)
	return elapsed.String()
}

// Start processing a pbf file to create a meshgrid
func Start(pbfFileName string, xSize, ySize int, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon bool) map[string]string {
	if basicGrid && (xSize < 360 || ySize < 360) {
		log.Fatal("\n\nBasic grid not possibly for xSize or ySize under 360.")
	}

	fmt.Printf("\nStarting processing of %s\n\n", pbfFileName)
	logging := make(map[string]string)
	pbfFileName = fmt.Sprintf("data/%s", pbfFileName)
	start := time.Now()

	// Read the pbf file
	var coastlineMap = make(map[int64][]int64)
	var nodeMap = make(map[int64][]float64)
	if lessMemory {
		var readTime = ReadFileLessMemory(pbfFileName, &coastlineMap, &nodeMap)
		logging["time_read"] = string(readTime)
		logging["filename"] += "_lm"
	} else {
		var readTime = ReadFile(pbfFileName, &coastlineMap, &nodeMap)
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
	var filenameAdditions = ""
	if basicPointInPolygon {
		logging["filename"] += "_bpip"
		filenameAdditions += "_bpip"
	}

	if basicGrid {
		logging["filename"] += "_bg"

		bg := make([][]bool, xSize) // inits with false!
		for i := range bg {
			bg[i] = make([]bool, ySize)
		}

		var meshgridTime string
		if noBoundingTree {
			meshgridTime = createMeshgridNBT(xSize, ySize, &allBoundingBoxes, &allCoastlines, &bg, basicPointInPolygon)
		} else {
			meshgridTime = createMeshgrid(xSize, ySize, &boundingTreeRoot, &allCoastlines, &bg, basicPointInPolygon)
		}
		logging["time_meshgrid"] = string(meshgridTime)

		var meshgridStoreTime = storeMeshgrid(&bg, fmt.Sprintf("data/output/meshgrid_%d_%d%s.json", xSize, ySize, filenameAdditions))
		logging["time_meshgrid_store"] = string(meshgridStoreTime)

	} else {
		var uniformGrid grids.UniformGrid

		var uniformGridTime string
		if noBoundingTree {
			uniformGridTime = createUniformGridNBT(xSize, ySize, &allBoundingBoxes, &allCoastlines, &uniformGrid, basicPointInPolygon)
		} else {
			uniformGridTime = createUniformGrid(xSize, ySize, &boundingTreeRoot, &allCoastlines, &uniformGrid, basicPointInPolygon)
		}
		logging["time_meshgrid"] = string(uniformGridTime)

		var uniformGridStoreTime = storeUniformGrid(&uniformGrid, fmt.Sprintf("data/output/uniformGrid_%v_%v%s.json", xSize, ySize, filenameAdditions))
		logging["time_meshgrid_store"] = string(uniformGridStoreTime)
	}

	// Create and safe coastline
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
