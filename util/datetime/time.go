package datetime

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Time struct {
	time.Time
}

func (t Time) Proto() *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t.Time)
	if err != nil {
		log.Println(err)
	}
	return ts
}

func WithTime(t time.Time) *Time {
	return &Time{t}
}

func FromProtoTime(ts *timestamp.Timestamp) *Time {
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		log.Println(err)
	}
	return &Time{t.Local()}
}

func Now() *Time {
	return &Time{time.Now()}
}
