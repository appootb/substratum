package cron

import (
	"context"
	"crypto/sha1"
	"fmt"
	"reflect"
	"time"

	"github.com/appootb/substratum/cron"
	"github.com/appootb/substratum/util/scheduler"
)

type Cron struct{}

func (c *Cron) Schedule(spec string, fn cron.JobFunc, opts ...cron.Option) error {
	options := cron.EmptyOptions()
	for _, o := range opts {
		o(options)
	}
	if options.Name == "" {
		t := reflect.TypeOf(fn)
		name := spec + t.PkgPath() + t.Name()
		options.Name = fmt.Sprintf("%x", sha1.Sum([]byte(name)))
	}
	schedule, err := scheduler.ParseStandard(spec)
	if err != nil {
		return err
	}
	go c.exec(schedule, fn, options)
	return nil
}

func (c *Cron) exec(schedule scheduler.Schedule, fn cron.JobFunc, opts *cron.Options) {
Reset:
	ctx := context.TODO()

	if opts.Singleton {
		err := cron.BackendImplementor().Lock(opts.Name)
		if err != nil {
			time.Sleep(time.Second)
			goto Reset
		}
		// Keep alive
		ctx = cron.BackendImplementor().KeepAlive(opts.Name)
	}

	for {
		now := time.Now()
		next := schedule.Next(now)

		select {
		case <-opts.Done():
			if opts.Singleton {
				cron.BackendImplementor().Unlock(opts.Name)
			}
			return

		case <-ctx.Done():
			// Keep alive failed, reset.
			goto Reset

		case <-time.After(next.Sub(now)):
			fn(opts.Argument)
		}
	}
}
