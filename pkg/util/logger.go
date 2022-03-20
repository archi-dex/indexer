package util

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = *zap.SugaredLogger

func InitLogger(debug bool) Logger {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.TimeKey = zapcore.OmitKey
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}

	return logger.Sugar()
}
