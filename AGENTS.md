# AGENTS.md

protomock is a Go mock server for HTTP/JSON and unary gRPC. Mocks are JavaScript files plus `.proto` definitions loaded from disk at startup.

## Commands

```bash
task run    # go run ./cmd/protomock (expects ./configs/protomock.yaml)
task test   # go test -v -failfast --tags=integration ./...
task lint   # golangci-lint run
task deps   # go mod tidy && go mod vendor && go mod verify
task tools  # install golangci-lint v2.14.0
```

- There are no `*_test.go` files yet; `go test ./...` is currently a compile check.
- CI (`.github/workflows/go.yml`) runs only `go test -v ./...`. Lint is not enforced in CI.
- `vendor/` exists locally but is gitignored and auto-selected by Go. After changing `go.mod`, run `task deps` or expect inconsistent-vendoring errors.
- `.golangci.yml` enables all linters with `fix: true`, so `task lint` may rewrite files.

## Run the example locally

From the repo root (relative mock dirs and `./include` are resolved from the current directory):

```bash
HTTP_SERVER_MOCKSDIR=example/mocks/http \
GRPC_SERVER_MOCKSDIR=example/mocks/grpc \
go run ./cmd/protomock -c example/configs/protomock.yaml
```

Do not run from `example/`: gRPC proto compilation hardcodes the import path `./include` relative to cwd (`internal/transport/grpc/mocks.go`).

## Layout

- `cmd/protomock/main.go` — entrypoint; builds enabled servers, middleware, and loads mocks once.
- `internal/config` — YAML config plus env overrides.
- `internal/container` — application lifecycle and server startup/shutdown.
- `internal/transport/{http,grpc}` — mock discovery, routing, and dynamic gRPC service registration.
- `pkg/` — reusable helpers (`js`, `protobuf/dynamic`, HTTP/gRPC middleware, logging, option, closer).
- `include/` — bundled Google well-known `.proto` files; the Docker image copies this next to the binary.
- `example/` — sample config, HTTP mocks, and gRPC mocks/protos.

## Behavior and gotchas

- The config file is always opened (default `./configs/protomock.yaml`, overridden by `-c`), then env vars override YAML. Env names: `LOG_LEVEL`, `HTTP_SERVER_ENABLED`, `HTTP_SERVER_PORT`, `HTTP_SERVER_MOCKSDIR`, and the matching `GRPC_SERVER_*` set.
- Mocks are loaded only at startup. Restart after editing `.js` or `.proto` files; there is no hot reload.
- HTTP mocks: `<mocksdir>/<path>/<METHOD>.js`; a `__name` directory becomes the route param `:name`. Response is `{status, body}`, where `body` is any JSON value.
- gRPC mocks: `<mocksdir>/<proto package path>/<Service>/<Method>.js`; matching uses the proto's declared package/service/method. Only unary RPCs are supported; streaming is not implemented.
- gRPC bodies go through `protojson`, so JS mocks see/return proto JSON names and rules (camelCase fields, enums as strings). `error` takes precedence over `body`; its `code` is a `google.golang.org/grpc/codes` value.
- Each request gets a fresh goja VM; globals do not persist between requests. `request` and `console` are injected, and `console.log` logs via slog at INFO.
- HTTP request header names are lowercased; values are passed through unchanged.
- Imports are grouped std / third-party / `github.com/sknv/protomock`; max line length is 120. Use `//nolint:name // reason` like the existing code.
