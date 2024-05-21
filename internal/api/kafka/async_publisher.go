package kafka

import (
	"github.com/IBM/sarama"
	"log"
	"rest_api/internal/api/model"
	"time"
)

type AsyncPublisher struct {
	producer sarama.AsyncProducer
}

func (p *AsyncPublisher) Create() *AsyncPublisher {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer([]string{brokerList}, config)
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

	return p
}

func (p *AsyncPublisher) Publish(movie *model.Movie) error {
	return nil
}
