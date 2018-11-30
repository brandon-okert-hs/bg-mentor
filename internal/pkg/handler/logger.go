package handler

import (
	"crypto/rand"
	"encoding/base64"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// SetRequestLogger registers the logger instance for all
// request handlers (everything in the handler package)
func SetRequestLogger(l *zap.SugaredLogger) {
	logger = l
}

// LogId creates a new unique string that be used to easily trace logs
func LogId() string {
	b := make([]byte, 128)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
