package main

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/agent/internal/storage"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	s := storage.New(storage.Config{
		ObjStoreAddr:  "localhost:9000",
		DataStoreAddr: "localhost:27017",
	})
	file, err := ioutil.ReadFile("./img.png")
	if err != nil {
		logrus.Fatal(err)
	}
	if err := s.Image("2.png", "image/png", file); err != nil {
		logrus.Error(err)
	}
	type Result struct {
		Code int
		Msg  string
	}
	if err := s.Data("2.png", Result{
		Code: 200,
		Msg:  "OK",
	}); err != nil {
		logrus.Error(err)
	}

	//t := transport.New(":7890", &api.AgentInfo{Id: os.Args[1], Service: "test"}, true)
	//go t.Start()
	//go func() {
	//	for cfg := range t.RecvConfig() {
	//		logrus.Info(cfg)
	//	}
	//}()
	//s := server.New(":9000")
	//go s.Start()
	//defer s.Stop()
	//imgCh := s.RecvImageReport()
	//p := process.New(imgCh, t)
	//p.Execute()
}
