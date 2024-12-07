package main

// @title           Absensi Karyawan API
// @version         1.0
// @description     API Server for Employee Attendance System
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@attendance.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"log"

	"golang-tes/config"
	"golang-tes/internal/delivery/http/attendance"
	"golang-tes/internal/delivery/http/user"
	"golang-tes/internal/repository"
	"golang-tes/internal/usecase"
	"golang-tes/pkg/db"

	_ "golang-tes/docs" // This will be auto-generated

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.NewDatabase(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewMySQLUserRepository(database)
	attendanceRepo := repository.NewMySQLAttendanceRepository(database)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, userRepo)

	// Initialize handlers
	userHandler := user.NewUserHandler(userUsecase)
	attendanceHandler := attendance.NewAttendanceHandler(attendanceUsecase)

	// Initialize Gin router with CORS middleware
	router := gin.Default()
	router.Use(corsMiddleware())

	// Setup routes
	setupRoutes(router, cfg, userHandler, attendanceHandler)

	// Start server
	log.Printf("Server starting on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
