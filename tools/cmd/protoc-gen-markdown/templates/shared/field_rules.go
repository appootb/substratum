package shared

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/appootb/substratum/v2/proto/go/validate"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (fns Funcs) valueCmp(v1, v2 reflect.Value) int {
	switch v1.Kind() {
	case reflect.Float32, reflect.Float64:
		return big.NewFloat(v1.Float()).Cmp(big.NewFloat(v2.Float()))
	case reflect.Int32, reflect.Int64:
		switch {
		case v1.Int() < v2.Int():
			return -1
		case v1.Int() > v2.Int():
			return 1
		}
	case reflect.Uint32, reflect.Uint64:
		switch {
		case v1.Uint() < v2.Uint():
			return -1
		case v1.Uint() > v2.Uint():
			return 1
		}
	default:
	}
	return 0
}

func (fns Funcs) ruleIs(v interface{}) string {
	return fmt.Sprintf("IS: `%v`", v)
}

func (fns Funcs) ruleIn(in interface{}, not bool) string {
	list := reflect.ValueOf(in)
	//
	var values []string
	for i := 0; i < list.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", list.Index(i).Interface()))
	}
	v := fmt.Sprintf("IN: `[%s]`", strings.Join(values, ", "))
	if not {
		v = "NOT " + v
	}
	return v
}

