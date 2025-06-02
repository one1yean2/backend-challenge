package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Container struct {
	UserDB *UserDB
	JWT    *JWT
}

type UserDB struct {
	URI string
}

type JWT struct {
	SecretKey []byte
}

func New() *Container {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return &Container{
		UserDB: &UserDB{
			URI: os.Getenv("MONGODB_URI"),
		},
		JWT: &JWT{
			SecretKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		},
	}
}
