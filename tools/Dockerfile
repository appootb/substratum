# syntax=docker/dockerfile:1
#
# grpc-runner: protoc + multi-language plugins used with substratum.
# Build from the repository root:
#   docker build -f tools/Dockerfile -t appootb/grpc-runner .
#
# Formerly: https://github.com/appootb/grpc-runner

#
# Go binaries (ecosystem plugins + substratum tools from this monorepo)
FROM --platform=$TARGETPLATFORM golang:1.25 AS go

ARG TARGETPLATFORM

ENV GO111MODULE=on
ENV CGO_ENABLED=0

# protoc-gen-go: https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go
ENV PROTOC_GEN_GO_VER=v1.36.11
# protoc-gen-go-grpc: https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc
ENV PROTOC_GEN_GO_GRPC_VER=v1.6.2
# grpc-gateway (aligned with substratum go.mod): https://github.com/grpc-ecosystem/grpc-gateway/releases
ENV GRPC_GATEWAY_VER=v2.30.0

# Full monorepo: tools/go.mod replace points at the parent framework module.
WORKDIR /go/src/github.com/appootb/substratum
COPY . .

WORKDIR /go/src/github.com/appootb/substratum/tools
RUN go install ./cmd/protoc-gen-ootb \
	./cmd/protoc-gen-validate \
	./cmd/protoc-gen-markdown \
	./cmd/protoc-gen-dart-export

# Well-known / third-party .proto include trees used by consumer protoc invocations.
RUN git clone --depth 1 https://github.com/googleapis/googleapis.git /go/src/github.com/googleapis/googleapis && \
	git clone --depth 1 -b ${GRPC_GATEWAY_VER} https://github.com/grpc-ecosystem/grpc-gateway.git /go/src/github.com/grpc-ecosystem/grpc-gateway

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VER} && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VER} && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@${GRPC_GATEWAY_VER} && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@${GRPC_GATEWAY_VER}

#
# Swift
FROM --platform=$TARGETPLATFORM swift:6.2-jammy AS swift

# protoc-gen-swift: https://github.com/apple/swift-protobuf/releases
ENV SWIFT_PLUGIN_VER=1.38.1

# Only build the protoc plugin product (default `swift build` also links the heavy
# `protoc` / protoc-tool product and often gets SIGKILL/OOM on CI / QEMU).
# -j 1 further caps peak memory during static stdlib linking.
RUN git clone --depth 1 -b ${SWIFT_PLUGIN_VER} https://github.com/apple/swift-protobuf && \
	cd swift-protobuf && \
	swift build --product protoc-gen-swift --static-swift-stdlib -c release -j 1

#
# CSharp
FROM --platform=$TARGETPLATFORM mcr.microsoft.com/dotnet/sdk:8.0-jammy AS csharp

ARG TARGETPLATFORM

# grpc_csharp_plugin: https://www.nuget.org/packages/Grpc.Tools
ENV CS_PLUGIN_VER=2.83.0

RUN case ${TARGETPLATFORM} in \
		"linux/amd64") CS_ARCH=x64 ;; \
		"linux/arm64") CS_ARCH=arm64 ;; \
		*) echo "unsupported TARGETPLATFORM=${TARGETPLATFORM}" >&2; exit 1 ;; \
	esac && \
	cd /root && dotnet new console && \
	dotnet add package Grpc.Tools --version ${CS_PLUGIN_VER} && \
	cp /root/.nuget/packages/grpc.tools/${CS_PLUGIN_VER}/tools/linux_${CS_ARCH}/grpc_csharp_plugin /usr/local/bin/grpc_csharp_plugin

#
# Downloads (protoc, Dart SDK, Java plugin)
FROM --platform=$TARGETPLATFORM ubuntu:jammy AS binary

ARG TARGETPLATFORM

