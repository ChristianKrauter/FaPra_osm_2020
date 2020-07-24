package main

import (
	"./evaluate"
	"./read"
	"./server"
	"flag"
	"fmt"
)

func main() {
	var mode int
	flag.IntVar(&mode, "m", 0,
		"Select Mode:\n  0: Test \n  1: create meshgrid\n  2: start server")

	var xSize int
	var ySize int
	var pbfFileName string
	var createTestGeoJSON bool
	var createCoastlineGeoJSON bool

	flag.StringVar(&pbfFileName, "f", "antarctica-latest.osm.pbf", "Name of the pbf file inside data/")
	flag.IntVar(&xSize, "x", 360, "Meshgrid size in x direction")
	flag.IntVar(&ySize, "y", 360, "Meshgrid size in y direction")
	flag.BoolVar(&createTestGeoJSON, "test", false, "Create test geoJSON?")
	flag.BoolVar(&createCoastlineGeoJSON, "coastline", false, "Create coastline geoJSON?")
	flag.Parse()

	switch mode {
	case 0:
		evaluate.Test()
	case 1:
		read.Main(pbfFileName)
	case 2:
		server.Run()
	default:
		fmt.Printf("Error: No mode %d specified", mode)
	}
}
