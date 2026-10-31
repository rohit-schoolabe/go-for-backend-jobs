// Package catalogsvc implements the gRPC catalog service
// on top of Shelf's store.
package catalogsvc

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
	catalogv1 "github.com/rohit-schoolabe/go-for-backend-jobs/ch18/gen/shelf/catalog/v1"
)

// Books is what the service needs from storage.
// *store.Store satisfies it.
type Books interface {
	Book(ctx context.Context, id int64) (store.Book, error)
	Books(ctx context.Context,
		author string) ([]store.Book, error)
	CreateBook(ctx context.Context,
		b store.Book) (store.Book, error)
}

// Server answers catalog RPCs.
type Server struct {
	catalogv1.UnimplementedCatalogServiceServer
	Books Books
}

// New builds a gRPC server with logging and auth.
func New(books Books, tokens auth.Tokens,
	log *slog.Logger) *grpc.Server {
	protected := map[string]bool{
		catalogv1.CatalogService_AddBook_FullMethodName: true,
	}
	gs := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			Logging(log), Auth(tokens, protected)),
		grpc.StreamInterceptor(StreamLogging(log)),
	)
	catalogv1.RegisterCatalogServiceServer(gs,
		&Server{Books: books})
	return gs
}

// GetBook returns one book.
func (s *Server) GetBook(ctx context.Context,
	req *catalogv1.GetBookRequest,
) (*catalogv1.GetBookResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument,
			"id must be positive")
	}
	b, err := s.Books.Book(ctx, req.GetId())
	if errors.Is(err, store.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound,
			"book %d not found", req.GetId())
	}
	if err != nil {
		return nil, internal(ctx, err)
	}
	return &catalogv1.GetBookResponse{Book: toProto(b)}, nil
}

// ListBooks sends each matching book as its own message.
func (s *Server) ListBooks(req *catalogv1.ListBooksRequest,
	stream catalogv1.CatalogService_ListBooksServer) error {
	ctx := stream.Context()
	books, err := s.Books.Books(ctx, req.GetAuthor())
	if err != nil {
		return internal(ctx, err)
	}
	for _, b := range books {
		msg := &catalogv1.ListBooksResponse{Book: toProto(b)}
		if err := stream.Send(msg); err != nil {
			return err // the client went away
		}
	}
	return nil
}

// AddBook creates a book owned by the calling member.
func (s *Server) AddBook(ctx context.Context,
	req *catalogv1.AddBookRequest,
) (*catalogv1.AddBookResponse, error) {
	title := strings.TrimSpace(req.GetTitle())
	author := strings.TrimSpace(req.GetAuthor())
	if title == "" || author == "" ||
		req.GetYear() < 1450 || req.GetYear() > 2100 {
		return nil, status.Error(codes.InvalidArgument,
			"title, author and a year 1450-2100 are required")
	}
	b, err := s.Books.CreateBook(ctx, store.Book{
		Title: title, Author: author,
		Year: int(req.GetYear()), AddedBy: UserID(ctx),
	})
	if err != nil {
		return nil, internal(ctx, err)
	}
	return &catalogv1.AddBookResponse{Book: toProto(b)}, nil
}

// internal logs the real error and returns a generic one,
// so SQL details never reach the client.
func internal(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return status.FromContextError(ctx.Err()).Err()
	}
	slog.Error("catalog", "err", err)
	return status.Error(codes.Internal, "internal error")
}

func toProto(b store.Book) *catalogv1.Book {
	a := catalogv1.Availability_AVAILABILITY_ON_SHELF
	if b.OnLoan {
		a = catalogv1.Availability_AVAILABILITY_ON_LOAN
	}
	return &catalogv1.Book{
		Id: b.ID, Title: b.Title, Author: b.Author,
		Year: int32(b.Year), Availability: a,
	}
}
