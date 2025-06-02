package services

import (
	"context"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, user domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id string, updateFields bson.M) error {
	return m.Called(ctx, id, updateFields).Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepository) GetUserCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := domain.User{
		Email:    "test@gmail.com",
		Password: "passwordkrub",
		Name:     "One1 yean",
	}

	mockRepo.On("GetUserByEmail", mock.Anything, user.Email).Return(nil, nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	jwtToken, err := service.Register(
		context.Background(),
		user,
		config.Container{
			JWT: &config.JWT{
				SecretKey: []byte("secret"),
			},
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)
}

func TestRegisterExistingUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	existingUser := &domain.User{
		Email: "test@gmail.com",
	}

	mockRepo.On("GetUserByEmail", mock.Anything, existingUser.Email).Return(existingUser, nil)

	_, err := service.Register(context.Background(), *existingUser, config.Container{})

	assert.Error(t, err)
	assert.Equal(t, "email already exist", err.Error())
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUser := domain.User{ID: "123", Name: "One1 yean", Email: "test@gmail.com"}
	mockRepo.On("GetUserByID", mock.Anything, "123").Return(expectedUser, nil)

	user, err := service.GetUserByID(context.Background(), "123")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	tests := []struct {
		Name        string
		User        domain.User
		ExpectError bool
	}{
		{
			Name: "valid update",
			User: domain.User{
				Email: "test@gmail.com",
				Name:  "One1 yean",
			},
			ExpectError: false,
		},
		{
			Name: "valid update 1 field provide",
			User: domain.User{
				Email: "test@gmail.com",
				Name:  "",
			},
			ExpectError: false,
		},
		{
			Name: "invalid update",
			User: domain.User{
				Email: "",
				Name:  "",
			},
			ExpectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mockRepo.On("GetUserByID", mock.Anything, "123").Return(domain.User{ID: "123"}, nil)
			mockRepo.On("GetUserByEmail", mock.Anything, test.User.Email).Return(nil, nil)
			mockRepo.On("UpdateUser", mock.Anything, "123", mock.Anything).Return(nil)

			err := service.UpdateUser(context.Background(), "123", test.User)

			if test.ExpectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		},
		)
	}

}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("GetUserByID", mock.Anything, "123").Return(domain.User{ID: "123"}, nil)
	mockRepo.On("DeleteUser", mock.Anything, "123").Return(nil)

	err := service.DeleteUser(context.Background(), "123")

	assert.NoError(t, err)
}

func TestLogTotalUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("GetUserCount", mock.Anything).Return(int64(100), nil)

	go service.LogTotalUser(context.Background())

	time.Sleep(12 * time.Second)

	mockRepo.AssertCalled(t, "GetUserCount", mock.Anything)
}
