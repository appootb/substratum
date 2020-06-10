package datetime

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
)

type Duration struct {
	time.Duration
}

func (d Duration) Proto() *duration.Duration {
	return ptypes.DurationProto(d.Duration)
}

func WithDuration(d time.Duration) *Duration {
	return &Duration{d}
}

func FromProtoDuration(d *duration.Duration) *Duration {
	dur, err := ptypes.Duration(d)
	if err != nil {
		log.Println(err)
	}
	return &Duration{dur}
}
