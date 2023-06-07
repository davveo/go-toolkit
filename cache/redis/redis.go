package redis

import (
	"context"
	"github.com/davveo/go-toolkit/cache"
	"time"
)

type AdapterRedis struct {
}

func NewAdapterRedis() cache.Adapter {
	return &AdapterRedis{}
}

func (c *AdapterRedis) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	return nil
}

func (c *AdapterRedis) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error) {
	return false, nil
}

func (c *AdapterRedis) SetIfNotExistFunc(ctx context.Context, key interface{}, f cache.Func, duration time.Duration) (ok bool, err error) {
	return false, nil
}
