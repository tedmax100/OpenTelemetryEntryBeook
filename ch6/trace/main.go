package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	// 創建標準輸出導出器
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to initialize stdouttrace exporter %v", err)
	}

	// 創建批量 Span 處理器
	batchSpanProcessor := traceSDK.NewBatchSpanProcessor(exporter)

	// 創建資源
	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNameKey.String("guide_to_observability"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// 創建 TracerProvider，並設置 TraceIdRatioBased 採樣器，採樣率設置為 50%
	tp := traceSDK.NewTracerProvider(
		traceSDK.WithSampler(traceSDK.TraceIDRatioBased(0.5)),
		traceSDK.WithSpanProcessor(batchSpanProcessor),
		traceSDK.WithResource(res),
	)

	// 註冊 TracerProvider 至OpenTelemetry框架中，此為全局可視的 Singleton 物件
	otel.SetTracerProvider(tp)

	// 使用 Tracer 的名稱來獲取 Tracer
	tracer := otel.Tracer("example-tracer")

	// 執行多次操作，生成不同的 Trace ID
	for i := 0; i < 5; i++ {
		// 創建父 Span
		ctx, parentSpan := tracer.Start(context.Background(), "operation")
		runChildOperation(tracer, ctx) // 调用 runChildOperation ，並傳入SpanContext
		parentSpan.End()               // 结束父 Span
		time.Sleep(10 * time.Millisecond)
	}

	// 等待所有 Spans 被導出
	err = tp.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failed to shutdown TracerProvider: %v", err)
	}
}

// runChildOperation 創建一個子 Span 並進行一些操作
func runChildOperation(tracer trace.Tracer, ctx context.Context) {
	// 創建子 Span
	_, childSpan := tracer.Start(ctx, "child-operation")
	defer childSpan.End()

	// 模擬一些操作
	time.Sleep(25 * time.Millisecond)

	// 隨機生成錯誤
	if rand.Float64() < 0.3 { // 30% 的概率生成錯誤
		// 透過 RecordError 記錄事件資訊
		// 該方法就是5.4小節裡介紹 Event 裡面的 RecordException
		childSpan.RecordError(errors.New("something went wrong"))
		// 設置 Span 狀態為錯誤
		childSpan.SetStatus(codes.Error, "Error occurred")

		// 加入 Span 的額外的資訊
		childSpan.SetAttributes(attribute.String("action", "do something"))
	} else {
		// 設置 Span 狀態為正常
		childSpan.SetStatus(codes.Ok, "")
	}
}
