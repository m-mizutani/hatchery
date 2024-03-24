package utils

import (
	"context"
	"log/slog"
	"time"

	"github.com/m-mizutani/hatchery/pkg/domain/model"
)

type ctxLoggerKey struct{}

// CtxWithLogger returns a new context with logger
func CtxWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

// CtxLoggerWith returns a new logger with additional fields
func CtxLoggerWith(ctx context.Context, attrs ...slog.Attr) context.Context {
	logger := CtxLogger(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	return CtxWithLogger(ctx, logger)
}

// CtxLogger returns logger from context. If logger is not set, return default logger
func CtxLogger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLoggerKey{}).(*slog.Logger); ok {
		return l
	}
	return logger
}

type ctxRequestIDKey struct{}

// CtxRequestID returns request ID from context. If request ID is not set, return new request ID and context with it
func CtxRequestID(ctx context.Context) (model.RequestID, context.Context) {
	if id, ok := ctx.Value(ctxRequestIDKey{}).(model.RequestID); ok {
		return id, ctx
	}

	newID := model.NewRequestID()
	ctx = CtxLoggerWith(ctx, slog.Any("request_id", newID))
	ctx = context.WithValue(ctx, ctxRequestIDKey{}, newID)
	return newID, ctx
}

type ctxNowKey struct{}

type nowFunc func() time.Time

func CtxNow(ctx context.Context) time.Time {
	f, ok := ctx.Value(ctxNowKey{}).(nowFunc)
	if !ok {
		return time.Now()
	}
	return f()
}

func CtxWithNow(ctx context.Context, f nowFunc) context.Context {
	return context.WithValue(ctx, ctxNowKey{}, f)
}
