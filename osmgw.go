package main

import (
	"./src/dataprocessing"
	"./src/evaluate"
	"./src/server"
	"./src/testserver"
	"flag"
	"fmt"
)

func main() {
	var mode int

	var pbfFileName string
	var note string
	var xSize, ySize int
	var createCoastlineGeoJSON bool

	var lessMemory bool
	var noBoundingTree bool
	var basicGrid bool
	var basicPointInPolygon bool

	var nRuns int
	var info string

	flag.IntVar(&mode, "m", 0,
		"Select Mode:\n  0: Start server \n  1: Evaluate dataprocessing\n  2: Evaluate wayfinding\n  3: Evaluate reading pbf\n  4: Evaluate ug neighbours\n  5: Test routes and neighbours\n  6: Add canals to grid")

	flag.StringVar(&pbfFileName, "f", "antarctica-latest.osm.pbf", "Name of the pbf file inside data/")
	flag.StringVar(&note, "n", "", "Additional note for evaluations.")
	flag.IntVar(&xSize, "x", 360, "Meshgrid size in x direction.")
	flag.IntVar(&ySize, "y", 360, "Meshgrid size in y direction.")
	flag.BoolVar(&createCoastlineGeoJSON, "coastline", false, "Create coastline geoJSON.")

	flag.BoolVar(&lessMemory, "lm", false, "Use memory efficient method to read unpruned pbf files.")
	flag.BoolVar(&noBoundingTree, "nbt", false, "Do not use a tree structure for the bounding boxes.")
	flag.BoolVar(&basicGrid, "bg", false, "Create a basic (non-uniform) grid.")
	flag.BoolVar(&basicPointInPolygon, "bpip", false, "Use the basic 2D point in polygon test.")

	flag.IntVar(&nRuns, "r", 1000, "Number of runs for wayfinding evaluation.")
	flag.Parse()

	if basicGrid {
		info += fmt.Sprintf("basic %vx%v grid ", xSize, ySize)
	} else {
		info += fmt.Sprintf("uniform %vx%v grid ", xSize, ySize)
	}
	if mode == 1 {
		info += fmt.Sprintf("\non the %s file", pbfFileName)
	}
	if mode == 0 || mode == 1 {
		if basicPointInPolygon {
			info += "\nwith the basic point in polygon test"
		} else {
			info += "\nwith the spherical point in polygon test"
		}
	}

	switch mode {
	case 0:
		fmt.Printf("Starting osmGW server with a %s\n", info)
		if basicGrid {
			server.RunBg(xSize, ySize, basicPointInPolygon)
		} else {
			server.Run(xSize, ySize, basicPointInPolygon)
		}
	case 1:
		fmt.Printf("Starting data processing for %s", info)
		evaluate.DataProcessing(pbfFileName, note, xSize, ySize, createCoastlineGeoJSON, lessMemory, noBoundingTree, basicGrid, basicPointInPolygon)
	case 2:
		fmt.Printf("Starting evaluation of wayfinding on %s\n", info)
		fmt.Printf("Averaging over %v random routings ", nRuns)
		if basicGrid {
			evaluate.WayFindingBg(xSize, ySize, nRuns, note)
		} else {
			evaluate.WayFinding(xSize, ySize, nRuns, note)
		}
	case 3:
		fmt.Printf("Starting evaluation of pbf reading for %s\n", pbfFileName)
		evaluate.ReadPBF(pbfFileName, note)
	case 4:
		evaluate.NeighboursUg(xSize, ySize, note)
	case 5:
		fmt.Printf("Starting test server with a %s\n", info)
		if basicGrid {
			testserver.StartBg(xSize, ySize)
		} else {
			testserver.Start(xSize, ySize)
		}
	case 6:
		fmt.Printf("Starting adding canals to %s\n", info)
		if basicGrid {
			dataprocessing.AddCanalsBg(xSize, ySize, basicPointInPolygon)
		} else {
			dataprocessing.AddCanals(xSize, ySize, basicPointInPolygon)
		}
	default:
		fmt.Printf("Error: No mode %d specified", mode)
	}
}
