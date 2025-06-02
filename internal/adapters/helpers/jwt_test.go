package helpers

import (
	"one1-be-chal/internal/adapters/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	mockConfig := config.Container{
		JWT: &config.JWT{SecretKey: []byte("secret")},
	}

	id := "123"
	name := "One1 yean"
	email := "test@gmail.com"

	token, err := GenerateJWT(id, name, email, mockConfig)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseJWT(t *testing.T) {
	mockConfig := config.Container{
		JWT: &config.JWT{SecretKey: []byte("secret")},
	}

	id := "123"
	name := "One1 yean"
	email := "test@gmail.com"

	token, err := GenerateJWT(id, name, email, mockConfig)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseJWT(token, mockConfig)
	assert.NoError(t, err)
	assert.Equal(t, id, claims.ID)
	assert.Equal(t, name, claims.Name)
	assert.Equal(t, email, claims.Email)
}

func TestParseJWTInvalidToken(t *testing.T) {
	mockConfig := config.Container{
		JWT: &config.JWT{SecretKey: []byte("secret")},
	}

	claims, err := ParseJWT("invalid.token.string", mockConfig)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
