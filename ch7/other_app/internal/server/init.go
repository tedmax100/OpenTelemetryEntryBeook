package server

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer       = otel.Tracer("rolldice")
	meter        = otel.Meter("rolldice")
	rollCnt      metric.Int64Counter
	rollDuration metric.Int64Histogram
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}
	rollDuration, err = meter.Int64Histogram("dice.roll_duration", metric.WithDescription("The duration of dice rolls"), metric.WithUnit("ms"))
	if err != nil {
		panic(err)
	}
}
