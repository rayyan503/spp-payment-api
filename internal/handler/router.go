package handler

import (
	"github.com/hiuncy/spp-payment-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine       *gin.Engine
	authHandler  AuthHandler
	adminHandler AdminHandler
	jwtSecretKey string
}

func NewRouter(engine *gin.Engine, authHandler AuthHandler, adminHandler AdminHandler, jwtSecretKey string) *Router {
	return &Router{engine, authHandler, adminHandler, jwtSecretKey}
}

func (r *Router) SetupRoutes() {
	api := r.engine.Group("/api/v1")

	// Auth routes
	api.POST("/login", r.authHandler.Login)
	api.GET("/me", middleware.AuthMiddleware(r.jwtSecretKey, "admin", "bendahara", "siswa"), r.authHandler.GetMe)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(r.jwtSecretKey, "admin"))
	{
		admin.POST("/users", r.adminHandler.CreateUser)
		admin.GET("/users", r.adminHandler.FindAllUsers)
		admin.GET("/users/:id", r.adminHandler.FindUserByID)
		admin.PUT("/users/:id", r.adminHandler.UpdateUser)
		admin.DELETE("/users/:id", r.adminHandler.DeleteUser)
		admin.POST("/class-levels", r.adminHandler.CreateClassLevel)
		admin.GET("/class-levels", r.adminHandler.FindAllClassLevels)
		admin.GET("/class-levels/:id", r.adminHandler.FindClassLevelByID)
		admin.PUT("/class-levels/:id", r.adminHandler.UpdateClassLevel)
		admin.DELETE("/class-levels/:id", r.adminHandler.DeleteClassLevel)
		admin.GET("/settings", r.adminHandler.FindAllSettings)
		admin.PUT("/settings", r.adminHandler.UpdateSettings)
	}

	// Treasurer routes
	treasurer := api.Group("/treasurer")
	treasurer.Use(middleware.AuthMiddleware(r.jwtSecretKey, "bendahara"))
	{
		// treasurer.GET("/students", ...)
	}

	// Student routes
	student := api.Group("/student")
	student.Use(middleware.AuthMiddleware(r.jwtSecretKey, "siswa"))
	{
		// student.GET("/bills", ...)
	}
}
