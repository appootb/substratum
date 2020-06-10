package jsonpb

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var Marshaler = &runtime.JSONPb{
	EnumsAsInts:  true,
	EmitDefaults: false,
	OrigName:     true,
}

func Marshal(v interface{}) ([]byte, error) {
	return Marshaler.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return Marshaler.Unmarshal(data, v)
}
