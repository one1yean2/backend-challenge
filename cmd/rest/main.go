package main

import (
	"context"
	"log"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/adapters/handlers"
	"one1-be-chal/internal/adapters/storages/mongo"
	"one1-be-chal/internal/adapters/storages/mongo/repositories"
	"one1-be-chal/internal/core/services"
	"os"
)

func main() {
	app := handlers.EchoMiddleware()
	app.Validator = handlers.NewRequestValidator()
	config := config.New()
	ctx := context.Background()
	userDBClient, err := mongo.New(ctx, config.UserDB)
	if err != nil {
		log.Printf("Error initializing MongoDB connection: %v\n", err)
		os.Exit(1)
	}
	defer userDBClient.Close(ctx)
	userDB := userDBClient.Client.Database("backend-challenge")

	userRepo := repositories.NewUserRepository(userDB, "users")
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewHttpUserHandler(userService, config)

	app.POST("/register", userHandler.Register)
	app.GET("/user/:id", userHandler.GetUserByID, handlers.JWTMiddleware(config))
	app.GET("/user", userHandler.GetAllUsers, handlers.JWTMiddleware(config))
	app.PATCH("/user/:id", userHandler.UpdateUser, handlers.JWTMiddleware(config))
	app.DELETE("/user/:id", userHandler.DeleteUser, handlers.JWTMiddleware(config))

	go userService.LogTotalUser(ctx)

	app.Logger.Fatal(app.Start(":8080"))
}
