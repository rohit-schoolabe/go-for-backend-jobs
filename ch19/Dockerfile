# Build context is the code/ directory, where go.mod is:
#   podman build -f ch19/Dockerfile -t shelf:dev .

# Stage 1: compile with the full Go toolchain.
FROM docker.io/library/golang:1.22.12-bookworm AS build
WORKDIR /src

# Modules first, so this layer stays cached until
# go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO off gives a static binary that needs no libc.
# modernc SQLite and pgx are both pure Go, so this works.
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/shelf ./ch19/shelf

# Stage 2: only the binary, on a base with no shell.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/shelf /shelf
# Writable home for a SQLite file, if one is used.
WORKDIR /home/nonroot
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/shelf"]
