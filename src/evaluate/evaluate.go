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

// WayFinding is evaluated
func WayFinding(xSize, ySize, nRuns int, note string) {
	var filename string
	var ug1D []bool
	var ug grids.UniformGrid
	var m runtime.MemStats
	var sum = make([]time.Duration, 5)
	var poppedSum = make([]int, 5)
	var count int
	var max = make([]time.Duration, 5)
	var min = make([]time.Duration, 5)
	var from = make([]int, nRuns)
	var to = make([]int, nRuns)
	var lengths = make([]Length, 5)
	const TOLERANCE = 0.001 // 0.1%#
	var outFilename string
	var timestamp string
	var Logs Log
	Logs.Parameters = make(map[string]string)
	Logs.Results = make(map[string]Algo)
	var exampleCount = 0
	var exampleCountMax = 10

	for i := 0; i < 5; i++ {
		max[i] = time.Duration(math.MinInt64)
		min[i] = time.Duration(math.MaxInt64)
	}

	fmt.Printf("\nwith a length tolerance of %v%%.\n", TOLERANCE*100)
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
		var start = make([]time.Time, 5)
		var end = make([]time.Time, 5)
		var popped = make([]int, 5)
		var dists = make([]float64, 5)

		start[0] = time.Now()
		_, popped[0], dists[0] = algorithms.Dijkstra(from[i], to[i], &ug)
		end[0] = time.Now()

		start[1] = time.Now()
		_, popped[1], dists[1] = algorithms.AStar(from[i], to[i], &ug)
		end[1] = time.Now()

		start[2] = time.Now()
		_, popped[2], dists[2] = algorithms.BiDijkstra(from[i], to[i], &ug)
		end[2] = time.Now()

		start[3] = time.Now()
		_, popped[3], dists[3] = algorithms.BiAStar(from[i], to[i], &ug)
		end[3] = time.Now()

		start[4] = time.Now()
		_, popped[4], dists[4] = algorithms.AStarJPS(from[i], to[i], &ug)
		end[4] = time.Now()

		if dists[0] == math.Inf(1) {
			fmt.Printf("\nThere was a combination which returned no route: %v, %v\n", from[i], to[i])
		} else {
			eq := true
			text := "Longer"
			for j := 0; j < 5; j++ {
				// PQ-Pops
				poppedSum[j] += popped[j]

				// Time
				var elapsed = end[j].Sub(start[j])
				if elapsed > max[j] {
					max[j] = elapsed
				}
				if elapsed < min[j] {
					min[j] = elapsed
				}
				sum[j] += elapsed

				// Length
				diff := dists[0] / dists[j]
				if diff <= 1+TOLERANCE && diff >= 1-TOLERANCE {
					lengths[j].Equal++
				} else if diff > 1+TOLERANCE {
					eq = false
					text = "Shorter"
					lengths[j].Shorter++
				} else {
					eq = false
					lengths[j].Longer++
				}
			}
			count++
			if !eq && exampleCount < exampleCountMax {
				exampleCount++
				fmt.Printf("\n%v: %v, %v", text, from[i], to[i])
			}
		}
	}

	if count > 0 {
		runtime.ReadMemStats(&m)
		Logs.Parameters["Note"] = note
		Logs.Parameters["Basic Grid"] = "false"
		Logs.Parameters["X-Size"] = strconv.Itoa(xSize)
		Logs.Parameters["Y-Size"] = strconv.Itoa(ySize)
		Logs.Parameters["CPU Cores"] = strconv.Itoa(runtime.NumCPU())
		Logs.Parameters["Lenght Tolerance (%)"] = strconv.FormatFloat(TOLERANCE*100, 'f', -1, 64)
		Logs.Parameters["Run Count"] = strconv.Itoa(count)

		for j := 0; j < 5; j++ {
			var a Algo
			a.Time = Time{
				Sum:         sum[j].String(),
				AVG:         (sum[j] / time.Duration(count)).String(),
				Min:         min[j].String(),
				Max:         max[j].String(),
				TimesFaster: math.Floor((float64(sum[0]/time.Duration(count))/float64(sum[j]/time.Duration(count)))*100) / 100,
			}
			a.PQPops = PQPops{
				AVG:     poppedSum[j] / count,
				Percent: math.Floor((float64(poppedSum[j]/count)/float64(poppedSum[0]/count))*10000) / 100,
			}
			a.Length = lengths[j]

			Logs.Results[algoStrPrint[j]] = a
		}

		for j := 0; j < 5; j++ {
			jsT, _ := json.MarshalIndent(Logs.Results[algoStrPrint[j]].Time, "", "    ")
			jsL, _ := json.MarshalIndent(Logs.Results[algoStrPrint[j]].Length, "", "    ")
			jsP, _ := json.MarshalIndent(Logs.Results[algoStrPrint[j]].PQPops, "", "    ")

			fmt.Printf("\n\n%v\n", algoStrPrint[j])
			fmt.Printf("Time:\n%v\n", string(jsT))
			fmt.Printf("Length:\n%v\n", string(jsL))
			fmt.Printf("PQ-Pops:\n%v\n", string(jsP))
		}

		jsonString, _ := json.MarshalIndent(Logs, "", "    ")
		timestamp = time.Now().Format("2006-01-02_15-04-05")
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_%s.json", Logs.Parameters["X-Size"], Logs.Parameters["Y-Size"], timestamp)
		saveLog(outFilename, jsonString)
	}
}

