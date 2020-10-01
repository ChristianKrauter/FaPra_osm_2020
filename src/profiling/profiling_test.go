package main

import (
	"../grids"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

// TestWayFinding for profiling
func TestWayFinding(t *testing.T) {
	xSize := 360
	ySize := 360
	nRuns := 100
	algorithm := 4

	var ug grids.UniformGrid
	var from = make([]int, nRuns)
	var to = make([]int, nRuns)
	filename := fmt.Sprintf("../../data/output/uniformgrid_%v_%v.json", xSize, ySize)

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		fmt.Printf("%v\n", errJSON)
		log.Fatal(fmt.Sprintf("\nThe meshgrid could not be found. Please create it first.\n"))
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &ug)

	ug1D := make([]bool, ug.N)
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
		WayFinding(from[i], to[i], algorithm, &ug)
	}
}
