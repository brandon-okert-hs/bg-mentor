package database

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// SetDBLogger registers the logger instance for all
// db interactions (everything in the database package)
func SetDBLogger(l *zap.SugaredLogger) {
	logger = l
}
