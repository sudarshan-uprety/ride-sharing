package kafka

import (
	"ride-sharing/config"

	"github.com/segmentio/kafka-go"
)

func NewProducerFromAppConfig(cfg *config.Config) *Producer {
	var balancer kafka.Balancer

	switch cfg.Kafka.Balancer {
	case "round-robin":
		balancer = &kafka.RoundRobin{}
	case "hash":
		balancer = &kafka.Hash{}
	default:
		balancer = &kafka.LeastBytes{}
	}

	return NewProducer(ProducerConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		Balancer: balancer,
	})
}
