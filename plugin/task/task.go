package task

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"time"

	sctx "github.com/appootb/substratum/context"
	ictx "github.com/appootb/substratum/internal/context"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/task"
	"github.com/appootb/substratum/util/scheduler"
	"github.com/appootb/substratum/util/timer"
)

const (
	DebugLog = "_TASK_.debug"
	ErrorLog = "_TASK_.error"

	LogName     = logger.LogTag + "name"
	LogExecutor = logger.LogTag + "executor"
	LogError    = logger.LogTag + "error"
)

type Task struct{}

func (c *Task) Schedule(spec string, exec task.Executor, opts ...task.Option) error {
	options := task.EmptyOptions()
	for _, o := range opts {
		o(options)
	}
	if options.Name == "" {
		options.Name = fmt.Sprintf("%x", sha1.Sum([]byte(spec+c.reflectName(exec))))
	}
	schedule, err := scheduler.ParseStandard(spec)
	if err != nil {
		return err
	}
	go c.exec(schedule, exec, options)
	return nil
}

func (c *Task) reflectName(exec task.Executor) string {
	t := reflect.TypeOf(exec)
	// reflect.Ptr's PkgPath and Name is empty
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return fmt.Sprintf("%s/%s", t.PkgPath(), t.Name())
}

func (c *Task) exec(schedule scheduler.Schedule, exec task.Executor, opts *task.Options) {
Reset:
	ctx := ictx.Context
	if opts.Singleton {
		// Blocked before acquired the locker.
		ctx = task.LockerImplementor().Lock(ctx, opts.Name)
	}

	for {
		now := time.Now()
		next := schedule.Next(now)

		select {
		case <-ctx.Done():
			select {
			case <-ictx.Context.Done():
				if opts.Singleton {
					task.LockerImplementor().Unlock(opts.Name)
				}
				return
			default:
				// Keep alive failed, reset.
				goto Reset
			}

		case <-timer.After(next.Sub(now)):
			err := exec.Execute(sctx.ServerContext(opts.Component), opts.Argument)
			if err != nil {
				logger.Error(ErrorLog, logger.Content{
					LogError:    err.Error(),
					LogName:     opts.Name,
					LogExecutor: c.reflectName(exec),
				})
			} else {
				logger.Debug(DebugLog, logger.Content{
					LogName:     opts.Name,
					LogExecutor: c.reflectName(exec),
				})
			}
		}
	}
}
