package producer

type Producer interface {
	Produce(topic string, key, value []byte) error
	Close() error
}
