package task

import (
	"context"
	"crypto/sha1"
	"fmt"
	"reflect"
	"time"

	"github.com/appootb/substratum/task"
	"github.com/appootb/substratum/util/scheduler"
)

type Task struct{}

func (c *Task) Schedule(spec string, fn task.JobFunc, opts ...task.Option) error {
	options := task.EmptyOptions()
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

func (c *Task) exec(schedule scheduler.Schedule, fn task.JobFunc, opts *task.Options) {
Reset:
	ctx := context.TODO()

	if opts.Singleton {
		err := task.BackendImplementor().Lock(opts.Name)
		if err != nil {
			time.Sleep(time.Second)
			goto Reset
		}
		// Keep alive
		ctx = task.BackendImplementor().KeepAlive(opts.Name)
	}

	for {
		now := time.Now()
		next := schedule.Next(now)

		select {
		case <-opts.Done():
			if opts.Singleton {
				task.BackendImplementor().Unlock(opts.Name)
			}
			return

		case <-ctx.Done():
			// Keep alive failed, reset.
			goto Reset

		case <-time.After(next.Sub(now)):
			err := fn(opts.Argument)
			if err != nil {
				// TODO succeed
			} else {
				// TODO failed
			}
		}
	}
}
