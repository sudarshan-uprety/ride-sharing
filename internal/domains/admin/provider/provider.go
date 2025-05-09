package provider

// import (
// 	"context"

// 	"ride-sharing/internal/domains/users/repository"
// 	"ride-sharing/internal/pkg/auth"
// )

// type UserProvider struct {
// 	repo repository.UserRepository
// }

// func NewUserProvider(repo repository.UserRepository) auth.UserProvider {
// 	return &UserProvider{repo: repo}
// }

// func (p *UserProvider) GetByID(ctx context.Context, id string) (interface{}, error) {
// 	return p.repo.GetByID(ctx, id)
// }
