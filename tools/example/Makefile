## Shipping package
PROJECT_NAME = example

_cleango:
	@rm -rf proto/*.go
	@rm -rf proto/*.md

_cleandart:
	@rm -f proto/export.dart

_gengo:
	@protoc -Iproto \
		-I${GOPATH}/src/github.com/appootb/substratum/proto \
		-I${GOPATH}/src/github.com/googleapis/googleapis \
		--go_out=paths=source_relative:proto \
		--go-grpc_out=paths=source_relative:proto \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:proto \
		--ootb_out=lang=go,paths=source_relative:proto \
		--validate_out=lang=go,paths=source_relative:proto \
		--markdown_out=paths=source_relative:proto \
		proto/*.proto

# dart-export only needs descriptors; it does not require protoc-gen-dart.
_gendart:
	@protoc -Iproto \
		-I${GOPATH}/src/github.com/appootb/substratum/proto \
		-I${GOPATH}/src/github.com/googleapis/googleapis \
		--dart-export_out=file=export.dart,prefix=./src/generated,paths=source_relative:proto \
		proto/*.proto

gen: clean _gengo _gendart

clean: _cleango _cleandart
