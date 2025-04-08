# substratum

A foundational framework library for Go that seamlessly integrates gRPC and gRPC-Gateway capabilities.

> [中文](README.zh.md)

## Features

- **Unified gRPC and REST API Support**
  - Implement gRPC methods once, get both gRPC and RESTful API endpoints automatically
  - gRPC Streaming methods are automatically mapped to WebSocket endpoints in gRPC-Gateway

- **Rich Context Features**
  - Built-in error handling
  - Monitoring and metrics
  - Authentication and authorization
  - Logging
  - Parameter validation
  - Service discovery
  - Data storage
  - Message queue support

- **Plugin Architecture**
  - All features are implemented as plugins
  - Easy integration with both open-source and custom services
  - Available plugins can be found at: https://github.com/appootb/plugins

## Quick Start

1. Create a new project:
```bash
mkdir my-service
cd my-service
go mod init my-service
```

2. Installation:
```bash
go get github.com/appootb/substratum/v2
```

3. Define your protobuf service:
```bash
mkdir -p protobuf/proto
cd protobuf/proto
```

```protobuf
syntax = "proto3";

package example;

import "appootb/api/websocket.proto";
import "appootb/permission/method.proto";
import "appootb/permission/service.proto";
import "appootb/validate/validate.proto";
import "google/api/annotations.proto";
import "google/protobuf/empty.proto";

option go_package = "go/example";

message Token {
  string token = 1;
}

message UpStream {
  string message = 1 [
    (appootb.validate.rules).string = {
      min_bytes: 1,
      max_bytes: 31,
    }
  ];
}

message DownStream {
  string message = 1;
}

service MyService {
  option (appootb.permission.service.visible) = CLIENT;

  rpc Login(google.protobuf.Empty) returns (Token) {
    option (google.api.http) = {
      post: "/my-service/v1/login"
      body: "*"
    };
  }

  rpc Stream(stream google.protobuf.Empty) returns (stream DownStream) {
    option (appootb.permission.method.required) = LOGGED_IN;

    option (appootb.api.websocket) = {
      url: "/my-service/v1/streaming"
    };
  }
}
```

4. Generate golang code:
```bash
cd ..
mkdir -p go # make sure $PWD is `my-service/protobuf`
docker run --rm -v $PWD:/mnt -it appootb/grpc-runner:master protoc \
	-Iproto \
	-I/usr/local/include \
	-I/go/src/github.com/googleapis/googleapis \
	-I/go/src/github.com/grpc-ecosystem/grpc-gateway \
	-I/go/src/github.com/appootb/substratum/proto \
	--go_out=paths=source_relative:go \
	--go-grpc_out=paths=source_relative:go \
	--grpc-gateway_out=logtostderr=true,paths=source_relative:go \
	--ootb_out=lang=go,paths=source_relative:go \
	--validate_out=lang=go,paths=source_relative:go \
	proto/*.proto
```

5. Generate MARKDOWN API document(OPTIONAL):
```bash
mkdir -p doc # make sure $PWD is `my-service/protobuf`
docker run --rm -v $PWD:/mnt -it appootb/grpc-runner:master protoc \
	-Iproto \
	-I/usr/local/include \
	-I/go/src/github.com/googleapis/googleapis \
	-I/go/src/github.com/grpc-ecosystem/grpc-gateway \
	-I/go/src/github.com/appootb/substratum/proto \
	--markdown_out=paths=source_relative:doc \
	proto/*.proto
```

6. Implement `MyService`(rpc/example.go):
```go
package rpc

import (
	"context"

	example "my-service/protobuf/go"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Example struct {
	example.UnimplementedMyServiceServer
}

func (s *Example) Login(ctx context.Context, _ *emptypb.Empty) (*example.Token, error) {
	panic("implement me")
}

func (s *Example) Stream(stream example.MyService_StreamServer) error {
	panic("implement me")
}
```

7. Create and register a component(component.go):
```go
package my_service

import (
  "context"

  example "my-service/protobuf/go"
  "my-service/rpc"

  "github.com/appootb/substratum/v2"
  "github.com/appootb/substratum/v2/configure"
  "github.com/appootb/substratum/v2/queue"
  "github.com/appootb/substratum/v2/service"
  "github.com/appootb/substratum/v2/storage"
  "github.com/appootb/substratum/v2/task"
)

type MyComponent struct {
  context.Context
}

func New(ctx context.Context) substratum.Component {
  return &MyComponent{
    Context: ctx,
  }
}

func (m MyComponent) Name() string {
  return "my_service" // A unique service name for service discovery
}

func (m MyComponent) Init(cfg configure.Configure) error {
  return nil
}

func (m MyComponent) InitStorage(s storage.Storage) error {
  return nil
}

func (m MyComponent) RegisterHandler(outer, inner service.HttpHandler) error {
  return nil
}

func (m MyComponent) RegisterService(auth service.Authenticator, srv service.Implementor) error {
  return example.RegisterMyServiceScopeServer(m.Name(), auth, srv, &rpc.Example{})
}

func (m MyComponent) RunQueueWorker(q queue.Queue) error {
  return nil
}

func (m MyComponent) ScheduleCronTask(t task.Task) error {
  return nil
}
```

8. Create main(bootstrap/main.go):
```go
package main

import (
	"fmt"
	"log"

	comp "my-service"

	"github.com/appootb/substratum/v2"
	"github.com/appootb/substratum/v2/context"
	"github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/proto/go/permission"
)

const (
	ClientRpcPort     = 6007
	ClientGatewayPort = 6009
	ServerRpcPort     = 8007
	ServerGatewayPort = 8009
)

func main() {
	// New server instance
	srv := substratum.NewServer(
		substratum.WithServeMux(permission.VisibleScope_CLIENT, ClientRpcPort, ClientGatewayPort),
		substratum.WithServeMux(permission.VisibleScope_SERVER, ServerRpcPort, ServerGatewayPort))

	// Register components
	if err := srv.Register(comp.New(context.Context())); err != nil {
		log.Panicf("register component failed, err: %v", err)
	}

	// Serve
	if err := srv.Serve(metadata.EnvDevelop == "local"); err != nil {
		fmt.Println("exiting...", err.Error())
	}
}
```

* Project Structure
```
my-service/
├── bootstrap/
│   └── main.go         # Service entry point
├── protobuf/
│   ├── doc/            # Generated MarkDown document
│   ├── go/             # Generated Go code
│   └── proto/          # Protocol buffer definitions
├── rpc/                # RPC implementations
└── component.go        # Component definition
```

9. Run the service:
```bash
go run main.go
```

Your service will now be available at:
- gRPC: localhost:6007
- REST API: http://localhost:6009
- Source code: [my-service](https://github.com/appootb/my-service)

## Plugin Usage(bootstrap/main.go)

### Enable database
```go
import _ "github.com/appootb/plugins/v2/storage/sql/mysql"
```

### Enable JSON log format
```go
import _ "github.com/appootb/plugins/v2/logger/json/console"
```

### Enable ETCD

1. Add import:
```go
import (
	_ "github.com/appootb/plugins/v2/configure/backend/etcd"
	_ "github.com/appootb/plugins/v2/configure/backend/etcd"
)
```

2. Add environments:
```bash
# Required environment variables
export COMPONENT="my_service"
export ETCD="http://username:password@etcd-0:2379,etcd-1:2379,etcd-2:2379/base_path"
```

### More [plugins](https://github.com/appootb/plugins)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.