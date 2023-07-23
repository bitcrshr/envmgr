package shared

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func Logger() *zap.Logger {
	if logger != nil {
		return logger
	}

	logger = zap.Must(zap.NewDevelopment())

	return logger
}
