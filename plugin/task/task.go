package task

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"time"

	"github.com/appootb/substratum/plugin/context"
	"github.com/appootb/substratum/task"
	"github.com/appootb/substratum/util/scheduler"
	"github.com/appootb/substratum/util/timer"
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
	ctx := opts.Context
	if opts.Singleton {
		// Blocked before acquired the locker.
		ctx = task.LockerImplementor().Lock(opts.Context, opts.Name)
	}

	ctx = context.WithImplementContext(opts.Context, opts.Component)

	for {
		now := time.Now()
		next := schedule.Next(now)

		select {
		case <-ctx.Done():
			select {
			case <-opts.Done():
				if opts.Singleton {
					task.LockerImplementor().Unlock(opts.Name)
				}
				return
			default:
				// Keep alive failed, reset.
				goto Reset
			}

		case <-timer.After(next.Sub(now)):
			err := fn(ctx, opts.Argument)
			if err != nil {
				// TODO succeed
			} else {
				// TODO failed
			}
		}
	}
}
