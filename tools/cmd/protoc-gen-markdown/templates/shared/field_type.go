package shared

import (
	"fmt"
	"strconv"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (fns Funcs) fieldElementType(el pgs.FieldTypeElem) string {
	if el.IsEmbed() {
		msg := el.Embed()
		return fmt.Sprintf("[%s](#%s)", fns.Name(msg), fns.anchorName(msg))
	} else if el.IsEnum() {
		enum := el.Enum()
		if enum.FullyQualifiedName() == ".google.protobuf.NullValue" {
			return "null"
		}
		return fmt.Sprintf("[%s](#%s)", fns.Name(enum), fns.anchorName(enum))
	}

	switch el.ProtoType() {
	case pgs.DoubleT, pgs.FloatT:
		return "float"
	case pgs.Int64T, pgs.UInt64T, pgs.SInt64, pgs.Fixed64T, pgs.SFixed64:
		return "int64"
	case pgs.Int32T, pgs.UInt32T, pgs.Fixed32T, pgs.SInt32, pgs.SFixed32:
		return "int"
	case pgs.BoolT:
		return "bool"
	case pgs.StringT:
		return "string"
	case pgs.BytesT:
		return "bytes"
	default:
		return "UNKNOWN"
	}
}

func (fns Funcs) embedJSONType(msg pgs.Message) string {
	switch msg.WellKnownType() {
	case pgs.DurationWKT:
		return `string ("1.000340012s")`
	case pgs.EmptyWKT:
		return `object "{}"`
	case pgs.TimestampWKT:
		return `string ("1972-01-01T10:00:20.021Z")`
	case pgs.ListValueWKT:
		return "array"
	case pgs.DoubleValueWKT, pgs.FloatValueWKT:
		return "number/string"
	case pgs.Int64ValueWKT, pgs.UInt64ValueWKT:
		return "string"
	case pgs.Int32ValueWKT, pgs.UInt32ValueWKT:
		return "number/string"
	case pgs.BoolValueWKT:
		return "true, false"
	case pgs.StringValueWKT:
		return "string"
	case pgs.BytesValueWKT:
		return "base64 string"
	case pgs.AnyWKT, pgs.StructWKT:
		return "object"
	default:
		return "object"
	}
}

func (fns Funcs) fieldType(field pgs.Field) (pbType, jsonType string) {
	switch field.Type().ProtoType() {
	case pgs.DoubleT, pgs.FloatT:
		pbType = "float"
		jsonType = "number/string"
	case pgs.Int64T, pgs.UInt64T, pgs.SInt64, pgs.Fixed64T, pgs.SFixed64:
		pbType = "int64"
		jsonType = "string"
	case pgs.Int32T, pgs.UInt32T, pgs.Fixed32T, pgs.SInt32, pgs.SFixed32:
		pbType = "int"
		jsonType = "number/string"
	case pgs.BoolT:
		pbType = "bool"
		jsonType = "true, false"
	case pgs.StringT:
		pbType = "string"
		jsonType = "string"
	case pgs.BytesT:
		pbType = "bytes"
		jsonType = "base64 string"
	case pgs.EnumT:
		var enum pgs.Enum
		if field.Type().IsRepeated() {
			enum = field.Type().Element().Enum()
		} else {
			enum = field.Type().Enum()
		}
		pbType = fmt.Sprintf("enum [%s](#%s)", fns.Name(enum), fns.anchorName(enum))
		jsonType = "string/integer"
		if enum.FullyQualifiedName() == ".google.protobuf.NullValue" {
			jsonType = "null"
		}
	case pgs.MessageT:
		if field.Type().IsMap() {
			key := fns.fieldElementType(field.Type().Key())
			value := fns.fieldElementType(field.Type().Element())
			pbType = fmt.Sprintf("map\\<%s, %s\\>", key, value)
			jsonType = "object"
		} else if field.Type().IsRepeated() {
			el := fns.fieldElementType(field.Type().Element())
			pbType = el
			jsonType = "array"
		} else {
			msg := field.Type().Embed()
			pbType = fmt.Sprintf("[%s](#%s)", fns.Name(msg), fns.anchorName(msg))
			jsonType = fns.embedJSONType(msg)
		}
	// TODO: deprecated
	case pgs.GroupT:
	}
	if field.InOneOf() {
		pbType = fmt.Sprintf("%s (oneof %s)", pbType, field.OneOf().Name())
	}
	return
}

func (fns Funcs) fieldDoc(field pgs.Field) string {
	columns := []string{
		field.Name().String(),
	}
	// type
	pbType, jsonType := fns.fieldType(field)
	if field.Type().IsRepeated() {
		pbType = fmt.Sprintf("array [%s]", pbType)
	}
	columns = append(columns, pbType, jsonType)
	// validation
	columns = append(columns, fns.fieldRules(field))
	// comment
	columns = append(columns, fns.trailingComment(field.SourceCodeInfo()))
	if field.Syntax().SupportsRequiredPrefix() {
		defaultValue := field.Descriptor().GetDefaultValue()
		if defaultValue == "" {
			defaultValue = "-"
		}
		columns = append(columns, defaultValue)
		columns = append(columns, strconv.FormatBool(field.Required()))
	}
	return "|" + strings.Join(columns, "|") + "|"
}
