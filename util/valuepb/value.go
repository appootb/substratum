// Package valuepb converts between Go scalars and google.protobuf.Value,
// encoding int64/uint64 as strings to preserve full precision in JSON/structpb.
package valuepb

import (
	"fmt"
	"strconv"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// New converts val to a structpb.Value. int64/uint64 use string encoding;
// other types delegate to structpb.NewValue.
func New(val interface{}) (*structpb.Value, error) {
	switch v := val.(type) {
	case int64:
		return Int64Value(v), nil
	case uint64:
		return Uint64Value(v), nil
	default:
		return structpb.NewValue(v)
	}
}

// BoolValue wraps b as a bool structpb.Value.
func BoolValue(b bool) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_BoolValue{
			BoolValue: b,
		},
	}
}

// Bool reads v as a boolean (non-zero numeric/string/bool kinds are true).
func Bool(v *structpb.Value) bool {
	return int64(Float64(v)) != 0
}

// BoolPtr returns a pointer to Bool(v).
func BoolPtr(v *structpb.Value) *bool {
	return proto.Bool(Bool(v))
}

// Int64Value encodes v as a decimal string Value to avoid JSON number precision loss.
func Int64Value(v int64) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: strconv.FormatInt(v, 10),
		},
	}
}

// Int64 parses v as int64 from number, decimal string, or bool kinds.
func Int64(v *structpb.Value) int64 {
	if v == nil {
		return 0
	}
	switch n := v.Kind.(type) {
	case *structpb.Value_NumberValue:
		return int64(n.NumberValue)
	case *structpb.Value_StringValue:
		if iv, err := strconv.ParseInt(n.StringValue, 10, 64); err == nil {
			return iv
		}
	case *structpb.Value_BoolValue:
		if n.BoolValue {
			return 1
		}
	}
	return 0
}

// Int64Ptr returns a pointer to Int64(v).
func Int64Ptr(v *structpb.Value) *int64 {
	return proto.Int64(Int64(v))
}

// Uint64Value encodes v as a decimal string Value.
func Uint64Value(v uint64) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: strconv.FormatUint(v, 10),
		},
	}
}

// Uint64 parses v as uint64 from number, decimal string, or bool kinds.
func Uint64(v *structpb.Value) uint64 {
	if v == nil {
		return 0
	}
	switch n := v.Kind.(type) {
	case *structpb.Value_NumberValue:
		return uint64(n.NumberValue)
	case *structpb.Value_StringValue:
		if uv, err := strconv.ParseUint(n.StringValue, 10, 64); err == nil {
			return uv
		}
	case *structpb.Value_BoolValue:
		if n.BoolValue {
			return 1
		}
	}
	return 0
}

// Uint64Ptr returns a pointer to Uint64(v).
func Uint64Ptr(v *structpb.Value) *uint64 {
	return proto.Uint64(Uint64(v))
}

// Int32Value encodes v via Int64Value.
func Int32Value(v int32) *structpb.Value {
	return Int64Value(int64(v))
}

// Int32 returns Int64(v) narrowed to int32.
func Int32(v *structpb.Value) int32 {
	return int32(Int64(v))
}

// Int32Ptr returns a pointer to Int32(v).
func Int32Ptr(v *structpb.Value) *int32 {
	return proto.Int32(Int32(v))
}

// Uint32Value encodes v via Uint64Value.
func Uint32Value(v uint32) *structpb.Value {
	return Uint64Value(uint64(v))
}

// Uint32 returns Uint64(v) narrowed to uint32.
func Uint32(v *structpb.Value) uint32 {
	return uint32(Uint64(v))
}

// Uint32Ptr returns a pointer to Uint32(v).
func Uint32Ptr(v *structpb.Value) *uint32 {
	return proto.Uint32(Uint32(v))
}

// Float64Value wraps v as a number Value.
func Float64Value(v float64) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_NumberValue{
			NumberValue: v,
		},
	}
}

// Float64 parses v as float64 from number, string, or bool kinds.
func Float64(v *structpb.Value) float64 {
	if v == nil {
		return 0
	}
	switch n := v.Kind.(type) {
	case *structpb.Value_NumberValue:
		return n.NumberValue
	case *structpb.Value_StringValue:
		if fv, err := strconv.ParseFloat(n.StringValue, 64); err == nil {
			return fv
		}
	case *structpb.Value_BoolValue:
		if n.BoolValue {
			return 1
		}
	}
	return 0
}

// Float64Ptr returns a pointer to Float64(v).
func Float64Ptr(v *structpb.Value) *float64 {
	return proto.Float64(Float64(v))
}

// Float32Value encodes v via Float64Value.
func Float32Value(v float32) *structpb.Value {
	return Float64Value(float64(v))
}

// Float32 returns Float64(v) narrowed to float32.
func Float32(v *structpb.Value) float32 {
	return float32(Float64(v))
}

// Float32Ptr returns a pointer to Float32(v).
func Float32Ptr(v *structpb.Value) *float32 {
	return proto.Float32(Float32(v))
}

// StringValue wraps s as a string Value.
func StringValue(s string) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: s,
		},
	}
}

// String formats v as a string from number, string, or bool kinds.
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

// StringPtr returns a pointer to String(v).
func StringPtr(v *structpb.Value) *string {
	return proto.String(String(v))
}
