package server

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/zaynjarvis/fyp/dc/api"

	"github.com/sirupsen/logrus"
)

func ListenConfig(addr string) chan *api.CollectionConfig {
	ch := make(chan *api.CollectionConfig)
	go func() {
		if err := http.ListenAndServe(addr, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer request.Body.Close()
			var cfg api.CollectionConfig
			if err := jsonpb.Unmarshal(request.Body, &cfg); err != nil {
				logrus.Error("unmarshal failed, err: ", err)
			}
			ch <- &cfg
		})); err != nil {
			logrus.Info("server closed, err: ", err)
		}
	}()
	return ch
}
