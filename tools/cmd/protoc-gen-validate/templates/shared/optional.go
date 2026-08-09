package shared

import pgs "github.com/lyft/protoc-gen-star/v2"

// Optional reports whether generated Go code should guard rules with a nil
// check for proto2 optional/presence fields (pointers). Repeated, map, oneof,
// and required fields are excluded.
func Optional(f pgs.Field) bool {
	if f.Syntax() != pgs.Proto2 || f.Required() || f.InOneOf() {
		return false
	}
	t := f.Type()
	return !t.IsRepeated() && !t.IsMap()
}
