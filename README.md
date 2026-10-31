# Go for Backend Jobs: example code

Every listing in the book comes from a file in this repository, and each one
was checked against these files before the book was built, so what you read
is what compiles.

## Requirements

Go 1.22 or newer. Check with `go version`. Chapter 12's race detector needs
cgo, which means a C compiler (gcc or clang) on your machine.

## Layout

One folder per chapter, all inside a single module, `github.com/rohit-schoolabe/go-for-backend-jobs`:

| Folder | What is in it |
|---|---|
| `ch01` to `ch06` | the first programs and the Shelf desk tool |
| `ch07` | saving the catalog as JSON |
| `ch08` | interfaces, `any` and generics |
| `ch09` | arrays, slices, maps and the catalog |
| `ch10` | functions as values, closures, recursion |
| `ch11/finereport` | the fine report batch job |
| `ch12` | goroutines, channels, the concurrent report |
| `ch13` | readers and writers, streaming JSON, slog, time zones, signals, embed |
| `ch14` | the REST API, first in memory, then with SQLite |
| `ch15/library-api` | users, logins, tokens and loans |
| `ch16/library-api` | a copy of chapter 15's API with fines, reminders, tests, fuzzing, benchmarks and pprof |
| `ch17/library-api` | worker pool, import pipeline and graceful shutdown; imports chapter 16's packages |
| `ch17/patterns` | ticker, `sync.Once`, `sync.Pool`, timeouts and a goroutine leak |
| `ch18` | the gRPC catalog service: `.proto`, generated code in `gen/`, server, interceptors, REST client, bufconn tests |
| `ch19/shelf` | the API ready for containers: SQLite or PostgreSQL, migrations, env config, health probes |
| `ch19/Dockerfile`, `ch19/compose.yaml` | the image and a Compose file with PostgreSQL |
| `.github/workflows/ci.yml` | the CI pipeline from chapter 20 |
| `Makefile` | `make ci`, `make proto`, `make docker` |
| `solutions` | longer answers from Appendix A |

Folders named `broken` hold the mistakes the book walks through. Each file
there starts with `//go:build ignore`, so the normal build skips it, and it
fails to compile on purpose.

## Running the examples

From this folder:

    go run ./ch01/hello
    go test ./...                       # every test in the book
    go test -race ./...                 # the same, with the race detector
    go vet ./...

Programs that read or write files expect to be run from their own folder.
For example:

    cd ch11/finereport && go run .

## The library API (chapters 14 to 17)

The server signs login tokens with a secret that you provide. It must be at
least 32 characters:

    export JWT_SECRET="$(openssl rand -hex 32)"
    go run ./ch15/library-api

It listens on port 8080 and stores data in `library.db` in the current
folder. Chapter 15 has the full `curl` session: sign up, log in, add a book,
borrow it and return it.

The chapter 17 version adds commands and stops cleanly on Ctrl+C:

    go run ./ch17/library-api                       # serve
    go run ./ch17/library-api remind                # overdue emails, logged
    go run ./ch17/library-api import FILE USER_ID   # old catalog export

Set `ADDR=:9090` to listen on another port. The fuzz tests run their seed
corpus with `go test`. To fuzz for real:

    go test -run '^$' -fuzz FuzzParseRupees -fuzztime 10s ./ch16/library-api/fines

## The catalog service (chapter 18)

    go run ./ch18/cmd/catalogd      # gRPC on :9090, data in library.db
    go run ./ch18/cmd/web           # REST on :8081, calls catalogd
    curl localhost:8081/books/1

The generated code in `ch18/gen` is committed. To regenerate it, install
`protoc-gen-go` and `protoc-gen-go-grpc`, then run `make proto` (with `buf`)
or `make proto-protoc` (with `protoc`).

## Containers (chapter 19)

    make docker                                     # or podman build ...
    cd ch19 && JWT_SECRET=$(openssl rand -hex 32) docker compose up

With no `DATABASE_URL` the server uses a SQLite file. With a `postgres://`
URL it uses PostgreSQL. To run the store and API tests on PostgreSQL:

    SHELF_TEST_DSN=postgres://... go test -race -p 1 ./ch19/...

## Checks (chapter 20)

`make ci` runs the same checks as the workflow: gofmt, `go vet`,
staticcheck on the service code, and `go test -race -cover ./...`.

## Dependencies

SQLite comes from `modernc.org/sqlite`, a pure Go driver, so you don't need
cgo to build the API. Tokens use `github.com/golang-jwt/jwt/v5` and passwords
use `golang.org/x/crypto/bcrypt`. Chapter 17 uses `golang.org/x/sync/errgroup`
and `golang.org/x/time/rate`. Chapter 18 uses `google.golang.org/grpc` and
`google.golang.org/protobuf`, and chapter 19 adds `github.com/jackc/pgx/v5`.
The versions are pinned in `go.mod`.

## License

The code is released under the MIT License. See `LICENSE`.
