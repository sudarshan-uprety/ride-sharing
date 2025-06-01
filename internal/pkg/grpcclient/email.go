package email

import (
	"context"
	"net"
	"ride-sharing/internal/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	client proto.NotificationServiceClient
	conn   *grpc.ClientConn
}

func NewNotificationClient(address string) (*NotificationClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new gRPC client connection
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "tcp", addr)
		}),
	)
	if err != nil {
		return nil, err
	}

	connectCtx, connectCancel := context.WithTimeout(ctx, 5*time.Second)
	defer connectCancel()

	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			break
		}
		if !conn.WaitForStateChange(connectCtx, state) {
			conn.Close()
			return nil, context.DeadlineExceeded
		}
	}

	client := proto.NewNotificationServiceClient(conn)
	return &NotificationClient{client: client, conn: conn}, nil
}

func (c *NotificationClient) Close() error {
	return c.conn.Close()
}

func (n *NotificationClient) SendRegisterEmail(ctx context.Context, to string, otp string) (*proto.StandardResponse, error) {
	req := &proto.RegisterEmailRequest{
		To:  to,
		Otp: otp,
	}
	return n.client.SendRegisterEmail(ctx, req)
}
