package logger

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"time"

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

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.zl.Fatal(msg, fields...)
}

func Interceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	next grpc.UnaryHandler,
) (interface{}, error) {
	guid := uuid.New().String()
	ctx = context.WithValue(ctx, RequestID, guid)
	GetLoggerFromCtx(ctx).Info(ctx,
		"request", zap.String("info", info.FullMethod),
		zap.Time("request time", time.Now()),
	)
	return next(ctx, req)
}
