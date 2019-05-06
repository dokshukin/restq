package main

import (
	"time"
)

type message struct {
	Body     string `json:"body"`
	UUID     string `json:"uuid"`
	TTL      int64  `json:"ttl"`
	Status   uint8
	Created  time.Time
	Modified time.Time
	Expires  time.Time
}

var (
	// QueueList is global list of all queues
	QueueList map[string][]message
)

func init() {
	QueueList = make(map[string][]message)
}

func main() {
	// All configration variables should be set from ENV (easy-to-use docker and systemd)

	// run routine to save data on disk
	go dbHandler()

	// run queue cleaner
	go queueCleaner()

	// run rest API with GIN
	restAPI()
}
