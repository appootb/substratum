package valuepb

import (
	"fmt"
	"strconv"

	structpb "github.com/golang/protobuf/ptypes/struct"
)

func Int64(v *structpb.Value) int64 {
	return int64(Float64(v))
}

func Int32(v *structpb.Value) int32 {
	return int32(Float64(v))
}

func Float64(v *structpb.Value) float64 {
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

func String(v *structpb.Value) string {
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
