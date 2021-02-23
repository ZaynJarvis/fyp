package transport

import (
	"context"
	"io"
	"log"

	"github.com/zaynjarvis/fyp/dc/api"
	"google.golang.org/grpc"
)

type PushModel struct {
	cloudAddr string
	info      *api.AgentInfo
	quit      chan struct{}
	configCh  chan *api.CollectionConfig
	eventCh   chan *api.CollectionEvent
}

func newPushModel(addr string, info *api.AgentInfo) *PushModel {
	return &PushModel{
		cloudAddr: addr,
		info:      info,
		quit:      make(chan struct{}),
		configCh:  make(chan *api.CollectionConfig),
		eventCh:   make(chan *api.CollectionEvent),
	}
}

func (p *PushModel) Start() {
	conn, err := grpc.Dial(p.cloudAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to start server err: ", err)
	}
	defer conn.Close()
	c := api.NewAgentPushServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		stream, err := c.ListenConfig(ctx, p.info)
		if err != nil {
			log.Fatal("cannot listen config, err: ", err)
		}
		for {
			select {
			case <-stream.Context().Done():
				return
			default:
			}
			config, err := stream.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("receive error on stream receive, err: ", err)
			}
			p.configCh <- config
		}
	}()

	go func() {
		notification, err := c.SendNotification(ctx)
		if err != nil {
			log.Println("setup send notification err: ", err)
			return
		}
		for event := range p.eventCh {
			if err := notification.Send(event); err != nil {
				log.Println("sending notification err: ", err)
			}
		}
	}()
	<-p.quit
}

func (p *PushModel) Stop() {
	close(p.quit)
}

func (p *PushModel) SendNotification(event *api.CollectionEvent) {
	p.eventCh <- event
}

func (p *PushModel) RecvConfig() <-chan *api.CollectionConfig {
	return p.configCh
}