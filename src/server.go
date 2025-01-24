package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "<html><body><h1>jmischler72</h1><ul><li><a href='https://www.youtube.com/@jmischler72'>yt</a></li><li><a href='https://jmischler72.github.io/cv/'>cv</a></li></ul></body></html>")
}

func main() {
    http.HandleFunc("/", hello)
    http.ListenAndServe(":8090", nil)
}