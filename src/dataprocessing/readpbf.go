package dataprocessing

import (
	"fmt"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

// ReadFile pbf
func ReadFile(pbfFileName string, coastlineMap *map[int64][]int64, nodeMap *map[int64][]float64) string {
	start := time.Now()

	// Read coastlines
	f, err := os.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)
	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(runtime.NumCPU()))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				(*nodeMap)[v.ID] = []float64{v.Lon, v.Lat}
			case *osmpbf.Way:
				for _, value := range v.Tags {
					if value == "coastline" {
						(*coastlineMap)[v.NodeIDs[0]] = v.NodeIDs
					}
				}
			case *osmpbf.Relation:
				continue
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Read file in                     : %s\n", elapsed)
	return elapsed.String()
}

// ReadFileLessMemory pbf
func ReadFileLessMemory(pbfFileName string, coastlineMap *map[int64][]int64, nodeMap *map[int64][]float64) string {
	start := time.Now()

	// Read coastlines
	f, err := os.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)
	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(runtime.NumCPU()))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				for _, value := range v.Tags {
					if value == "coastline" {
						(*coastlineMap)[v.NodeIDs[0]] = v.NodeIDs
						for _, id := range v.NodeIDs {
							(*nodeMap)[id] = []float64{}
						}
					}
				}
			default:
				continue
			}
		}
	}

	// Read nodes
	f, err = os.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d = osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)
	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(runtime.NumCPU()))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				if _, ok := (*nodeMap)[v.ID]; ok {
					(*nodeMap)[v.ID] = []float64{v.Lon, v.Lat}
				}
			default:
				continue
			}
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Read file in                     : %s\n", elapsed)
	return elapsed.String()
}
