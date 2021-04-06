package transport

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/zaynjarvis/fyp/dc/api"
	"google.golang.org/grpc"
)

// should apply rate limiter in upstream, mtx is only used for precaution
type server struct {
	mtx sync.Mutex
	api.UnimplementedAgentPushServiceServer
	retainedCfg map[string]*api.CollectionConfig
	svcs        map[string]struct{}
	notyCh      chan *api.CollectionEvent
	agents      sync.Map // map[string]chan api.CollectionConfig
}

func (s *server) ListenConfig(agentInfo *api.AgentInfo, stream api.AgentPushService_ListenConfigServer) error {
	stored, loaded := s.agents.LoadOrStore(agentInfo.Id, make(chan *api.CollectionConfig))
	if loaded {
		return errors.New("agent with id " + agentInfo.Id + " already exists")
	}
	logrus.Infof("received connection from %s", agentInfo)
	s.svcs[agentInfo.Service] = struct{}{}
	cfgCh := stored.(chan *api.CollectionConfig)
	for k, v := range s.retainedCfg {
		logrus.Debugf("sending config for service: %v", k)
		if err := stream.Send(v); err != nil {
			logrus.Error("stream send config error: ", err)
		}
	}
	for {
		select {
		case <-stream.Context().Done():
			s.agents.Delete(agentInfo.Id)
			return errors.New("stream context is closed")
		case cfg, ok := <-cfgCh:
			if !ok {
				logrus.Info("config channel is closed")
				return errors.New("config channel is closed by the cloud")
			}
			if err := stream.Send(cfg); err != nil {
				logrus.Error("stream send config error: ", err)
			}
		}
	}
}

func (s *server) GetConfig(ctx context.Context, agentInfo *api.AgentInfo) (*api.CollectionConfig, error) {
	if s.retainedCfg[agentInfo.Service] == nil {
		return nil, errors.New("service config not found")
	}
	return s.retainedCfg[agentInfo.Service], nil
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
		if err == context.Canceled {
			logrus.Debug("client left")
			return nil
		} else if err == io.EOF {
			logrus.Debug("client done")
			return nil
		} else if err != nil {
			return err
		}
		s.notyCh <- recv
	}
}

func (s *server) update(cfg *api.CollectionConfig) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.retainedCfg[cfg.Service] = cfg
	s.agents.Range(func(_, v interface{}) bool {
		ch := v.(chan *api.CollectionConfig)
		ch <- cfg
		return true
	})
}

type PushModel struct {
	port string
	quit chan struct{}
	svr  *server
}

func newPushModel(port string) *PushModel {
	return &PushModel{
		port: port,
		quit: make(chan struct{}),
		svr: &server{
			notyCh:      make(chan *api.CollectionEvent),
			retainedCfg: make(map[string]*api.CollectionConfig),
			svcs:        make(map[string]struct{}),
		},
	}
}

func (p PushModel) Start() {
	lis, err := net.Listen("tcp", p.port)
	if err != nil {
		logrus.Fatal("cannot listen on port, err: ", err)
	}
	s := grpc.NewServer()
	api.RegisterAgentPushServiceServer(s, p.svr)
	go func() {
		<-p.quit
		s.GracefulStop()
	}()
	if err := s.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}

func (p PushModel) Stop() {
	close(p.quit)
}

func (p PushModel) Services() []string {
	svcs := make([]string, 0, len(p.svr.svcs))
	for k := range p.svr.svcs {
		svcs = append(svcs, k)
	}
	return svcs
}

func (p PushModel) RecvNotification() <-chan *api.CollectionEvent {
	return p.svr.notyCh
}

func (p PushModel) SendConfig(cfg *api.CollectionConfig) {
	p.svr.update(cfg)
}
