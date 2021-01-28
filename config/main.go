package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"sync"

	"github.com/zaynjarvis/fyp/config/api"
	"google.golang.org/grpc"
)

var (
	address = flag.String("addr", "localhost:3700", "address of the config center")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *address)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterConfigCenterServer(grpcServer, New())
	if grpcServer.Serve(lis) != nil {
		return
	}
}

type ServiceID struct {
	name    string
	version uint32
}

type Watcher func(service *api.ServiceConfig) error
type ConfigCenter struct {
	api.UnimplementedConfigCenterServer
	configs  sync.Map
	watchers sync.Map
}

func (r *ConfigCenter) Set(ctx context.Context, cfg *api.ServiceConfig) (*api.Result, error) {
	if cfg == nil || cfg.Name == "" {
		return nil, errors.New("invalid request")
	}
	id := ServiceID{name: cfg.Name, version: cfg.Version}
	if cfg.Config == nil {
		r.configs.LoadAndDelete(id)
	} else {
		r.configs.Store(id, cfg)
	}

	mappedWI, ok := r.watchers.Load(id)
	if !ok {
		return &api.Result{Status: 0, Message: "OK. no watcher"}, nil
	}
	mappedW := mappedWI.([]Watcher)
	for _, w := range mappedW {
		if err := w(cfg); err != nil {
			log.Printf("cfg update not send, err: %v\n", err)
		}
	}
	return &api.Result{Status: 0, Message: "OK"}, nil
}

func (r *ConfigCenter) Watch(service *api.ServiceConfig, stream api.ConfigCenter_WatchServer) error {
	if service == nil || service.Name == "" {
		return errors.New("invalid request")
	}
	id := ServiceID{name: service.Name, version: service.Version}
	cfg, ok := r.configs.Load(id)
	if ok {
		if err := stream.Send(cfg.(*api.ServiceConfig)); err != nil {
			return err
		}
	}

	loadWI, _ := r.watchers.LoadOrStore(id, make([]Watcher, 0, 1))
	loadW := loadWI.([]Watcher)
	loadW = append(loadW, func(service *api.ServiceConfig) error {
		if err := stream.Send(service); err != nil {
			return err
		}
		return nil
	})
	r.watchers.Store(id, loadW)

	<-stream.Context().Done()
	r.watchers.LoadAndDelete(id)
	return nil
}

func New() api.ConfigCenterServer {
	return &ConfigCenter{}
}