// WayFindingBg is evaluated
func WayFindingBg(xSize, ySize, nRuns int, note string) {
	var filename string
	var bg grids.BasicGrid
	var bg2D [][]bool
	var m runtime.MemStats
	var sum = make([]time.Duration, 5)
	var poppedSum = make([]int, 5)
	var count int
	var max = make([]time.Duration, 5)
	var min = make([]time.Duration, 5)
	from := make([]int, nRuns)
	to := make([]int, nRuns)
	var lengths = make([]Length, 5)
	const TOLERANCE = 0.001 // 0.1%#
	var outFilename string
	var timestamp string
	var Logs Log
	Logs.Parameters = make(map[string]string)
	Logs.Results = make(map[string]Algo)
	var exampleCount = 0
	var exampleCountMax = 10

	bg.XSize = xSize
	bg.YSize = ySize
	bg.XFactor = float64(xSize) / 360.0
	bg.YFactor = float64(ySize) / 360.0

	for i := 0; i < 5; i++ {
		max[i] = time.Duration(math.MinInt64)
		min[i] = time.Duration(math.MaxInt64)
	}

	fmt.Printf("\nwith a length tolerance of %v%%.\n", TOLERANCE*100)
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
		var start = make([]time.Time, 5)
		var end = make([]time.Time, 5)
		var popped = make([]int, 5)
		var dists = make([]float64, 5)

		start[0] = time.Now()
		_, popped[0], dists[0] = algorithms.DijkstraBg(from[i], to[i], &bg)
		end[0] = time.Now()

		start[1] = time.Now()
		_, popped[1], dists[1] = algorithms.AStarBg(from[i], to[i], &bg)
		end[1] = time.Now()

		start[2] = time.Now()
		_, popped[2], dists[2] = algorithms.BiDijkstraBg(from[i], to[i], &bg)
		end[2] = time.Now()

		start[3] = time.Now()
		_, popped[3], dists[3] = algorithms.BiAStarBg(from[i], to[i], &bg)
		end[3] = time.Now()

		start[4] = time.Now()
		_, popped[4], dists[4] = algorithms.AStarJPSBg(from[i], to[i], &bg)
		end[4] = time.Now()

		if dists[0] == math.Inf(1) {
			fmt.Printf("\nThere was a combination which returned no route: %v, %v\n", from[i], to[i])
		} else {
			eq := true
			text := "Longer"
			for j := 0; j < 5; j++ {
				// PQ-Pops
				poppedSum[j] += popped[j]

				// Time
				var elapsed = end[j].Sub(start[j])
				if elapsed > max[j] {
					max[j] = elapsed
				}
				if elapsed < min[j] {
					min[j] = elapsed
				}
				sum[j] += elapsed

				// Length
				diff := dists[0] / dists[j]
				if diff <= 1+TOLERANCE && diff >= 1-TOLERANCE {
					lengths[j].Equal++
				} else if diff > 1+TOLERANCE {
					eq = false
					text = "Shorter"
					lengths[j].Shorter++
				} else {
					eq = false
					lengths[j].Longer++
				}
			}
			count++
			if !eq && exampleCount < exampleCountMax {
				exampleCount++
				fmt.Printf("\n%v: %v, %v", text, from[i], to[i])
			}
		}
	}

	if count > 0 {
		runtime.ReadMemStats(&m)
		Logs.Parameters["Note"] = note
		Logs.Parameters["Basic Grid"] = "true"
		Logs.Parameters["X-Size"] = strconv.Itoa(xSize)
		Logs.Parameters["Y-Size"] = strconv.Itoa(ySize)
		Logs.Parameters["CPU Cores"] = strconv.Itoa(runtime.NumCPU())
		Logs.Parameters["Lenght Tolerance (%)"] = strconv.FormatFloat(TOLERANCE*100, 'f', -1, 64)
		Logs.Parameters["Run Count"] = strconv.Itoa(count)

		for j := 0; j < 5; j++ {
			var a Algo
			a.Time = Time{
				Sum:         sum[j].String(),
				AVG:         (sum[j] / time.Duration(count)).String(),
				Min:         min[j].String(),
				Max:         max[j].String(),
				TimesFaster: math.Floor((float64(sum[0]/time.Duration(count))/float64(sum[j]/time.Duration(count)))*100) / 100,
			}
			a.PQPops = PQPops{
				AVG:     poppedSum[j] / count,
				Percent: math.Floor((float64(poppedSum[j]/count)/float64(poppedSum[0]/count))*10000) / 100,
			}
			a.Length = lengths[j]

			Logs.Results[algoStrPrint[j]] = a
		}

		for j := 0; j < 5; j++ {
			jsT, _ := json.MarshalIndent(Logs.Results[algoStrPrint[j]].Time, "", "    ")
			jsL, _ := json.MarshalIndent(Logs.Results[algoStrPrint[j]].Length, "", "    ")
			jsP, _ := json.MarshalIndent(Logs.Results[algoStrPrint[j]].PQPops, "", "    ")

			fmt.Printf("\n\n%v\n", algoStrPrint[j])
			fmt.Printf("Time:\n%v\n", string(jsT))
			fmt.Printf("Length:\n%v\n", string(jsL))
			fmt.Printf("PQ-Pops:\n%v\n", string(jsP))
		}

		jsonString, _ := json.MarshalIndent(Logs, "", "    ")
		timestamp = time.Now().Format("2006-01-02_15-04-05")
		outFilename = fmt.Sprintf("data/evaluation/wf_%s_%s_bg_%s.json", Logs.Parameters["X-Size"], Logs.Parameters["Y-Size"], timestamp)
		saveLog(outFilename, jsonString)
	}
}

