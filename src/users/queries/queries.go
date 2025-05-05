package userQueries

import (
	"context"
	"ride-sharing/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	PhoneExists(ctx context.Context, phone string) (bool, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Implement EmailExists
func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).
		Error
	return count > 0, err
}

// Implement PhoneExists
func (r *userRepository) PhoneExists(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("phone = ?", phone).
		Count(&count).
		Error
	return count > 0, err
}

// Implement CreateUser
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
