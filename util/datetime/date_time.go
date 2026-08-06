package datetime

import (
	"encoding/json"
	"time"
)

// DefaultDateTimeLayout is the default layout for DateTime (YYYY-MM-DD HH:MM:SS).
const (
	DefaultDateTimeLayout = "2006-01-02 15:04:05"
)

// DateTime is a wall-clock timestamp with a configurable format for JSON/URL codecs.
type DateTime struct {
	payload time.Time
	layout  string
}

// NewDateTime wraps ts as a DateTime. An optional format overrides DefaultDateTimeLayout.
func NewDateTime(ts time.Time, format ...string) *DateTime {
	layout := DefaultDateTimeLayout
	if len(format) > 0 && format[0] != "" {
		layout = format[0]
	}
	return &DateTime{
		payload: ts,
		layout:  layout,
	}
}

// Time returns the underlying time.Time.
func (dt *DateTime) Time() time.Time {
	return dt.payload
}

// String formats the timestamp with the configured layout (or DefaultDateTimeLayout).
func (dt *DateTime) String() string {
	layout := dt.layout
	if layout == "" {
		layout = DefaultDateTimeLayout
	}
	return dt.payload.Format(layout)
}

func (dt *DateTime) parse(v string) (err error) {
	if v == "" {
		return
	}
	if dt.layout == "" {
		dt.layout = DefaultDateTimeLayout
	}
	dt.payload, err = time.ParseInLocation(dt.layout, v, time.Local)
	return
}

// MarshalURL encodes the timestamp using the same layout as String.
func (dt *DateTime) MarshalURL() (string, error) {
	return dt.String(), nil
}

// UnmarshalURL parses v using the configured (or default) layout.
func (dt *DateTime) UnmarshalURL(v string) error {
	return dt.parse(v)
}

// MarshalJSON encodes the timestamp as a JSON string using the layout.
func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}

// UnmarshalJSON decodes a JSON string into the timestamp; empty string is a no-op.
func (dt *DateTime) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	return dt.parse(s)
}
