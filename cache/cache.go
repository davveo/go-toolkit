package cache

import (
	"context"
	"time"
)

type Func func(ctx context.Context) (value interface{}, err error)

type Adapter interface {
	Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error
	SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error)
	SetIfNotExistFunc(ctx context.Context, key interface{}, f Func, duration time.Duration) (ok bool, err error)
}