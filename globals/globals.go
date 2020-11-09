package globals

import (
	"github.com/utsavgupta/go-leader-election/config"
	"github.com/utsavgupta/go-leader-election/config/utils"
	"github.com/utsavgupta/go-leader-election/logger"
)

// Logger is a global logger used throughout the application
var Logger logger.Logger

// Config stores the application configuration
var Config *config.Config

// InitConfig should be the first initializer that should be called.
// It reads values from the environment and builds the configuration
// settings for the application.
func InitConfig() {
	Config = &config.Config{}

	Config.Stage = utils.GetString("APP_STAGE")
	Config.Port = utils.GetInt("APP_PORT")
}

// InitLogger creates a new default logger for the application
func InitLogger(application, env string) {
	Logger = logger.NewZapLogger(application, env)
}
