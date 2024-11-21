package main

import (
	"context"
	"log/slog"
	"runtime"
)

type CustLogger struct {
	base *slog.Logger
}

func NewCustLoggerLogger(base *slog.Logger) *CustLogger {
	return &CustLogger{base: base}
}

func getCallerInfo() (file string, line int, funcName string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown", 0, "unknown"
	}
	funcName = runtime.FuncForPC(pc).Name()
	return file, line, funcName
}

func (l *CustLogger) InfoContext(ctx context.Context, msg string, keysAndValues ...any) {
	file, line, funcName := getCallerInfo()
	extendedKeysAndValues := append(keysAndValues,
		"caller.file", file,
		"caller.line", line,
		"caller.func", funcName,
	)
	l.base.InfoContext(ctx, msg, extendedKeysAndValues...)
}

func (l *CustLogger) ErrorContext(ctx context.Context, err error, keysAndValues ...any) {
	file, line, funcName := getCallerInfo()
	extendedKeysAndValues := append(keysAndValues,
		"caller.file", file,
		"caller.line", line,
		"caller.func", funcName,
	)
	l.base.ErrorContext(ctx, err.Error(), extendedKeysAndValues...)
}
