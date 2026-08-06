package valuepb

import (
	"github.com/appootb/substratum/v2/util/jsonpb"
	"google.golang.org/protobuf/types/known/structpb"
)

// StructValue builds a struct Value from a Go map via JSON marshal/unmarshal.
// On failure it returns an empty struct Value rather than an error.
func StructValue(kvs map[string]interface{}) *structpb.Value {
	var wrapper *structpb.Struct
	val, err := jsonpb.Marshal(kvs)
	if err != nil {
		goto Error
	}
	if err = jsonpb.Unmarshal(val, &wrapper); err != nil {
		goto Error
	}
	return &structpb.Value{
		Kind: &structpb.Value_StructValue{
			StructValue: wrapper,
		},
	}

Error:
	return &structpb.Value{
		Kind: &structpb.Value_StructValue{
			StructValue: &structpb.Struct{
				Fields: map[string]*structpb.Value{},
			},
		},
	}
}

// StructFields returns the field map of a struct Value, or an empty map if v is not a struct.
func StructFields(v *structpb.Value) map[string]*structpb.Value {
	sv, ok := v.GetKind().(*structpb.Value_StructValue)
	if !ok || sv.StructValue == nil {
		return map[string]*structpb.Value{}
	}
	return sv.StructValue.GetFields()
}
