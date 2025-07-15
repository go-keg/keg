package schedule

import (
	"context"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	cron *cron.Cron
	log  *log.Helper
}

func NewSchedule(logger log.Logger) *Schedule {
	cronLog := cronLogger{log.NewHelper(log.With(logger, "module", "schedule/cron"))}
	return &Schedule{
		log: log.NewHelper(log.With(logger, "module", "schedule")),
		cron: cron.New(cron.WithChain(
			cron.Recover(cronLog),
			cron.DelayIfStillRunning(cronLog),
		)),
	}
}

type ArgOption func(*scheduleArgs)

type scheduleArgs struct {
	lock func(ctx context.Context) (*AutoLock, error)
}

func OnOneServer(rdb *redis.Client, key string, ttl time.Duration) ArgOption {
	return func(args *scheduleArgs) {
		args.lock = func(ctx context.Context) (*AutoLock, error) {
			return TryLock(ctx, rdb, key, ttl)
		}
	}
}

func (s Schedule) Add(name string, spec string, job func() error, opts ...ArgOption) (cron.EntryID, error) {
	s.log.Infof("add schedule: %s", name)
	var args scheduleArgs
	for _, opt := range opts {
		opt(&args)
	}

	return s.cron.AddFunc(spec, func() {
		s.log.Infof("run schedule: %s", name)
		if args.lock != nil {
			ctx := context.Background()
			lock, err := args.lock(ctx)
			if err != nil {
				return
			}
			defer func(lock *AutoLock, ctx context.Context) {
				_ = lock.Release(ctx)
			}(lock, ctx)
		}
		err := job()
		if err != nil {
			s.log.Errorw(
				"method", "schedule_err",
				"name", name,
				"err", err,
			)
		}
	})
}

func (s Schedule) AddCtx(ctx context.Context, name string, spec string, job func(ctx2 context.Context) error, opts ...ArgOption) (cron.EntryID, error) {
	s.log.Infof("add schedule: %s", name)
	var args scheduleArgs
	for _, opt := range opts {
		opt(&args)
	}
	return s.cron.AddFunc(spec, func() {
		s.log.Infof("run schedule: %s", name)
		if args.lock != nil {
			lock, err := args.lock(ctx)
			if err != nil {
				return
			}
			defer func(lock *AutoLock, ctx context.Context) {
				_ = lock.Release(ctx)
			}(lock, ctx)
		}
		err := job(ctx)
		if err != nil {
			s.log.Errorw(
				"method", "schedule_err",
				"name", name,
				"err", err,
			)
		}
	})
}

func (s Schedule) Start() error {
	if os.Getenv("SCHEDULE_ENABLE") != "false" {
		s.cron.Run()
	}
	return nil
}

func (s Schedule) Stop() error {
	s.cron.Stop()
	return nil
}