# Protocol Buffers: https://github.com/protocolbuffers/protobuf/releases
ENV PROTOC_VER=35.1
# protoc-gen-grpc-java: https://mvnrepository.com/artifact/io.grpc/protoc-gen-grpc-java
ENV JAVA_GRPC_VER=1.82.2
# Dart SDK: https://dart.dev/get-dart/archive
ENV DART_SDK_VER=3.12.2

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl unzip && \
	rm -rf /var/lib/apt/lists/*

# protoc
RUN case ${TARGETPLATFORM} in \
		"linux/amd64") PROTOC_ARCH=x86_64 ;; \
		"linux/arm64") PROTOC_ARCH=aarch_64 ;; \
		*) echo "unsupported TARGETPLATFORM=${TARGETPLATFORM}" >&2; exit 1 ;; \
	esac && \
	curl -fsSL -o protoc.zip \
		"https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-${PROTOC_ARCH}.zip" && \
	unzip -q protoc.zip -d protoc3 && \
	rm protoc.zip

# Dart SDK
RUN case ${TARGETPLATFORM} in \
		"linux/amd64") DART_ARCH=x64 ;; \
		"linux/arm64") DART_ARCH=arm64 ;; \
		*) echo "unsupported TARGETPLATFORM=${TARGETPLATFORM}" >&2; exit 1 ;; \
	esac && \
	curl -fsSL -o dart.zip \
		"https://storage.googleapis.com/dart-archive/channels/stable/release/${DART_SDK_VER}/sdk/dartsdk-linux-${DART_ARCH}-release.zip" && \
	unzip -q dart.zip -d dart && \
	rm dart.zip

# protoc-gen-grpc-java
RUN case ${TARGETPLATFORM} in \
		"linux/amd64") JAVA_ARCH=x86_64 ;; \
		"linux/arm64") JAVA_ARCH=aarch_64 ;; \
		*) echo "unsupported TARGETPLATFORM=${TARGETPLATFORM}" >&2; exit 1 ;; \
	esac && \
	curl -fsSL -o /usr/local/bin/protoc-gen-grpc-java \
		"https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/${JAVA_GRPC_VER}/protoc-gen-grpc-java-${JAVA_GRPC_VER}-linux-${JAVA_ARCH}.exe" && \
	chmod +x /usr/local/bin/protoc-gen-grpc-java

#
# Runner
FROM --platform=$TARGETPLATFORM ubuntu:jammy

# protoc-gen-dart: https://pub.dev/packages/protoc_plugin
ENV DART_GRPC_VER=25.0.0

ENV GOPATH=/go
ENV PATH=/usr/lib/dart/bin:$PATH

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates make && \
	rm -rf /var/lib/apt/lists/*

# Go binaries (ecosystem + substratum tools)
COPY --from=go /go/bin/ /usr/local/bin/

# Proto include trees expected by docs / consumer Makefiles
COPY --from=go /go/src/github.com/googleapis/googleapis /go/src/github.com/googleapis/googleapis
COPY --from=go /go/src/github.com/appootb/substratum/proto/appootb /go/src/github.com/appootb/substratum/proto/appootb
COPY --from=go /go/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2/options /go/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2/options

# protoc-gen-swift
COPY --from=swift /swift-protobuf/.build/release/protoc-gen-swift /usr/local/bin/protoc-gen-swift

# grpc_csharp_plugin
COPY --from=csharp /usr/local/bin/grpc_csharp_plugin /usr/local/bin/grpc_csharp_plugin

# protoc
COPY --from=binary /protoc3/bin/ /usr/local/bin/
COPY --from=binary /protoc3/include/ /usr/local/include/

# protoc-gen-grpc-java
COPY --from=binary /usr/local/bin/protoc-gen-grpc-java /usr/local/bin/protoc-gen-grpc-java

# protoc-gen-dart
COPY --from=binary /dart/dart-sdk/ /usr/lib/dart/

RUN dart pub global activate protoc_plugin ${DART_GRPC_VER} && \
	ln -sf /root/.pub-cache/bin/protoc-gen-dart /usr/local/bin/protoc-gen-dart

WORKDIR /mnt
