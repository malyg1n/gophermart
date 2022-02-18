package logger

import "go.uber.org/zap"

var instance *zap.SugaredLogger

// GetLogger returns instance of SugaredLogger.
func GetLogger() *zap.SugaredLogger {
	if instance != nil {
		return instance
	}

	logger, _ := zap.NewDevelopment()
	instance = logger.Sugar()

	return instance
}
