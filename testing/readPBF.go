package main

import (
	"fmt"
	"github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"os"
	"runtime"
	"time"
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

func main() {
	start := time.Now()

	var pbfFileName = "../data/antarctica-latest.osm.pbf"
	// pbfFileName = "../data/planet-coastlines.pbf"

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

	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				nodeMap[v.ID] = v
				// Process Node v.
				nc++
			case *osmpbf.Way:
				// Process Way v.
				for _, value := range v.Tags {
					if value == "coastline" {
						coastlineMap[v.NodeIDs[0]] = v
						wc++
					}
				}
			case *osmpbf.Relation:
				// Process Relation v.
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Read file in          : %s\n", elapsed)

	var allPolygonsID [][]int64
	var allPolygonsCoord [][][]float64

	var coastlineID []int64
	var coastlineCoord [][]float64

	for len(coastlineMap) > 0 {
		var key = get_some_key(coastlineMap)
		var nodeids = coastlineMap[key].NodeIDs
		coastlineID = nodeids[:len(nodeids)-1]
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
				coastlineID = append(coastlineID, nodeids[:len(nodeids)-1]...)
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

		allPolygonsID = append(allPolygonsID, coastlineID)
		allPolygonsCoord = append(allPolygonsCoord, coastlineCoord)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Created polygons in   : %s\n\n", elapsed)

	var rawJson []byte
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

	fmt.Printf("Number of:\n")
	fmt.Printf("Nodes                 : %d\n", nc)
	fmt.Printf("Ways                  : %d\n", wc)
	fmt.Printf("Relations             : %d\n", rc)
	fmt.Printf("Coastline Polygons    : %d\n\n", len(allPolygonsID))
	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Printf("Program finished after: %s\n", elapsed)
}
