# The same checks CI runs, for your own machine.
# Run `make ci` before you push.

GO      ?= go
PKGS    := ./...
# Service code. Chapters 1 to 13 hold small demos, some
# with deliberate bugs that a linter would flag.
LINT    := ./ch14/... ./ch15/... ./ch16/... ./ch17/... \
           ./ch18/... ./ch19/...
IMAGE   ?= localhost/shelf:dev
ENGINE  ?= $(shell command -v docker || command -v podman)

.PHONY: ci fmt-check vet lint test proto proto-protoc docker

ci: fmt-check vet lint test

fmt-check:
	@out=$$(find . -name '*.go' -not -path '*/broken/*' \
	  | xargs gofmt -l); if [ -n "$$out" ]; then \
	  echo "gofmt needed:"; echo "$$out"; exit 1; fi

vet:
	$(GO) vet $(PKGS)

lint:
	staticcheck $(LINT)

test:
	$(GO) test -race -cover $(PKGS)

# Regenerate the gRPC code from ch18/proto with buf.
proto:
	cd ch18 && buf lint && buf generate

# The same, with plain protoc.
proto-protoc:
	cd ch18 && ./scripts/protoc.sh

docker:
	$(ENGINE) build -f ch19/Dockerfile -t $(IMAGE) .
