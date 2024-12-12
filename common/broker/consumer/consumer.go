package consumer

import (
	"time"
)

type EventMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Timestamp time.Time
}

type EventHandler func(EventMessage) error

type Consumer interface {
	Start(map[string]EventHandler) error
	Close() error
}
