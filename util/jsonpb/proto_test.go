package jsonpb

import (
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestMarshalUnmarshalStruct(t *testing.T) {
	in := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"k": structpb.NewStringValue("v"),
		},
	}
	b, err := Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out structpb.Struct
	if err := Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Fields["k"].GetStringValue() != "v" {
		t.Fatalf("out=%v", out.Fields)
	}
}
