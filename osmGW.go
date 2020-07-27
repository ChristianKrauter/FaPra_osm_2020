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
		"Select Mode:\n  0: Evaluate data processing \n  1: Evaluate wayfinding \n  2: Start dataprocessing\n  3: Start server")

	var note string
	var xSize int
	var ySize int
	var pbfFileName string
	var lessMemory bool
	var noBoundingTree bool
	var createTestGeoJSON bool
	var createCoastlineGeoJSON bool

	flag.StringVar(&pbfFileName, "f", "antarctica-latest.osm.pbf", "Name of the pbf file inside data/")
	flag.IntVar(&xSize, "x", 360, "Meshgrid size in x direction.")
	flag.IntVar(&ySize, "y", 360, "Meshgrid size in y direction.")
	flag.BoolVar(&createTestGeoJSON, "test", false, "Create test geoJSON?")
	flag.BoolVar(&createCoastlineGeoJSON, "coastline", false, "Create coastline geoJSON?")
	flag.StringVar(&note, "n", "", "Additional note for evaluations.")
	flag.BoolVar(&lessMemory, "lm", false, "Use memory efficient method to read unpruned pbf files.")
	flag.BoolVar(&noBoundingTree, "nbt", false, "Do not use a tree structure for the bounding boxes.")
	flag.Parse()

	switch mode {
	case 0:
		fmt.Printf("Starting evaluation of data processing")
		evaluate.DataProcessing(pbfFileName, note, xSize, ySize, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree)
	case 1:
		fmt.Printf("Evaluation of wayfinding not implemented")
	case 2:
		fmt.Printf("Starting data processing")
		dataprocessing.Start(pbfFileName, xSize, ySize, createTestGeoJSON, createCoastlineGeoJSON, lessMemory, noBoundingTree)
	case 3:
		fmt.Printf("Starting osmGW server with %dx%d grid", xSize, ySize)
		server.Run(xSize, ySize)
	default:
		fmt.Printf("Error: No mode %d specified", mode)
	}
}
