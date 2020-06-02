package main

import (
	"fmt"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"time"
	"encoding/json"
	"github.com/paulmach/go.geojson"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getSomeKey(m map[int64]*osmpbf.Way) int64 {
	for k := range m {
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

// https://github.com/paulmach/orb
func polygonContains(polygon [][]float64, point []float64) bool {
	b, on := rayCast(point, polygon[0], polygon[len(polygon)-1])
	if on {
		return true
	}

	for i := 0; i < len(polygon)-1; i++ {
		inter, on := rayCast(point, polygon[i], polygon[i+1])
		if on {
			return true
		}
		if inter {
			b = !b
		}
	}
	return b
}

func boundingContains(bounding map[string]float64, point []float64) bool{
	if (bounding["minX"] <= point[0] && point[0] <= bounding["maxX"]) {
		if (bounding["minY"] <= point[1] && point[1] <= bounding["maxY"]) {
			return true
		}
	}
	return false
}

func createBoundingBox(polygon [][]float64) map[string]float64 {
	minX := math.Inf(1)
	maxX := math.Inf(-1)
	minY := math.Inf(1)
	maxY := math.Inf(-1)
	for _, coord := range polygon {
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


	//boundingBox := map[string]float64{"minX":minX, "maxX":maxX, "minY":minY, "maxY":maxY}
	/*var boundingBox map[string]float64
	boundingBox = make(map[string]float64)
	boundingBox["minX"]=minX
	boundingBox["maxX"]=maxX
	boundingBox["minY"]=minY
	boundingBox["maxY"]=maxY

	return boundingBox
	*/
	return map[string]float64{"minX":minX, "maxX":maxX, "minY":minY, "maxY":maxY}
}

func main() {
	start := time.Now()

	var pbfFileName = "../data/antarctica-latest.osm.pbf"
	//var pbfFileName = "../data/planet-coastlines.pbf"

	//fs, err := os.Stat(pbfFileName)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Printf("\nStarting processing of %s (%d KB)\n\n", pbfFileName, fs.Size()/1000)
	fmt.Printf("\nStarting processing of %s\n\n", pbfFileName)

	f, err := os.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	//var coastlineMap map[int64]*osmpbf.Way
	//var nodeMap map[int64]*osmpbf.Node
	var coastlineMap = make(map[int64]*osmpbf.Way)
	var nodeMap = make(map[int64]*osmpbf.Node)

	// start decoding with several goroutines, it is faster
	//err = d.Start(runtime.GOMAXPROCS(-1))
	//print(runtime.NumCPU())
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
				nodeMap[v.ID] = v
			case *osmpbf.Way:
				for _, value := range v.Tags {
					if value == "coastline" {
						coastlineMap[v.NodeIDs[0]] = v
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
	fmt.Printf("Done Reading file after: %s\n", elapsed)

	var allCoastlines [][][]float64
	var allBoundingBoxes []map[string]float64
	var coastline [][]float64

	for len(coastlineMap) > 0 {
		var key = getSomeKey(coastlineMap)
		var nodeIDs = coastlineMap[key].NodeIDs
		coastline = nil
		for _, x := range nodeIDs {
			//var coord []float64
			//coord = append(coord, nodeMap[x].Lon)
			//coord = append(coord, nodeMap[x].Lat)
			coastline = append(coastline, []float64{nodeMap[x].Lon, nodeMap[x].Lat})
		}
		delete(coastlineMap, key)
		key = nodeIDs[len(nodeIDs)-1]
		for {
			if val, ok := coastlineMap[key]; ok {
				var nodeIDs = val.NodeIDs
				for _, x := range nodeIDs {
					//var coord []float64
					//coord = append(coord, nodeMap[x].Lon)
					//coord = append(coord, nodeMap[x].Lat)
					coastline = append(coastline, []float64{nodeMap[x].Lon, nodeMap[x].Lat})
				}
				delete(coastlineMap, key)
				key = nodeIDs[len(nodeIDs)-1]
			} else {
				break
			}

		}
		allCoastlines = append(allCoastlines, coastline)
	}
	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Made all polygons: %s\n", elapsed)

	for _, i := range allCoastlines {
		//boundingBox := createBoundingBox(i)
		allBoundingBoxes = append(allBoundingBoxes, createBoundingBox(i))
	}
	fmt.Printf("\n")

	fmt.Printf("Creating meshgrid:\n")
	var testGeoJSON [][]float64
	var meshgrid [360][360]bool
	for x := 0.0; x < 360; x++ {
		for y := 0.0; y < 360; y++ {
			var xs = x - 180
			var ys = (y/2) -90
			isWater := true
			for i, j := range allBoundingBoxes {
				if boundingContains(j, []float64{xs, ys}) {
					if polygonContains(allCoastlines[i], []float64{xs, ys}) {

						// For test geojson
						//var coord []float64
						//coord = append(coord, x-180)
						//coord = append(coord, (y/2)-90)
						testGeoJSON = append(testGeoJSON, []float64{xs, ys})

						//testGeoJSON = append(testGeoJSON, coord)
						// End for test geojson

						isWater = false
						break
					}
				}
			}
			meshgrid[int(x)][int(y)] = isWater
		}
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Meshrid finished after: %s\n", elapsed)

	// Save meshgrid to disk
	var meshgridBytes []byte
	meshgridBytes, err1 := json.Marshal(meshgrid)
	check(err1)
	var filename = fmt.Sprintf("tmp/meshgrid.json")
	f, err2 := os.Create(filename)
	check(err2)
	_, err3 := f.Write(meshgridBytes)
	check(err3)
	f.Sync()

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Saved Meshrid to disc after: %s\n", elapsed)

	//fmt.Printf("Points in test geojson: %d\n", len(testGeoJSON))
	fmt.Printf("creating test geojson\n")
	var rawJson []byte
	g := geojson.NewMultiPointGeometry(testGeoJSON...)
	rawJson, err4 := g.MarshalJSON()
	check(err4)
	var testgeojsonFilename = fmt.Sprintf("tmp/datatestgeojson.geojson")
	f, err5 := os.Create(testgeojsonFilename)
	check(err5)
	_, err6 := f.Write(rawJson)
	check(err6)
	f.Sync()

	// Create coastline geojson file
	/*var rawJson []byte
	for j, i := range allCoastlines {
		var polygon [][][]float64
		polygon = append(polygon, i)
		g := geojson.NewPolygonGeometry(polygon)
		rawJson, err = g.MarshalJSON()
		var filename = fmt.Sprintf("tmp/data%d.geojson", j)
		f, err := os.Create(filename)
		check(err)
		_, err1 := f.Write(rawJson)
		check(err1)
		f.Sync()
	}
	fmt.Printf("Created coastline geojson in   : %s\n\n", elapsed)
	*/

	//fmt.Printf("\nNumber of:\n")
	//fmt.Printf("Nodes                 : %d\n", nc)
	//fmt.Printf("Ways                  : %d\n", wc)
	//fmt.Printf("Relations             : %d\n", rc)
	//fmt.Printf("Coastline Polygons    : %d\n\n", len(allPolygonsID))

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Program finished after: %s\n", elapsed)
}
