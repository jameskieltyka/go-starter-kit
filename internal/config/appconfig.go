package config

import (
	"fmt"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AppConfig defines the variables available for the application
type AppConfig struct {
	ServiceName string `flag:"servicename s" desc:"service name"`
	BasePath    string `flag:"basepath b" desc:"base path for this service"`
	ServerPort  int    `flag:"serverport p" desc:"port for the server to listen on"`
	Env         string `flag:"env e" desc:"environment"`
	LogLevel    string `flag:"loglevel v" desc:"log level"`
	LogFormat   string `flag:"logformat l" desc:"log format"`
}

//Configure settings from multiple possible sources
// Flags
// ENV variables
// Config File
var config *AppConfig

func initConfig() {

	// Use config file from the flag.
	cfgFile := viper.Get("config").(string)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	}

	err := viper.ReadInConfig()
	if cfgFile != "" && err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Read the env variables
	viper.AutomaticEnv()
}

func initInputFlags() {
	pflag.String("config", "./config/dev/config.yaml", "location of config file")

	err := gpflag.ParseTo(config, pflag.CommandLine)
	if err != nil {
		//panic
	}
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		//panic
	}

}

func NewApplicationConfig() *AppConfig {

	config = &AppConfig{}
	initInputFlags()
	initConfig()

	err := viper.Unmarshal(&config)
	if err != nil {
		//panic
		fmt.Println(err)
	}

	pflag.Parse()
	return config
}

//RegisterWatchFields registers application configuration that should be watched overtime for changes
func RegisterWatchFields() {
	//TODO
}
