package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sunzhqr/sharedlibprompt/internal/service"
	test "github.com/sunzhqr/sharedlibprompt/pkg/api/test/api"
	"github.com/sunzhqr/sharedlibprompt/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %w", err)
	}

	srv := service.New()
	server := grpc.NewServer()
	test.RegisterOrderServiceServer(server, srv)
	if err := server.Serve(listener); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed ro serve", zap.Error(err))
	}
}
