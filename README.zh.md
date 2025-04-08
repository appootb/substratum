# substratum

一个基础 Go 框架库，无缝集成 gRPC 和 gRPC-Gateway 功能。

> [English](README.md)

## 特性

- **统一的 gRPC 和 REST API 支持**
  - 只需实现一次 gRPC 方法，即可自动获得 gRPC 和 RESTful API 端点
  - gRPC Streaming 方法自动映射为 gRPC-Gateway 中的 WebSocket 端点

- **丰富的上下文功能**
  - 内置错误处理
  - 监控和指标
  - 认证和授权
  - 日志记录
  - 参数验证
  - 服务发现
  - 数据存储
  - 消息队列支持

- **插件架构**
  - 所有功能都通过插件实现
  - 轻松集成开源和自定义服务
  - 可用插件请访问：https://github.com/appootb/plugins

## 快速开始

1. 创建新项目：
```bash
mkdir my-service
cd my-service
go mod init my-service
```

2. 安装：
```bash
go get github.com/appootb/substratum/v2
```

3. 定义 protobuf 服务：
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

4. 生成 Go 代码：
```bash
cd ..
mkdir -p go # 确保当前目录是 `my-service/protobuf`
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

5. 生成 MARKDOWN API 文档（可选）：
```bash
mkdir -p doc # 确保当前目录是 `my-service/protobuf`
docker run --rm -v $PWD:/mnt -it appootb/grpc-runner:master protoc \
	-Iproto \
	-I/usr/local/include \
	-I/go/src/github.com/googleapis/googleapis \
	-I/go/src/github.com/grpc-ecosystem/grpc-gateway \
	-I/go/src/github.com/appootb/substratum/proto \
	--markdown_out=paths=source_relative:doc \
	proto/*.proto
```

6. 实现 `