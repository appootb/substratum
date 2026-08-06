package storage

import (
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func TestIsEmpty(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"redis.Nil", redis.Nil, true},
		{"gorm not found", gorm.ErrRecordNotFound, true},
		{"wrapped redis", errors.Join(errors.New("x"), redis.Nil), true},
		{"wrapped gorm", errors.Join(errors.New("x"), gorm.ErrRecordNotFound), true},
		{"other", errors.New("boom"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsEmpty(tc.err); got != tc.want {
				t.Fatalf("IsEmpty(%v)=%v want %v", tc.err, got, tc.want)
			}
		})
	}
}
