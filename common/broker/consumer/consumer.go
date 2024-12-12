package consumer

type ConsumedMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
}

type Consumer interface {
	Consume() (ConsumedMessage, error)
	Close() error
}
