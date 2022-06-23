package jsonpb

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

var Marshaler = &runtime.JSONPb{
	MarshalOptions: protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: false,
	},
	UnmarshalOptions: protojson.UnmarshalOptions{
		DiscardUnknown: true,
	},
}

func Marshal(v interface{}) ([]byte, error) {
	return Marshaler.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return Marshaler.Unmarshal(data, v)
}
