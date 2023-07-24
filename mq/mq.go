package mq

import "context"

type MQ interface {
	Send(ctx context.Context)
	SendDelay(ctx context.Context)
	Consume(ctx context.Context)
}

func Init() error {
	return nil
}
