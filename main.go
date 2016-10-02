package main

import (
	"io"
	"net/http"
	//	"github.com/devincarr/timequeue"
)

func push(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func main() {
	//	tq := timequeue.NewTimeQueue()
	http.HandleFunc("/", push)
	http.HandleFunc("/push", push)
	http.ListenAndServe(":8090", nil)
}
