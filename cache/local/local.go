package local

import (
	"context"
	"github.com/davveo/go-toolkit/cache"
	"time"
)

type AdapterLocal struct {
}

func NewAdapterRedis() cache.Adapter {
	return &AdapterLocal{}
}

func (a AdapterLocal) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (a AdapterLocal) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AdapterLocal) SetIfNotExistFunc(ctx context.Context, key interface{}, f cache.Func, duration time.Duration) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}
