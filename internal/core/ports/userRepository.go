package ports

import (
	"context"
	"one1-be-chal/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id string, user bson.M) error
	DeleteUser(ctx context.Context, id string) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserCount(ctx context.Context) (int64, error)
}
