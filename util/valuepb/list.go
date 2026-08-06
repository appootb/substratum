package valuepb

import "google.golang.org/protobuf/types/known/structpb"

// BoolListValue encodes l as a list Value of bools.
func BoolListValue(l []bool) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, b := range l {
		values = append(values, BoolValue(b))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// BoolList decodes a list Value into []bool; non-list kinds yield an empty slice.
func BoolList(v *structpb.Value) []bool {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []bool{}
	}
	//
	list := make([]bool, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, Bool(b))
	}
	return list
}

// Int64ListValue encodes l as a list Value of int64 string numbers.
func Int64ListValue(l []int64) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, v := range l {
		values = append(values, Int64Value(v))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// Int64List decodes a list Value into []int64.
func Int64List(v *structpb.Value) []int64 {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []int64{}
	}
	//
	list := make([]int64, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, Int64(b))
	}
	return list
}

// Uint64ListValue encodes l as a list Value of uint64 string numbers.
func Uint64ListValue(l []uint64) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, v := range l {
		values = append(values, Uint64Value(v))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// Uint64List decodes a list Value into []uint64.
func Uint64List(v *structpb.Value) []uint64 {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []uint64{}
	}
	//
	list := make([]uint64, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, Uint64(b))
	}
	return list
}

// Int32ListValue encodes l as a list Value of int32 values.
func Int32ListValue(l []int32) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, v := range l {
		values = append(values, Int32Value(v))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// Int32List decodes a list Value into []int32.
func Int32List(v *structpb.Value) []int32 {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []int32{}
	}
	//
	list := make([]int32, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, Int32(b))
	}
	return list
}

// Uint32ListValue encodes l as a list Value of uint32 values.
func Uint32ListValue(l []uint32) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, v := range l {
		values = append(values, Uint32Value(v))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// Uint32List decodes a list Value into []uint32.
func Uint32List(v *structpb.Value) []uint32 {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []uint32{}
	}
	//
	list := make([]uint32, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, Uint32(b))
	}
	return list
}

// Float64ListValue encodes l as a list Value of float64 numbers.
func Float64ListValue(l []float64) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, v := range l {
		values = append(values, Float64Value(v))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// Float64List decodes a list Value into []float64.
func Float64List(v *structpb.Value) []float64 {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []float64{}
	}
	//
	list := make([]float64, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, Float64(b))
	}
	return list
}

// Float32ListValue encodes l as a list Value of float32 numbers.
func Float32ListValue(l []float32) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, v := range l {
		values = append(values, Float32Value(v))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// Float32List decodes a list Value into []float32.
func Float32List(v *structpb.Value) []float32 {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []float32{}
	}
	//
	list := make([]float32, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, Float32(b))
	}
	return list
}

// StringListValue encodes l as a list Value of strings.
func StringListValue(l []string) *structpb.Value {
	values := make([]*structpb.Value, 0, len(l))
	for _, v := range l {
		values = append(values, StringValue(v))
	}
	return &structpb.Value{
		Kind: &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		},
	}
}

// StringList decodes a list Value into []string.
func StringList(v *structpb.Value) []string {
	kind, ok := v.GetKind().(*structpb.Value_ListValue)
	if !ok || kind.ListValue == nil {
		return []string{}
	}
	//
	list := make([]string, 0, len(kind.ListValue.GetValues()))
	for _, b := range kind.ListValue.GetValues() {
		list = append(list, String(b))
	}
	return list
}
