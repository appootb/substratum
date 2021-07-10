package datetime

import (
	"encoding/json"
	"time"
)

type Date time.Time

func (dt Date) Time() time.Time {
	return time.Time(dt)
}

func (dt Date) String() string {
	return time.Time(dt).Format("2006-01-02")
}

func (dt *Date) parse(v string) error {
	ts, err := time.ParseInLocation("2006-01-02", v, time.Local)
	if err != nil {
		return err
	}
	*dt = Date(ts)
	return nil
}

func (dt Date) MarshalURL() (string, error) {
	return dt.String(), nil
}

func (dt *Date) UnmarshalURL(v string) error {
	return dt.parse(v)
}

func (dt Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}

func (dt *Date) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	return dt.parse(s)
}
