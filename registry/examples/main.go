package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/zaynjarvis/fyp/registry/api"
	"google.golang.org/grpc"
)

var (
	registry = flag.String("registry", "localhost:3900", "address of the registry")
	watch    = flag.Bool("watch", false, "act as a watcher (default to a actor)")
	name     = flag.String("name", "default", "service name")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*registry, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := api.NewRegistryClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if *watch {
		watch, err := client.Watch(ctx, &api.Service{Name: *name})
		if err != nil {
			panic(err)
		}
		for {
			recv, err := watch.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				panic(err)
			}
			log.Printf("%s\n", recv)
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := client.Register(ctx, &api.Service{Name: *name, Endpoint: fmt.Sprintf("localhost:%v", rand.Int()%65535)})
				if err != nil {
					log.Printf("client request register err: %v\n", err)
					return
				}
				log.Printf("%s\n", res.String())
			}()
		}
		wg.Wait()
	}
}
