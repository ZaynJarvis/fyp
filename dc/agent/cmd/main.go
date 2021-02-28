package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/agent/internal/process"
	"github.com/zaynjarvis/fyp/dc/agent/internal/server"
	"github.com/zaynjarvis/fyp/dc/agent/internal/transport"
	"github.com/zaynjarvis/fyp/dc/api"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	t := transport.New(":7890", &api.AgentInfo{Id: os.Args[1], Service: "test"}, true)
	go t.Start()
	s := server.New(":7000")
	go s.Start()
	defer s.Stop()
	imgCh := s.RecvImageReport()
	p := process.New(imgCh, t)
	p.Execute()
}
