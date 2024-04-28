package server

import (
	otelsdk "ch7/internal/otel"
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func rolldice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()
	other_app_handle(ctx, span)

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

func other_app_handle(ctx context.Context, span trace.Span) {
	// 初始化OpenTelemetry包装的HTTP客户端
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	// 构造要发起的HTTP请求
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://other_app:8081/handler", nil)

	// 发起HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		otelsdk.Logger.ErrorContext(ctx, "HTTP request failed", zap.Error(err))
		span.SetStatus(codes.Error, err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// 处理响应...
	otelsdk.Logger.DebugContext(ctx, "API response", zap.ByteString("response", body))

}
