package grpcserver

import (
	"fmt"
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Done             chan bool
	Stop             chan os.Signal
	StreamMiddleware []grpc.StreamServerInterceptor
	UnaryMiddleware  []grpc.UnaryServerInterceptor
}

func NewServer() *Server {
	return &Server{}
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

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panicf("failed to listen: %v", err)
	}
	zap.S().Infof("Starting server %s", lis.Addr().String())

	return s.Serve(lis)
}
