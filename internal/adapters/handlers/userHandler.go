package handlers

import (
	"context"
	"net/http"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/core/domain"
	"one1-be-chal/internal/core/ports"

	"github.com/labstack/echo"
)

type HttpUserHandler struct {
	service ports.UserService
	config  *config.Container
}

func NewHttpUserHandler(service ports.UserService, config *config.Container) *HttpUserHandler {
	return &HttpUserHandler{
		service: service,
		config:  config,
	}
}

func (u *HttpUserHandler) Register(c echo.Context) error {
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(
			echo.ErrBadRequest.Code,
			echo.Map{"error": err.Error()},
		)
	}
	jwt, err := u.service.Register(context.Background(), user, *u.config)
	if err != nil {
		return c.JSON(
			echo.ErrInternalServerError.Code,
			echo.Map{"error": err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		echo.Map{"jwToken": jwt},
	)
}

func (u *HttpUserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := u.service.GetUserByID(context.Background(), id)
	if err != nil {
		return c.JSON(
			echo.ErrInternalServerError.Code,
			echo.Map{"error": err.Error()},
		)
	}
	return c.JSON(http.StatusOK, user)
}

func (u *HttpUserHandler) GetAllUsers(c echo.Context) error {
	users, err := u.service.GetAllUsers(context.Background())
	if err != nil {
		return c.JSON(
			echo.ErrInternalServerError.Code,
			echo.Map{"error": err.Error()},
		)
	}
	return c.JSON(http.StatusOK, users)
}

func (u *HttpUserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(
			echo.ErrBadRequest.Code,
			echo.Map{"error": err.Error()},
		)
	}

	if err := u.service.UpdateUser(context.Background(), id, user); err != nil {
		return c.JSON(
			echo.ErrInternalServerError.Code,
			echo.Map{"error": err.Error()},
		)
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User updated successfully"})
}

func (u *HttpUserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if err := u.service.DeleteUser(context.Background(), id); err != nil {
		return c.JSON(
			echo.ErrInternalServerError.Code,
			echo.Map{"error": err.Error()},
		)
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}
