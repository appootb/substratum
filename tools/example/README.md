# tools/example

End-to-end sample protos for the substratum generators.

## Generate

Install plugins first (`cd tools && go install ./cmd/...`), then:

```bash
make gen
```

Produces Go (`*.pb.go`, `*_grpc.pb.go`, `*.pb.gw.go`, `*.pb.ootb.go`, `*.pb.validate.go`),
Markdown docs, and `proto/export.dart`.

## Proto coverage

| File | Syntax | Focus |
|------|--------|--------|
| [`proto/example.proto`](proto/example.proto) | proto2 | optional presence, required fields, custom `error`, map/repeated, oneof required, disabled message, CLIENT/SERVER services, HTTP, bidi + server-stream WebSocket |
| [`proto/example_v3.proto`](proto/example_v3.proto) | proto3 | `optional` fields, `ignore_empty`, wrappers/timestamp/duration, nested message, HTTP, client-stream WebSocket, LOGGED_IN/GUEST/CLIENT/ANY subjects |

### validate highlights

- `(appootb.validate.error)` custom failure text (proto2 `Request.id`)
- `ignore_empty` on string/int/bytes
- `enum.defined_only` / `not_in`
- `message.required`, `oneof` required, `validate.disabled`
- well-known: `email`, `uri`, `uuid`, `timestamp`, `duration`
- generated APIs: `Validate()` (fail-fast) and `ValidateAll()` (`MultiError`)

### ootb highlights

- `FullMethod` / subjects keys use **proto method names** (e.g. `/example.example/test1`)
- VisibleScope `CLIENT` vs `SERVER`
- HTTP via `google.api.http`
- WebSocket via `(appootb.api.websocket).url` on streaming RPCs only

Use `appootb.validate` rules (not upstream `validate.rules`).
