package grpcunary

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ExampleInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h, err := handler(ctx, req)
		zap.S().Infof("this is an example")
		return h, err
	}
}
