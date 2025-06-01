package grpcclient

import (
	"context"
	"log"
	"ride-sharing/config"
	"ride-sharing/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	client proto.NotificationServiceClient
	conn   *grpc.ClientConn
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

func (n *NotificationClient) SendRegisterEmail(ctx context.Context, to string, otp string) (*proto.StandardResponse, error) {
	req := &proto.RegisterEmailRequest{
		To:  to,
		Otp: otp,
	}

	resp, err := n.client.SendRegisterEmail(ctx, req)
	if err != nil {
		log.Printf("gRPC SendRegisterEmail failed: %v", err)
		return nil, err
	}

	return resp, nil
}
