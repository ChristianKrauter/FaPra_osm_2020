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

func saveLog(filename string, jsonString []byte) {
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

// ReadPBF is evaluated
func ReadPBF(pbfFileName, note string) {
	log := make(map[string]string)
	log["pbfFileName"] = pbfFileName
	pbfFileName = fmt.Sprintf("data/%s", pbfFileName)

	var coastlineMap = make(map[int64][]int64)
	var nodeMap = make(map[int64][]float64)
	var readTime = dataprocessing.ReadFile(pbfFileName, &coastlineMap, &nodeMap)

	coastlineMap = make(map[int64][]int64)
	nodeMap = make(map[int64][]float64)
	var readLessMemoryTime = dataprocessing.ReadFileLessMemory(pbfFileName, &coastlineMap, &nodeMap)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log["note"] = note
	log["numCPU"] = strconv.Itoa(runtime.NumCPU())
	log["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	log["time_read"] = string(readTime)
	log["time_readLessMemory"] = string(readLessMemoryTime)

	jsonString, _ := json.MarshalIndent(log, "", "    ")
	var filename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	filename = fmt.Sprintf("data/evaluation/rf_%s_%s.json", strings.Split(log["pbfFileName"], ".")[0], timestamp)
	saveLog(filename, jsonString)
}

// WayFindingBG is evaluated
func WayFindingBG(xSize, ySize, nRuns int, basicPointInPolygon bool, note string) {
	var filename string
	var meshgrid []bool
	var meshgrid2d [][]bool
	var m runtime.MemStats
	var sum time.Duration
	var count int
	var max = time.Duration(math.MinInt64)
	var min = time.Duration(math.MaxInt64)
	from := make([]int, nRuns)
	to := make([]int, nRuns)
	var wg sync.WaitGroup

	log := make(map[string]string)
	log["note"] = note
	log["basicGrid"] = "true"
	log["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	log["xSize"] = strconv.Itoa(xSize)
	log["ySize"] = strconv.Itoa(ySize)
	log["numCPU"] = strconv.Itoa(runtime.NumCPU())

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v_bpip.json", xSize, ySize)
		log["filename"] += "_bpip"
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

	wg.Wait()
	fmt.Printf("Total Time: %v\nNumber of routings: %v\nAverage duration: %v\nMin, Max: %v, %v", sum, count, sum/time.Duration(count), min, max)

	runtime.ReadMemStats(&m)
	log["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	log["time_sum"] = sum.String()
	log["count_runs"] = strconv.Itoa(count)
	log["time_avg"] = (sum / time.Duration(count)).String()
	log["time_min"] = min.String()
	log["time_max"] = max.String()

	jsonString, _ := json.MarshalIndent(log, "", "    ")
	var outFilename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := log["filename"]; ok {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg_dij%s_%s.json", log["xSize"], log["ySize"], val, timestamp)
	} else {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg_dij_%s.json", log["xSize"], log["ySize"], timestamp)
	}
	saveLog(outFilename, jsonString)
}

// WayFinding is evaluated
func WayFinding(xSize, ySize, nRuns, algorithm int, basicPointInPolygon bool, note string) {
	var filename string
	var ug1D []bool
	var ug algorithms.UniformGrid
	var m runtime.MemStats
	var sum time.Duration
	var count int
	var max = time.Duration(math.MinInt64)
	var min = time.Duration(math.MaxInt64)
	from := make([]int, nRuns)
	to := make([]int, nRuns)
	var wg sync.WaitGroup

	log := make(map[string]string)
	log["note"] = note
	log["basicGrid"] = "false"
	log["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	log["xSize"] = strconv.Itoa(xSize)
	log["ySize"] = strconv.Itoa(ySize)
	log["numCPU"] = strconv.Itoa(runtime.NumCPU())

	var algostring string
	switch algorithm {
	case 0:
		algostring = "_dij"
	default:
		algostring = "_dij"
	}

	log["filename"] = algostring

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/uniformGrid_%v_%v_bpip.json", xSize, ySize)
		log["filename"] += "_bpip"
	} else {
		filename = fmt.Sprintf("data/output/uniformGrid_%v_%v.json", xSize, ySize)
	}

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &ug)
	for i := 0; i < len(ug.VertexData); i++ {
		for j := 0; j < len(ug.VertexData[i]); j++ {
			ug1D = append(ug1D, ug.VertexData[i][j])
		}
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nRuns; i++ {
		for {
			from[i] = rand.Intn(len(ug1D))
			to[i] = rand.Intn(len(ug1D))
			if !ug1D[from[i]] && !ug1D[to[i]] {
				break
			}
		}
	}

	ug.XSize = xSize
	ug.YSize = ySize
	ug.BigN = xSize * ySize
	ug.A = 4.0 * math.Pi / float64(ug.BigN)
	ug.D = math.Sqrt(ug.A)
	ug.MTheta = math.Round(math.Pi / ug.D)
	ug.DTheta = math.Pi / ug.MTheta
	ug.DPhi = ug.A / ug.DTheta

	for i := 0; i < len(from); i++ {
		var x = from[i]
		var y = to[i]
		wg.Add(1)
		go func(x, y int) {
			defer wg.Done()
			var start = time.Now()
			var a = ug.IDToGrid(x)
			var b = ug.IDToGrid(y)
			var _ = algorithms.UniformDijkstra(a[0], a[1], b[0], b[1], &ug)
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

	wg.Wait()
	fmt.Printf("Total Time: %v\nNumber of routings: %v\nAverage duration: %v\nMin, Max: %v, %v", sum, count, sum/time.Duration(count), min, max)

	runtime.ReadMemStats(&m)
	log["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	log["time_sum"] = sum.String()
	log["count_runs"] = strconv.Itoa(count)
	log["time_avg"] = (sum / time.Duration(count)).String()
	log["time_min"] = min.String()
	log["time_max"] = max.String()

	jsonString, _ := json.MarshalIndent(log, "", "    ")
	var outFilename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s%s_%s.json", log["xSize"], log["ySize"], log["filename"], timestamp)
	saveLog(outFilename, jsonString)
}

// DataProcessing is evaluated
func DataProcessing(pbfFileName, note string, xSize, ySize int, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon bool) {
	var log = dataprocessing.Start(pbfFileName, xSize, ySize, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log["pbfFileName"] = pbfFileName
	log["note"] = note
	log["basicGrid"] = strconv.FormatBool(basicGrid)
	log["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	log["xSize"] = strconv.Itoa(xSize)
	log["ySize"] = strconv.Itoa(ySize)
	log["numCPU"] = strconv.Itoa(runtime.NumCPU())
	log["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)

	jsonString, _ := json.MarshalIndent(log, "", "    ")
	var filename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := log["filename"]; ok {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s%s_%s.json", strings.Split(log["pbfFileName"], ".")[0], log["xSize"], log["ySize"], val, timestamp)
	} else {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s_%s.json", strings.Split(log["pbfFileName"], ".")[0], log["xSize"], log["ySize"], timestamp)
	}
	saveLog(filename, jsonString)
}
