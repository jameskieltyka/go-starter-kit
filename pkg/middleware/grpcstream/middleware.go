package grpcstream

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ExampleInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		zap.S().Infof("this is an example")
		return nil
	}
}
