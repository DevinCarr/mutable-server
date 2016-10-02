package main

import (
	"fmt"
	"net/http"
	"time"
)

type Items struct {
	consumerId     string
	consumerSecret string
	userId         string
	expire         int64
}

type TimeQueue struct {
	count uint64
}

// NewTimeQueue creates a new TimeQueue object
func NewTimeQueue() *TimeQueue {
	return &TimeQueue{
		count: 0,
	}
}

// wait will wait for the duration of the store and then call the
// provided url.
func wait(tq *TimeQueue, i *Items, address string) {
	time.Sleep(time.Duration(i.expire) * time.Minute)
	req, _ := http.NewRequest("POST", UnmuteAddress, nil)

	q := req.URL.Query()
	q.Add("consumerId", i.consumerId)
	q.Add("consumerSecret", i.consumerSecret)
	q.Add("userId", i.userId)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Printf("Unmute: address: %s\nstatus: %s\n", req.URL.RawQuery, resp.Status)
	tq.Done()
}

// Push puts a new item in the queue with a time d in minutes
func (tq *TimeQueue) Push(i *Items, address string) {
	tq.count += 1
	fmt.Println("Put new item")
	go wait(tq, i, address)
}

// Done will remove the item in the queue and return the count of
// remaining items.
func (tq *TimeQueue) Done() uint64 {
	tq.count -= 1
	return tq.count
}

// Count returns the current amount of items in the queue
func (tq *TimeQueue) Count() uint64 {
	return tq.count
}
