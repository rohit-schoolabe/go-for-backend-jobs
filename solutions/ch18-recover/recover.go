// Package recoverrpc turns a panic in a gRPC handler into
// an Internal status instead of a crashed server.
package recoverrpc

import (
	"context"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Recover catches a panic in anything after it in the
// chain, so put it first.
func Recover(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if p := recover(); p != nil {
				log.Error("panic", "method", info.FullMethod,
					"value", p, "stack", string(debug.Stack()))
				err = status.Error(codes.Internal,
					"internal error")
			}
		}()
		return next(ctx, req)
	}
}
