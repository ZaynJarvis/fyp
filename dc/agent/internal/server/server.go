package server

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/api"
	"google.golang.org/grpc"
)

type Server struct {
	api.UnimplementedLocalServer
	port  string
	quit  chan struct{}
	imgCh chan *api.ImageReport
}

func New(port string) *Server {
	return &Server{
		port:  port,
		quit:  make(chan struct{}),
		imgCh: make(chan *api.ImageReport),
	}
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		logrus.Fatal("cannot listen on port, err: ", err)
	}
	svr := grpc.NewServer()
	api.RegisterLocalServer(svr, s)
	go func() {
		<-s.quit
		svr.GracefulStop()
		close(s.imgCh)
	}()
	if err := svr.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}
func (s *Server) Stop() {
	close(s.quit)
}

func (s *Server) Image(stream api.Local_ImageServer) error {
	for {
		rpt, err := stream.Recv()
		if err == context.Canceled {
			logrus.Debug("client left")
			return nil
		} else if err == io.EOF {
			logrus.Debug("client done")
			return nil
		} else if err != nil {
			return err
		}
		select {
		case s.imgCh <- rpt:
		case <-time.After(time.Millisecond):
			logrus.Debug("image report channel full")
		}
	}
}

func (s *Server) RecvImageReport() chan *api.ImageReport {
	return s.imgCh
}
