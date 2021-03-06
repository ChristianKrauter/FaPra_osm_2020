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
		for _, x := range nodeIDs {
			poly.Points = append(poly.Points, []float64{(*nodeMap)[x][0], (*nodeMap)[x][1]})
		}
		delete(*coastlineMap, key)
		key = nodeIDs[len(nodeIDs)-1]
		for {
			if val, ok := (*coastlineMap)[key]; ok {
				for i, x := range val {
					if i != 0 {
						poly.Points = append(poly.Points, []float64{(*nodeMap)[x][0], (*nodeMap)[x][1]})
					}
				}
				delete(*coastlineMap, key)
				key = val[len(val)-1]
			} else {
				break
			}
		}

		if !basicPointInPolygon {
			poly.EoWNext = make([]int, len(poly.Points))
			poly.LngTNext = make([]float64, len(poly.Points))
			poly.BtoX = make([]int, len(poly.Points))

			for i, x := range poly.Points {
				var nortPole = []float64{0.0, 90.0}
				var aT = x[0]
				var bT = poly.Points[(i+1)%len(poly.Points)][0]
				if aT == bT {
					x[0] -= 0.000000001
					nortPole = []float64{0.1, 89.9}
					aT = transformLon(nortPole, x)
					bT = transformLon(nortPole, poly.Points[(i+1)%len(poly.Points)])
				}

				poly.EoWNext[i] = eastOrWest(aT, bT)
				poly.LngTNext[i] = transformLon(x, poly.Points[(i+1)%len(poly.Points)])
				poly.BtoX[i] = eastOrWest((poly.LngTNext)[i], transformLon(x, nortPole))
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
	var mpg = make([][][]float64, len(*polygons))
	for i, j := range *polygons {
		mpg[i] = j.Points
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

	if noBoundingTree {
		fmt.Printf("\nwithout using a bounding tree structure")
	}
	if lessMemory {
		fmt.Printf("\noptimized for unpruned pbf files")
	}
	fmt.Printf("\n\n")

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
	var allBoundingBoxes = make([]map[string]float64, len(polygons))
	if noBoundingTree {
		logging["filename"] += "_nbt"

		start := time.Now()
		for i, j := range polygons {
			allBoundingBoxes[i] = createBoundingBox(&j.Points)
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
