package main

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/sdk"
)

type Result struct {
	Conf float64 `json:"confidence"`
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	s := sdk.New("localhost:8000", "example", sdk.WithTag("mode", "test"))
	go func() {
		if err := s.Connect(context.Background()); err != nil {
			logrus.Error("example failed, err: ", err)
		}
	}()
	img, err := ioutil.ReadFile("./test/assets/img.png")
	if err != nil {
		logrus.Error(err)
		return
	}
	s.Image("a.png", img, map[string]float64{"confidence": 0.8}, sdk.Tag{K: "extra", V: "001"})
	time.Sleep(time.Millisecond)
	s.Image("xyzw.png", img, Result{Conf: 0.99}, sdk.Tag{K: "extra", V: "002"})
	select {}
}
