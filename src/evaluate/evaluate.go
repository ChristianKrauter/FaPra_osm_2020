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

	jsonString, _ := json.MarshalIndent(logging, "", "    ")
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
func WayFindingBG(xSize, ySize, nRuns int, basicPointInPolygon bool, note string) {
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
		logging["filename"] += "_bpip"
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

	var sum time.Duration
	var count int
	var max = time.Duration(math.MinInt64)
	var min = time.Duration(math.MaxInt64)
	from := make([]int, nRuns)
	to := make([]int, nRuns)
	var wg sync.WaitGroup

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nRuns; i++ {
		for {
			from[i] = rand.Intn(len(meshgrid))
			to[i] = rand.Intn(len(meshgrid))
			if !meshgrid[from[i]] && !meshgrid[to[i]] {
				break
			}
		}
	}

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

	/*for i := 0; i < len(meshgrid)/1000; i++ {
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

	logging["time_sum"] = sum.String()
	logging["count_runs"] = strconv.Itoa(count)
	logging["time_avg"] = (sum / time.Duration(count)).String()
	logging["time_min"] = min.String()
	logging["time_max"] = max.String()

	jsonString, _ := json.MarshalIndent(logging, "", "    ")
	var outFilename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := logging["filename"]; ok {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg_dij%s_%s.json", logging["xSize"], logging["ySize"], val, timestamp)
	} else {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg_dij_%s.json", logging["xSize"], logging["ySize"], timestamp)
	}

	f, err := os.Create(outFilename)
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

// WayFinding is evaluated
func WayFinding(xSize, ySize, nRuns, algorithm int, basicPointInPolygon bool, note string) {
	var filename string
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
	var algostring string
	switch algorithm {
	case 0:
		algostring = "_dij"
	default:
		algostring = "_dij"
	}

	logging["filename"] = algostring

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/uniformGrid_%v_%v_bpip.json", xSize, ySize)
		logging["filename"] += "_bpip"
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

	var sum time.Duration
	var count int
	var max = time.Duration(math.MinInt64)
	var min = time.Duration(math.MaxInt64)
	from := make([]int, nRuns)
	to := make([]int, nRuns)
	var wg sync.WaitGroup

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nRuns; i++ {
		from[i] = rand.Intn(len(uniformgrid))
		to[i] = rand.Intn(len(uniformgrid))
	}

	for i := 0; i < len(from); i++ {
		if !uniformgrid[from[i]] && !uniformgrid[to[i]] {
			var x = from[i]
			var y = to[i]
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				var start = time.Now()
				var a = uniformgrid2d.IDToGrid(x)
				var b = uniformgrid2d.IDToGrid(y)
				var _ = algorithms.UniformDijkstra(a[0], a[1], b[0], b[1], xSize, ySize, &uniformgrid2d)
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

	wg.Wait()
	fmt.Printf("Total Time: %v\nNumber of routings: %v\nAverage duration: %v\nMin, Max: %v, %v", sum, count, sum/time.Duration(count), min, max)

	logging["time_sum"] = sum.String()
	logging["count_runs"] = strconv.Itoa(count)
	logging["time_avg"] = (sum / time.Duration(count)).String()
	logging["time_min"] = min.String()
	logging["time_max"] = max.String()

	jsonString, _ := json.MarshalIndent(logging, "", "    ")
	var outFilename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s%s_%s.json", logging["xSize"], logging["ySize"], logging["filename"], timestamp)

	f, err := os.Create(outFilename)
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

	jsonString, _ := json.MarshalIndent(logging, "", "    ")
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
