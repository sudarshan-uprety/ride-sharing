package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

type ProducerConfig struct {
	Brokers  []string
	Topic    string
	Balancer kafka.Balancer
}

func NewProducer(cfg ProducerConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			Balancer:     cfg.Balancer,
			BatchTimeout: 50 * time.Millisecond,
			Async:        true,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) Produce(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: jsonValue,
			Time:  time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("kafka produce failed: %w", err)
	}

	return nil
}
