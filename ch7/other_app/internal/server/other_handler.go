package server

import (
	otelsdk "ch7_other_app/internal/otel"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

func other_handler(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()
	start := time.Now()
	roll := 1 + rand.Intn(6)
	duration := time.Since(start).Milliseconds()

	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))
	rollDuration.Record(ctx, duration, metric.WithAttributes(rollValueAttr))

	otelsdk.Logger.DebugContext(ctx, "roll dice", zap.Int("roll_value", roll))

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		otelsdk.Logger.ErrorContext(ctx, "Write failed", zap.Error(err))
		span.SetStatus(codes.Error, err.Error())
	}
	span.SetStatus(codes.Ok, "")
}
