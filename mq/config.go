package mq

import "time"

type Msg struct {
	Topic string
	Tag   string
	Body  []byte
	Delay time.Duration
}
