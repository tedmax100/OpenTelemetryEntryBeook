package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	ctx := context.Background()

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to initialize otlp trace exporter: %v\n", err)
		return
	}

	// 创建TracerProvider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("CacheService"),
		)),
	)
	otel.SetTracerProvider(tp)

	// 获取Meter和Tracer

	tracer := otel.Tracer("ex.com/basic")

	// 创建直方图
	cacheAccessDurations := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "demo",
			Name:      "cache_access_durations",
			Help:      "Histogram of cache access durations.",
			Buckets:   prometheus.LinearBuckets(20, 5, 5), // 示例：从20毫秒开始，每5毫秒一个桶，共5个桶。
		},
	)

	// 注册直方图
	prometheus.MustRegister(cacheAccessDurations)

	// Create non-global registry.
	registry := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		cacheAccessDurations,
	)

	// 启动 Prometheus HTTP 服务器
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

	go func() {
		fmt.Println("Serving metrics on :9090")
		_ = http.ListenAndServe(":9090", nil)
	}()

	// 模拟缓存命中和未命中的事件
	for i := 0; i < 10; i++ {
		// 开始追踪
		_, span := tracer.Start(ctx, "CacheOperation")
		spanCtx := span.SpanContext()

		// 将追踪信息附加为 Exemplar
		exemplarLabels := prometheus.Labels{
			"trace_id":  spanCtx.TraceID().String(),
			"span_id":   spanCtx.SpanID().String(),
			"operation": "cache_access",
		}

		// 模拟一次缓存访问的时间
		cacheAccessTime := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(cacheAccessTime)
		// 记录缓存访问时间
		cacheAccessDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
			float64(cacheAccessTime.Milliseconds()),
			exemplarLabels,
		)
		fmt.Printf("Cache accessed for %v\n", cacheAccessTime)

		time.Sleep(500 * time.Millisecond)

		// 结束追踪
		span.End()
	}

	// Shutdown the MeterProvider and TracerProvider, this will flush any remaining data
	err = tp.Shutdown(ctx)
	if err != nil {
		fmt.Printf("failed to stop tracer provider: %v\n", err)
	}
}
