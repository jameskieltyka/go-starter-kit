package main

import (
	"fmt"

	"github.com/jkieltyka/go-starter-kit/internal"
	"github.com/jkieltyka/go-starter-kit/internal/config"
	"github.com/jkieltyka/go-starter-kit/pkg/logger"
	"go.uber.org/zap"

	"net/http"
	_ "net/http/pprof"
)

func main() {

	//start the profiler

	cfg := config.NewConfig()

	go func() {
		zap.S().Info(http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppConfig.ProfilerPort), nil))
	}()

	cfg.RegisterWatchFields(map[string]config.UpdateFunc{
		"ServerPort": config.DefaultUpdater(),
	})

	//update global loggers based on configuration
	logger.NewLogger(cfg)

	//Start Config Watcher if fields should be updated dynamically
	if len(cfg.WatchFields) != 0 {
		go cfg.StartWatch(15)
	}

	//Start the HTTP Server
	internal.SetupHTTPServer(cfg.AppConfig).
		Start(cfg.AppConfig.ServerPort)

	//Start the GRPC server
	//internal.SetupGRPCServer(cfg.AppConfig).Start(cfg.AppConfig.ServerPort)
}
