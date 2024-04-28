package main

import (
	"ch7/internal/otel"
	"ch7/internal/server"

	"go.uber.org/zap"
)

func main() {
	logger := otel.Logger
	if err := server.Run(); err != nil {
		logger.Error("Server run failed", zap.Error(err))
	}
}
