package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPStore struct {
	cli *redis.Client
}

func NewOTPStore(client *Client) *OTPStore {
	return &OTPStore{cli: client.cli}
}

func (s *OTPStore) SetOTP(ctx context.Context, email, otp string) error {
	return s.cli.Set(ctx,
		"otp:"+email,
		otp,
		2*time.Minute,
	).Err()
}

func (s *OTPStore) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	storedOTP, err := s.cli.Get(ctx, "otp:"+email).Result()
	if err != nil {
		return false, err
	}
	return storedOTP == otp, nil
}
