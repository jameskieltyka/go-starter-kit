package main

import (
	"fmt"

	"github.com/jkieltyka/go-starter-kit/internal/config"
	"github.com/spf13/viper"
)

func main() {
	cfg := config.NewApplicationConfig()
	fmt.Println(cfg, viper.Get("server-port"))
	//server := internal.Setup(cfg)
	//server.Start("", cfg.ServerPort())
}
