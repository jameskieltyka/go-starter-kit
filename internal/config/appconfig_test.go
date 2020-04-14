package config

import (
	"os"
	"reflect"
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
	type fields struct {
		AppConfig   *AppConfig
		WatchFields map[string]UpdateFunc
	}
	type args struct {
		watch map[string]UpdateFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				AppConfig:   tt.fields.AppConfig,
				WatchFields: tt.fields.WatchFields,
			}
			c.RegisterWatchFields(tt.args.watch)
		})
	}
}

func TestDefaultUpdater(t *testing.T) {
	tests := []struct {
		name string
		want UpdateFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultUpdater(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultUpdater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_StartWatch(t *testing.T) {
	type fields struct {
		AppConfig   *AppConfig
		WatchFields map[string]UpdateFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				AppConfig:   tt.fields.AppConfig,
				WatchFields: tt.fields.WatchFields,
			}
			c.StartWatch()
		})
	}
}
