package catalogsvc

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
)

// Logging logs every unary call with its result code.
func Logging(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := next(ctx, req)
		log.Info("rpc", "method", info.FullMethod,
			"code", status.Code(err).String(),
			"ms", time.Since(start).Milliseconds())
		return resp, err
	}
}

// StreamLogging does the same for streaming calls.
func StreamLogging(
	log *slog.Logger,
) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		next grpc.StreamHandler) error {
		start := time.Now()
		err := next(srv, ss)
		log.Info("rpc", "method", info.FullMethod,
			"code", status.Code(err).String(),
			"ms", time.Since(start).Milliseconds())
		return err
	}
}

type userKey struct{}

// UserID is the member a protected call was made for.
func UserID(ctx context.Context) int64 {
	id, _ := ctx.Value(userKey{}).(int64)
	return id
}

// Auth checks the bearer token on protected methods and
// puts the member's ID in the context.
func Auth(tokens auth.Tokens,
	protected map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler) (any, error) {
		if !protected[info.FullMethod] {
			return next(ctx, req)
		}
		md, _ := metadata.FromIncomingContext(ctx)
		vals := md.Get("authorization")
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated,
				"missing token")
		}
		tok, ok := strings.CutPrefix(vals[0], "Bearer ")
		id, err := tokens.Verify(tok)
		if !ok || err != nil {
			return nil, status.Error(codes.Unauthenticated,
				"invalid token")
		}
		ctx = context.WithValue(ctx, userKey{}, id)
		return next(ctx, req)
	}
}
