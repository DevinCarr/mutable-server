package main

import (
	"github.com/devincarr/timequeue"
	"io"
	"net/http"
)

var tq = timequeue.NewTimeQueue()

var count = 0

func push(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
	tq.Push(0.1, 123, "http://getmuted.com/db/check")
}

func check(w http.ResponseWriter, r *http.Request) {
	count += 1
	io.WriteString(w, string(count))
}

func main() {
	http.HandleFunc("/check", check)
	http.HandleFunc("/push", push)
	http.ListenAndServe(":8090", nil)
}
