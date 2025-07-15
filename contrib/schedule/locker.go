package schedule

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type AutoLock struct {
	lock    *redislock.Lock
	ttl     time.Duration
	cancel  context.CancelFunc
	stopped chan struct{}
}

// TryLock attempts to acquire a distributed lock and starts automatic renewal.
func TryLock(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) (*AutoLock, error) {
	locker := redislock.New(rdb)

	// Try to obtain the lock
	lock, err := locker.Obtain(ctx, key, ttl, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return nil, errors.New("lock already held")
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	al := &AutoLock{
		lock:    lock,
		ttl:     ttl,
		cancel:  cancel,
		stopped: make(chan struct{}),
	}

	go al.autoRefresh(ctx)
	return al, nil
}

// autoRefresh periodically refreshes the lock until cancellation or failure.
func (al *AutoLock) autoRefresh(ctx context.Context) {
	ticker := time.NewTicker(al.ttl / 2)
	defer func() {
		ticker.Stop()
		close(al.stopped)
	}()

	for {
		select {
		case <-ticker.C:
			err := al.lock.Refresh(ctx, al.ttl, nil)
			if err != nil {
				al.cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// Release stops the auto-renewal process and releases the lock.
func (al *AutoLock) Release(ctx context.Context) error {
	al.cancel()
	<-al.stopped
	return al.lock.Release(ctx)
}
