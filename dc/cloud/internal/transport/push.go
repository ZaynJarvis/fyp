package transport

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/zaynjarvis/fyp/dc/api"
	"google.golang.org/grpc"
)

// should apply rate limiter in upstream, mtx is only used for precaution
type server struct {
	mtx sync.Mutex
	api.UnimplementedAgentPushServiceServer
	retainedCfg *api.CollectionConfig
	notyCh      chan *api.CollectionEvent
	agents      sync.Map // map[string]chan api.CollectionConfig
}

func (s *server) ListenConfig(agentInfo *api.AgentInfo, stream api.AgentPushService_ListenConfigServer) error {
	stored, loaded := s.agents.LoadOrStore(agentInfo.Id, make(chan *api.CollectionConfig))
	if loaded {
		return errors.New("agent with id " + agentInfo.Id + " already exists")
	}
	cfgCh := stored.(chan *api.CollectionConfig)
	if s.retainedCfg != nil {
		if err := stream.Send(s.retainedCfg); err != nil {
			log.Println("stream send config error: ", err)
		}
	}
	for {
		select {
		case <-stream.Context().Done():
			s.agents.Delete(agentInfo.Id)
			return errors.New("stream context is closed")
		case cfg, ok := <-cfgCh:
			if !ok {
				log.Println("config channel is closed")
				return errors.New("config channel is closed by the cloud")
			}
			if err := stream.Send(cfg); err != nil {
				log.Println("stream send config error: ", err)
			}
		}
	}
}

func (s *server) GetConfig(ctx context.Context, agentInfo *api.AgentInfo) (*api.CollectionConfig, error) {
	return s.retainedCfg, nil
}

func (s *server) SendNotification(stream api.AgentPushService_SendNotificationServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return errors.New("stream context is closed")
		default:
		}
		recv, err := stream.Recv()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		s.notyCh <- recv
	}
}

func (s *server) update(cfg *api.CollectionConfig) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.retainedCfg = cfg
	s.agents.Range(func(_, v interface{}) bool {
		ch := v.(chan *api.CollectionConfig)
		ch <- cfg
		return true
	})
}

type PushModel struct {
	port        string
	quit        chan struct{}
	retainedCfg *api.CollectionConfig
	svr         *server
}

func newPushModel(port string) *PushModel {
	return &PushModel{
		port: port,
		quit: make(chan struct{}),
		svr: &server{
			notyCh: make(chan *api.CollectionEvent),
		},
	}
}

func (p PushModel) Start() {
	lis, err := net.Listen("tcp", p.port)
	if err != nil {
		log.Fatal("cannot listen on port, err: ", err)
	}
	s := grpc.NewServer()
	api.RegisterAgentPushServiceServer(s, p.svr)
	go func() {
		<-p.quit
		s.GracefulStop()
	}()
	if err := s.Serve(lis); err != nil {
		log.Println(err)
	}
}

func (p PushModel) Stop() {
	close(p.quit)
}

func (p PushModel) RecvNotification() <-chan *api.CollectionEvent {
	return p.svr.notyCh
}

func (p PushModel) SendConfig(cfg *api.CollectionConfig) {
	p.retainedCfg = cfg
	p.svr.update(cfg)
}
