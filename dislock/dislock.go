package dislock

import (
	"context"
)

/*
https://github.com/ggicci/distlock/blob/main/redis.go
*/

type Lock interface {
	Lock(ctx context.Context, key string)
	UnLock(ctx context.Context, key string)
}
