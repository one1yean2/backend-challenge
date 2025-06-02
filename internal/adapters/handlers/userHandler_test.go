package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/core/domain"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, user domain.User, config config.Container) (string, error) {
	args := m.Called(ctx, user, config)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, user domain.User) error {
	return m.Called(ctx, id, user).Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockUserService) LogTotalUser(ctx context.Context) {
	m.Called(ctx)
}

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	e.Validator = NewRequestValidator()
	mockService := new(MockUserService)
	mockConfig := &config.Container{JWT: &config.JWT{SecretKey: []byte("secret")}}
	handler := NewHttpUserHandler(mockService, mockConfig)

	userJSON := `{
		"name": "Wannueng Krub",
		"email" : "one123@gmail.com",
		"password" : "123456Test!"    
	}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("Register", mock.Anything, mock.Anything, mock.Anything).Return("token", nil)

	err := handler.Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRegisterUserBadRequest(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	mockConfig := &config.Container{}
	handler := NewHttpUserHandler(mockService, mockConfig)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserByID(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	mockConfig := &config.Container{}
	handler := NewHttpUserHandler(mockService, mockConfig)

	expectedUser := domain.User{ID: "123", Name: "One1 Yean", Email: "test@gmail.com"}
	mockService.On("GetUserByID", mock.Anything, "123").Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	err := handler.GetUserByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateUser(t *testing.T) {
	e := echo.New()
	e.Validator = NewRequestValidator()
	mockService := new(MockUserService)
	mockConfig := &config.Container{}
	handler := NewHttpUserHandler(mockService, mockConfig)

	userJSON := `{"name": "One3", "email": "test@gmail.com"}`
	req := httptest.NewRequest(http.MethodPut, "/user/123", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	mockService.On("UpdateUser", mock.Anything, "123", mock.Anything).Return(nil)

	err := handler.UpdateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteUser(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	mockConfig := &config.Container{}
	handler := NewHttpUserHandler(mockService, mockConfig)

	req := httptest.NewRequest(http.MethodDelete, "/user/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	mockService.On("DeleteUser", mock.Anything, "123").Return(nil)

	err := handler.DeleteUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
