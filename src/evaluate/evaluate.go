package evaluate

import (
	"../algorithms"
	"../dataprocessing"
	//"../server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ReadPBF is evaluated
func ReadPBF(pbfFileName, note string) {
	logging := make(map[string]string)
	logging["pbfFileName"] = pbfFileName
	pbfFileName = fmt.Sprintf("data/%s", pbfFileName)

	var coastlineMap = make(map[int64][]int64)
	var nodeMap = make(map[int64][]float64)
	var readTime = dataprocessing.ReadFile(pbfFileName, &coastlineMap, &nodeMap)

	coastlineMap = make(map[int64][]int64)
	nodeMap = make(map[int64][]float64)
	var readLessMemoryTime = dataprocessing.ReadFileLessMemory(pbfFileName, &coastlineMap, &nodeMap)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logging["note"] = note
	logging["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logging["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	logging["time_read"] = string(readTime)
	logging["time_readLessMemory"] = string(readLessMemoryTime)

	jsonString, _ := json.Marshal(logging)
	var filename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	filename = fmt.Sprintf("data/evaluation/rf_%s_%s.json", strings.Split(logging["pbfFileName"], ".")[0], timestamp)

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

// WayFindingBG is evaluated
func WayFindingBG(xSize, ySize, algorithm int, basicPointInPolygon bool, note string) {
	var filename string
	var meshgrid []bool
	var meshgrid2d [][]bool

	logging := make(map[string]string)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logging["note"] = note
	logging["basicGrid"] = "true"
	logging["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	logging["xSize"] = strconv.Itoa(xSize)
	logging["ySize"] = strconv.Itoa(ySize)
	logging["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logging["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)

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

	for i := 0; i < len(meshgrid2d[0]); i++ {
		for j := 0; j < len(meshgrid2d); j++ {
			meshgrid = append(meshgrid, meshgrid2d[j][i])
		}
	}

	// TODO: More efficient average
	var sum time.Duration
	var count int
	var max = time.Duration(math.MinInt64)
	var min = time.Duration(math.MaxInt64)

	var from [100000]int
	var to [100000]int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		from[i] = rand.Intn(len(meshgrid))
		to[i] = rand.Intn(len(meshgrid))
	}

	var wg sync.WaitGroup

	fmt.Printf("\n%v\n", len(meshgrid))
	for i := 0; i < len(from); i++ {
		if !meshgrid[from[i]] && !meshgrid[to[i]] {
			var x = from[i]
			var y = to[i]
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				var start = time.Now()
				var a = algorithms.ExpandIndex(x, xSize)
				var b = algorithms.ExpandIndex(y, xSize)
				var _ = algorithms.Dijkstra(a[0], a[1], b[0], b[1], xSize, ySize, &meshgrid)
				t := time.Now()
				var elapsed = t.Sub(start)
				if elapsed > max {
					max = elapsed
				}
				if elapsed < min {
					min = elapsed
				}
				sum += elapsed
				count++
			}(x, y)
		}
	}

	/*fmt.Printf("\n%v\n", len(meshgrid))
	for i := 0; i < len(meshgrid)/1000; i++ {
		if !meshgrid[i] {
			for j := 0; j < len(meshgrid)/1000; j++ {
				if !meshgrid[j] {
					wg.Add(1)
					go func(i, j int) {
						defer wg.Done()
						var start = time.Now()
						var a = algorithms.ExpandIndex(i, xSize)
						var b = algorithms.ExpandIndex(j, xSize)
						var _ = algorithms.Dijkstra(a[0], a[1], b[0], b[1], xSize, ySize, &meshgrid)
						t := time.Now()
						var elapsed = t.Sub(start)
						sum += elapsed
						count++
					}(i, j)
				}
			}
		}
	}*/

	wg.Wait()
	fmt.Printf("Total Time: %v\nNumber of routings: %v\nAverage duration: %v\nMin, Max: %v, %v", sum, count, sum/time.Duration(count), min, max)
	//fmt.Printf("%v\n", testCount)
}

// WayFinding is evaluated
func WayFinding(xSize, ySize, algorithm int, basicPointInPolygon bool, note string) {
	var filename string
	var meshgrid []bool
	var meshgrid2d [][]bool
	var uniformgrid []bool
	var uniformgrid2d algorithms.UniformGrid

	logging := make(map[string]string)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logging["note"] = note
	logging["basicGrid"] = "false"
	logging["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	logging["xSize"] = strconv.Itoa(xSize)
	logging["ySize"] = strconv.Itoa(ySize)
	logging["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logging["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/uniformGrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/uniformGrid_%v_%v.json", xSize, ySize)
	}

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &uniformgrid2d)
	for i := 0; i < len(uniformgrid2d.VertexData); i++ {
		for j := 0; j < len(uniformgrid2d.VertexData[i]); j++ {
			uniformgrid = append(uniformgrid, uniformgrid2d.VertexData[i][j])
		}
	}

	//var testCount int
	// TODO: More efficient average
	var sum time.Duration
	var count int

	for i := 0; i < len(meshgrid2d); i++ {
		for j := 0; j < len(meshgrid2d[i]); j++ {
			for k := 0; k < len(meshgrid2d); k++ {
				for l := 0; l < len(meshgrid2d[k]); l++ {
					if !meshgrid2d[i][j] && !meshgrid2d[k][l] {
						//testCount += 1
						// fmt.Printf("%v:%v - %v:%v", j, i, l, k)
						var start = time.Now()
						var _ = algorithms.Dijkstra(int(i), int(j), int(k), int(l), int(xSize), int(ySize), &meshgrid)
						t := time.Now()
						var elapsed = t.Sub(start)
						sum += elapsed
						count++
					}
				}
			}
		}
	}

	fmt.Printf("%v, %v, %v", sum, count, float64(sum)/float64(count))
	//fmt.Printf("%v\n", testCount)
}

// DataProcessing is evaluated
func DataProcessing(pbfFileName, note string, xSize, ySize int, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon bool) {
	var logging = dataprocessing.Start(pbfFileName, xSize, ySize, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)

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
