package kafka

const brokerList = "localhost:9092"

type Publisher interface {
	Publish(msg string) error
}
