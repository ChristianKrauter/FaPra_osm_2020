package main

import ("fmt"
		"os"
		"log"
		"runtime"
		"io"
		"time"
		"github.com/qedus/osmpbf"
		//"github.com/paulmach/go.geojson"
		)

func get_some_key(m map[int64]*osmpbf.Way) int64 {
    for k := range m {
        return k
    }
    return 0
}

func main() {
	start := time.Now()
	f, err := os.Open("antarctica-latest.osm.pbf")
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

var test map[int64]*osmpbf.Way
test = make(map[int64]*osmpbf.Way)
var coastlines []*osmpbf.Way
var nodeIDs []int64
var nodes []*osmpbf.Node
var nc, wc, rc uint64
for {
    if v, err := d.Decode(); err == io.EOF {
        break
    } else if err != nil {
        log.Fatal(err)
    } else {
        switch v := v.(type) {
        case *osmpbf.Node:
        	nodes = append(nodes,v)
            // Process Node v.
            nc++
        case *osmpbf.Way:
            // Process Way v.
            for _,value := range v.Tags {
            	if(value == "coastline"){
            		test[v.NodeIDs[0]] = v
            		coastlines = append(coastlines,v)
            		nodeIDs = append(nodeIDs,v.NodeIDs[:len(v.NodeIDs)-1]...)
            		wc++
            	}
            }
        case *osmpbf.Relation:
            // Process Relation v.
            rc++
        default:
            log.Fatalf("unknown type %T\n", v)
        }
    }
}



var coastline []int64

for len(test) > 0 {
	var key = get_some_key(test)
	var nodeids = test[key].NodeIDs
	coastline = nodeids[:len(nodeids)-1]
	delete(test,key)
	key = coastline[len(coastline)-1]
	for {
		if val,ok := test[key]; ok {
    		var nodeids = val.NodeIDs
			coastline = nodeids[:len(nodeids)-1]
			delete(test,key)
			key = coastline[len(coastline)-1]
		} else {
			break
		}

	}
}
/*var coastline []*osmpbf.Node
for _,id := range nodeIDs {
	for _,node := range nodes{
		if node.ID == id {
			coastline = append(coastline,node)
			//fmt.Printf("%d\n", id)
			break
		}
	}
}*/
//var coastlineIDs []int64
//for _,x := range coastlines{
	//for _,y := range x.NodeIDs[:len(x.NodeIDs)-1]{
	//	coastlineIDs = append(coastlineIDs,y)
	//}
//}
fmt.Printf("%d\n", len(coastlines))
fmt.Printf("%d\n", len(nodeIDs))
fmt.Printf("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)
t := time.Now()
elapsed := t.Sub(start)
fmt.Printf("%s\n", elapsed)
}