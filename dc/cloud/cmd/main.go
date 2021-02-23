package main

import (
	"log"
	"time"

	"github.com/zaynjarvis/fyp/dc/api"
	"github.com/zaynjarvis/fyp/dc/cloud/internal/transport"
)

func main() {
	t := transport.New(":7890", true)
	go t.Start()
	ch := t.RecvNotification()
	go func() {
		for e := range ch {
			log.Printf("received %#v", e.Message)
		}
	}()
	time.Sleep(5 * time.Second)
	log.Println("sending config")
	t.SendConfig(&api.CollectionConfig{
		Version:             "1",
		Service:             "2",
		ObjectStoragePath:   "abc.com",
		DocumentStoragePath: "opq.com",
		TextIndexPath:       "xyz.com",
	})
	select {}
}
