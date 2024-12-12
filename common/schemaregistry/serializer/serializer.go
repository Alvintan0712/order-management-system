package serializer

type Serializer interface {
	Serialize(topic string, value interface{}) ([]byte, error)
}
