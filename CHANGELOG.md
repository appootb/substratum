# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Work on `master` after `v2.0.4` (dependency upgrades, security hardening, and framework performance). Not yet tagged.

### Added

- Framework interceptor benchmarks (`rpc`, `metadata`, `logger`, `auth`)
- `context.WithImplementors` and a single interceptor for client/discovery/storage/queue/task injection
- Process-local JWT verification cache (`plugin/token` ParseRaw), with `SetParseCacheEnabled` / TTL controls
- Access-log and metrics global switches: `logger.SetAccessLogEnabled`, `monitor.SetMetricsEnabled` (default on)
- Panic recovery sanitization: generic Internal messages outside `DEVELOPMENT`
- REST token via header: drop query `token` for REST; `Authorization: Bearer`; `WebsocketQueryMetadata`
- Access-log secret redaction (`logger.RedactSecret`)
- JWT `cty` must match request platform; token plugin on `github.com/appootb/jwt`
- Util improvements: LFU cache, timers, pools, and expanded tests
- Nested `tools/` module (from `appootb/grpc-gen`): `protoc-gen-ootb`, `protoc-gen-validate`, `protoc-gen-markdown`, `protoc-gen-dart-export`
- `tools/Dockerfile` + `.github/workflows/docker-publish.yml` (from `appootb/grpc-runner`): multi-arch `appootb/grpc-runner` image built from this monorepo
- `protoc-gen-validate` aligned with PGV v1.3.3 Go features: `ValidateAll` / `MultiError`, `ignore_empty`, proto3 optional (`SyntheticOneOfFields`), compile-time rule checks
- `appootb.validate` proto gains `ignore_empty` on scalar/string/bytes/repeated/map/any rules (keeps custom `error` field option)
- `protoc-gen-ootb`: FullMethod/subjects align with go-grpc proto method names; `FEATURE_PROTO3_OPTIONAL`; `DEBUG_OOTB`; sorted targets; skip files without services; websocket url/streaming checks
- `tools/example`: richer proto2/proto3 demos for validate + ootb (HTTP/WebSocket, subjects, ignore_empty, wrappers/WKT)
- `protoc-gen-markdown` / `protoc-gen-dart-export`: declare `FEATURE_PROTO3_OPTIONAL`

### Changed

- Dependency upgrades (gorm, grpc-ecosystem middleware/prometheus, related stack)
- Interceptor chain composition uses a single per-request continuation
- Metadata parsing is lazy (holder + first `IncomingMetadata` parse); stream parses once
- Errors interceptor: success path only sets `timestamp` header (code/message on errors)
- WebSocket registration codegen uses `runtime.ServeMux.HandlePath` (grpc-gateway v2 public API)
- `tools/` layout: generators under `tools/cmd/protoc-gen-*`, sample under `tools/example/`

## [2.0.4] - 2026-08-07

### Added

- Phase-based gRPC interceptor registration: `rpc.RegisterUnary` / `RegisterStream` (and `Set*`) with slots merged at `plugin.Register` via `BuildInterceptorChains`
- Duplicate plugin implementor guards: second `RegisterImplementor` panics instead of silently overwriting

### Changed

- Default unary/stream chain assembly moved from option helpers into phase merge (`rpc/phase.go`)

## [2.0.3] - 2025-04-08

### Added

- Locale-aware error messages: `errors.Prompter.Translate` now takes a language code and resolves locale from request metadata
- `valuepb.New()` helper to construct `structpb.Value`, with dedicated paths for `int64` / `uint64`
- Service method context: inject `service.Method` (including `FullMethod`) into HTTP/gRPC request context; `model.Base.ServiceMethod()` accessor
- Configurable `Date` / `DateTime` wrappers with custom layout support and default layouts (`2006-01-02`, `2006-01-02 15:04:05`)

### Changed

- `dao.Base` and `model.Base` methods use pointer receivers
- `ClientSeed.Lock`: when `duration <= 0`, default lock duration is ~100 years
- Internal credential / resolver errors use `status.Error` instead of the custom error helper
- Discovery empty-address report code: `NotFound` → `Unavailable`
- Dependency bumps: `golang.org/x/net`, `google.golang.org/protobuf`
- Expanded English README and added Chinese README (`README.zh.md`)

### Fixed

- Date/DateTime parse: fall back to default layout when layout is empty
- Date/DateTime parse: empty string no longer returns an error

## [2.0.2] - 2022-11-15

### Added

- Default `round_robin` balancer package and plugin
- Built-in configure / discovery plugin implementations (removed runtime dependency on embedding `grc` wrappers)
- Empty storage dialect registration for safer defaults
- Discovery refreshes resolver addresses on init and on service updates

### Changed

- Plugin initialize sequence: balancer → discovery → resolver (and related ordering tweaks)
- `discovery.GetAddresses` returns `[]resolver.Address` instead of `[]string`
- Config/discovery path options moved into plugins; path-related options removed from public configure/discovery options
- `github.com/appootb/grc` upgraded to v1.0.0, then decoupled from the core module path as plugins took over config/discovery backends
- Context helpers consolidated under `context` (server context wiring)

### Removed

- In-tree `grc` wrapper packages under plugins (configure / discovery)

## [2.0.1] - 2022-09-28

### Changed

- Queue initialize logic simplified: removed `Queue.Initialize` / backend `Init` from the public queue interfaces

### Removed

- `pulsar+ssl` configure schema constant
- Explicit queue backend initialize API surface (`Init` / `Initialize` / `InitBackend`)

## [2.0.0] - 2022-09-21

Major release: module path moves to `github.com/appootb/substratum/v2`.

### Changed

- **Upgrade `grpc-gateway` to v2.x** (`github.com/grpc-ecosystem/grpc-gateway` v1.16.0 → `github.com/grpc-ecosystem/grpc-gateway/v2` v2.10.3), and align related gRPC / protobuf stack dependencies
- Module path `/v2` (`github.com/appootb/substratum/v2`)
- Protos and related APIs aligned to v2
- Metadata naming updates
- DAO model updates
- Gateway JSON marshaling and HTTP/gRPC gateway wiring for the v2 gateway stack
- Client context / connection pool behavior for the v2 stack
- Bump companion deps: `google.golang.org/grpc`, `google.golang.org/protobuf`, `prometheus/client_golang`, `gorm.io/gorm`, etc.

### Added

- Queue initialize sequence and `Queue.Initialize` (later refined in 2.0.1)
- `valuepb` utility package
- DAO option helpers and `dao.Base` transaction support
- Expanded queue backend interface upgrades

### Fixed

- HTTP error status handling
- Resolver `Target` empty-address edge case
- LRU load issues (pre-tag fix included in the v2 line)

### Removed

- Direct dependencies no longer required at the core module boundary (e.g. embedded elasticsearch client packages from v1 `go.mod`)

[Unreleased]: https://github.com/appootb/substratum/compare/v2.0.4...HEAD
[2.0.4]: https://github.com/appootb/substratum/compare/v2.0.3...v2.0.4
[2.0.3]: https://github.com/appootb/substratum/compare/v2.0.2...v2.0.3
[2.0.2]: https://github.com/appootb/substratum/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/appootb/substratum/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/appootb/substratum/releases/tag/v2.0.0
