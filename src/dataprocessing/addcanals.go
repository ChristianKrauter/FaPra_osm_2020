package dataprocessing

import (
	"../algorithms"
	"../grids"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var fromCoord = [][]float64{
	{-79.35200520015435188270, 8.39506599016100807376},
	{-6.40559952102847685040, 35.95154462115006310796},
	{24.90517595851296661635, 40.22370962194782606502},
	{31.86675754180648212355, 31.99357994086424383795},
}
var toCoord = [][]float64{
	{-80.14068065781292204974, 9.89468610374165891130},
	{-4.77168700638261888969, 36.01010019501850933921},
	{29.37309462345675115102, 41.59221573205132216344},
	{34.38728084670030682446, 27.24810015724757406019},
}

// AddCanals to uniform grid
func AddCanals(xSize, ySize int, basicPointInPolygon bool) {
	var filename string
	var ug grids.UniformGrid
	var from int
	var to int

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/uniformgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/uniformgrid_%v_%v.json", xSize, ySize)
	}

	uniformgridRaw, errJSON := os.Open(filename)
	if errJSON != nil {
		log.Fatal(fmt.Sprintf("\nThe meshgrid '%s'\ncould not be found. Please create it first.\n", filename))
	}
	defer uniformgridRaw.Close()
	byteValue, _ := ioutil.ReadAll(uniformgridRaw)
	json.Unmarshal(byteValue, &ug)

	for i := 0; i < len(fromCoord); i++ {
		from = ug.GridToID(ug.CoordToGrid(fromCoord[i][0], fromCoord[i][1]))
		to = ug.GridToID(ug.CoordToGrid(toCoord[i][0], toCoord[i][1]))
		fmt.Printf("\nfrom: %v\n", from)
		fmt.Printf("to: %v\n", to)

		fmt.Printf("Route:\n")
		var route *[][][]float64
		route, _, _ = algorithms.DijkstraCanal(from, to, &ug)
		for j := 0; j < len(*route); j++ {
			for k := 0; k < len((*route)[j]); k++ {
				grid := ug.CoordToGrid((*route)[j][k][0], (*route)[j][k][1])
				id := ug.GridToID(grid)
				fmt.Printf("%v\n", id)
				ug.VertexData[grid[0]][grid[1]] = false
			}
		}

	}
	_ = storeUniformGrid(&ug, filename)
}

// AddCanalsBg to basic grid
func AddCanalsBg(xSize, ySize int, basicPointInPolygon bool) {
	var filename string
	var bg grids.BasicGrid
	var bg2D [][]bool
	var from int
	var to int
	bg.XSize = xSize
	bg.YSize = ySize
	bg.XFactor = float64(xSize) / 360.0
	bg.YFactor = float64(ySize) / 360.0

	if basicPointInPolygon {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v_bpip.json", xSize, ySize)
	} else {
		filename = fmt.Sprintf("data/output/meshgrid_%v_%v.json", xSize, ySize)
	}

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

	for i := 0; i < len(fromCoord); i++ {
		from = bg.GridToID(bg.CoordToGrid(fromCoord[i]))
		to = bg.GridToID(bg.CoordToGrid(toCoord[i]))
		fmt.Printf("\nfrom: %v\n", from)
		fmt.Printf("to: %v\n", to)

		fmt.Printf("Route:\n")
		var route [][][]int
		route, _, _ = algorithms.DijkstraCanalBg(from, to, &bg)
		for j := 0; j < len(route); j++ {
			for k := 0; k < len(route[j]); k++ {
				grid := route[j][k]
				id := bg.GridToID(grid)
				fmt.Printf("%v\n", id)
				bg2D[grid[0]][grid[1]] = false
			}
		}

	}
	_ = storeMeshgrid(&bg2D, filename)
}
