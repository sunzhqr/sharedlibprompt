package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	Key       = "logger"
	RequestID = "request_id"
)

type Logger struct {
	zl *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, err
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.zl.Info(msg, fields...)
}
