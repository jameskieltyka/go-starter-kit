package main

import (
	"fmt"
	"os"

	"github.com/jkieltyka/go-starter-kit/internal/config"
	"github.com/jkieltyka/go-starter-kit/pkg/logger"
)

func init() {
	fmt.Println(os.Args[0])
}

func main() {
	cfg := config.NewConfig()
	cfg.RegisterWatchFields(map[string]config.UpdateFunc{
		"ServerPort": config.DefaultUpdater(),
	})

	//update global loggers based on configuration
	logger.NewLogger(cfg)

	//Start Config Watcher if fields should be updated dynamically
	if len(cfg.WatchFields) != 0 {
		go cfg.StartWatch()
	}

	select {}

	//server := internal.Setup(cfg)
	//server.Start("", cfg.ServerPort())
}