func (fns Funcs) isUnsigned(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func (fns Funcs) minRule(gt, gte interface{}) (string, bool) {
	v, eqV := reflect.ValueOf(gt), reflect.ValueOf(gte)
	if v.IsNil() && eqV.IsNil() {
		if fns.isUnsigned(v.Type().Elem()) {
			return "0", true
		}
		return "-INF", false
	}
	if v.IsNil() {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	} else if eqV.IsNil() {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	}
	//
	if fns.valueCmp(v.Elem(), eqV.Elem()) > 0 {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	} else {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	}
}

func (fns Funcs) maxRule(lt, lte interface{}) (string, bool) {
	v, eqV := reflect.ValueOf(lt), reflect.ValueOf(lte)
	if v.IsNil() && eqV.IsNil() {
		return "+INF", false
	}
	if v.IsNil() {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	} else if eqV.IsNil() {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	}
	//
	if fns.valueCmp(v.Elem(), eqV.Elem()) > 0 {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	} else {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	}
}

func (fns Funcs) ruleInterval(lower, upper string, lowerEq, upperEq bool) string {
	if lowerEq {
		lower = "[" + lower
	} else {
		lower = "(" + lower
	}
	if upperEq {
		upper = upper + "]"
	} else {
		upper = upper + ")"
	}
	return fmt.Sprintf("RANGE: `%s, %s`", lower, upper)
}

func (fns Funcs) floatRules(f *validate.FloatRules) string {
	var rules []string
	//
	if f.Gt != nil || f.Gte != nil || f.Lt != nil || f.Lte != nil {
		lo, loEq := fns.minRule(f.Gt, f.Gte)
		hi, hiEq := fns.maxRule(f.Lt, f.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if f.Const != nil {
		rules = append(rules, fns.ruleIs(f.GetConst()))
	}
	if len(f.In) > 0 {
		rules = append(rules, fns.ruleIn(f.In, false))
	}
	if len(f.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(f.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) doubleRules(d *validate.DoubleRules) string {
	var rules []string
	//
	if d.Gt != nil || d.Gte != nil || d.Lt != nil || d.Lte != nil {
		lo, loEq := fns.minRule(d.Gt, d.Gte)
		hi, hiEq := fns.maxRule(d.Lt, d.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if d.Const != nil {
		rules = append(rules, fns.ruleIs(d.GetConst()))
	}
	if len(d.In) > 0 {
		rules = append(rules, fns.ruleIn(d.In, false))
	}
	if len(d.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(d.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) int32Rules(i *validate.Int32Rules) string {
	var rules []string
	//
	if i.Gt != nil || i.Gte != nil || i.Lt != nil || i.Lte != nil {
		lo, loEq := fns.minRule(i.Gt, i.Gte)
		hi, hiEq := fns.maxRule(i.Lt, i.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if i.Const != nil {
		rules = append(rules, fns.ruleIs(i.GetConst()))
	}
	if len(i.In) > 0 {
		rules = append(rules, fns.ruleIn(i.In, false))
	}
	if len(i.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(i.NotIn, true))
	}
	if i.GetIgnoreEmpty() {
		rules = append(rules, "IGNORE_EMPTY")
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) int64Rules(i *validate.Int64Rules) string {
	var rules []string
	//
	if i.Gt != nil || i.Gte != nil || i.Lt != nil || i.Lte != nil {
		lo, loEq := fns.minRule(i.Gt, i.Gte)
		hi, hiEq := fns.maxRule(i.Lt, i.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if i.Const != nil {
		rules = append(rules, fns.ruleIs(i.GetConst()))
	}
	if len(i.In) > 0 {
		rules = append(rules, fns.ruleIn(i.In, false))
	}
	if len(i.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(i.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) uint32Rules(u *validate.UInt32Rules) string {
	var rules []string
	//
	if u.Gt != nil || u.Gte != nil || u.Lt != nil || u.Lte != nil {
		lo, loEq := fns.minRule(u.Gt, u.Gte)
		hi, hiEq := fns.maxRule(u.Lt, u.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if u.Const != nil {
		rules = append(rules, fns.ruleIs(u.GetConst()))
	}
	if len(u.In) > 0 {
		rules = append(rules, fns.ruleIn(u.In, false))
	}
	if len(u.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(u.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) uint64Rules(u *validate.UInt64Rules) string {
	var rules []string
	//
	if u.Gt != nil || u.Gte != nil || u.Lt != nil || u.Lte != nil {
		lo, loEq := fns.minRule(u.Gt, u.Gte)
		hi, hiEq := fns.maxRule(u.Lt, u.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if u.Const != nil {
		rules = append(rules, fns.ruleIs(u.GetConst()))
	}
	if len(u.In) > 0 {
		rules = append(rules, fns.ruleIn(u.In, false))
	}
	if len(u.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(u.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sInt32Rules(s *validate.SInt32Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		lo, loEq := fns.minRule(s.Gt, s.Gte)
		hi, hiEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sInt64Rules(s *validate.SInt64Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		lo, loEq := fns.minRule(s.Gt, s.Gte)
		hi, hiEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) fixed32Rules(f *validate.Fixed32Rules) string {
	var rules []string
	//
	if f.Gt != nil || f.Gte != nil || f.Lt != nil || f.Lte != nil {
		lo, loEq := fns.minRule(f.Gt, f.Gte)
		hi, hiEq := fns.maxRule(f.Lt, f.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if f.Const != nil {
		rules = append(rules, fns.ruleIs(f.GetConst()))
	}
	if len(f.In) > 0 {
		rules = append(rules, fns.ruleIn(f.In, false))
	}
	if len(f.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(f.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) fixed64Rules(f *validate.Fixed64Rules) string {
	var rules []string
	//
	if f.Gt != nil || f.Gte != nil || f.Lt != nil || f.Lte != nil {
		lo, loEq := fns.minRule(f.Gt, f.Gte)
		hi, hiEq := fns.maxRule(f.Lt, f.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if f.Const != nil {
		rules = append(rules, fns.ruleIs(f.GetConst()))
	}
	if len(f.In) > 0 {
		rules = append(rules, fns.ruleIn(f.In, false))
	}
	if len(f.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(f.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sFixed32Rules(s *validate.SFixed32Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		lo, loEq := fns.minRule(s.Gt, s.Gte)
		hi, hiEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sFixed64Rules(s *validate.SFixed64Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		lo, loEq := fns.minRule(s.Gt, s.Gte)
		hi, hiEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) boolRules(b *validate.BoolRules) string {
	if b.Const != nil {
		return fns.ruleIs(b.GetConst())
	}
	return "-"
}

func (fns Funcs) uintRangeRule(label string, lower, upper *uint64) string {
	lo, hi := "-INF", "+INF"
	loEq, hiEq := false, false
	if lower != nil {
		lo = fmt.Sprintf("%d", *lower)
		loEq = true
	}
	if upper != nil {
		hi = fmt.Sprintf("%d", *upper)
		hiEq = true
	}
	return label + " " + fns.ruleInterval(lo, hi, loEq, hiEq)
}

func (fns Funcs) stringRules(s *validate.StringRules) string {
	var rules []string
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if s.Len != nil {
		rules = append(rules, fmt.Sprintf("LEN: `%d`", s.GetLen()))
	}
	if s.MinLen != nil || s.MaxLen != nil {
		rules = append(rules, fns.uintRangeRule("CHARS", s.MinLen, s.MaxLen))
	}
	if s.LenBytes != nil {
		rules = append(rules, fmt.Sprintf("LEN_BYTES: `%d`", s.GetLenBytes()))
	}
	if s.MinBytes != nil || s.MaxBytes != nil {
		rules = append(rules, fns.uintRangeRule("BYTES", s.MinBytes, s.MaxBytes))
	}
	if s.Pattern != nil {
		rules = append(rules, fmt.Sprintf("PATTERN: `%s`", s.GetPattern()))
	}
	if s.Prefix != nil {
		rules = append(rules, fmt.Sprintf("PREFIX: `%s`", s.GetPrefix()))
	}
	if s.Suffix != nil {
		rules = append(rules, fmt.Sprintf("SUFFIX: `%s`", s.GetSuffix()))
	}
	if s.Contains != nil {
		rules = append(rules, fmt.Sprintf("CONTAINS: `%s`", s.GetContains()))
	}
	if s.NotContains != nil {
		rules = append(rules, fmt.Sprintf("NOT CONTAINS: `%s`", s.GetNotContains()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	switch {
	case s.GetEmail():
		rules = append(rules, "WELL_KNOWN: `email`")
	case s.GetHostname():
		rules = append(rules, "WELL_KNOWN: `hostname`")
	case s.GetIp():
		rules = append(rules, "WELL_KNOWN: `ip`")
	case s.GetIpv4():
		rules = append(rules, "WELL_KNOWN: `ipv4`")
	case s.GetIpv6():
		rules = append(rules, "WELL_KNOWN: `ipv6`")
	case s.GetUri():
		rules = append(rules, "WELL_KNOWN: `uri`")
	case s.GetUriRef():
		rules = append(rules, "WELL_KNOWN: `uri_ref`")
	case s.GetAddress():
		rules = append(rules, "WELL_KNOWN: `address`")
	case s.GetUuid():
		rules = append(rules, "WELL_KNOWN: `uuid`")
	case s.GetWellKnownRegex() != validate.KnownRegex_UNKNOWN:
		rules = append(rules, fmt.Sprintf("WELL_KNOWN_REGEX: `%s`", s.GetWellKnownRegex().String()))
	}
	if s.GetIgnoreEmpty() {
		rules = append(rules, "IGNORE_EMPTY")
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) bytesRules(b *validate.BytesRules) string {
	var rules []string
	if len(b.Const) > 0 {
		rules = append(rules, fns.ruleIs(fmt.Sprintf("%x", b.GetConst())))
	}
	if b.Len != nil {
		rules = append(rules, fmt.Sprintf("LEN: `%d`", b.GetLen()))
	}
	if b.MinLen != nil || b.MaxLen != nil {
		rules = append(rules, fns.uintRangeRule("BYTES", b.MinLen, b.MaxLen))
	}
	if b.Pattern != nil {
		rules = append(rules, fmt.Sprintf("PATTERN: `%s`", b.GetPattern()))
	}
	if len(b.Prefix) > 0 {
		rules = append(rules, fmt.Sprintf("PREFIX: `%x`", b.GetPrefix()))
	}
	if len(b.Suffix) > 0 {
		rules = append(rules, fmt.Sprintf("SUFFIX: `%x`", b.GetSuffix()))
	}
	if len(b.Contains) > 0 {
		rules = append(rules, fmt.Sprintf("CONTAINS: `%x`", b.GetContains()))
	}
	if len(b.In) > 0 {
		vals := make([]string, 0, len(b.In))
		for _, v := range b.In {
			vals = append(vals, fmt.Sprintf("%x", v))
		}
		rules = append(rules, fmt.Sprintf("IN: `[%s]`", strings.Join(vals, ", ")))
	}
	if len(b.NotIn) > 0 {
		vals := make([]string, 0, len(b.NotIn))
		for _, v := range b.NotIn {
			vals = append(vals, fmt.Sprintf("%x", v))
		}
		rules = append(rules, fmt.Sprintf("NOT IN: `[%s]`", strings.Join(vals, ", ")))
	}
	switch {
	case b.GetIp():
		rules = append(rules, "WELL_KNOWN: `ip`")
	case b.GetIpv4():
		rules = append(rules, "WELL_KNOWN: `ipv4`")
	case b.GetIpv6():
		rules = append(rules, "WELL_KNOWN: `ipv6`")
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) enumRules(enum *validate.EnumRules) string {
	var rules []string
	if enum.Const != nil {
		rules = append(rules, fns.ruleIs(enum.GetConst()))
	}
	if enum.GetDefinedOnly() {
		rules = append(rules, "DEFINED_ONLY")
	}
	if len(enum.In) > 0 {
		rules = append(rules, fns.ruleIn(enum.In, false))
	}
	if len(enum.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(enum.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) repeatedRules(r *validate.RepeatedRules) string {
	var rules []string
	if r.MinItems != nil || r.MaxItems != nil {
		rules = append(rules, fns.uintRangeRule("ITEMS", r.MinItems, r.MaxItems))
	}
	if r.GetUnique() {
		rules = append(rules, "UNIQUE")
	}
	if r.Items != nil {
		rules = append(rules, "ITEMS: `(nested rules)`")
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) mapRules(m *validate.MapRules) string {
	var rules []string
	if m.MinPairs != nil || m.MaxPairs != nil {
		rules = append(rules, fns.uintRangeRule("PAIRS", m.MinPairs, m.MaxPairs))
	}
	if m.GetNoSparse() {
		rules = append(rules, "NO_SPARSE")
	}
	if m.Keys != nil {
		rules = append(rules, "KEYS: `(nested rules)`")
	}
	if m.Values != nil {
		rules = append(rules, "VALUES: `(nested rules)`")
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) anyRules(any *validate.AnyRules) string {
	var rules []string
	if any.GetRequired() {
		rules = append(rules, "REQUIRED")
	}
	if len(any.In) > 0 {
		rules = append(rules, fns.ruleIn(any.In, false))
	}
	if len(any.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(any.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) durationRules(dur *validate.DurationRules) string {
	var rules []string
	if dur.GetRequired() {
		rules = append(rules, "REQUIRED")
	}
	if dur.Const != nil {
		rules = append(rules, fns.ruleIs(dur.GetConst().AsDuration()))
	}
	if dur.Gt != nil || dur.Gte != nil || dur.Lt != nil || dur.Lte != nil {
		lo, loEq := "-INF", false
		hi, hiEq := "+INF", false
		if dur.Gte != nil {
			lo, loEq = dur.GetGte().AsDuration().String(), true
		} else if dur.Gt != nil {
			lo, loEq = dur.GetGt().AsDuration().String(), false
		}
		if dur.Lte != nil {
			hi, hiEq = dur.GetLte().AsDuration().String(), true
		} else if dur.Lt != nil {
			hi, hiEq = dur.GetLt().AsDuration().String(), false
		}
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	if len(dur.In) > 0 {
		vals := make([]string, 0, len(dur.In))
		for _, v := range dur.In {
			vals = append(vals, v.AsDuration().String())
		}
		rules = append(rules, fmt.Sprintf("IN: `[%s]`", strings.Join(vals, ", ")))
	}
	if len(dur.NotIn) > 0 {
		vals := make([]string, 0, len(dur.NotIn))
		for _, v := range dur.NotIn {
			vals = append(vals, v.AsDuration().String())
		}
		rules = append(rules, fmt.Sprintf("NOT IN: `[%s]`", strings.Join(vals, ", ")))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) timestampRules(ts *validate.TimestampRules) string {
	var rules []string
	if ts.GetRequired() {
		rules = append(rules, "REQUIRED")
	}
	if ts.Const != nil {
		rules = append(rules, fns.ruleIs(ts.GetConst().AsTime().UTC().Format("2006-01-02T15:04:05.999999999Z07:00")))
	}
	if ts.Gt != nil || ts.Gte != nil || ts.Lt != nil || ts.Lte != nil {
		lo, loEq := "-INF", false
		hi, hiEq := "+INF", false
		if ts.Gte != nil {
			lo, loEq = ts.GetGte().AsTime().UTC().Format(time.RFC3339Nano), true
		} else if ts.Gt != nil {
			lo, loEq = ts.GetGt().AsTime().UTC().Format(time.RFC3339Nano), false
		}
		if ts.Lte != nil {
			hi, hiEq = ts.GetLte().AsTime().UTC().Format(time.RFC3339Nano), true
		} else if ts.Lt != nil {
			hi, hiEq = ts.GetLt().AsTime().UTC().Format(time.RFC3339Nano), false
		}
		rules = append(rules, fns.ruleInterval(lo, hi, loEq, hiEq))
	}
	if ts.GetLtNow() {
		rules = append(rules, "LT_NOW")
	}
	if ts.GetGtNow() {
		rules = append(rules, "GT_NOW")
	}
	if ts.Within != nil {
		rules = append(rules, fmt.Sprintf("WITHIN: `%s`", ts.GetWithin().AsDuration()))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) fieldRules(field pgs.Field) string {
	var rules validate.FieldRules
	if ok, _ := field.Extension(validate.E_Rules, &rules); !ok {
		return "-"
	}

	switch r := rules.GetType().(type) {
	case *validate.FieldRules_Float:
		return fns.floatRules(r.Float)
	case *validate.FieldRules_Double:
		return fns.doubleRules(r.Double)
	case *validate.FieldRules_Int32:
		return fns.int32Rules(r.Int32)
	case *validate.FieldRules_Int64:
		return fns.int64Rules(r.Int64)
	case *validate.FieldRules_Uint32:
		return fns.uint32Rules(r.Uint32)
	case *validate.FieldRules_Uint64:
		return fns.uint64Rules(r.Uint64)
	case *validate.FieldRules_Sint32:
		return fns.sInt32Rules(r.Sint32)
	case *validate.FieldRules_Sint64:
		return fns.sInt64Rules(r.Sint64)
	case *validate.FieldRules_Fixed32:
		return fns.fixed32Rules(r.Fixed32)
	case *validate.FieldRules_Fixed64:
		return fns.fixed64Rules(r.Fixed64)
	case *validate.FieldRules_Sfixed32:
		return fns.sFixed32Rules(r.Sfixed32)
	case *validate.FieldRules_Sfixed64:
		return fns.sFixed64Rules(r.Sfixed64)
	case *validate.FieldRules_Bool:
		return fns.boolRules(r.Bool)
	case *validate.FieldRules_String_:
		return fns.stringRules(r.String_)
	case *validate.FieldRules_Bytes:
		return fns.bytesRules(r.Bytes)
	case *validate.FieldRules_Enum:
		return fns.enumRules(r.Enum)
	case *validate.FieldRules_Repeated:
		return fns.repeatedRules(r.Repeated)
	case *validate.FieldRules_Map:
		return fns.mapRules(r.Map)
	case *validate.FieldRules_Any:
		return fns.anyRules(r.Any)
	case *validate.FieldRules_Duration:
		return fns.durationRules(r.Duration)
	case *validate.FieldRules_Timestamp:
		return fns.timestampRules(r.Timestamp)
	//case nil: // TODO
	default:
		return "-"
	}
}
