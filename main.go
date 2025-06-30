package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/hiuncy/spp-payment-api/internal/config"
	"github.com/hiuncy/spp-payment-api/internal/handler"
	"github.com/hiuncy/spp-payment-api/internal/repository"
	"github.com/hiuncy/spp-payment-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Repository
	userRepo := repository.NewUserRepository(db)

	// Service
	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey)
	userService := service.NewUserService(userRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService, userService)
	adminHandler := handler.NewAdminHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())
	apiRouter := handler.NewRouter(router, authHandler, adminHandler, cfg.JWTSecretKey)
	apiRouter.SetupRoutes()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
