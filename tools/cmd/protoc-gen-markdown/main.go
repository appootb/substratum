package main

import (
	"github.com/appootb/substratum/v2/tools/cmd/protoc-gen-markdown/generator"
	pgs "github.com/lyft/protoc-gen-star/v2"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.
		Init(pgs.DebugEnv("DEBUG_MARKDOWN"), pgs.SupportedFeatures(&optional)).
		RegisterModule(generator.New()).
		Render()
}
