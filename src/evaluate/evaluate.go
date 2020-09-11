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
	"sort"
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
func WayFindingBG(xSize, ySize, nRuns int, basicPointInPolygon bool, note string) {
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
		var a = bg.ExpandIndex(from[i])
		var b = bg.ExpandIndex(to[i])
		var _, popped = algorithms.DijkstraBg(a, b, &bg)
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

	for i := 0; i < len(from); i++ {
		var start = time.Now()
		var from = ug.IDToGrid(from[i])
		var to = ug.IDToGrid(to[i])
		var _, popped = algorithms.Dijkstra(from, to, &ug)
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

// NeighboursUg is evaluated and tested
func NeighboursUg(xSize, ySize int, note string) {
	var filename string
	var ug grids.UniformGrid
	var failedIDS = make(map[int]struct{})
	var exists = struct{}{}
	var sumSN time.Duration
	var sumN time.Duration

	filename = fmt.Sprintf("data/output/uniformGrid_%v_%v.json", xSize, ySize)

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		panic(errJSON)
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &ug)

	k := 0
	for i := 0; i < len(ug.VertexData); i++ {
		for j := 0; j < len(ug.VertexData[i]); j++ {
			var start = time.Now()
			var sn = algorithms.SimpleNeighboursUg(k, &ug)
			t := time.Now()
			sumSN += t.Sub(start)

			start = time.Now()
			var n = algorithms.NeighboursUg(k, &ug)
			t = time.Now()
			sumN += t.Sub(start)

			if len(sn) != len(n) {
				failedIDS[k] = exists
			} else {
				sort.Ints(sn)
				sort.Ints(n)

				for i := 0; i < len(sn); i++ {
					if sn[i] != n[i] {
						failedIDS[k] = exists
						continue
					}
				}
			}
			k++
		}
	}
	fmt.Printf("\nAVG SN Time : %v\nAVG N Time  : %v\nErrors      : %v / %v (%.3f %%)\nSpeed up    : Simple neighbours is %.2f times faster\n\n",
		sumSN/time.Duration(ug.N), sumN/time.Duration(ug.N), len(failedIDS), k, float64(len(failedIDS))/float64(k)*100, float64(sumN)/float64(sumSN))

	if len(failedIDS) > 0 {

		var uniqueFailedIDS []int
		for i := range failedIDS {
			uniqueFailedIDS = append(uniqueFailedIDS, i)
		}
		sort.Ints(uniqueFailedIDS)

		fmt.Printf("Error IDXs:\n%v\n", uniqueFailedIDS)
		sn := algorithms.SimpleNeighboursUg(uniqueFailedIDS[0], &ug)
		n := algorithms.NeighboursUg(uniqueFailedIDS[0], &ug)
		fmt.Printf("\nerror for IDX %v:\n", uniqueFailedIDS[0])
		sort.Ints(sn)
		sort.Ints(n)
		fmt.Printf("sn: %v\n", sn)
		fmt.Printf("n : %v\n", n)
	}
}
