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
		"Select Mode:\n  0: Start server \n  1: Start dataprocessing \n  2: Evaluate data processing\n  3: Evaluate wayfinding")

	var pbfFileName, note string
	var xSize, ySize int
	var createTestGeoJSON, createCoastlineGeoJSON bool

	var lessMemory bool
	var noBoundingTree bool
	var basicGrid bool
	var basicPointInPolygon bool

	flag.StringVar(&pbfFileName, "f", "antarctica-latest.osm.pbf", "Name of the pbf file inside data/")
	flag.IntVar(&xSize, "x", 360, "Meshgrid size in x direction.")
	flag.IntVar(&ySize, "y", 360, "Meshgrid size in y direction.")
	flag.BoolVar(&createTestGeoJSON, "test", false, "Create test geoJSON?")
	flag.BoolVar(&createCoastlineGeoJSON, "coastline", false, "Create coastline geoJSON?")
	flag.StringVar(&note, "n", "", "Additional note for evaluations.")

	flag.BoolVar(&lessMemory, "lm", false, "Use memory efficient method to read unpruned pbf files.")
	flag.BoolVar(&noBoundingTree, "nbt", false, "Do not use a tree structure for the bounding boxes.")
	flag.BoolVar(&basicGrid, "bg", false, "Create a basic (non-unidistant) grid.")
	flag.BoolVar(&basicPointInPolygon, "bpip", false, "Use a simple 2D point in polygon test.")
	flag.Parse()

	var info string
	if basicGrid {
		info += fmt.Sprintf("basic %vx%v grid ", xSize, ySize)
	} else {
		info += fmt.Sprintf("uniform %vx%v grid ", xSize, ySize)
	}
	if basicPointInPolygon {
		info += "\nwith the simple point in polygon test"
	} else {
		info += "\nwith the advanced point in polygon test"
	}

	switch mode {
	case 0:
		fmt.Printf("Starting osmGW server with a %s", info)
		if basicGrid {
			server.Run(xSize, ySize, basicPointInPolygon)
		} else {
			server.RunUnidistant(xSize, ySize, basicPointInPolygon)
		}
	case 1:
		fmt.Printf("Starting data processing for a %s", info)
		dataprocessing.Start(pbfFileName, xSize, ySize, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)
	case 2:
		fmt.Printf("Starting evaluation of data processing for %s", info)
		evaluate.DataProcessing(pbfFileName, note, xSize, ySize, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)
	case 3:
		fmt.Printf("Evaluation of wayfinding not implemented")
	default:
		fmt.Printf("Error: No mode %d specified", mode)
	}
}
