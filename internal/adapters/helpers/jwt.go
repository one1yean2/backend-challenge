package helpers

import (
	"one1-be-chal/internal/adapters/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(id, name, email string, config config.Container) (string, error) {
	claims := Claims{
		ID:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //Use HMAC (HS256) with a secret key.
	return token.SignedString(config.JWT.SecretKey)
}

func ParseJWT(tokenStr string, config config.Container) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return config.JWT.SecretKey, nil
		},
	)
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
