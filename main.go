package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var tq = NewTimeQueue()

type Response struct {
	Status  uint64
	Content string
	Error   string
}

func callMute(i *Items) error {
	req, _ := http.NewRequest("POST", MuteAddress, nil)

	q := req.URL.Query()
	q.Add("consumerId", i.consumerId)
	q.Add("consumerSecret", i.consumerSecret)
	q.Add("userId", i.userId)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("Mute: address: %s\nstatus: %s\n", req.URL.RawQuery, resp.Status)
	resp.Body.Close()

	return nil
}

func returnResponse(w http.ResponseWriter, s string, e error) {
	res := Response{0, s, ""}
	if e != nil {
		res = Response{1, "", e.Error()}
	}

	js, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Error: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func push(w http.ResponseWriter, r *http.Request) {
	// Check to make sure the method is a post
	if r.Method != "POST" {
		fmt.Printf("INVALID METHOD: %s", r.Method)
		returnResponse(w, "", errors.New("Invalid Method"))
		return
	}

	// parse the form
	ci := r.FormValue("consumerId")
	cs := r.FormValue("consumerSecret")
	ui := r.FormValue("userId")
	ex, err := strconv.ParseInt(r.FormValue("expire"), 10, 64)

	// verify inputs
	if ci == "" || cs == "" || ui == "" || ex < 0 || err != nil {
		fmt.Printf("Error: missing payload (ci: %s, cs: %s, ui: %s)", ci, cs, ui)
		returnResponse(w, "", errors.New("Missing POST payload items"))
		return
	}

	items := &Items{ci, cs, ui, ex}

	// Mute the user
	errMute := callMute(items)
	if errMute != nil {
		returnResponse(w, "", errors.New("Unable to Mute: http error"))
		return
	}
	// Push the item in the db timequeue
	tq.Push(items, UnmuteAddress)
	returnResponse(w, "", nil)
}

func check(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("%d count", tq.Count())
	returnResponse(w, s, nil)
}

func main() {
	http.HandleFunc("/db/check", check)
	http.HandleFunc("/db/new", push)
	http.ListenAndServe(":8090", nil)
}
