package handler

import (
	"github.com/hiuncy/spp-payment-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine           *gin.Engine
	authHandler      AuthHandler
	adminHandler     AdminHandler
	treasurerHandler TreasurerHandler
	studentHandler   StudentHandler
	midtransHandler  MidtransHandler
	jwtSecretKey     string
}

func NewRouter(engine *gin.Engine, authHandler AuthHandler, adminHandler AdminHandler, treasurerHandler TreasurerHandler, studentHandler StudentHandler, midtransHandler MidtransHandler, jwtSecretKey string) *Router {
	return &Router{engine, authHandler, adminHandler, treasurerHandler, studentHandler, midtransHandler, jwtSecretKey}
}

func (r *Router) SetupRoutes() {
	api := r.engine.Group("/api/v1")

	// Auth routes
	api.POST("/login", r.authHandler.Login)
	api.GET("/me", middleware.AuthMiddleware(r.jwtSecretKey, "admin", "bendahara", "siswa"), r.authHandler.GetMe)

	// Midtrans routes
	api.POST("/payments/midtrans-notification", r.midtransHandler.HandleNotification)

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
		admin.POST("/classes", r.adminHandler.CreateClass)
		admin.GET("/classes/:id", r.adminHandler.FindClassByID)
		admin.PUT("/classes/:id", r.adminHandler.UpdateClass)
		admin.DELETE("/classes/:id", r.adminHandler.DeleteClass)
		admin.GET("/settings", r.adminHandler.FindAllSettings)
		admin.PUT("/settings", r.adminHandler.UpdateSettings)
	}

	// Treasurer routes
	treasurer := api.Group("/treasurer")
	treasurer.Use(middleware.AuthMiddleware(r.jwtSecretKey, "bendahara", "admin"))
	{
		treasurer.GET("/classes", r.adminHandler.FindAllClasses)
		treasurer.POST("/students", r.treasurerHandler.CreateStudent)
		treasurer.GET("/students", r.treasurerHandler.FindAllStudents)
		treasurer.GET("/students/:id", r.treasurerHandler.FindStudentByID)
		treasurer.PUT("/students/:id", r.treasurerHandler.UpdateStudent)
		treasurer.DELETE("/students/:id", r.treasurerHandler.DeleteStudent)
		treasurer.POST("/periods", r.treasurerHandler.CreatePeriod)
		treasurer.GET("/periods", r.treasurerHandler.FindAllPeriods)
		treasurer.GET("/periods/:id", r.treasurerHandler.FindPeriodByID)
		treasurer.PUT("/periods/:id", r.treasurerHandler.UpdatePeriod)
		treasurer.DELETE("/periods/:id", r.treasurerHandler.DeletePeriod)
		treasurer.POST("/periods/:id/generate-bills", r.treasurerHandler.GenerateBills)
		treasurer.GET("/bills", r.treasurerHandler.FindAllBills)
		treasurer.GET("/bills/:id", r.treasurerHandler.FindBillByID)
		treasurer.PUT("/bills/:id", r.treasurerHandler.UpdateBill)
		treasurer.DELETE("/bills/:id", r.treasurerHandler.DeleteBill)
		reports := treasurer.Group("/reports")
		{
			reports.GET("/per-student", r.treasurerHandler.GetLaporanSiswa)
			reports.GET("/per-class", r.treasurerHandler.GetLaporanKelas)
			reports.GET("/overall", r.treasurerHandler.GetLaporanKeseluruhan)
		}
	}

	// Student routes
	student := api.Group("/student")
	student.Use(middleware.AuthMiddleware(r.jwtSecretKey, "siswa"))
	{
		student.GET("/profile", r.studentHandler.GetProfile)
		student.GET("/bills", r.studentHandler.FindMyBills)
		student.POST("/bills/:id/pay", r.studentHandler.InitiatePayment)
		student.GET("/payment-history", r.studentHandler.GetPaymentHistory)
	}
}
