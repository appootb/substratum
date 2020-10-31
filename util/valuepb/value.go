package valuepb

import (
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func BoolValue(b bool) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_BoolValue{
			BoolValue: b,
		},
	}
}

func Bool(v *structpb.Value) bool {
	return int64(Float64(v)) != 0
}

func BoolPtr(v *structpb.Value) *bool {
	return proto.Bool(Bool(v))
}

func Int64Value(v int64) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_NumberValue{
			NumberValue: float64(v),
		},
	}
}

func Int64(v *structpb.Value) int64 {
	return int64(Float64(v))
}

func Int64Ptr(v *structpb.Value) *int64 {
	return proto.Int64(Int64(v))
}

func Uint64Value(v uint64) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_NumberValue{
			NumberValue: float64(v),
		},
	}
}

func Uint64(v *structpb.Value) uint64 {
	return uint64(Float64(v))
}

func Uint64Ptr(v *structpb.Value) *uint64 {
	return proto.Uint64(Uint64(v))
}

func Int32Value(v int32) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_NumberValue{
			NumberValue: float64(v),
		},
	}
}

func Int32(v *structpb.Value) int32 {
	return int32(Float64(v))
}

func Int32Ptr(v *structpb.Value) *int32 {
	return proto.Int32(Int32(v))
}

func Uint32Value(v uint32) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_NumberValue{
			NumberValue: float64(v),
		},
	}
}

func Uint32(v *structpb.Value) uint32 {
	return uint32(Float64(v))
}

func Uint32Ptr(v *structpb.Value) *uint32 {
	return proto.Uint32(Uint32(v))
}

func Float64Value(v float64) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_NumberValue{
			NumberValue: v,
		},
	}
}

func Float64(v *structpb.Value) float64 {
	if v == nil {
		return 0
	}
	switch n := v.Kind.(type) {
	case *structpb.Value_NumberValue:
		return n.NumberValue
	case *structpb.Value_BoolValue:
		if n.BoolValue {
			return 1
		}
	}
	return 0
}

func Float64Ptr(v *structpb.Value) *float64 {
	return proto.Float64(Float64(v))
}

func StringValue(s string) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: s,
		},
	}
}

func String(v *structpb.Value) string {
	if v == nil {
		return ""
	}
	switch n := v.Kind.(type) {
	case *structpb.Value_NumberValue:
		return fmt.Sprintf("%v", n.NumberValue)
	case *structpb.Value_StringValue:
		return n.StringValue
	case *structpb.Value_BoolValue:
		return strconv.FormatBool(n.BoolValue)
	}
	return ""
}

func StringPtr(v *structpb.Value) *string {
	return proto.String(String(v))
}
