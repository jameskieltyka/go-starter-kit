package internal

import (
	"fmt"
	"strings"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jkieltyka/go-starter-kit/internal/config"
	"github.com/jkieltyka/go-starter-kit/internal/version"
	"github.com/jkieltyka/go-starter-kit/pkg/grpcserver"
	"github.com/jkieltyka/go-starter-kit/pkg/httpserver"
	"github.com/jkieltyka/go-starter-kit/pkg/middleware/grpcstream"
	"github.com/jkieltyka/go-starter-kit/pkg/middleware/grpcunary"
	httpmiddleware "github.com/jkieltyka/go-starter-kit/pkg/middleware/http"
	versioner "github.com/jkieltyka/go-starter-kit/proto"
	"go.uber.org/zap"
)

func SetupGRPCServer(c *config.AppConfig) *grpcserver.Server {
	server := grpcserver.NewServer().
		WithStreamMiddleware(
			grpc_zap.StreamServerInterceptor(zap.L()),
			grpcstream.ExampleInterceptor(),
		).
		WithUnaryMiddleware(
			grpc_zap.UnaryServerInterceptor(zap.L()),
			grpcunary.ExampleInterceptor(),
		).
		Configure()

	versioner.RegisterVersionerServer(server.Server, version.NewVersionGRPCServer(c))

	return server
}

func SetupHTTPServer(c *config.AppConfig) *httpserver.Server {
	versionPath := fmt.Sprintf("/%v", strings.Split(version.BuildVersion, ".")[0])
	server := httpserver.NewServer().
		WithBasePath(versionPath).
		WithMiddleware(httpmiddleware.ExampleMiddleware).
		WithBaseRoutes(version.VersionRoute()).
		WithRoutes(
			version.VersionRoute(),
		)
	return server
}
