package kafka

import (
	"github.com/IBM/sarama"
	"log"
	"rest_api/internal/api/config"
	"time"
)

type AsyncPublisher struct {
	producer sarama.AsyncProducer
}

func (p *AsyncPublisher) Configure(topic string) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.DefaultVersion
	cfg.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	cfg.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer([]string{config.BrokerLink}, cfg)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to produce message:", err)
		}
	}()

	p.producer = producer
}

func (p *AsyncPublisher) Publish(msg string) error {
	// no key, the message with end up in random partition
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: config.Topic,
		Value: sarama.StringEncoder(msg),
	}
	return nil
}
