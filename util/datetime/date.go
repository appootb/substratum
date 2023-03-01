package datetime

import (
	"encoding/json"
	"time"
)

const (
	DefaultDateLayout = "2006-01-02"
)

type Date struct {
	payload time.Time
	layout  string
}

func NewDate(ts time.Time, format ...string) *Date {
	layout := DefaultDateLayout
	if len(format) > 0 && format[0] != "" {
		layout = format[0]
	}
	return &Date{
		payload: ts,
		layout:  layout,
	}
}

func (d *Date) Time() time.Time {
	return d.payload
}

func (d *Date) String() string {
	return d.payload.Format(d.layout)
}

func (d *Date) parse(v string) (err error) {
	if d.layout == "" {
		d.layout = DefaultDateLayout
	}
	d.payload, err = time.ParseInLocation(d.layout, v, time.Local)
	return nil
}

func (d *Date) MarshalURL() (string, error) {
	return d.payload.String(), nil
}

func (d *Date) UnmarshalURL(v string) error {
	return d.parse(v)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	return d.parse(s)
}
