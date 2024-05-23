package kafka

import (
	"github.com/IBM/sarama"
	"log"
	"rest_api/internal/api/config"
)

type SyncPublisher struct {
	producer sarama.SyncProducer
	topic    string
}

func (p *SyncPublisher) Configure(topic string) {
	p.topic = topic

	cfg := sarama.NewConfig()
	cfg.Version = sarama.DefaultVersion
	cfg.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	cfg.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{config.BrokerLink}, cfg)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	p.producer = producer
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
