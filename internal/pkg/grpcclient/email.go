package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"math"
	"ride-sharing/config"
	"ride-sharing/internal/pkg/constants"
	"ride-sharing/internal/pkg/kafka"
	"ride-sharing/internal/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	client proto.NotificationServiceClient
	kafka  *kafka.Producer
	conn   *grpc.ClientConn
}

func NewNotificationClient(cfg *config.Config, producer *kafka.Producer) (*NotificationClient, error) {
	conn, err := grpc.NewClient(
		cfg.Notification.Host+":"+cfg.Notification.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &NotificationClient{
		client: proto.NewNotificationServiceClient(conn),
		conn:   conn,
		kafka:  producer,
	}, nil
}

func (c *NotificationClient) Close() error {
	return c.conn.Close()
}

func (n *NotificationClient) SendRegisterEmail(ctx context.Context, to string, otp string) (bool, error) {
	req := &proto.RegisterEmailRequest{
		To:  to,
		Otp: otp,
	}

	var lastErr error
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, err := n.client.SendRegisterEmail(ctx, req)
		if err == nil {
			return true, nil
		}

		lastErr = err

		if attempt < maxAttempts {
			time.Sleep(exponentialBackoff(attempt))
		}
	}

	// Fallback to Kafka
	kafkaErr := n.kafka.Produce(ctx, string(constants.OTPUserRegister), map[string]string{
		"type": string(constants.OTPUserRegister),
		"to":   to,
		"otp":  otp,
	})

	if kafkaErr != nil {
		return false, errors.Join(
			fmt.Errorf("gRPC attempts failed: %w", lastErr),
			fmt.Errorf("kafka fallback failed: %w", kafkaErr),
		)
	}

	return true, nil
}

func (n *NotificationClient) SendForgetPasswordEmail(ctx context.Context, to string, otp string) (bool, error) {
	req := &proto.ForgetPasswordEmailRequest{
		To:  to,
		Otp: otp,
	}

	var lastErr error
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, err := n.client.SendForgetPasswordEmail(ctx, req)
		if err == nil {
			return true, nil
		}

		lastErr = err

		if attempt < maxAttempts {
			time.Sleep(exponentialBackoff(attempt))
		}
	}

	// Fallback to Kafka
	kafkaErr := n.kafka.Produce(ctx, string(constants.OTPUserRegister), map[string]string{
		"type": string(constants.OTPUserRegister),
		"to":   to,
		"otp":  otp,
	})

	if kafkaErr != nil {
		return false, errors.Join(
			fmt.Errorf("gRPC attempts failed: %w", lastErr),
			fmt.Errorf("kafka fallback failed: %w", kafkaErr),
		)
	}

	return true, nil
}

func exponentialBackoff(attempt int) time.Duration {
	base := math.Pow(2, float64(attempt))
	scale := 3.0 / (2 + 4)
	scaledDelay := base * scale
	return time.Duration(scaledDelay * float64(time.Second))
}
