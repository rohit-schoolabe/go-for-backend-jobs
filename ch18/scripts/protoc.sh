#!/bin/sh
# Regenerate the Go code with plain protoc, for machines
# without buf. Run from the ch18 directory.
set -eu
out=${1:-gen}
mkdir -p "$out"
protoc -I proto \
  --go_out="$out" --go_opt=paths=source_relative \
  --go-grpc_out="$out" --go-grpc_opt=paths=source_relative \
  proto/shelf/catalog/v1/catalog.proto
