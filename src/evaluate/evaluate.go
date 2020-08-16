package evaluate

import (
	"../algorithms"
	"../dataprocessing"
	//"../server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// WayFinding is evaluated
func WayFinding(xSize, ySize, algorithm int, basicPointInPolygon bool) {
	var filename string
	var meshgrid []bool
	var meshgrid2d [][]bool

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v.json", xSize, ySize)
	}

	meshgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &meshgrid2d)

	//var trueCount int
	for i := 0; i < len(meshgrid2d[0]); i++ {
		for j := 0; j < len(meshgrid2d); j++ {
			meshgrid = append(meshgrid, meshgrid2d[j][i])
			//if !meshgrid2d[j][i] {
			//	trueCount += 1
			//}
		}
	}
	//fmt.Printf("\n%v\n", trueCount)
	//var testCount int
	// TODO: More efficient average
	var sum time.Duration
	var count int64
	for i := 0; i < len(meshgrid2d); i++ {
		for j := 0; j < len(meshgrid2d[i]); j++ {
			for k := 0; k < len(meshgrid2d); k++ {
				for l := 0; l < len(meshgrid2d[k]); l++ {
					if !meshgrid2d[i][j] && !meshgrid2d[k][l] {
						//testCount += 1
						// fmt.Printf("%v:%v - %v:%v", j, i, l, k)
						var start = time.Now()
						var _, _ = algorithms.DijkstraAllNodes(int64(i), int64(j), int64(k), int64(l), int64(xSize), int64(ySize), &meshgrid)
						t := time.Now()
						var elapsed = t.Sub(start)
						sum += elapsed
						count += 1
					}
				}
			}
		}
	}
	fmt.Printf("%v, %v, %v", sum, count, float64(sum)/float64(count))
	//fmt.Printf("%v\n", testCount)
}

// DataProcessing is evaluated
func DataProcessing(pbfFileName, note string, xSize, ySize int, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon bool) {
	var logging = dataprocessing.Start(pbfFileName, xSize, ySize, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logging["pbfFileName"] = pbfFileName
	logging["note"] = note
	logging["basicGrid"] = strconv.FormatBool(basicGrid)
	logging["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	logging["xSize"] = strconv.Itoa(xSize)
	logging["ySize"] = strconv.Itoa(ySize)
	logging["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logging["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)

	jsonString, _ := json.Marshal(logging)
	var filename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := logging["filename"]; ok {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s%s_%s.json", strings.Split(logging["pbfFileName"], ".")[0], logging["xSize"], logging["ySize"], val, timestamp)
	} else {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s_%s.json", strings.Split(logging["pbfFileName"], ".")[0], logging["xSize"], logging["ySize"], timestamp)
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(string(jsonString))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
