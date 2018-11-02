package handler

import "go.uber.org/zap"

var logger *zap.SugaredLogger

// SetRequestLogger registers the logger instance for all
// request handlers (everything in the handler package)
func SetRequestLogger(l *zap.SugaredLogger) {
	logger = l
}
