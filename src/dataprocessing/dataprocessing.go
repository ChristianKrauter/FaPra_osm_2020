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

func getSomeKey(m *map[int64][]int64) int64 {
	for k := range *m {
		return k
	}
	return 0
}

func createPolygons(polygons *Polygons, coastlineMap *map[int64][]int64, nodeMap *map[int64][]float64, basicPointInPolygon bool) string {
	start := time.Now()
	var poly Polygon
	for len(*coastlineMap) > 0 {
		var key = getSomeKey(coastlineMap)
		var nodeIDs = (*coastlineMap)[key]
		poly = Polygon{}
		for i, x := range nodeIDs {
			var point = []float64{(*nodeMap)[x][0], (*nodeMap)[x][1]}
			poly.Points = append(poly.Points, point)
			if !basicPointInPolygon {
				poly.LngTNorth = append(poly.LngTNorth, transformLon(point, []float64{0.0, 90.0}))
				var nIDX = nodeIDs[(i+1)%len(nodeIDs)]
				var nPoint = []float64{(*nodeMap)[nIDX][0], (*nodeMap)[nIDX][1]}
				poly.LngTNext = append(poly.LngTNext, transformLon(point, nPoint))
			}
		}
		delete(*coastlineMap, key)
		key = nodeIDs[len(nodeIDs)-1]
		for {
			if val, ok := (*coastlineMap)[key]; ok {
				for i, x := range val {
					if i != 0 {
						var point = []float64{(*nodeMap)[x][0], (*nodeMap)[x][1]}
						poly.Points = append(poly.Points, point)
						if !basicPointInPolygon {
							poly.LngTNorth = append(poly.LngTNorth, transformLon(point, []float64{0.0, 90.0}))
							var nIDX = nodeIDs[(i+1)%len(nodeIDs)]
							var nPoint = []float64{(*nodeMap)[nIDX][0], (*nodeMap)[nIDX][1]}
							poly.LngTNext = append(poly.LngTNext, transformLon(point, nPoint))
						}
					}
				}
				delete(*coastlineMap, key)
				key = val[len(val)-1]
			} else {
				break
			}
		}
		*polygons = append(*polygons, poly)
	}

	sort.Sort(Polygons(*polygons))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Made all polygons in             : %s\n", elapsed)
	return elapsed.String()
}

func createAndStoreCoastlineGeoJSON(polygons *Polygons, filename string) string {
	start := time.Now()
	var mpg [][][]float64
	for _, i := range *polygons {
		mpg = append(mpg, i.Points)
	}
	var fc = geojson.NewMultiPolygonGeometry(mpg)
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
	var polygons Polygons
	var polygonTime = createPolygons(&polygons, &coastlineMap, &nodeMap, basicPointInPolygon)
	logging["time_poly"] = string(polygonTime)

	// Create bounding boxes
	var boundingTreeRoot boundingTree
	var allBoundingBoxes []map[string]float64
	if noBoundingTree {
		logging["filename"] += "_nbt"

		start := time.Now()
		for _, i := range polygons {
			allBoundingBoxes = append(allBoundingBoxes, createBoundingBox(&i.Points))
		}
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("Created bounding boxes           : %s\n", elapsed)
		logging["time_boundingBoxes"] = elapsed.String()

	} else {
		var boundingTreeTime = createBoundingTree(&boundingTreeRoot, &polygons)
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
			meshgridTime = createMeshgridNBT(xSize, ySize, &allBoundingBoxes, &polygons, &bg, basicPointInPolygon)
		} else {
			meshgridTime = createMeshgrid(xSize, ySize, &boundingTreeRoot, &polygons, &bg, basicPointInPolygon)
		}
		logging["time_meshgrid"] = string(meshgridTime)

		var meshgridStoreTime = storeMeshgrid(&bg, fmt.Sprintf("data/output/meshgrid_%d_%d%s.json", xSize, ySize, filenameAdditions))
		logging["time_meshgrid_store"] = string(meshgridStoreTime)

	} else {
		var uniformGrid = grids.NewUG(xSize, ySize)

		var uniformGridTime string
		if noBoundingTree {
			uniformGridTime = createUniformGridNBT(xSize, ySize, &allBoundingBoxes, &polygons, uniformGrid, basicPointInPolygon)
		} else {
			uniformGridTime = createUniformGrid(xSize, ySize, &boundingTreeRoot, &polygons, uniformGrid, basicPointInPolygon)
		}
		logging["time_meshgrid"] = string(uniformGridTime)

		var uniformGridStoreTime = storeUniformGrid(uniformGrid, fmt.Sprintf("data/output/uniformGrid_%v_%v%s.json", xSize, ySize, filenameAdditions))
		logging["time_meshgrid_store"] = string(uniformGridStoreTime)
	}

	// Create and safe coastline
	if createCoastlineGeoJSON {
		var coastlineGeoJSONTime = createAndStoreCoastlineGeoJSON(&polygons, fmt.Sprintf("data/output/coastlines_for_%d_%d.geojson", xSize, ySize))
		logging["time_coastlineGeoJSON"] = string(coastlineGeoJSONTime)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("\nProgram finished after           : %s\n", elapsed)
	logging["time_total"] = elapsed.String()
	return logging
}
