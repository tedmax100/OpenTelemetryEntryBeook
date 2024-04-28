package main

import (
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogEntity struct {
	App      string        `json:"app"`
	Duration time.Duration `json:"duration"`
	Status   int32         `json:"status"`
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			duration := time.Duration(rand.Intn(300)+10) * time.Millisecond
			time.Sleep(duration)
			status := int32(200)   // 預設為 2xx
			if rand.Intn(2) == 0 { // 隨機選擇是否改為 4xx
				status = int32(400)
			}

			entity := LogEntity{
				App:      "Ch10App",
				Duration: duration,
				Status:   status,
			}

			if duration.Milliseconds()%2 == 0 {
				Logger.Info("response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
			} else {
				Logger.Error("err response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
			}

		}
	}()

	go func() {
		for {
			duration := time.Duration(rand.Intn(300)+10) * time.Millisecond
			time.Sleep(duration)
			status := int32(200)   // 預設為 2xx
			if rand.Intn(2) == 0 { // 隨機選擇是否改為 4xx
				status = int32(400)
			}

			entity := LogEntity{
				App:      "Ch10App2",
				Duration: duration,
				Status:   status,
			}

			if duration.Milliseconds()%2 == 0 {
				Logger.Info("response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
			} else {
				Logger.Error("err response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
			}

		}
	}()

	<-sigs
}

func formatDurationWithMS(d time.Duration) string {
	return d.String() // `String` method already formats as "72ms", "2h45m", etc.
}

var Logger *otelzap.Logger

func init() {
	// Wrap zap logger to extend Zap with API that accepts a context.Context.
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = EncodeTime
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := config.Build(zap.AddCaller(), zap.Fields(zap.String("service", "ch10")))
	if err != nil {
		slog.Error(err.Error())
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
