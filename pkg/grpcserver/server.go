package grpcserver

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Done             chan bool
	StopSignal       chan os.Signal
	StreamMiddleware []grpc.StreamServerInterceptor
	UnaryMiddleware  []grpc.UnaryServerInterceptor
}

func NewServer() *Server {
	return &Server{
		Done:       make(chan bool, 1),
		StopSignal: make(chan os.Signal, 1),
	}
}

func (s *Server) WithStreamMiddleware(middlewares ...grpc.StreamServerInterceptor) *Server {
	s.StreamMiddleware = append(s.StreamMiddleware, middlewares...)
	return s
}

func (s *Server) WithUnaryMiddleware(middlewares ...grpc.UnaryServerInterceptor) *Server {
	s.UnaryMiddleware = append(s.UnaryMiddleware, middlewares...)
	return s
}

func (s *Server) Configure() *Server {
	s.Server = grpc.NewServer(
		grpc.ChainStreamInterceptor(s.StreamMiddleware...),
		grpc.ChainUnaryInterceptor(s.UnaryMiddleware...),
	)
	return s
}

func (s *Server) WaitShutdown() {
	signal.Notify(s.StopSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-s.Done:
		zap.S().Info("server stopping due to finish request")
	case sig := <-s.StopSignal:
		zap.S().Infof("server stopped due to signal %v", sig)
	}

	//give the server 2 seconds to finish any outstanding requests
	timer := time.NewTimer(2 * time.Second)

	finished := make(chan bool)
	go func() {
		s.GracefulStop()
		close(finished)
	}()

	select {
	case <-timer.C:
		s.Stop()
		zap.S().Infof("server failed to gracefully shut down")
	case <-finished:
		timer.Stop()
		zap.S().Infof("server gracefully shut down")
	}
}

func (s *Server) Start(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panicf("failed to listen: %v", err)
	}
	zap.S().Infof("Starting server %s", lis.Addr().String())

	go func() {
		err := s.Serve(lis)
		if err != nil {
			zap.S().Errorf("server closed due to error %v", err)
			s.Done <- true
		}
	}()

	// wait shutdown
	s.WaitShutdown()
}
