package main

import (
	"encoding/json"
	"fmt"

	"github.com/enriquebris/goconcurrentqueue"
)

type OutGoingEvent struct {
	Type         string `json:"eventType"`
	JIDConcerned string `json:"jid"`
	Content      any    `json:"content"`
}

type ReturnFunc struct {
	ID      string `json:"returnId"`
	Content any    `json:"returnData"`
}

type EQ struct {
	backend *goconcurrentqueue.FIFO
}

func (queue *EQ) Enqueue(event OutGoingEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = queue.backend.Enqueue(data)
	if err != nil {
		return err
	}
	return nil
}

func (queue *EQ) Recoverer() {
	if r := recover(); r != nil {
		defer func(r1 interface{}) {
			if r := recover(); r != nil {
				fmt.Println("An Error happened while handling another error, the library is likely to break soon (would recommend checking both errors and shutting it down).", r1, r)
			}
		}(r)
		queue.Enqueue(OutGoingEvent{
			Type:    "error",
			Content: r,
		})
	}
}
