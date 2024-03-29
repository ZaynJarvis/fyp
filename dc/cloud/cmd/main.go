package main

import (
	"log"

	"github.com/sirupsen/logrus"

	"github.com/zaynjarvis/fyp/dc/cloud/internal/server"
	"github.com/zaynjarvis/fyp/dc/cloud/internal/transport"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	t := transport.New(":7890", true)
	go t.Start()
	notyCh := t.RecvNotification()
	go func() {
		for e := range notyCh {
			log.Printf("received %#v", e.Message)
		}
	}()
	cfgCh := server.ListenConfig(":8900", t.Services)
	for cfg := range cfgCh {
		logrus.Info("received a config and sent to agents")
		t.SendConfig(cfg)
	}
}
