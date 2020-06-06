package main

import (
	"fmt"
	"github.com/paulmach/go.geojson"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"net/http/httputil"
)

var port int = 8081

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toGeojson(route [][]float64) []byte {
	var rawJson []byte
	g := geojson.NewLineStringGeometry(route)
	rawJson, err4 := g.MarshalJSON()
	check(err4)
	return rawJson
}

func dijkstra(startLng float64, startLat float64, endLng float64, endLat float64) [][]float64 {
	var route [][]float64
	route = append(route, []float64{startLng, startLat})
	route = append(route, []float64{endLng, endLat})

	return route
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/point") {
			query := r.URL.Query()
			var startLat, err = strconv.ParseFloat(query.Get("startLat"), 10)
			if err != nil {
				panic(err)
			}
			var startLng, err1 = strconv.ParseFloat(query.Get("startLng"), 10)
			if err1 != nil {
				panic(err1)
			}
			var endLat, err2 = strconv.ParseFloat(query.Get("endLat"), 10)
			if err2 != nil {
				panic(err2)
			}
			var endLng, err3 = strconv.ParseFloat(query.Get("endLng"), 10)
			if err3 != nil {
				panic(err3)
			}

			fmt.Printf("start: %d / %d ", int(startLat), int(startLng))
			fmt.Printf("end: %d / %d\n", int(endLat), int(endLng))

			var route = dijkstra(startLng, startLat, endLng, endLat)
			var result = toGeojson(route)

			w.Write(result)
		} else {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	})

	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
