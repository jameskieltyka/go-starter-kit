package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_initConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		config   string
		expected int
	}{
		{"no env or config", make(map[string]string), "", 0},
		{"port env", map[string]string{"SERVERPORT": "9090"}, "", 9090},
		{"port config", make(map[string]string), "../../config/dev/config.yaml", 8080},
		{"port config and env", map[string]string{"SERVERPORT": "9090"}, "../../config/dev/config.yaml", 9090},
	}

	for _, tt := range tests {
		for k, v := range tt.env {
			os.Setenv(k, v)
		}
		viper.Set("config", tt.config)

		t.Run(tt.name, func(t *testing.T) {
			initConfig()
			port := viper.GetInt("serverport")
			assert.Equal(t, tt.expected, port)
		})

		for k := range tt.env {
			os.Unsetenv(k)
		}
	}
}

func Test_initInputFlags(t *testing.T) {
	t.Run("test input flags", func(t *testing.T) {
		initInputFlags(&AppConfig{})
		_ = pflag.CommandLine.Set("config", "../test.yaml")
		_ = pflag.CommandLine.Set("serverport", "8080")

		config := viper.Get("config")
		port := viper.GetInt("serverport")
		assert.Equal(t, "../test.yaml", config)
		assert.Equal(t, 8080, port)
	})
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		port   int
	}{
		{"set some config via env variable", map[string]string{"SERVERPORT": "9090"}, 9090},
	}

	for _, tt := range tests {
		for k, v := range tt.fields {
			os.Setenv(k, v)
		}

		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig()
			assert.Equal(t, tt.port, config.AppConfig.ServerPort)
		})

		for k := range tt.fields {
			os.Unsetenv(k)
		}
	}
}

func TestConfig_RegisterWatchFields(t *testing.T) {

	t.Run("add field to be watched", func(t *testing.T) {
		c := &Config{
			AppConfig:   &AppConfig{},
			WatchFields: make(map[string]UpdateFunc),
		}

		c.RegisterWatchFields(map[string]UpdateFunc{
			"ServerPort": func(config *AppConfig, field string, newConfig *AppConfig) { config.ServerPort = 8080 },
		})

		c.WatchFields["ServerPort"](c.AppConfig, "ServerPort", c.AppConfig)

		assert.Equal(t, 8080, c.AppConfig.ServerPort)
	})

}

func TestDefaultUpdater(t *testing.T) {
	tests := []struct {
		name        string
		initialPort int
		updatedPort int
	}{
		{"test updating server port", 8080, 9090},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := DefaultUpdater()
			config := &AppConfig{ServerPort: tt.initialPort}
			update(config, "ServerPort", &AppConfig{ServerPort: tt.updatedPort})
			assert.Equal(t, tt.updatedPort, config.ServerPort)
		})
	}
}
