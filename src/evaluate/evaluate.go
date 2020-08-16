package evaluate

import (
	"../dataprocessing"
	//"../algorithms"
	//"../server"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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