// ReadPBF is evaluated
func ReadPBF(pbfFileName, note string) {
	var filename string
	var coastlineMap = make(map[int64][]int64)
	var nodeMap = make(map[int64][]float64)
	var readTime string
	var readLessMemoryTime string
	var m runtime.MemStats
	var totalAllocNormal uint64
	var timestamp string

	var logs = make(map[string]string)
	logs["pbfFileName"] = pbfFileName
	pbfFileName = fmt.Sprintf("data/%s", pbfFileName)

	fmt.Printf("\nNormal:\n")
	readTime = dataprocessing.ReadFile(pbfFileName, &coastlineMap, &nodeMap)
	runtime.ReadMemStats(&m)
	totalAllocNormal = m.TotalAlloc

	fmt.Printf("\nLess memory:\n")
	coastlineMap = make(map[int64][]int64)
	nodeMap = make(map[int64][]float64)
	readLessMemoryTime = dataprocessing.ReadFileLessMemory(pbfFileName, &coastlineMap, &nodeMap)
	runtime.ReadMemStats(&m)

	logs["Note"] = note
	logs["CPU Cores"] = strconv.Itoa(runtime.NumCPU())
	logs["Mem alloc normal"] = strconv.FormatUint(totalAllocNormal/1024/1024, 10)
	logs["Mem alloc less memory"] = strconv.FormatUint((m.TotalAlloc-totalAllocNormal)/1024/1024, 10)
	logs["Time"] = string(readTime)
	logs["Time less memory"] = string(readLessMemoryTime)

	jsonString, _ := json.MarshalIndent(logs, "", "    ")
	timestamp = time.Now().Format("2006-01-02_15-04-05")
	filename = fmt.Sprintf("data/evaluation/rf_%s_%s.json", strings.Split(logs["pbfFileName"], ".")[0], timestamp)
	saveLog(filename, jsonString)
}

