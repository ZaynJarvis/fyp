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
	id := os.Getenv("DRAIS_AGENT_ID")
	svc := os.Getenv("DRAIS_AGENT_SVC")
	if id == "" || svc == "" {
		panic("should configure env vars for DRAIS_AGENT_ID and DRAIS_AGENT_SVC")
	}

	t := transport.New("cloud:7890", &api.AgentInfo{Id: id, Service: svc}, true)
	go t.Start()

	s := server.New(":7000")
	go s.Start()
	defer s.Stop()
	imgCh := s.RecvImageReport()
	p := process.New(imgCh, t)
	p.Execute()
}
