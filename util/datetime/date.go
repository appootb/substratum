package datetime

import (
	"encoding/json"
	"time"
)

// DefaultDateLayout is the default layout for Date formatting and parsing (YYYY-MM-DD).
const (
	DefaultDateLayout = "2006-01-02"
)

// Date is a date-only value with a configurable format layout for JSON/URL codecs.
// The payload is always normalized to 00:00:00 in the source location.
type Date struct {
	payload time.Time
	layout  string
}

// NewDate wraps ts as a Date at local midnight of that calendar day.
// An optional format overrides DefaultDateLayout.
func NewDate(ts time.Time, format ...string) *Date {
	layout := DefaultDateLayout
	if len(format) > 0 && format[0] != "" {
		layout = format[0]
	}
	y, m, d := ts.Date()
	return &Date{
		payload: time.Date(y, m, d, 0, 0, 0, 0, ts.Location()),
		layout:  layout,
	}
}

// Time returns the underlying time.Time (always 00:00:00 when built via NewDate/parse).
func (d *Date) Time() time.Time {
	return d.payload
}

// String formats the date with the configured layout (or DefaultDateLayout).
func (d *Date) String() string {
	layout := d.layout
	if layout == "" {
		layout = DefaultDateLayout
	}
	return d.payload.Format(layout)
}

func (d *Date) parse(v string) (err error) {
	if v == "" {
		return
	}
	if d.layout == "" {
		d.layout = DefaultDateLayout
	}
	d.payload, err = time.ParseInLocation(d.layout, v, time.Local)
	return err
}

// MarshalURL encodes the date using the same layout as String.
func (d *Date) MarshalURL() (string, error) {
	return d.String(), nil
}

// UnmarshalURL parses v using the configured (or default) layout.
func (d *Date) UnmarshalURL(v string) error {
	return d.parse(v)
}

// MarshalJSON encodes the date as a JSON string using the layout.
func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON decodes a JSON string into the date; empty string is a no-op.
func (d *Date) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	return d.parse(s)
}
