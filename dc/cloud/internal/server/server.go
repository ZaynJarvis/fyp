package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/api"
)

func ListenConfig(addr string, getSvc func() []string) chan *api.CollectionConfig {
	ch := make(chan *api.CollectionConfig)
	cfgs := make(map[string]*api.CollectionConfig)
	r := mux.NewRouter()
	r.HandleFunc("/services/{svc}", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		defer request.Body.Close()
		var cfg api.CollectionConfig
		if err := jsonpb.Unmarshal(request.Body, &cfg); err != nil {
			logrus.Error("unmarshal failed, err: ", err)
			writer.WriteHeader(500)
			return
		}
		select {
		case ch <- &cfg:
		case <-time.After(time.Millisecond):
			writer.WriteHeader(503)
		}
		cfgs[cfg.Service] = &cfg
		marshal, err := json.Marshal(map[string]string{"status": "success"})
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		writer.WriteHeader(200)
		writer.Write(marshal)
	}).Methods("POST")

	r.HandleFunc("/services", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		for _, svc := range getSvc() {
			if _, ok := cfgs[svc]; !ok {
				cfgs[svc] = &api.CollectionConfig{
					Version: "0",
					Service: svc,
				}
			}
		}
		// should be unnecessary, svcs is keys
		keys := make([]string, 0, len(cfgs))
		for k := range cfgs {
			keys = append(keys, k)
		}
		marshal, err := json.Marshal(keys)
		if err != nil {
			writer.WriteHeader(500)
		}
		writer.WriteHeader(200)
		writer.Write(marshal)
	}).Methods("GET")

	r.HandleFunc("/services/{svc}", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(request)
		cfg, ok := cfgs[vars["svc"]]
		if !ok {
			writer.WriteHeader(404)
			return
		}
		marshal, err := json.Marshal(cfg)
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		writer.WriteHeader(200)
		writer.Write(marshal)
	}).Methods("GET")

	go http.ListenAndServe(addr, cors.Default().Handler(r))

	return ch
}
