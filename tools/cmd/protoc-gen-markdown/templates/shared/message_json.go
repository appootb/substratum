package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (fns Funcs) enumJSON(enum pgs.Enum) string {
	var val []string

	for _, v := range enum.Values() {
		val = append(val, fmt.Sprintf("%d (%s)", v.Value(), v.Name().String()))
	}
	return strings.Join(val, " | ")
}

func (fns Funcs) fieldElementJSON(el pgs.FieldTypeElem, messages map[string]int, mapKey bool) string {
	if el.IsEmbed() {
		messages[el.Embed().FullyQualifiedName()]++
		return fns.messageJSON(el.Embed(), messages)
	} else if el.IsEnum() {
		if el.Enum().FullyQualifiedName() == ".google.protobuf.NullValue" {
			return "null"
		}
		return fmt.Sprintf(`"%s"`, fns.enumJSON(el.Enum()))
	}

	switch el.ProtoType() {
	case pgs.DoubleT, pgs.FloatT:
		if mapKey {
			return `"3.1415926"`
		} else {
			return "3.1415926"
		}
	case pgs.Int64T, pgs.UInt64T, pgs.SInt64, pgs.Fixed64T, pgs.SFixed64:
		return `"string($int64)"`
	case pgs.Int32T, pgs.UInt32T, pgs.Fixed32T, pgs.SInt32, pgs.SFixed32:
		if mapKey {
			return `"0"`
		} else {
			return "0"
		}
	case pgs.BoolT:
		if mapKey {
			return `"true"`
		} else {
			return "true"
		}
	case pgs.StringT:
		return `"string"`
	case pgs.BytesT:
		return `"YmFzZTY0IHN0cmluZw=="`
	default:
		return `"UNKNOWN"`
	}
}

func (fns Funcs) fieldJSON(field pgs.Field, messages map[string]int) string {
	switch field.Type().ProtoType() {
	case pgs.DoubleT, pgs.FloatT:
		return `3.1415926`
	case pgs.Int64T, pgs.UInt64T, pgs.SInt64, pgs.Fixed64T, pgs.SFixed64:
		return `"string($int64)"`
	case pgs.Int32T, pgs.UInt32T, pgs.Fixed32T, pgs.SInt32, pgs.SFixed32:
		return "0"
	case pgs.BoolT:
		return "true"
	case pgs.StringT:
		return `"string"`
	case pgs.BytesT:
		return `"YmFzZTY0IHN0cmluZw=="`
	case pgs.EnumT:
		var enum pgs.Enum
		if field.Type().IsRepeated() {
			enum = field.Type().Element().Enum()
		} else {
			enum = field.Type().Enum()
		}
		if enum.FullyQualifiedName() == ".google.protobuf.NullValue" {
			return "null"
		}
		return fmt.Sprintf(`"%s"`, fns.enumJSON(enum))
	case pgs.MessageT:
		if field.Type().IsMap() {
			key := fns.fieldElementJSON(field.Type().Key(), messages, true)
			value := fns.fieldElementJSON(field.Type().Element(), messages, false)
			return fmt.Sprintf(`{%s:%s}`, key, value)
		} else if field.Type().IsRepeated() {
			return fns.fieldElementJSON(field.Type().Element(), messages, false)
		} else {
			return fns.embedJSON(field.Type().Embed(), messages)
		}
	// TODO: deprecated
	//case pgs.GroupT:
	default:
		return ""
	}
}

func (fns Funcs) embedJSON(message pgs.Message, messages map[string]int) string {
	switch message.WellKnownType() {
	case pgs.DurationWKT:
		return `"1.000340012s"`
	case pgs.EmptyWKT:
		return `{}`
	case pgs.TimestampWKT:
		return `"1972-01-01T10:00:20.021Z"`
	case pgs.ListValueWKT:
		return `["foo", "bar"]`
	case pgs.DoubleValueWKT, pgs.FloatValueWKT:
		return `3.1415926`
	case pgs.Int64ValueWKT, pgs.UInt64ValueWKT:
		return `"string($int64)"`
	case pgs.Int32ValueWKT, pgs.UInt32ValueWKT:
		return "0"
	case pgs.BoolValueWKT:
		return "true"
	case pgs.StringValueWKT:
		return `"string"`
	case pgs.BytesValueWKT:
		return `"YmFzZTY0IHN0cmluZw=="`
	case pgs.AnyWKT, pgs.StructWKT:
		return `{"foo": "bar"}`
	default:
		messages[message.FullyQualifiedName()]++
		return fns.messageJSON(message, messages)
	}
}

func (fns Funcs) messageJSON(message pgs.Message, messages map[string]int) string {
	var lines []string
	if messages[message.FullyQualifiedName()] > 2 {
		return "{}"
	}

	for _, field := range message.Fields() {
		val := fns.fieldJSON(field, messages)
		if field.Type().IsRepeated() {
			val = fmt.Sprintf(`[%s]`, val)
		}
		lines = append(lines, fmt.Sprintf(`"%s":%s`, field.Name(), val))
	}
	return fmt.Sprintf(`{%s}`, strings.Join(lines, ","))
}

func (fns Funcs) messageJSONDemo(message pgs.Message) string {
	var lines []string
	for _, field := range message.Fields() {
		val := fns.fieldJSON(field, map[string]int{
			message.FullyQualifiedName(): 1,
		})
		if field.Type().IsRepeated() {
			val = fmt.Sprintf(`[%s]`, val)
		}
		lines = append(lines, fmt.Sprintf("%q:%s", field.Name(), val))
	}
	jsonVal := fmt.Sprintf("{%s}", strings.Join(lines, ","))
	return fns.prettyJSON(jsonVal)
}

func (fns Funcs) wellKnownJSONDemo(message pgs.Message) string {
	jsonVal := "{}"
	switch message.WellKnownType() {
	case pgs.ListValueWKT:
		jsonVal = `["foo", "bar"]`
	case pgs.AnyWKT, pgs.StructWKT:
		jsonVal = `{"foo": "bar"}`
	}
	return fns.prettyJSON(jsonVal)
}

func (fns Funcs) prettyJSON(v string) string {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, []byte(v), "", "  "); err != nil {
		return "json.Indent err:" + err.Error()
	}
	return fmt.Sprintf("```json\n%s\n```", string(pretty.Bytes()))
}
