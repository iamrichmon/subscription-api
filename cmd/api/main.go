package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamrichmon/subscription-api/internal/config"
	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/repository"
)

func main() {

	cfg := config.Load()

	db := repository.Connect(cfg)

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to auto-migrate database: %v", err)
	}

	r := gin.Default()

	r.Run()
}
