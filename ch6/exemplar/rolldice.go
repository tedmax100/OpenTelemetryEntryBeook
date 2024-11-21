package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "go.opentelemetry.io/otel/example/dice"

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
	logger = NewCustLoggerLogger(otelslog.NewLogger(name))
	// logger  = otelslog.NewLogger(name)
	rollCnt          metric.Int64Counter
	rollCntHistogram metric.Int64Histogram
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}

	rollCntHistogram, err = meter.Int64Histogram("dice.rolls",
		metric.WithDescription("The histogram of rolls by roll value"),
		metric.WithUnit("{roll}"),
	)
	if err != nil {
		panic(err)
	}
}

func rolldice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	roll := 1 + rand.Intn(6)

	// 記錄 Exemplar 的 TraceID 和 SpanID
	// exemplarAttributes := []attribute.KeyValue{
	// 	attribute.String("trace_id", span.SpanContext().TraceID().String()),
	// 	attribute.String("span_id", span.SpanContext().SpanID().String()),
	// }

	var msg string
	if player := r.PathValue("player"); player != "" {
		msg = fmt.Sprintf("%s is rolling the dice", player)
	} else {
		msg = "Anonymous player is rolling the dice"
	}

	logger.InfoContext(ctx, msg,
		"result", roll,
	)

	// 記錄 Roll Value
	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)

	span.SetAttributes(rollValueAttr)
	// 增加 Counter 並附加 Exemplar
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	// 記錄到 Histogram 並附加 Exemplar
	rollCntHistogram.Record(ctx, int64(roll), metric.WithAttributes(rollValueAttr))

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}

	logger.ErrorContext(ctx, errors.New("errrrr"),
		"argName", "1212",
		"argName2", 4444)
}
