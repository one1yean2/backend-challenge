package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mypassword123!"
	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "correctPass"
	hashedPassword, _ := HashPassword(password)

	t.Run("Valid password match", func(t *testing.T) {
		match := CheckPasswordHash(password, hashedPassword)
		assert.True(t, match)
	})

	t.Run("Invalid password mismatch", func(t *testing.T) {
		match := CheckPasswordHash("wrongPass", hashedPassword)
		assert.False(t, match)
	})
}
