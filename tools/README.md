# substratum tools

Protocol buffer code generators for the [substratum](https://github.com/appootb/substratum) framework.

These plugins previously lived in [`appootb/grpc-gen`](https://github.com/appootb/grpc-gen) and are now
versioned with the framework under a nested Go module so runtime consumers are not forced to pull
generator-only dependencies.

## Layout

```text
tools/
  cmd/                          # installable protoc plugins
    protoc-gen-ootb/
    protoc-gen-validate/
    protoc-gen-markdown/
    protoc-gen-dart-export/
  example/                      # end-to-end sample protos + Makefile
  Dockerfile                    # appootb/grpc-runner image (build from repo root)
  go.mod                        # nested module: github.com/appootb/substratum/v2/tools
```

## Plugins

| Binary | Purpose |
|--------|---------|
| `protoc-gen-ootb` | Service registration, gateway/WebSocket wrappers, permission metadata; debug with `DEBUG_OOTB=1` |
| `protoc-gen-validate` | Message validation from `appootb.validate` rules (`Validate` / `ValidateAll`, `ignore_empty`, proto2/proto3 optional); debug with `DEBUG_VALIDATE=1` |
| `protoc-gen-markdown` | API documentation (Markdown); debug with `DEBUG_MARKDOWN=1` |
| `protoc-gen-dart-export` | Dart barrel/export helpers (`file`, `prefix`) |

## Install

From a published tag (nested module tags use the `tools/` prefix, e.g. `tools/v0.1.0`):

```bash
go install github.com/appootb/substratum/v2/tools/cmd/protoc-gen-ootb@latest
go install github.com/appootb/substratum/v2/tools/cmd/protoc-gen-validate@latest
go install github.com/appootb/substratum/v2/tools/cmd/protoc-gen-markdown@latest
go install github.com/appootb/substratum/v2/tools/cmd/protoc-gen-dart-export@latest
```

From a local checkout of this repository:

```bash
cd tools
go install ./cmd/protoc-gen-ootb ./cmd/protoc-gen-validate ./cmd/protoc-gen-markdown ./cmd/protoc-gen-dart-export
```

### `protoc-gen-dart-export` parameters

| Param | Default | Meaning |
|-------|---------|---------|
| `file` | `export.dart` | Output barrel filename (`.dart` appended when omitted) |
| `prefix` | `./src/generated` | Relative import prefix for generated Dart files |

Example:

```bash
protoc ... --dart-export_out=file=export.dart,prefix=./src/generated:out_dir
```

Debug with `DEBUG_DART_EXPORT=1`.

### `protoc-gen-ootb` notes

- Permission subjects/roles keys and interceptor `FullMethod` use the **proto method name** (same as `protoc-gen-go-grpc`, e.g. `/pkg.Service/test1`).
- `(appootb.api.websocket)` requires a non-empty `url` and a streaming RPC.
- Files without services produce no `*.ootb.go` output.

## Example

```bash
cd tools/example
make gen
```

See [`example/README.md`](example/README.md) for the proto2/proto3 coverage matrix
(`Validate` / `ValidateAll`, HTTP, WebSocket, permission subjects).

Rules must use the `appootb.validate` package (not upstream `validate.rules`).

## Docker image (`appootb/grpc-runner`)

The former [`appootb/grpc-runner`](https://github.com/appootb/grpc-runner) image is built from
[`tools/Dockerfile`](Dockerfile) in this repository (CI: `.github/workflows/docker-publish.yml`).

It bundles `protoc`, official/language plugins, and the four substratum generators above, plus
include trees for `appootb` / `googleapis` / `grpc-gateway` protos.

Local build (from the **repository root**):

```bash
docker build -f tools/Dockerfile -t appootb/grpc-runner:local .
```

Usage (same as before):

```bash
docker run --rm -v "$PWD":/mnt -it appootb/grpc-runner:master protoc \
  -Iproto \
  -I/usr/local/include \
  -I/go/src/github.com/googleapis/googleapis \
  -I/go/src/github.com/appootb/substratum/proto \
  --go_out=paths=source_relative:go \
  --go-grpc_out=paths=source_relative:go \
  --grpc-gateway_out=logtostderr=true,paths=source_relative:go \
  --ootb_out=lang=go,paths=source_relative:go \
  --validate_out=lang=go,paths=source_relative:go \
  proto/*.proto
```

Published tags: `master` / `latest` on pushes to `master`, plus semver from `v*.*.*` and `tools/v*.*.*` tags.

## Module layout

```text
github.com/appootb/substratum/v2              # runtime framework (repo root)
github.com/appootb/substratum/v2/tools        # this module (generators only)
github.com/appootb/substratum/v2/tools/cmd/*  # plugin main packages
```

`tools/go.mod` depends on the parent framework module and uses a local `replace` for development.
Published installs resolve the required framework version from the module proxy.
