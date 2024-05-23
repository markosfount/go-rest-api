package kafka

type Publisher interface {
	Publish(msg string) error
	Configure(topic string)
}
