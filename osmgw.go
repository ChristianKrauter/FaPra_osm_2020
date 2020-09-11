package main

import (
	"./src/dataprocessing"
	"./src/evaluate"
	"./src/server"
	"flag"
	"fmt"
)

func main() {
	var mode int
	flag.IntVar(&mode, "m", 0,
		"Select Mode:\n  0: Start server \n  1: Start dataprocessing \n  2: Evaluate data processing\n  3: Evaluate wayfinding\n  4: Evaluate reading pbf\n  4: Evaluate & Test UG neighbours")

	var pbfFileName, note string
	var xSize, ySize int
	var createCoastlineGeoJSON bool

	var lessMemory bool
	var noBoundingTree bool
	var basicGrid bool
	var basicPointInPolygon bool

	var algo int
	var nRuns int

	flag.StringVar(&pbfFileName, "f", "antarctica-latest.osm.pbf", "Name of the pbf file inside data/")
	flag.IntVar(&xSize, "x", 360, "Meshgrid size in x direction.")
	flag.IntVar(&ySize, "y", 360, "Meshgrid size in y direction.")
	flag.BoolVar(&createCoastlineGeoJSON, "coastline", false, "Create coastline geoJSON?")
	flag.StringVar(&note, "n", "", "Additional note for evaluations.")

	flag.BoolVar(&lessMemory, "lm", false, "Use memory efficient method to read unpruned pbf files.")
	flag.BoolVar(&noBoundingTree, "nbt", false, "Do not use a tree structure for the bounding boxes.")
	flag.BoolVar(&basicGrid, "bg", false, "Create a basic (non-unidistant) grid.")
	flag.BoolVar(&basicPointInPolygon, "bpip", false, "Use a basic 2D point in polygon test.")

	flag.IntVar(&algo, "a", 0, "Select Algorithm:\n  0: Dijkstra\n  1: A*\n  2: Bi-Dijkstra\n  3: Bi-A*\n  4: A*-JPS")
	flag.IntVar(&nRuns, "r", 1000, "Number of runs for wayfinding evaluation.")
	flag.Parse()

	var info string
	if basicGrid {
		info += fmt.Sprintf("basic %vx%v grid ", xSize, ySize)
	} else {
		info += fmt.Sprintf("uniform %vx%v grid ", xSize, ySize)
	}
	if mode == 1 || mode == 2 {
		info += fmt.Sprintf("\non the %s file", pbfFileName)

	}
	if basicPointInPolygon {
		info += "\nwith the basic point in polygon test"
	} else {
		info += "\nwith the spherical point in polygon test"
	}

	switch mode {
	case 0:
		fmt.Printf("Starting osmGW server with a %s\n", info)
		if basicGrid {
			server.Run(xSize, ySize, basicPointInPolygon)
		} else {
			server.RunUnidistant(xSize, ySize, basicPointInPolygon)
		}
	case 1:
		fmt.Printf("Starting data processing for a %s", info)
		dataprocessing.Start(pbfFileName, xSize, ySize, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)
	case 2:
		fmt.Printf("Starting evaluation of data processing for %s", info)
		evaluate.DataProcessing(pbfFileName, note, xSize, ySize, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)
	case 3:
		fmt.Printf("Starting evaluation of wayfinding for %s\n", info)
		fmt.Printf("Averaging over %v routings ", nRuns)
		if basicGrid {
			evaluate.WayFindingBG(xSize, ySize, nRuns, algo, basicPointInPolygon, note)
		} else {
			evaluate.WayFinding(xSize, ySize, nRuns, algo, basicPointInPolygon, note)
		}
	case 4:
		fmt.Printf("Starting evaluation of pbf reading for %s\n", pbfFileName)
		evaluate.ReadPBF(pbfFileName, note)
	case 5:
		evaluate.NeighboursUg(xSize, ySize, note)
	default:
		fmt.Printf("Error: No mode %d specified", mode)
	}
}
