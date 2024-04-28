package server

import (
	otelsdk "ch7/internal/otel"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
)

func register(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "register")
	defer span.End()

	otelsdk.Logger.DebugContext(ctx, "register",
		zap.String("pwd", "aa1234"),
		zap.String("address", "taipei neihu"))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rng.Intn(2) == 0 { // 有一半的機率生成錯誤
		err := errors.New("random error occurred")
		otelsdk.Logger.ErrorContext(ctx, "Error during register",
			zap.Error(err),
			zap.Int("severity_number", int(log.SeverityError)), // 記錄 SeverityNumber 為 17
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	resp := "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		otelsdk.Logger.ErrorContext(ctx, "Write failed", zap.Error(err))
		span.SetStatus(codes.Error, err.Error())
	}
	span.SetStatus(codes.Ok, "")
}
