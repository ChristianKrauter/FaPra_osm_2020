package main
import (
	"os"
  "log"
  "io"
  "fmt"
  "runtime"
  //"github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
)

func main(){
	f, err := os.Open("../data/antarctica-latest.osm.pbf")
	if err != nil {
	    log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
	    log.Fatal(err)
	}

	var nodeIDs []int64
	//var coordinateList [][]float64

	var nc, wc, rc, cc uint64
	for {
	    if v, err := d.Decode(); err == io.EOF {
	        break
	    } else if err != nil {
	        log.Fatal(err)
	    } else {
	        switch v := v.(type) {
	        case *osmpbf.Node:
	            // Process Node v.
	            nc++
	        case *osmpbf.Way:
	            // Process Way v.
	            wc++
	            for _, e := range v.Tags {
	            	if e == "coastline" {
	            		cc++
			            fmt.Printf("%d", v.NodeIDs)
			            nodeIDs = append(nodeIDs, v.NodeIDs...)
	            	}
	            }
	        case *osmpbf.Relation:
	            // Process Relation v.
	            rc++
	        default:
	            log.Fatalf("unknown type %T\n", v)
	        }
	    }
	    if nc == 3 {
	    	break
	    }
	}


/*
	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
	    log.Fatal(err)
	}

	var nc, wc, rc, cc uint64
	for {
	    if v, err := d.Decode(); err == io.EOF {
	        break
	    } else if err != nil {
	        log.Fatal(err)
	    } else {
	        switch v := v.(type) {
	        case *osmpbf.Node:
	            // Process Node v.
	            nc++
	           	if  v.ID
							coordinateList = append(coordinateList,[]float64{v.Lon, v.Lat})
	        case *osmpbf.Way:
	            // Process Way v.
	            wc++
	        case *osmpbf.Relation:
	            // Process Relation v.
	            rc++
	        default:
	            log.Fatalf("unknown type %T\n", v)
	        }
	    }
	    if nc == 3 {
	    	break
	    }
	}*/






















  //fc := geojson.NewFeatureCollection()
  //fc.AddFeature(geojson.NewLineStringFeature(coordinateList))
  //rawJSON, err := fc.MarshalJSON()
  //fmt.Printf("%s", rawJSON)

	fmt.Printf("\nNodes: %d, Ways: %d, Relations: %d, Coastlines: %d\n", nc, wc, rc, cc)
}