package redis

import (
	"context"
	"ride-sharing/internal/pkg/errors"
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
	key := "otp:" + email

	// Check if the OTP already exists
	exists, err := s.cli.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		return errors.NewConflictError("OTP already exists for this email")
	}

	// Set the new OTP with 2-minute expiration
	return s.cli.Set(ctx, key, otp, 2*time.Minute).Err()
}

func (s *OTPStore) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	storedOTP, err := s.cli.Get(ctx, "otp:"+email).Result()
	if err != nil {
		return false, err
	}
	return storedOTP == otp, nil
}
