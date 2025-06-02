package services

import (
	"context"
	"errors"
	"log"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/adapters/helpers"
	"one1-be-chal/internal/core/domain"
	"one1-be-chal/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	UserRepository ports.UserRepository
}

func NewUserService(userRepository ports.UserRepository) ports.UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}

func (s *UserServiceImpl) Register(
	ctx context.Context,
	user domain.User,
	config config.Container,
) (string, error) {
	existUser, err := s.UserRepository.GetUserByEmail(ctx, user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}
	if existUser != nil {
		return "", errors.New("email already exist")
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	user.ID = uuid.NewString()
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	if err := s.UserRepository.Save(ctx, user); err != nil {
		return "", err
	}

	jwToken, err := helpers.GenerateJWT(user.ID, user.Name, user.Email, config)
	if err != nil {
		return "", err
	}

	return jwToken, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return s.UserRepository.GetUserByID(ctx, id)
}

func (s *UserServiceImpl) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.UserRepository.GetAllUsers(ctx)
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, user domain.User) error {
	updateFields := bson.M{}
	if err := user.ValidateEmailAndName(); err != nil {
		return err
	}

	_, err := s.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Email != "" {
		existUser, _ := s.UserRepository.GetUserByEmail(ctx, user.Email)
		if existUser != nil && existUser.ID != id {
			return errors.New("email already exist")
		}
		updateFields["email"] = user.Email
	}
	if user.Name != "" {
		updateFields["name"] = user.Name
	}

	return s.UserRepository.UpdateUser(ctx, id, updateFields)
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	_, err := s.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	return s.UserRepository.DeleteUser(ctx, id)
}
func (s *UserServiceImpl) LogTotalUser(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		count, err := s.UserRepository.GetUserCount(context.Background())
		if err != nil {
			log.Println("Error getting user count:", err)
			continue
		}
		log.Println("Total number of users:", count)
	}
}
