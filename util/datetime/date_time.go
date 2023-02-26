package datetime

import (
	"encoding/json"
	"time"
)

const (
	DefaultDateTimeLayout = "2006-01-02 15:04:05"
)

type DateTime struct {
	payload time.Time
	layout  string
}

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

func (dt *DateTime) Time() time.Time {
	return dt.payload
}

func (dt *DateTime) String() string {
	return dt.payload.Format(dt.layout)
}

func (dt *DateTime) parse(v string) (err error) {
	dt.payload, err = time.ParseInLocation(dt.layout, v, time.Local)
	return
}

func (dt *DateTime) MarshalURL() (string, error) {
	return dt.String(), nil
}

func (dt *DateTime) UnmarshalURL(v string) error {
	return dt.parse(v)
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}

func (dt *DateTime) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	return dt.parse(s)
}
