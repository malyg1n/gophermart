package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	zapLogger, _ := zap.NewDevelopment()
	logger = zapLogger.Sugar()
}

// NewDefaultLogger returns instance of SugaredLogger.
func NewDefaultLogger() *zap.SugaredLogger {
	return logger
}
