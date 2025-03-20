package main

import (
	"context"
	"fmt"
	"github.com/sunzhqr/sharedlibprompt/internal/config"
	"github.com/sunzhqr/sharedlibprompt/internal/service"
	test "github.com/sunzhqr/sharedlibprompt/pkg/api/test/api"
	"github.com/sunzhqr/sharedlibprompt/pkg/logger"
	"github.com/sunzhqr/sharedlibprompt/pkg/postgres"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to load config", zap.Error(err))
	}
	pool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to database", zap.Error(err))
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %w", err)
	}

	srv := service.New()
	server := grpc.NewServer(grpc.UnaryInterceptor(logger.Interceptor))
	test.RegisterOrderServiceServer(server, srv)
	go func() {
		if err := server.Serve(listener); err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "failed ro serve", zap.Error(err))
		}
	}()

	select {
	case <-ctx.Done():
		server.Stop()
		pool.Close()
		logger.GetLoggerFromCtx(ctx).Info(ctx, "server stopped")
	}
}
