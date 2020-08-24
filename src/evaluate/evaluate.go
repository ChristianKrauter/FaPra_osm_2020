package evaluate

import (
	"../algorithms"
	"../dataprocessing"
	"../grids"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
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
func WayFindingBG(xSize, ySize, nRuns, algorithm int, basicPointInPolygon bool, note string) {
	var filename string
	var bg grids.BasicGrid
	var bg2D [][]bool
	var m runtime.MemStats
	var sum time.Duration
	var poppedSum int
	var count int
	var max = time.Duration(math.MinInt64)
	var min = time.Duration(math.MaxInt64)
	from := make([]int, nRuns)
	to := make([]int, nRuns)

	bg.XSize = xSize
	bg.YSize = ySize

	log := make(map[string]string)
	log["note"] = note
	log["basicGrid"] = "true"
	log["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	log["xSize"] = strconv.Itoa(xSize)
	log["ySize"] = strconv.Itoa(ySize)
	log["numCPU"] = strconv.Itoa(runtime.NumCPU())

	var algoStr, algoStrPrint string
	switch algorithm {
	case 0:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	case 1:
		algoStr = "_as"
		algoStrPrint = "A-Star"
	default:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	}

	fmt.Printf("using %s.\n", algoStrPrint)
	log["filename"] = algoStr

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
	json.Unmarshal(byteValue, &bg2D)

	for i := 0; i < len(bg2D[0]); i++ {
		for j := 0; j < len(bg2D); j++ {
			bg.VertexData = append(bg.VertexData, bg2D[j][i])
		}
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nRuns; i++ {
		for {
			from[i] = rand.Intn(len(bg.VertexData))
			to[i] = rand.Intn(len(bg.VertexData))
			if !bg.VertexData[from[i]] && !bg.VertexData[to[i]] {
				break
			}
		}
	}

	for i := 0; i < len(from); i++ {
		var start = time.Now()

		var popped int
		switch algorithm {
		case 0:
			_, popped = algorithms.DijkstraBg(bg.ExpandIndex(from[i]), bg.ExpandIndex(to[i]), &bg)
		case 1:
			_, popped = algorithms.AStarBg(bg.ExpandIndex(from[i]), bg.ExpandIndex(to[i]), &bg)
		default:
			_, popped = algorithms.DijkstraBg(bg.ExpandIndex(from[i]), bg.ExpandIndex(to[i]), &bg)
		}

		poppedSum += popped
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
	}

	fmt.Printf("Total Time: %v\nNumber of routings: %v\nAverage duration: %v\nMin, Max: %v, %v", sum, count, sum/time.Duration(count), min, max)

	runtime.ReadMemStats(&m)
	log["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	log["time_sum"] = sum.String()
	log["count_runs"] = strconv.Itoa(count)
	log["time_avg"] = (sum / time.Duration(count)).String()
	log["time_min"] = min.String()
	log["time_max"] = max.String()
	log["nodes_popped_avg"] = strconv.Itoa(poppedSum / count)

	jsonString, _ := json.MarshalIndent(log, "", "    ")
	var outFilename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := log["filename"]; ok {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg%s_%s.json", log["xSize"], log["ySize"], val, timestamp)
	} else {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg_%s.json", log["xSize"], log["ySize"], timestamp)
	}
	saveLog(outFilename, jsonString)
}

// WayFinding is evaluated
func WayFinding(xSize, ySize, nRuns, algorithm int, basicPointInPolygon bool, note string) {
	var filename string
	var ug1D []bool
	var ug grids.UniformGrid
	var m runtime.MemStats
	var sum time.Duration
	var poppedSum int
	var count int
	var max = time.Duration(math.MinInt64)
	var min = time.Duration(math.MaxInt64)
	from := make([]int, nRuns)
	to := make([]int, nRuns)

	log := make(map[string]string)
	log["note"] = note
	log["basicGrid"] = "false"
	log["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	log["xSize"] = strconv.Itoa(xSize)
	log["ySize"] = strconv.Itoa(ySize)
	log["numCPU"] = strconv.Itoa(runtime.NumCPU())

	var algoStr, algoStrPrint string
	switch algorithm {
	case 0:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	case 1:
		algoStr = "_as"
		algoStrPrint = "A-Star"
	default:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	}

	fmt.Printf("using %s.\n", algoStrPrint)
	log["filename"] = algoStr

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

	for i := 0; i < len(from); i++ {
		var start = time.Now()
		var from = ug.IDToGrid(from[i])
		var to = ug.IDToGrid(to[i])
		var popped int
		switch algorithm {
		case 0:
			_, popped = algorithms.Dijkstra(from, to, &ug)
		case 1:
			_, popped = algorithms.AStar(from, to, &ug)
		default:
			_, popped = algorithms.Dijkstra(from, to, &ug)
		}
		poppedSum += popped
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
	}

	fmt.Printf("Total Time: %v\nNumber of routings: %v\nAverage duration: %v\nMin, Max: %v, %v", sum, count, sum/time.Duration(count), min, max)

	runtime.ReadMemStats(&m)
	log["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	log["time_sum"] = sum.String()
	log["count_runs"] = strconv.Itoa(count)
	log["time_avg"] = (sum / time.Duration(count)).String()
	log["time_min"] = min.String()
	log["time_max"] = max.String()
	log["nodes_popped"] = strconv.Itoa(poppedSum / count)

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
