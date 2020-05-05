package main

import (
    "fmt"
    "log"
    "net/http"
)

var port int = 8081

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
    var portStr = fmt.Sprintf(":%d", port)
    fmt.Printf("Starting server on %s", portStr)
    log.Fatal(http.ListenAndServe(portStr, nil))
}