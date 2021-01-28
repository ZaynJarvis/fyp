package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

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
	svc      map[string][]*api.Service
	watchers map[string][]Watcher
}

func (r *Registry) Register(ctx context.Context, service *api.Service) (*api.Result, error) {
	if service == nil || service.Name == "" || service.Endpoint == "" {
		return nil, errors.New("invalid request")
	}
	for _, svc := range r.svc[service.Name] {
		if svc.Endpoint == service.Endpoint {
			return nil, errors.New("duplicate endpoint update")
		}
	}
	r.svc[service.Name] = append(r.svc[service.Name], service)
	mappedW, ok := r.watchers[service.Name]
	if !ok {
		return &api.Result{Status: 0, Message: "OK. no watcher"}, nil
	}
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
	if _, ok := r.watchers[service.Name]; !ok {
		r.watchers[service.Name] = make([]Watcher, 0, 1)
	}
	for _, svc := range r.svc[service.Name] {
		if err := stream.Send(svc); err != nil {
			return fmt.Errorf("initialize error, err: %w", err)
		}
	}
	r.watchers[service.Name] = append(r.watchers[service.Name], func(service *api.Service) error {
		if err := stream.Send(service); err != nil {
			return err
		}
		return nil
	})
	<-stream.Context().Done()
	return nil
}

func New() api.RegistryServer {
	return &Registry{
		svc:      make(map[string][]*api.Service),
		watchers: make(map[string][]Watcher),
	}
}
