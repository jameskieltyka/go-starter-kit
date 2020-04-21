package main

import (
	"github.com/jkieltyka/go-starter-kit/internal"
	"github.com/jkieltyka/go-starter-kit/internal/config"
	"github.com/jkieltyka/go-starter-kit/pkg/logger"

	"net/http"
	_ "net/http/pprof"
)

func main() {

	//start the profiler
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	cfg := config.NewConfig()
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
	// internal.SetupHTTPServer(cfg.AppConfig).
	// 	Start(cfg.AppConfig.ServerPort)

	//Start the GRPC server
	internal.SetupGRPCServer(cfg.AppConfig).Start(cfg.AppConfig.ServerPort)
}
