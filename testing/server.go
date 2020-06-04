package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"net/http/httputil"
)

var port int = 8081

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("loaded something\n")
		// Save a copy of this request for debugging.
		if strings.Contains(r.URL.Path, "/point") {
			/*requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
			  fmt.Println(err)
			}
			fmt.Println(string(requestDump))*/

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

			w.Write([]byte("test response"))

		} else {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	})
	var portStr = fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
