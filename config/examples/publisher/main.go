package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"

	"github.com/zaynjarvis/fyp/config/api"
	"google.golang.org/grpc"
)

type Config struct {
	Name    string
	Version int
	Desc    string
}

func main() {
	var (
		cfg            = Config{Name: "service-name", Version: 1, Desc: "hello world"}
		name           = "default"
		version uint32 = 1
		addr           = flag.String("addr", "localhost:3700", "address of config center")
	)
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := api.NewConfigCenterClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	data, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	res, err := client.Set(ctx, &api.ServiceConfig{Name: name, Version: version, Config: data})
	if err != nil {
		panic(err)
	}
	log.Println(res)
}
