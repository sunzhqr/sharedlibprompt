package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sunzhqr/sharedlibprompt/internal/config"
	"github.com/sunzhqr/sharedlibprompt/internal/service"
	"github.com/sunzhqr/sharedlibprompt/pkg/api/test/api"
	"github.com/sunzhqr/sharedlibprompt/pkg/logger"
	"github.com/sunzhqr/sharedlibprompt/pkg/postgres"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	api.RegisterOrderServiceServer(server, srv)

	rt := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = test.RegisterOrderServiceHandlerFromEndpoint(ctx, rt, "localhost:"+strconv.Itoa(cfg.GRPCPort), opts)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to register handler server", zap.Error(err))
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort), rt); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to serve REST", zap.Error(err))
		}
	}()

	go func() {
		if err := server.Serve(listener); err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "failed ro serve", zap.Error(err))
		}
	}()
	fmt.Println("Server launched")

	select {
	case <-ctx.Done():
		server.GracefulStop()
		pool.Close()
		logger.GetLoggerFromCtx(ctx).Info(ctx, "server stopped")
	}
}
