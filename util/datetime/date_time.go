package datetime

import (
	"encoding/json"
	"time"
)

type DateTime time.Time

func (dt DateTime) Time() time.Time {
	return time.Time(dt)
}

func (dt DateTime) String() string {
	return time.Time(dt).Format("2006-01-02 15:04:05")
}

func (dt *DateTime) parse(v string) error {
	ts, err := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	if err != nil {
		return err
	}
	*dt = DateTime(ts)
	return nil
}

func (dt DateTime) MarshalURL() (string, error) {
	return dt.String(), nil
}

func (dt *DateTime) UnmarshalURL(v string) error {
	return dt.parse(v)
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(dt.String()), nil
}

func (dt *DateTime) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	return dt.parse(s)
}
