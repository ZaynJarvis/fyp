package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/sdk"
)

type Result struct {
	Conf float64 `json:"confidence"`
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	s := sdk.New("localhost:9000", "example", sdk.WithTag("mode", "test"))
	go func() {
		if err := s.Connect(context.Background()); err != nil {
			logrus.Error("example failed, err: ", err)
		}
	}()
	s.Image([]byte("img"), map[string]float64{"confidence": 0.8}, sdk.Tag{K: "imageID", V: "001"})
	time.Sleep(time.Millisecond)
	s.Image([]byte("img"), Result{Conf: 0.99}, sdk.Tag{K: "imageID", V: "002"})
	select {}
}
