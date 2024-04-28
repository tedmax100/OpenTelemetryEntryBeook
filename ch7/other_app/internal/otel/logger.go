package otel

import (
	"log"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *otelzap.Logger

func init() {
	// Wrap zap logger to extend Zap with API that accepts a context.Context.
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = EncodeTime
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := config.Build(zap.AddCaller(), zap.Fields(zap.String("service", "ch7")))
	if err != nil {
		log.Println(err.Error())
	}
	logger.Sync()

	Logger = otelzap.New(logger, otelzap.WithTraceIDField(true), otelzap.WithMinLevel(zapcore.DebugLevel))

}

func EncodeTime(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendInt64(t.Unix())
}

func CapitalLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(level.String())
}
