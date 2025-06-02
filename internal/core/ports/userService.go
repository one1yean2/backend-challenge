package ports

import (
	"context"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/core/domain"
)

type UserService interface {
	Register(ctx context.Context, user domain.User, config config.Container) (string, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id string, user domain.User) error
	DeleteUser(ctx context.Context, id string) error
	LogTotalUser(ctx context.Context)
}
