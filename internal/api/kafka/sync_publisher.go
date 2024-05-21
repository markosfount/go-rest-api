package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

type SyncPublisher struct {
	producer sarama.SyncProducer
	topic    string
}

func (p *SyncPublisher) Configure(topic string) *SyncPublisher {
	p.topic = topic

	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{brokerList}, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	p.producer = producer

	return p
}

func (p *SyncPublisher) Publish(msg string) error {
	pm := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(msg),
	}
	_, _, err := p.producer.SendMessage(pm)
	if err != nil {
		return err
	}
	return nil
}
