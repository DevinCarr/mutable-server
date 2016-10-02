package main

import (
	"fmt"
	"net/http"
	"net/url"
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
	resp, err := http.PostForm(address,
		url.Values{"consumerId": {i.consumerId}})

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
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
