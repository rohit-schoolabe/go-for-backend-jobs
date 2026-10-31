package catalogsvc

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
	catalogv1 "github.com/rohit-schoolabe/go-for-backend-jobs/ch18/gen/shelf/catalog/v1"
)

var tokens = auth.Tokens{
	Secret: []byte("test-secret"), TTL: time.Hour,
}

// dial starts the service on an in-memory listener and
// returns a client connected to it.
func dial(t *testing.T,
	books Books) catalogv1.CatalogServiceClient {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := New(books, tokens, slog.New(
		slog.NewTextHandler(io.Discard, nil)))
	go srv.Serve(lis)
	t.Cleanup(srv.Stop)

	dialer := func(ctx context.Context,
		_ string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return catalogv1.NewCatalogServiceClient(conn)
}

// seeded opens an in-memory store with one member and
// three books.
func seeded(t *testing.T) (*store.Store, int64) {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	ctx := context.Background()
	u, err := st.CreateUser(ctx, "asha@example.com", "x")
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range []store.Book{
		{Title: "Godaan", Author: "Premchand", Year: 1936},
		{Title: "Nirmala", Author: "Premchand", Year: 1927},
		{Title: "Malgudi Days", Author: "R. K. Narayan",
			Year: 1943},
	} {
		b.AddedBy = u.ID
		if _, err := st.CreateBook(ctx, b); err != nil {
			t.Fatal(err)
		}
	}
	return st, u.ID
}

func TestGetBook(t *testing.T) {
	st, _ := seeded(t)
	c := dial(t, st)
	ctx := context.Background()

	resp, err := c.GetBook(ctx,
		&catalogv1.GetBookRequest{Id: 1})
	if err != nil {
		t.Fatal(err)
	}
	if got := resp.GetBook().GetTitle(); got != "Godaan" {
		t.Errorf("title = %q", got)
	}

	cases := []struct {
		id   int64
		want codes.Code
	}{
		{99, codes.NotFound},
		{0, codes.InvalidArgument},
	}
	for _, tc := range cases {
		_, err := c.GetBook(ctx,
			&catalogv1.GetBookRequest{Id: tc.id})
		if got := status.Code(err); got != tc.want {
			t.Errorf("id %d: code %v, want %v",
				tc.id, got, tc.want)
		}
	}
}

func TestListBooksStreams(t *testing.T) {
	st, _ := seeded(t)
	c := dial(t, st)
	stream, err := c.ListBooks(context.Background(),
		&catalogv1.ListBooksRequest{Author: "premchand"})
	if err != nil {
		t.Fatal(err)
	}
	var titles []string
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		titles = append(titles, msg.GetBook().GetTitle())
	}
	if len(titles) != 2 {
		t.Errorf("got %v, want the two Premchand books",
			titles)
	}
}

func TestAddBookNeedsToken(t *testing.T) {
	st, uid := seeded(t)
	c := dial(t, st)
	req := &catalogv1.AddBookRequest{
		Title: "Train to Pakistan", Author: "Khushwant Singh",
		Year: 1956,
	}

	_, err := c.AddBook(context.Background(), req)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("no token: %v", err)
	}

	tok, err := tokens.Issue(uid, "asha@example.com")
	if err != nil {
		t.Fatal(err)
	}
	ctx := metadata.AppendToOutgoingContext(
		context.Background(), "authorization", "Bearer "+tok)
	resp, err := c.AddBook(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetBook().GetId() != 4 {
		t.Errorf("id = %d, want 4", resp.GetBook().GetId())
	}
}

// slowBooks waits until its context ends. Embedding the
// real store fills in the methods the test never calls.
type slowBooks struct{ *store.Store }

func (slowBooks) Book(ctx context.Context,
	id int64) (store.Book, error) {
	<-ctx.Done()
	return store.Book{}, ctx.Err()
}

func TestDeadlineReachesServer(t *testing.T) {
	c := dial(t, slowBooks{})
	ctx, cancel := context.WithTimeout(
		context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := c.GetBook(ctx, &catalogv1.GetBookRequest{Id: 1})
	if status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("got %v, want DeadlineExceeded", err)
	}
}
