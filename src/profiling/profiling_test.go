package main

import (
	"../grids"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

var xSize, ySize, algorithm, nRuns int
var basicGrid bool

func TestMain(m *testing.M) {
	flag.IntVar(&xSize, "xs", 1000, "Meshgrid size in x direction.")
	flag.IntVar(&ySize, "ys", 500, "Meshgrid size in y direction.")
	flag.IntVar(&algorithm, "alg", 0, "Select Algorithm:\n  0: Dijkstra\n  1: A*\n  2: Bi-Dijkstra\n  3: Bi-A*\n  4: A*-JPS")
	flag.IntVar(&nRuns, "r", 100, "Number of runs for wayfinding evaluation.")

	flag.Parse()
	os.Exit(m.Run())
}

// TestWayFinding for profiling
func TestWayFinding(t *testing.T) {

	fmt.Printf("Starting profilining of")

	switch algorithm {
	case 0:
		fmt.Printf(" %v Dijkstra runs ", nRuns)
	case 1:
		fmt.Printf(" %v A* runs ", nRuns)
	case 2:
		fmt.Printf(" %v Bi-Dijkstra runs ", nRuns)
	case 3:
		fmt.Printf(" %v Bi-A* runs ", nRuns)
	case 4:
		fmt.Printf(" %v A* JPS runs ", nRuns)
	}

	fmt.Printf("on the %vx%v uniform grid\n\n", xSize, ySize)

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
