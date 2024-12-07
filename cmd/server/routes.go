package main

import (
	"golang-tes/config"
	"golang-tes/internal/delivery/http/attendance"
	"golang-tes/internal/delivery/http/user"
	"golang-tes/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRoutes(router *gin.Engine, cfg *config.Config, userHandler *user.UserHandler, attendanceHandler *attendance.AttendanceHandler) {
	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.POST("/api/users/register", userHandler.Register)
	router.POST("/api/users/login", userHandler.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(authMiddleware.AuthRequired())
	{
		// User routes
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)

		// Attendance routes
		protected.POST("/attendance", attendanceHandler.MarkAttendance)
		protected.GET("/attendance", attendanceHandler.GetAttendance)
		protected.GET("/attendance/user", attendanceHandler.GetUserAttendance)
	}
}
