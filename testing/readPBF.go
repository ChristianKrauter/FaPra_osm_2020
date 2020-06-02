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

func get_some_key(m map[int64]*osmpbf.Way) int64 {
	for k := range m {
		return k
	}
	return 0
}

// After https://github.com/paulmach/orb
func rayCast(p, s, e []float64) (bool, bool) {
	if s[0] > e[0] {
		s, e = e, s
	}

	if p[0] == s[0] {
		if p[1] == s[1] {
			// p == start
			return false, true
		} else if s[0] == e[0] {
			// vertical segment (s -> e)
			// return true if within the line, check to see if start or end is greater.
			if s[1] > e[1] && s[1] >= p[1] && p[1] >= e[1] {
				return false, true
			}

			if e[1] > s[1] && e[1] >= p[1] && p[1] >= s[1] {
				return false, true
			}
		}

		// Move the y coordinate to deal with degenerate case
		p[0] = math.Nextafter(p[0], math.Inf(1))
	} else if p[0] == e[0] {
		if p[1] == e[1] {
			// matching the end point
			return false, true
		}

		p[0] = math.Nextafter(p[0], math.Inf(1))
	}

	if p[0] < s[0] || p[0] > e[0] {
		return false, false
	}

	if s[1] > e[1] {
		if p[1] > s[1] {
			return false, false
		} else if p[1] < e[1] {
			return true, false
		}
	} else {
		if p[1] > e[1] {
			return false, false
		} else if p[1] < s[1] {
			return true, false
		}
	}

	rs := (p[1] - s[1]) / (p[0] - s[0])
	ds := (e[1] - s[1]) / (e[0] - s[0])

	if rs == ds {
		return false, true
	}

	return rs <= ds, false
}

// https://github.com/paulmach/orb
func polygon_contains(polygon [][]float64, point []float64) bool {
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

func bounding_contains(bounding map[string]float64, point []float64) bool{
	if (bounding["minX"] <= point[0] && point[0] <= bounding["maxX"]) {
		if (bounding["minY"] <= point[1] && point[1] <= bounding["maxY"]) {
			return true
		}
	}
	return false
}

func create_bounding_Box(polygon [][]float64) map[string]float64 {
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

	var bounding_box map[string]float64
	bounding_box = make(map[string]float64)
	bounding_box["minX"]=minX
	bounding_box["maxX"]=maxX
	bounding_box["minY"]=minY
	bounding_box["maxY"]=maxY

	return bounding_box
}

func main() {
	start := time.Now()

	var pbfFileName = "../data/antarctica-latest.osm.pbf"
	//var pbfFileName = "../data/planet-coastlines.pbf"

	fs, err := os.Stat(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nStarting processing of %s (%d KB)\n\n", pbfFileName, fs.Size()/1000)

	f, err := os.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	var coastlineMap map[int64]*osmpbf.Way
	coastlineMap = make(map[int64]*osmpbf.Way)

	var nodeMap map[int64]*osmpbf.Node
	nodeMap = make(map[int64]*osmpbf.Node)

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
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Done Reading file after: %s\n", elapsed)

	//var allPolygonsID [][]int64
	var allPolygonsCoord [][][]float64

	var allPolygonsBounding []map[string]float64

	//var coastlineID []int64
	var coastlineCoord [][]float64

	for len(coastlineMap) > 0 {
		var key = get_some_key(coastlineMap)
		var nodeids = coastlineMap[key].NodeIDs
		//coastlineID = nodeids[:len(nodeids)-1]
		coastlineCoord = nil
		for _, x := range nodeids {
			var coord []float64
			coord = append(coord, nodeMap[x].Lon)
			coord = append(coord, nodeMap[x].Lat)
			coastlineCoord = append(coastlineCoord, coord)
		}
		delete(coastlineMap, key)
		key = nodeids[len(nodeids)-1]
		//fmt.Printf("1")
		for {
			//fmt.Printf("2")
			if val, ok := coastlineMap[key]; ok {
				var nodeids = val.NodeIDs
				//coastlineID = append(coastlineID, nodeids[:len(nodeids)-1]...)
				for _, x := range nodeids {
					var coord []float64
					coord = append(coord, nodeMap[x].Lon)
					coord = append(coord, nodeMap[x].Lat)
					coastlineCoord = append(coastlineCoord, coord)
				}
				delete(coastlineMap, key)
				key = nodeids[len(nodeids)-1]
			} else {
				break
			}

		}
		//allPolygonsID = append(allPolygonsID, coastlineID)
		allPolygonsCoord = append(allPolygonsCoord, coastlineCoord)
	}
	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Made all polygons: %s\n", elapsed)

	for _, i := range allPolygonsCoord {
		bounding_Box := create_bounding_Box(i)
		allPolygonsBounding = append(allPolygonsBounding, bounding_Box)
	}
	fmt.Printf("\n")

	//var allTestPolygon [][][]float64
	//var testPolygon [][]float64
	//var landCount int64

	fmt.Printf("Creating Meshgrid:\n")
	var Meshgrid [360][360]bool
	var MeshgridString = ""
	for x := 0.0; x < 360; x++ {
		for y := 0.0; y < 360; y++ {
			isWater := true
			for i, j := range allPolygonsBounding {
				//var printString = ""
				if bounding_contains(j, []float64{x-180, (y/2)-90}) {
					//printString = "in bounding"
					if polygon_contains(allPolygonsCoord[i], []float64{x-180, (y/2)-90}) {

						// For test polygon
						//var coord []float64
						//coord = append(coord, x-180)
						//coord = append(coord, y-180)
						// fmt.Printf("%f/%f", x, y)
						//testPolygon = append(testPolygon, coord)
						// End for test polygon

						isWater = false
						//landCount++
						//printString = "in poly"
					}
				}
				//if len(printString) > 0 {fmt.Printf("%s\n", printString)}
			}
			Meshgrid[int(x)][int(y)] = isWater
			if isWater {
				MeshgridString = MeshgridString + "w"
			} else {
				MeshgridString = MeshgridString + "o"
			}
		}
		MeshgridString = MeshgridString + "\n"
	}

	// Save meshgrid to disk
	var meshgrid_bytes []byte
	meshgrid_bytes, err := json.Marshal(m)
	var filename = fmt.SPrintf("tmp/meshgrid.json")
	f, err2 := os.Create(filename)
	check(err2)
	_, err3 := f.Write(meshgrid_bytes)
	check(err3)
	f.Sync()

	fmt.Printf("\nNot water: %d\n", landCount)
	/*
	fmt.Printf("creating test polygon\n")
	var rawJson []byte
	testPolygon = append(testPolygon, testPolygon[0])
	allTestPolygon = append(allTestPolygon, testPolygon)
	g := geojson.NewPolygonGeometry(allTestPolygon)
	rawJson, err = g.MarshalJSON()
	var filename = fmt.Sprintf("tmp/dataTestPolygon.geojson")
	f, err2 := os.Create(filename)
	check(err2)
	_, err3 := f.Write(rawJson)
	check(err3)
	f.Sync()
	*/

	// Create GeoJSON file
	/*var rawJson []byte
	for j, i := range allPolygonsCoord {
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
	fmt.Printf("Created polygons in   : %s\n\n", elapsed)
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
