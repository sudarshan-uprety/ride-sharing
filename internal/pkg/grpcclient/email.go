package grpcclient

import (
	"context"
	"errors"
	"math"
	"ride-sharing/config"
	"ride-sharing/internal/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	client proto.NotificationServiceClient
	// kafka  *kafka.Producer
	conn *grpc.ClientConn
}

func NewNotificationClient(cfg *config.Config) (*NotificationClient, error) {
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

	for attempt := 1; attempt <= 3; attempt++ {
		_, err := n.client.SendRegisterEmail(ctx, req)
		if err == nil {
			return true, nil
		}
		// lastErr := err
		time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
	}

	// err := n.kafka.Produce("notification-topic", map[string]string{
	// 	"type": "register_email",
	// 	"to":   to,
	// 	"otp":  otp,
	// })
	// if err != nil {
	// 	log.Printf("Failed to enqueue message to Kafka: %v", err)
	// 	return false, errors.New("all gRPC attempts failed and Kafka fallback also failed")
	// }

	return false, errors.New("email queued in Kafka due to notification service downtime")

}

func (n *NotificationClient) SendForgetPasswordEmail(ctx context.Context, to string, otp string) (bool, error) {
	req := &proto.ForgetPasswordEmailRequest{
		To:  to,
		Otp: otp,
	}

	for attempt := 1; attempt <= 3; attempt++ {
		_, err := n.client.SendForgetPasswordEmail(ctx, req)
		if err == nil {
			return true, nil
		}
		// lastErr := err
		time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
	}

	// err := n.kafka.Produce("notification-topic", map[string]string{
	// 	"type": "register_email",
	// 	"to":   to,
	// 	"otp":  otp,
	// })
	// if err != nil {
	// 	log.Printf("Failed to enqueue message to Kafka: %v", err)
	// 	return false, errors.New("all gRPC attempts failed and Kafka fallback also failed")
	// }

	return false, errors.New("email queued in Kafka due to notification service downtime")
}
