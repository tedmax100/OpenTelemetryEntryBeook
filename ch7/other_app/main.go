package main

import (
	"ch7_other_app/internal/otel"
	"ch7_other_app/internal/server"

	"go.uber.org/zap"
)

func main() {
	logger := otel.Logger
	if err := server.Run(); err != nil {
		logger.Error("Server run failed", zap.Error(err))
	}
}
