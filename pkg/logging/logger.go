package logger

import (
	"fmt"

	"github.com/gcamlicali/tradeshopExample/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logger with the given log level
func NewLogger(config *config.Config) {
	logLevel, err := zapcore.ParseLevel(config.Logger.Level)
	if err != nil {
		panic(fmt.Sprintf("Unknown log level %v", logLevel))
	}

	var cfg zap.Config
	if config.Logger.Development {
		// That config includes developing config options
		cfg = zap.NewDevelopmentConfig()
		// Colorful output
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// More Light config options
		cfg = zap.NewProductionConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	//We could choose return logger but this function set this variable as global
	zap.ReplaceGlobals(logger)
}

func Close() {
	// Last actions before close. Like "Flush"
	defer zap.L().Sync()
}
