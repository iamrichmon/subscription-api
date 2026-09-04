package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamrichmon/subscription-api/internal/config"
	"github.com/iamrichmon/subscription-api/internal/handler"
	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/repository"
	"github.com/iamrichmon/subscription-api/internal/service"
)

func main() {

	cfg := config.Load()

	if cfg.JWTSECRET == "" {
		log.Fatal("JWT secret is required. Please set the JWT_SECRET environment variable.")
	}

	db := repository.Connect(cfg)

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to auto-migrate database: %v", err)
	}

	// wire dependencies

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo, cfg.JWTSECRET)

	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/auth/register", userHandler.Register)
		v1.POST("/auth/login", userHandler.Login)
	}

	r.Run()
}
