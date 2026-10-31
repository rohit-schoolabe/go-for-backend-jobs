package recoverrpc

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRecover(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	icpt := Recover(log)
	boom := func(context.Context, any) (any, error) {
		panic("nil map write")
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/t/Boom"}
	_, err := icpt(context.Background(), nil, info, boom)
	if status.Code(err) != codes.Internal {
		t.Fatalf("got %v, want Internal", err)
	}
}
