package valuepb

import "testing"

func TestStructValueAndFields(t *testing.T) {
	v := StructValue(map[string]interface{}{
		"name": "alice",
		"age":  30,
	})
	fields := StructFields(v)
	if fields["name"] == nil || String(fields["name"]) != "alice" {
		t.Fatalf("fields=%v", fields)
	}
}

func TestStructFieldsWrongKind(t *testing.T) {
	fields := StructFields(StringValue("x"))
	if len(fields) != 0 {
		t.Fatalf("want empty, got %v", fields)
	}
}
