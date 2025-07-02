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
	classLevelRepo := repository.NewClassLevelRepository(db)
	classRepo := repository.NewClassRepository(db)
	settingRepo := repository.NewSettingRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	periodRepo := repository.NewPeriodRepository(db)
	billRepo := repository.NewBillRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Service
	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey)
	userService := service.NewUserService(userRepo)
	classLevelService := service.NewClassLevelService(classLevelRepo)
	classService := service.NewClassService(classRepo)
	settingService := service.NewSettingService(settingRepo, db)
	studentService := service.NewStudentService(studentRepo, userRepo, db)
	periodService := service.NewPeriodService(periodRepo)
	billService := service.NewBillService(billRepo)
	reportService := service.NewReportService(reportRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService, userService)
	adminHandler := handler.NewAdminHandler(userService, classLevelService, classService, settingService)
	treasurerHandler := handler.NewTreasurerHandler(studentService, periodService, billService, reportService)

	router := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))
	apiRouter := handler.NewRouter(router, authHandler, adminHandler, treasurerHandler, cfg.JWTSecretKey)
	apiRouter.SetupRoutes()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
