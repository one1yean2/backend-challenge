package handlers

import (
	"net/http"
	"net/http/httptest"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/adapters/helpers"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func mockHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func TestJWTMiddleware(t *testing.T) {
	e := echo.New()
	mockConfig := &config.Container{
		JWT: &config.JWT{SecretKey: []byte("secret")},
	}
	middleware := JWTMiddleware(mockConfig)

	validToken, _ := helpers.GenerateJWT("123", "one1", "test@gmail.com", *mockConfig)

	tests := []struct {
		Name           string
		Authorization  string
		ExpectedStatus int
	}{
		{
			Name:           "Valid Token",
			Authorization:  "Bearer " + validToken,
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Missing Token",
			Authorization:  "",
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Name:           "Invalid Token Format",
			Authorization:  "InvalidTokenString",
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Name:           "Malformed Token",
			Authorization:  "Bearer invalid.token.string",
			ExpectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", test.Authorization)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := middleware(mockHandler)
			err := handler(c)

			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedStatus, rec.Code)
		})
	}
}
