package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// AppConfig defines the variables available for the application
type AppConfig struct {
	ServiceName string `flag:"servicename s" desc:"service name"`
	ServerPort  int    `flag:"serverport p" desc:"port for the server to listen on"`
	Env         string `flag:"env e" desc:"environment"`
	LogLevel    string `flag:"loglevel v" desc:"log level"`
}

//Config Contains the application level configuration and the dynamic fields to watch
type Config struct {
	AppConfig   *AppConfig
	WatchFields map[string]UpdateFunc
}

func initConfig() {
	// Use config file from the flag.
	cfgFile := viper.Get("config").(string)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	}

	err := viper.ReadInConfig()
	if cfgFile != "" && err == nil {
		zap.S().Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	// Read the env variables
	viper.AutomaticEnv()
}

func initInputFlags(config *AppConfig) {
	pflag.String("config", "./config/dev/config.yaml", "location of config file")

	err := gpflag.ParseTo(config, pflag.CommandLine)
	if err != nil {
		zap.L().Panic(fmt.Sprintf("unable to generate flags from configuration struct, %v", err))
	}
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		zap.L().Panic(fmt.Sprintf("unable to bind commandline flags, %v", err))
	}

}

// NewConfig loads the current config settings for the application
// Config is loaded from multiple sources (highest to lowest precedence):
// Flags
// ENV variables
// Config File
func NewConfig() *Config {
	config := &AppConfig{}
	defer func() {
		zap.L().Sync()
		zap.S().Sync()
	}()
	initInputFlags(config)
	initConfig()

	err := viper.Unmarshal(&config)
	if err != nil {
		zap.L().Panic(fmt.Sprintf("unable to set startup configuration, %v", err))
	}

	pflag.Parse()
	return &Config{AppConfig: config}
}

//RegisterWatchFields registers application configuration that should be watched overtime for changes
type UpdateFunc func(config *AppConfig, field string, newConfig *AppConfig)

//RegisterWatchFields is used to register what config values to update on changes
func (c *Config) RegisterWatchFields(watch map[string]UpdateFunc) {
	c.WatchFields = watch
}

//DefaultUpdater updates the config value to the new value and logs
func DefaultUpdater() UpdateFunc {
	return func(config *AppConfig, field string, newConfig *AppConfig) {
		newFieldValue := reflect.ValueOf(newConfig).Elem().FieldByName(field)
		if reflect.ValueOf(config).Elem().FieldByName(field).Kind() == reflect.Invalid ||
			newFieldValue.Kind() == reflect.Invalid {
			//log error and continue
			zap.L().Warn(fmt.Sprintf("field not found in configuration %s", field))
			return
		}

		if reflect.ValueOf(config).Elem().FieldByName(field).Interface() != newFieldValue.Interface() {
			reflect.ValueOf(config).Elem().FieldByName(field).Set(newFieldValue)
			zap.S().Infof("updated field %s to value %v", field, newFieldValue)
		}
	}
}

//StartWatch Starts polling registered fields for changes in value
func (c *Config) StartWatch(delay time.Duration) {
	go func() {
		for {
			time.Sleep(time.Second * delay)
			var tempConfig AppConfig

			err := viper.Unmarshal(&tempConfig)
			if err != nil {
				//log error
				zap.L().Warn(fmt.Sprintf("could not unmarshal new configuration %v", err))
				continue
			}

			for field, updateFunc := range c.WatchFields {
				//run update function on field
				updateFunc(c.AppConfig, field, &tempConfig)
			}
		}

	}()
}
