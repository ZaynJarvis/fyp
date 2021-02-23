package main

import (
	"log"
	"os"

	"github.com/zaynjarvis/fyp/dc/api"

	"github.com/zaynjarvis/fyp/dc/agent/internal/transport"
)

func main() {
	t := transport.New(":7890", &api.AgentInfo{Id: os.Args[1], Service: "test"}, true)
	go t.Start()
	go func() {
		for cfg := range t.RecvConfig() {
			log.Printf("%s", cfg.String())
		}
	}()
	t.SendNotification(&api.CollectionEvent{
		Type:    api.ContentType_Text,
		Url:     "localhost:/api/con",
		Message: "OK",
	})
	select {}
}
