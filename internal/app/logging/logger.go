package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func GetLogger() (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"logfile.log", "stdout"}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
