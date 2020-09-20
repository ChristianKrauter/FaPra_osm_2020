package evaluate

import (
	"../algorithms"
	"../dataprocessing"
	"../grids"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	logs := make(map[string]string)
	logs["pbfFileName"] = pbfFileName
	pbfFileName = fmt.Sprintf("data/%s", pbfFileName)

	fmt.Printf("\nNormal:\n")
	var coastlineMap = make(map[int64][]int64)
	var nodeMap = make(map[int64][]float64)
	var readTime = dataprocessing.ReadFile(pbfFileName, &coastlineMap, &nodeMap)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	var totalAllocNormal = m.TotalAlloc

	fmt.Printf("\nLess memory:\n")
	coastlineMap = make(map[int64][]int64)
	nodeMap = make(map[int64][]float64)
	var readLessMemoryTime = dataprocessing.ReadFileLessMemory(pbfFileName, &coastlineMap, &nodeMap)

	runtime.ReadMemStats(&m)
	logs["note"] = note
	logs["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logs["totalAllocNormal"] = strconv.FormatUint(totalAllocNormal/1024/1024, 10)
	logs["totalAllocLessMemory"] = strconv.FormatUint((m.TotalAlloc-totalAllocNormal)/1024/1024, 10)
	logs["time_read"] = string(readTime)
	logs["time_readLessMemory"] = string(readLessMemoryTime)

	jsonString, _ := json.MarshalIndent(logs, "", "    ")
	var filename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	filename = fmt.Sprintf("data/evaluation/rf_%s_%s.json", strings.Split(logs["pbfFileName"], ".")[0], timestamp)
	saveLog(filename, jsonString)
}

// WayFindingBg is evaluated
func WayFindingBg(xSize, ySize, nRuns, algorithm int, note string) {
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
	logs := make(map[string]string)

	bg.XSize = xSize
	bg.YSize = ySize
	bg.XFactor = float64(xSize) / 360.0
	bg.YFactor = float64(ySize) / 360.0

	var algoStr, algoStrPrint string
	switch algorithm {
	case 0:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	case 1:
		algoStr = "_as"
		algoStrPrint = "A-Star"
	case 2:
		algoStr = "_bidij"
		algoStrPrint = "Bi-Dijkstra"
	case 3:
		algoStr = "_bias"
		algoStrPrint = "Bi-A-Star"
	default:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	}

	fmt.Printf("using %s.\n", algoStrPrint)
	logs["filename"] = algoStr

	filename = fmt.Sprintf("data/output/meshgrid_%v_%v.json", xSize, ySize)

	meshgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		log.Fatal(fmt.Sprintf("\nThe meshgrid '%s'\ncould not be found. Please create it first.\n", filename))
	}
	defer meshgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(meshgridRaw)
	json.Unmarshal(byteValue, &bg2D)

	bg.VertexData = make([]bool, xSize*ySize)
	k := 0
	for i := 0; i < len(bg2D[0]); i++ {
		for j := 0; j < len(bg2D); j++ {
			bg.VertexData[k] = bg2D[j][i]
			k++
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
		var start time.Time
		var end time.Time
		var popped int
		switch algorithm {
		case 0:
			start = time.Now()
			_, popped, _ = algorithms.DijkstraBg(from[i], to[i], &bg)
			end = time.Now()
		case 1:
			start = time.Now()
			_, popped, _ = algorithms.AStarBg(from[i], to[i], &bg)
			end = time.Now()
		case 2:
			start = time.Now()
			_, popped, _ = algorithms.BiDijkstraBg(from[i], to[i], &bg)
			end = time.Now()
		case 3:
			start = time.Now()
			_, popped, _ = algorithms.BiAStarBg(from[i], to[i], &bg)
			end = time.Now()
		default:
			start = time.Now()
			_, popped, _ = algorithms.DijkstraBg(from[i], to[i], &bg)
			end = time.Now()
		}

		poppedSum += popped
		var elapsed = end.Sub(start)
		if elapsed > max {
			max = elapsed
		}
		if elapsed < min {
			min = elapsed
		}
		sum += elapsed
		count++
	}

	fmt.Printf("\nNumber of routings: %v\nTotal Time        : %v\nAVG Time          : %v\nAVG nodes popped  : %v\nMin, Max Time     : %v, %v",
		count, sum, sum/time.Duration(count), strconv.Itoa(poppedSum/count), min, max)

	runtime.ReadMemStats(&m)
	logs["note"] = note
	logs["basicGrid"] = "true"
	logs["xSize"] = strconv.Itoa(xSize)
	logs["ySize"] = strconv.Itoa(ySize)
	logs["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logs["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	logs["time_sum"] = sum.String()
	logs["count_runs"] = strconv.Itoa(count)
	logs["time_avg"] = (sum / time.Duration(count)).String()
	logs["time_min"] = min.String()
	logs["time_max"] = max.String()
	logs["nodes_popped_avg"] = strconv.Itoa(poppedSum / count)

	jsonString, _ := json.MarshalIndent(logs, "", "    ")
	var outFilename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := logs["filename"]; ok {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg%s_%s.json", logs["xSize"], logs["ySize"], val, timestamp)
	} else {
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg_%s.json", logs["xSize"], logs["ySize"], timestamp)
	}
	saveLog(outFilename, jsonString)
}

// WayFinding is evaluated
func WayFinding(xSize, ySize, nRuns, algorithm int, note string) {
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
	logs := make(map[string]string)

	var algoStr, algoStrPrint string
	switch algorithm {
	case 0:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	case 1:
		algoStr = "_as"
		algoStrPrint = "A-Star"
	case 2:
		algoStr = "_bidij"
		algoStrPrint = "Bi-Dijkstra"
	case 3:
		algoStr = "_bias"
		algoStrPrint = "Bi-A-Star"
	case 4:
		algoStr = "_asjps"
		algoStrPrint = "AStar-JPS"
	default:
		algoStr = "_dij"
		algoStrPrint = "Dijkstra"
	}

	fmt.Printf("using %s.\n", algoStrPrint)
	logs["filename"] = algoStr

	filename = fmt.Sprintf("data/output/uniformGrid_%v_%v.json", xSize, ySize)
	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		log.Fatal(fmt.Sprintf("\nThe meshgrid '%s'\ncould not be found. Please create it first.\n", filename))
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &ug)

	ug1D = make([]bool, ug.N)
	k := 0
	for i := 0; i < len(ug.VertexData); i++ {
		for j := 0; j < len(ug.VertexData[i]); j++ {
			ug1D[k] = ug.VertexData[i][j]
			k++
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
		var start time.Time
		var end time.Time
		var popped int

		switch algorithm {
		case 0:
			start = time.Now()
			_, popped, _ = algorithms.Dijkstra(from[i], to[i], &ug)
			end = time.Now()
		case 1:
			start = time.Now()
			_, popped, _ = algorithms.AStar(from[i], to[i], &ug)
			end = time.Now()
		case 2:
			start = time.Now()
			_, popped, _ = algorithms.BiDijkstra(from[i], to[i], &ug)
			end = time.Now()
		case 3:
			start = time.Now()
			_, popped, _ = algorithms.BiAStar(from[i], to[i], &ug)
			end = time.Now()
		case 4:
			start = time.Now()
			_, popped, _ = algorithms.AStarJPS(from[i], to[i], &ug)
			end = time.Now()
		default:
			start = time.Now()
			_, popped, _ = algorithms.Dijkstra(from[i], to[i], &ug)
			end = time.Now()
		}
		poppedSum += popped
		var elapsed = end.Sub(start)
		if elapsed > max {
			max = elapsed
		}
		if elapsed < min {
			min = elapsed
		}
		sum += elapsed
		count++
	}

	fmt.Printf("\nNumber of routings: %v\nTotal Time        : %v\nAVG Time          : %v\nAVG nodes popped  : %v\nMin, Max Time     : %v, %v",
		count, sum, sum/time.Duration(count), strconv.Itoa(poppedSum/count), min, max)

	runtime.ReadMemStats(&m)
	logs["note"] = note
	logs["basicGrid"] = "false"
	logs["xSize"] = strconv.Itoa(xSize)
	logs["ySize"] = strconv.Itoa(ySize)
	logs["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logs["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)
	logs["time_sum"] = sum.String()
	logs["count_runs"] = strconv.Itoa(count)
	logs["time_avg"] = (sum / time.Duration(count)).String()
	logs["time_min"] = min.String()
	logs["time_max"] = max.String()
	logs["nodes_popped_avg"] = strconv.Itoa(poppedSum / count)

	jsonString, _ := json.MarshalIndent(logs, "", "    ")
	var outFilename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s%s_%s.json", logs["xSize"], logs["ySize"], logs["filename"], timestamp)
	saveLog(outFilename, jsonString)
}

// DataProcessing is evaluated
func DataProcessing(pbfFileName, note string, xSize, ySize int, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon bool) {
	var logs = dataprocessing.Start(pbfFileName, xSize, ySize, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logs["pbfFileName"] = pbfFileName
	logs["note"] = note
	logs["basicGrid"] = strconv.FormatBool(basicGrid)
	logs["basicPointInPolygon"] = strconv.FormatBool(basicPointInPolygon)
	logs["xSize"] = strconv.Itoa(xSize)
	logs["ySize"] = strconv.Itoa(ySize)
	logs["numCPU"] = strconv.Itoa(runtime.NumCPU())
	logs["totalAlloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)

	jsonString, _ := json.MarshalIndent(logs, "", "    ")
	var filename string
	var timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := logs["filename"]; ok {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s%s_%s.json", strings.Split(logs["pbfFileName"], ".")[0], logs["xSize"], logs["ySize"], val, timestamp)
	} else {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s_%s.json", strings.Split(logs["pbfFileName"], ".")[0], logs["xSize"], logs["ySize"], timestamp)
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
		log.Fatal(fmt.Sprintf("\nThe meshgrid '%s'\ncould not be found. Please create it first.\n", filename))
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
