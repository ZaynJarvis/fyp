package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/zaynjarvis/fyp/registry/api"
	"google.golang.org/grpc"
)

var (
	address = flag.String("addr", "localhost:3900", "address of the registry")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *address)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterRegistryServer(grpcServer, New())
	if grpcServer.Serve(lis) != nil {
		return
	}
}

type Watcher func(service *api.Service) error
type Registry struct {
	api.UnimplementedRegistryServer
	svc      sync.Map
	watchers sync.Map
}

func (r *Registry) Register(ctx context.Context, service *api.Service) (*api.Result, error) {
	if service == nil || service.Name == "" || service.Endpoint == "" {
		return nil, errors.New("invalid request")
	}
	svcI, _ := r.svc.LoadOrStore(service.Name, make([]*api.Service, 0, 1))
	store := svcI.([]*api.Service)
	for _, svc := range store {
		if svc.Endpoint == service.Endpoint {
			return nil, errors.New("duplicate endpoint update")
		}
	}
	store = append(store, service)
	r.svc.Store(service.Name, store)

	mappedWI, ok := r.watchers.Load(service.Name)
	if !ok {
		return &api.Result{Status: 0, Message: "OK. no watcher"}, nil
	}
	mappedW := mappedWI.([]Watcher)
	for _, w := range mappedW {
		if err := w(service); err != nil {
			log.Printf("service update not send, err: %v\n", err)
		}
	}
	return &api.Result{Status: 0, Message: "OK"}, nil
}

func (r *Registry) Watch(service *api.Service, stream api.Registry_WatchServer) error {
	if service == nil || service.Name == "" {
		return errors.New("invalid request")
	}
	load, ok := r.svc.Load(service.Name)
	if ok {
		for _, svc := range load.([]*api.Service) {
			if err := stream.Send(svc); err != nil {
				return fmt.Errorf("initialize error, err: %w", err)
			}
		}
	}

	loadWI, _ := r.watchers.LoadOrStore(service.Name, make([]Watcher, 0, 1))
	loadW := loadWI.([]Watcher)
	loadW = append(loadW, func(service *api.Service) error {
		if err := stream.Send(service); err != nil {
			return err
		}
		return nil
	})
	r.watchers.Store(service.Name, loadW)

	<-stream.Context().Done()
	r.watchers.LoadAndDelete(service.Name)
	return nil
}

func New() api.RegistryServer {
	return &Registry{}
}
