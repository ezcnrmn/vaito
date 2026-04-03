package app

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *App) recoverPanic(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			a.log.Error("panic", "method", info.FullMethod, "err", r)
			err = status.Error(codes.Internal, "Internal error")
		}
	}()

	return handler(ctx, req)
}
