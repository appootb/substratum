package shared

import (
	"strings"
	"testing"

	"github.com/appootb/substratum/v2/proto/go/validate"
)

func TestStringRulesDocs(t *testing.T) {
	var fns Funcs
	minLen, maxLen := uint64(1), uint64(10)
	pattern := "^[a-z]+$"
	s := &validate.StringRules{
		MinLen:    &minLen,
		MaxLen:    &maxLen,
		Pattern:   &pattern,
		WellKnown: &validate.StringRules_Email{Email: true},
	}
	got := fns.stringRules(s)
	for _, want := range []string{"CHARS", "PATTERN", "WELL_KNOWN: `email`"} {
		if !strings.Contains(got, want) {
			t.Fatalf("stringRules=%q missing %q", got, want)
		}
	}
}

func TestEnumRulesDocs(t *testing.T) {
	var fns Funcs
	defined := true
	got := fns.enumRules(&validate.EnumRules{DefinedOnly: &defined, In: []int32{1, 2}})
	if !strings.Contains(got, "DEFINED_ONLY") || !strings.Contains(got, "IN:") {
		t.Fatalf("enumRules=%q", got)
	}
}

func TestUintRangeRule(t *testing.T) {
	var fns Funcs
	minItems := uint64(2)
	got := fns.uintRangeRule("ITEMS", &minItems, nil)
	if got != "ITEMS RANGE: `[2, +INF)`" {
		t.Fatalf("got %q", got)
	}
}
