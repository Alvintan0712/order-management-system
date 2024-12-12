package deserializer

type Deserializer interface {
	Deserialize(topic string, payload []byte, value interface{}) error
}
