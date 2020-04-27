package logger

import (
	"encoding/json"
	"fmt"

	"github.com/jkieltyka/go-starter-kit/internal/config"
	"go.uber.org/zap"
)

func init() {
	//Create default global logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create default logger")
	}
	zap.ReplaceGlobals(logger)
}

func NewLogger(c *config.Config) {
	loggerConfig := []byte(fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "msg",
		  "timeKey": "ts",
		  "callerKey": "caller",
		  "levelKey": "level",
		  "levelEncoder": "lowercase",
		  "timeEncoder": "seconds",
		  "callerEncoder": "short"
		}
	  }`, c.AppConfig.LogLevel))

	var cfg zap.Config
	if err := json.Unmarshal(loggerConfig, &cfg); err != nil {
		zap.L().Panic(fmt.Sprintf("unable to build configured logger %v", err))
	}

	logger, err := cfg.Build()
	if err != nil {
		zap.L().Panic(fmt.Sprintf("unable to build configured logger %v", err))
	}

	zap.ReplaceGlobals(logger)
	zap.L().Info("updated logger using application configuration")
}
