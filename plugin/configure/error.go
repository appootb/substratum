package configure

import "reflect"

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "substratum: config type(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "substratum: config type(non-pointer " + e.Type.String() + ")"
	}
	return "substratum: config type(nil " + e.Type.String() + ")"
}
