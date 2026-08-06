// Package jsonpb wraps grpc-gateway JSONPb for protobuf JSON encode/decode.
package jsonpb

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// Marshaler is the shared JSONPb codec (proto field names, enum numbers).
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

// Marshal encodes v as JSON using Marshaler.
func Marshal(v interface{}) ([]byte, error) {
	return Marshaler.Marshal(v)
}

// Unmarshal decodes JSON data into v using Marshaler.
func Unmarshal(data []byte, v interface{}) error {
	return Marshaler.Unmarshal(data, v)
}
