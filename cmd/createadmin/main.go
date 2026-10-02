package main

import (
	"log"
	"os"

	"github.com/iamrichmon/subscription-api/internal/config"
	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/repository"
	"github.com/iamrichmon/subscription-api/internal/service"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("usage: go run ./cmd/createadmin <Richmon> <richmonalbon@gmail.com> <superlongpassword12>")
	}

	cfg := config.Load()
	db := repository.Connect(cfg)

	if err := db.AutoMigrate(&model.Admin{}); err != nil {
		log.Fatalf("failed to auto-migrate database: %v", err)
	}

	adminRepo := repository.NewAdminRepository(db)
	adminService := service.NewAdminService(adminRepo, cfg.JWTSECRET)

	admin, err := adminService.RegisterAdmin(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("admin created: id=%d email=%s role=%s", admin.ID, admin.Email, admin.Role)
}