// DataProcessing is evaluated
func DataProcessing(pbfFileName, note string, xSize, ySize int, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon bool) {
	var logs = dataprocessing.Start(pbfFileName, xSize, ySize, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)
	var filename string
	var m runtime.MemStats
	var timestamp string

	runtime.ReadMemStats(&m)
	logs["PBF filename"] = pbfFileName
	logs["Note"] = note
	logs["Basic Grid"] = strconv.FormatBool(basicGrid)
	logs["Basic point in polygon"] = strconv.FormatBool(basicPointInPolygon)
	logs["X-Size"] = strconv.Itoa(xSize)
	logs["Y-Size"] = strconv.Itoa(ySize)
	logs["CPU Cores"] = strconv.Itoa(runtime.NumCPU())
	logs["Mem alloc"] = strconv.FormatUint(m.TotalAlloc/1024/1024, 10)

	jsonString, _ := json.MarshalIndent(logs, "", "    ")
	timestamp = time.Now().Format("2006-01-02_15-04-05")
	if val, ok := logs["filename"]; ok {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s%s_%s.json", strings.Split(logs["pbfFileName"], ".")[0], logs["X-Size"], logs["Y-Size"], val, timestamp)
	} else {
		filename = fmt.Sprintf("data/evaluation/dp_%s_%s_%s_%s.json", strings.Split(logs["pbfFileName"], ".")[0], logs["X-Size"], logs["Y-Size"], timestamp)
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
	var logs = make(map[string]string)
	var timestamp string

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

	logs["Note"] = note
	logs["X-Size"] = strconv.Itoa(xSize)
	logs["Y-Size"] = strconv.Itoa(ySize)
	logs["CPU Cores"] = strconv.Itoa(runtime.NumCPU())
	logs["AVG time Simple Neighbours"] = (sumSN / time.Duration(ug.N)).String()
	logs["AVG time Neighbours"] = (sumN / time.Duration(ug.N)).String()
	logs["Times faster"] = strconv.FormatFloat(math.Floor((float64(sumN)/float64(sumSN))*100)/100, 'f', -1, 64)
	logs["Errors"] = strconv.Itoa(len(failedIDS))
	logs["Errors (%)"] = strconv.FormatFloat(math.Floor(float64(len(failedIDS))/float64(k)*10000)/100, 'f', -1, 64)

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

	jsonString, _ := json.MarshalIndent(logs, "", "    ")
	timestamp = time.Now().Format("2006-01-02_15-04-05")
	filename = fmt.Sprintf("data/evaluation/nb_%s_%s_%s.json", logs["X-Size"], logs["Y-Size"], timestamp)
	saveLog(filename, jsonString)

}

// Time log
type Time struct {
	Sum         string
	AVG         string
	Min         string
	Max         string
	TimesFaster float64
}

// Length log
type Length struct {
	Shorter int
	Equal   int
	Longer  int
}

// PQPops log
type PQPops struct {
	AVG     int
	Percent float64
}

// Algo log
type Algo struct {
	Time   Time
	Length Length
	PQPops PQPops
}

// Log structure
type Log struct {
	Parameters map[string]string
	Results    map[string]Algo
}

var algoStrPrint = []string{"Dijkstra", "A-Star", "Bi-Dijkstra", "Bi-A-Star", "AStar-JPS"}

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
