package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/zaynjarvis/fyp/config/api"
	"google.golang.org/grpc"
)

func GetConfig(cfg sync.Locker, name string, version uint32, cfgCenterAddr string) error {
	conn, err := grpc.Dial(cfgCenterAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := api.NewConfigCenterClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watch, err := client.Watch(ctx, &api.ServiceConfig{Name: name, Version: version})
	if err != nil {
		return err
	}
	for {
		recv, err := watch.Recv()
		if err != nil {
			return err
		} else if recv.Config == nil {
			return errors.New("received nil config")
		}
		cfg.Lock()
		if err := json.Unmarshal(recv.Config, cfg); err != nil {
			cfg.Unlock()
			return err
		}
		cfg.Unlock()
	}
}
